package main

import (
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sort"
	"sync"
	"syscall"
	"time"

	"github.com/scorum/scorum-go"
	"github.com/scorum/scorum-go/apis/account_history"
	"github.com/scorum/scorum-go/sign"
	"github.com/scorum/scorum-go/transport/http"
	"github.com/scorum/scorum-go/types"
	"github.com/shopspring/decimal"
)

const (
	testNet        = "https://testnet.scorum.com"
	paymentAccount = "sheldon"
	paymentWIF     = "5JwWJ2m2jGG9RPcpDix5AvkDzQZJoZvpUQScsDzzXWAKMs8Q6jH"
)

var chain = sign.TestChain

var (
	// mapping between deposit and balance
	balances map[string]decimal.Decimal
	// blockchain client
	client *scorumgo.Client
	// sync balances changes
	mutex sync.Mutex
	// history seq cursor
	seq uint32
)

func main() {
	// accounts and their balances, in a real-world app should be loaded from database
	balances = map[string]decimal.Decimal{
		"noelle":   decimal.NewFromFloat(0),
		"gina":     decimal.NewFromFloat(0),
		"margy":    decimal.NewFromFloat(0),
		"leonarda": decimal.NewFromFloat(0),
	}

	// create a blockchain rcp client
	client = scorumgo.NewClient(http.NewTransport(testNet))

	// seq is sequence number of the last processed history item
	seq = 0

	// listen for incoming payments
	go Monitor()
	// make payouts
	go Payout()

	// wait for signal to exit
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	<-exit
}

// Monitor payment account transactions
func Monitor() {
	for {
		var recentSeq uint32

		recent, err := client.AccountHistory.GetAccountHistory(paymentAccount, -1, 0)
		if err != nil {
			log.Printf("failed to get recent account history: %s\n", err)
			goto Step
		}

		// recent contain only one item, take it sequence number
		for recentSeq = range recent {
			break
		}

		if recentSeq > seq {
			limit := recentSeq - seq - 1
			// retrieve transactions created since the last step
			history, err := client.AccountHistory.GetAccountHistory(paymentAccount, int32(recentSeq), int32(limit))
			if err != nil {
				log.Printf("failed to get recent account history: %s\n", err)
				goto Step
			}

			mutex.Lock()
			processHistory(history)
			seq = recentSeq
			mutex.Unlock()
		}

	Step:
		<-time.After(10 * time.Second)
	}
}

func processHistory(history account_history.AccountHistory) {
	// order keys (seq numbers), to process transaction in chronological order
	keys := make([]uint32, len(history))
	index := 0
	for k := range history {
		keys[index] = k
		index++
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[j] > keys[i]
	})

	// process the transfers
	for _, seq := range keys {
		trx := history[seq]
		for _, op := range trx.Operations {
			switch body := op.(type) {
			case *types.TransferOperation:
				log.Printf("transfer: %+v\n", op)
				processTransfer(seq, trx, body)
			default:
				log.Printf("operation %s: %+v\n", op.Type(), op)
			}
		}
	}
}

func processTransfer(seq uint32, trx *types.OperationObject, op *types.TransferOperation) {
	// transaction memo is a deposit
	deposit := op.Memo
	balance, ok := balances[deposit]
	if !ok {
		//unrecognized deposit, save it somewhere for later review
		log.Printf("unrecognized deposit: `%s`\n", op.Memo)
		return
	}

	// increase deposit balance
	balances[deposit] = balance.Add(op.Amount.Decimal())

	log.Printf("%d %+v transfer from %s to deposit %s processed\n", seq, trx, op.From, deposit)
}

// makes random payout to the existing accounts
func Payout() {
	for {
		mutex.Lock()

		// get a random deposit
		deposits := make([]string, len(balances))
		idx := 0
		for d := range balances {
			deposits[idx] = d
			idx++
		}
		deposit := deposits[rand.Intn(len(deposits))]

		amount, _ := types.AssetFromString("0.000001 SCR")

		// check the balance
		if balances[deposit].LessThan(amount.Decimal()) {
			log.Printf("not enough SCR to transfer to %s\n", deposit)
		} else {
			// broadcast the transfer operation
			_, err := client.Broadcast(chain, []string{paymentWIF}, &types.TransferOperation{
				From:   paymentAccount,
				To:     deposit,
				Amount: *amount,
				Memo:   "payout from", //specify needed memo
			})

			if err != nil {
				log.Printf("failed to transfer %s to %s", amount, deposit)
			} else {
				// decrease deposit balance
				balances[deposit] = balances[deposit].Sub(amount.Decimal())
			}
		}

		mutex.Unlock()

		<-time.After(time.Second * 5)
	}
}
