package scorumgo

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/scorum/scorum-go/apis/account_history"
	"github.com/scorum/scorum-go/apis/betting"
	"github.com/scorum/scorum-go/apis/blockchain_history"
	"github.com/scorum/scorum-go/apis/chain"
	"github.com/scorum/scorum-go/apis/database"
	"github.com/scorum/scorum-go/apis/network_broadcast"
	"github.com/scorum/scorum-go/caller"
	"github.com/scorum/scorum-go/key"
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

	getReferenceBlock getReferenceBlock
}

type reference struct {
	time   types.Time
	number uint16
	prefix uint32
}

type getReferenceBlock func(ctx context.Context) (reference, error)

type Option func(c *Client)

func WithHeadBlockReferenceSign() Option {
	return func(c *Client) {
		c.getReferenceBlock = c.getHeadBlockReference
	}
}

// NewClient creates a new RPC client that use the given CallCloser internally.
func NewClient(cc caller.CallCloser, opts ...Option) *Client {
	client := &Client{
		cc: cc,
	}
	client.Database = database.NewAPI(client.cc)
	client.Chain = chain.NewAPI(client.cc)
	client.AccountHistory = account_history.NewAPI(client.cc)
	client.NetworkBroadcast = network_broadcast.NewAPI(client.cc)
	client.BlockchainHistory = blockchain_history.NewAPI(client.cc)
	client.Betting = betting.NewAPI(client.cc)
	client.getReferenceBlock = client.getLastIrreversibleBlockReference

	for _, opt := range opts {
		opt(client)
	}

	return client
}

// Close should be used to close the client when no longer needed.
// It simply calls Close() on the underlying CallCloser.
func (client *Client) Close() error {
	return client.cc.Close()
}

func (client *Client) BroadcastTransactionSynchronous(ctx context.Context, chainID []byte, operations []types.Operation, keys ...*key.PrivateKey) (*network_broadcast.BroadcastResponse, error) {
	stx, err := client.createSignedTransaction(ctx, chainID, operations, keys...)
	if err != nil {
		return nil, err
	}
	return client.NetworkBroadcast.BroadcastTransactionSynchronous(ctx, stx.Transaction)
}

func (client *Client) BroadcastTransaction(ctx context.Context, chainID []byte, operations []types.Operation, keys ...*key.PrivateKey) (string, error) {
	stx, err := client.createSignedTransaction(ctx, chainID, operations, keys...)
	if err != nil {
		return "", err
	}

	id, _ := stx.ID()

	return hex.EncodeToString(id), client.NetworkBroadcast.BroadcastTransaction(ctx, stx.Transaction)
}

func (client *Client) createSignedTransaction(ctx context.Context, chainID []byte, operations []types.Operation, keys ...*key.PrivateKey) (*sign.SignedTransaction, error) {
	refBlock, err := client.getReferenceBlock(ctx)
	if err != nil {
		return nil, err
	}

	expiration := refBlock.time.Add(10 * time.Minute)
	stx := sign.NewSignedTransaction(&types.Transaction{
		Operations:     operations,
		RefBlockNum:    refBlock.number,
		RefBlockPrefix: refBlock.prefix,
		Expiration:     &types.Time{Time: &expiration},
	})

	if err = stx.Sign(chainID, keys...); err != nil {
		return nil, fmt.Errorf("sign transaction: %w", err)
	}

	return stx, nil
}

func (client *Client) getLastIrreversibleBlockReference(ctx context.Context) (reference, error) {
	props, err := client.Chain.GetChainProperties(ctx)
	if err != nil {
		return reference{}, fmt.Errorf("get chainID properties: %w", err)
	}

	block, err := client.BlockchainHistory.GetBlock(ctx, props.LastIrreversibleBlockNumber)
	if err != nil {
		return reference{}, fmt.Errorf("blockchain history get block: %w", err)
	}

	number, prefix, err := sign.ParseBlockID(block.Previous)
	if err != nil {
		return reference{}, fmt.Errorf("ref block prefix: %w", err)
	}

	return reference{
		time:   props.Time,
		number: number,
		prefix: prefix,
	}, nil
}

func (client *Client) getHeadBlockReference(ctx context.Context) (reference, error) {
	props, err := client.Chain.GetChainProperties(ctx)
	if err != nil {
		return reference{}, fmt.Errorf("get chainID properties: %w", err)
	}

	number, prefix, err := sign.ParseBlockID(props.HeadBlockID)
	if err != nil {
		return reference{}, fmt.Errorf("ref block prefix: %w", err)
	}

	return reference{
		time:   props.Time,
		number: number,
		prefix: prefix,
	}, nil
}
