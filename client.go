package rpc

import (
	"github.com/scorum/scorum-go/apis/account_history"
	"github.com/scorum/scorum-go/apis/database"
	"github.com/scorum/scorum-go/apis/network_broadcast"
	"github.com/scorum/scorum-go/caller"
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
}

// NewClient creates a new RPC client that use the given CallCloser internally.
func NewClient(cc caller.CallCloser) *Client {
	client := &Client{cc: cc}
	client.Database = database.NewAPI(client.cc)
	client.AccountHistory = account_history.NewAPI(client.cc)
	client.NetworkBroadcast = network_broadcast.NewAPI(client.cc)
	return client
}

// Close should be used to close the client when no longer needed.
// It simply calls Close() on the underlying CallCloser.
func (client *Client) Close() error {
	return client.cc.Close()
}
