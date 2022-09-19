package scorumgo

import (
	"context"
	"testing"
	"time"

	"github.com/scorum/scorum-go/rpc"
	"github.com/scorum/scorum-go/rpc/protocol"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"

	"github.com/scorum/scorum-go/key"
	"github.com/scorum/scorum-go/sign"
	"github.com/scorum/scorum-go/types"
)

const (
	nodeWSS         = "wss://testnet.scorum.work"
	nodeHTTPS       = "https://testnet.scorum.work"
	mainNetNodeHTTP = "https://prodnet.scorum.com"
)

// test accounts available at https://github.com/scorum/scorum/wiki/Testnet-existent-accounts

func newWebsocketClient(t *testing.T) *Client {
	transport := rpc.NewWebSocketTransport(nodeWSS, websocket.DefaultDialer)
	require.NoError(t, transport.Dial(context.Background()))

	return NewClient(transport)
}

func newHTTPClient() *Client {
	client := NewClient(rpc.NewHTTPTransport(nodeHTTPS))
	return client
}

func newMainNetHTTPClient() *Client {
	client := NewClient(rpc.NewHTTPTransport(mainNetNodeHTTP))
	return client
}

func TestGetConfigViaWS(t *testing.T) {
	client := newWebsocketClient(t)
	defer client.Close()

	config, err := client.Database.GetConfig(context.Background())
	require.NoError(t, err)
	require.Equal(t, "SCR", config.ScorumAddressPrefix)
}

func TestGetDynamicGlobalProperties(t *testing.T) {
	client := newWebsocketClient(t)
	defer client.Close()

	config, err := client.Database.GetDynamicGlobalProperties(context.Background())
	require.NoError(t, err)
	t.Logf("dynamic properties: %+v", config)
}

func TestGetAccounts(t *testing.T) {
	client := newHTTPClient()
	defer client.Close()

	accounts, err := client.Database.GetAccounts(context.Background(), "leonarda", "kristie")
	require.NoError(t, err)

	require.Len(t, accounts, 2)
	require.Equal(t, "leonarda", accounts[0].Name)
	require.Equal(t, "kristie", accounts[1].Name)
}

func TestGetAccountHistory(t *testing.T) {
	client := newHTTPClient()

	history, err := client.AccountHistory.GetAccountHistory(context.Background(), "leonarda", -1, 3)
	require.NoError(t, err)
	require.True(t, len(history) > 0)
	spew.Dump(history)
}

func TestClient_Broadcast_AccountWitnessVoteOperation(t *testing.T) {
	client := newHTTPClient()

	roselle, err := key.PrivateKeyFromString("5JwWJ2m2jGG9RPcpDix5AvkDzQZJoZvpUQScsDzzXWAKMs8Q6jH")
	require.NoError(t, err)

	ops := []types.Operation{
		&types.AccountWitnessVoteOperation{
			Account: "roselle",
			Witness: "scorumwitness1",
			Approve: true,
		},
	}
	_, err = client.BroadcastTransactionSynchronous(context.Background(), sign.TestNetChainID, ops, roselle)
	require.NotNil(t, err)

	perr, ok := err.(*protocol.RPCError)
	require.True(t, ok)
	require.Equal(t, "assert_exception", perr.Data.Name)
	require.Equal(t, int(10), perr.Data.Code)
}

func TestClient_Broadcast_Transfer(t *testing.T) {
	client := newHTTPClient()
	amount, _ := types.AssetFromString("0.000009 SCR")

	azucena, err := key.PrivateKeyFromString("5J7FEcpqc1sZ7ZbKx2kVvBHx2oTjWG2wMU2e2FYX85sGA2qu8KT")
	require.NoError(t, err)
	ops := []types.Operation{
		&types.TransferOperation{
			From:   "azucena",
			To:     "leonarda",
			Amount: *amount,
			Memo:   "1",
		},
	}
	resp, err := client.BroadcastTransactionSynchronous(context.Background(), sign.TestNetChainID, ops, azucena)
	require.NoError(t, err)
	require.NotEmpty(t, resp.ID)
	require.NotEmpty(t, resp.BlockNum)
	require.False(t, resp.Expired)
}

func TestSetBlockAppliedCallback(t *testing.T) {
	client := newWebsocketClient(t)
	defer client.Close()

	var called bool
	err := client.Database.SetBlockAppliedCallback(func(block *types.BlockHeader, err error) {
		t.Log("block:", block, "error", err)
		called = true
	})
	require.NoError(t, err)
	time.Sleep(10 * time.Second)
	require.True(t, called)
}

func TestAccountUpdateOperation(t *testing.T) {
	t.Skip()
	client := newHTTPClient()

	blockIDWithAcountUpdateOp := uint32(1799)

	block, err := client.BlockchainHistory.GetBlock(context.Background(), blockIDWithAcountUpdateOp)

	require.NoError(t, err)
	require.Len(t, block.Transactions, 1)
	require.Len(t, block.Transactions[0].Operations, 1)

	op := block.Transactions[0].Operations[0]
	require.Equal(t, op.Type(), types.AccountUpdateOpType)

	accUpdOpt, ok := op.(*types.AccountUpdateOperation)
	require.True(t, ok)
	require.Equal(t, accUpdOpt.Account, "lizzette")

	require.EqualValues(t, accUpdOpt.Active.WeightThreshold, 1)
	require.Len(t, accUpdOpt.Active.AccountAuths, 0)
	require.Len(t, accUpdOpt.Active.KeyAuths, 2)

	v, ok := accUpdOpt.Active.KeyAuths.Get("SCR6W2AjgDsuYCmeaaMsZUU2Aa8wXxetZY7LEsuYEKEYf5ddMDY48")
	require.True(t, ok)
	require.EqualValues(t, v, 1)

	v, ok = accUpdOpt.Active.KeyAuths.Get("SCR7bRd3xQLCozabeBTXkxPWYzMQgHP3Aorj1h81WK68ovr83muoo")
	require.True(t, ok)
	require.EqualValues(t, v, 1)

	require.Equal(t, accUpdOpt.MemoKey, "SCR6W2AjgDsuYCmeaaMsZUU2Aa8wXxetZY7LEsuYEKEYf5ddMDY48")
	require.Equal(t, accUpdOpt.JsonMetadata, "{\"created_at\": \"GENESIS\"}")
}

func TestDelegateScorumpowerOperation(t *testing.T) {
	client := newMainNetHTTPClient()

	blockIDDelegateSCP := uint32(5709822)

	block, err := client.BlockchainHistory.GetBlock(context.Background(), blockIDDelegateSCP)

	require.NoError(t, err)
	require.Len(t, block.Transactions, 1)
	require.Len(t, block.Transactions[0].Operations, 1)

	op := block.Transactions[0].Operations[0]
	require.Equal(t, op.Type(), types.DelegateScorumpower)

	delegateScpOpt, ok := op.(*types.DelegateScorumpowerOperation)
	require.True(t, ok)
	require.Equal(t, delegateScpOpt.Delegator, "cali488")
	require.Equal(t, delegateScpOpt.Delegatee, "showtenseven")
	require.Equal(t, delegateScpOpt.Scorumpower, "2.182693663 SP")
}
