package scorumgo

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/scorum/scorum-go/apis/account_history"
	"github.com/scorum/scorum-go/apis/betting"
	"github.com/scorum/scorum-go/apis/blockchain_history"
	"github.com/scorum/scorum-go/apis/chain"
	"github.com/scorum/scorum-go/apis/database"
	"github.com/scorum/scorum-go/apis/network_broadcast"
	"github.com/scorum/scorum-go/caller"
	"github.com/scorum/scorum-go/sign"
	"github.com/scorum/scorum-go/types"
)

// Client can be used to access Scorum remote APIs.
//
// There is a public field for every Scorum API available,
// e.g. Client.Database corresponds to database_api.
type Client struct {
	cc caller.CallCloser

	// Database represents database_api
	Database *database.API

	// AccountHistory represents account_history_api
	AccountHistory *account_history.API

	// NetworkBroadcast represents network_broadcast_api
	NetworkBroadcast *network_broadcast.API

	// BlockchainHistory represents blockchain_history_api
	BlockchainHistory *blockchain_history.API

	// Betting represents betting_api
	Betting *betting.API

	// Chain represents chain_api
	Chain *chain.API
}

// NewClient creates a new RPC client that use the given CallCloser internally.
func NewClient(cc caller.CallCloser) *Client {
	client := &Client{cc: cc}
	client.Database = database.NewAPI(client.cc)
	client.Chain = chain.NewAPI(client.cc)
	client.AccountHistory = account_history.NewAPI(client.cc)
	client.NetworkBroadcast = network_broadcast.NewAPI(client.cc)
	client.BlockchainHistory = blockchain_history.NewAPI(client.cc)
	client.Betting = betting.NewAPI(client.cc)
	return client
}

// Close should be used to close the client when no longer needed.
// It simply calls Close() on the underlying CallCloser.
func (client *Client) Close() error {
	return client.cc.Close()
}

// Broadcast Sign the given operations with the wifs and broadcast them as one transaction
func (client *Client) Broadcast(chain *sign.Chain, wifs []string, operations ...types.Operation) (*network_broadcast.BroadcastResponse, error) {
	return client.BroadcastContext(context.Background(), chain, wifs, operations...)
}

func (client *Client) BroadcastContext(ctx context.Context, chain *sign.Chain, wifs []string, operations ...types.Operation) (*network_broadcast.BroadcastResponse, error) {
	props, err := client.Chain.GetChainProperties(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get dynamic global properties")
	}

	block, err := client.BlockchainHistory.GetBlock(ctx, props.LastIrreversibleBlockNumber)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get block")
	}

	refBlockPrefix, err := sign.RefBlockPrefix(block.Previous)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign block prefix")
	}

	expiration := props.Time.Add(10 * time.Minute)
	stx := sign.NewSignedTransaction(&types.Transaction{
		RefBlockNum:    sign.RefBlockNum(props.LastIrreversibleBlockNumber - 1&0xffff),
		RefBlockPrefix: refBlockPrefix,
		Expiration:     &types.Time{Time: &expiration},
	})

	for _, op := range operations {
		stx.PushOperation(op)
	}

	if err = stx.Sign(wifs, chain); err != nil {
		return nil, errors.Wrap(err, "failed to sign the transaction")
	}

	return client.NetworkBroadcast.BroadcastTransactionSynchronous(ctx, stx.Transaction)
}
