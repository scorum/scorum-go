package scorumgo

import (
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/scorum/scorum-go/sign"
	rpc "github.com/scorum/scorum-go/transport"
	"github.com/scorum/scorum-go/transport/http"
	"github.com/scorum/scorum-go/transport/websocket"
	"github.com/scorum/scorum-go/types"
	"github.com/stretchr/testify/require"
)

const (
	nodeWSS   = "wss://testnet.scorum.com"
	nodeHTTPS = "https://testnet.scorum.com"
)

// test accounts available at https://github.com/scorum/scorum/wiki/Testnet-existent-accounts

func newWebsocketClient(t *testing.T) *Client {
	transport, err := websocket.NewTransport(nodeWSS)
	require.NoError(t, err)
	client := NewClient(transport)
	return client
}

func newHTTPClient() *Client {
	transport := http.NewTransport(nodeHTTPS)
	client := NewClient(transport)
	return client
}

func TestGetConfigViaWS(t *testing.T) {
	client := newWebsocketClient(t)
	defer client.Close()

	config, err := client.Database.GetConfig()
	require.NoError(t, err)
	require.Equal(t, "SCR", config.ScorumAddressPrefix)
}

func TestGetDynamicGlobalProperties(t *testing.T) {
	client := newWebsocketClient(t)
	defer client.Close()

	config, err := client.Database.GetDynamicGlobalProperties()
	require.NoError(t, err)
	t.Logf("dynamic properties: %+v", config)
}

func TestGetAccounts(t *testing.T) {
	client := newHTTPClient()
	defer client.Close()

	accounts, err := client.Database.GetAccounts("leonarda", "kristie")
	require.NoError(t, err)

	require.Len(t, accounts, 2)
	require.Equal(t, "leonarda", accounts[0].Name)
	require.Equal(t, "kristie", accounts[1].Name)
}

func TestGetAccountHistory(t *testing.T) {
	client := newHTTPClient()

	history, err := client.AccountHistory.GetAccountHistory("leonarda", -1, 3)
	require.NoError(t, err)
	require.True(t, len(history) > 0)
	spew.Dump(history)
}

func TestClient_Broadcast_AccountWitnessVoteOperation(t *testing.T) {
	client := newHTTPClient()

	roselle := "5JwWJ2m2jGG9RPcpDix5AvkDzQZJoZvpUQScsDzzXWAKMs8Q6jH"
	_, err := client.Broadcast(sign.TestChain, []string{roselle}, &types.AccountWitnessVoteOperation{
		Account: "roselle",
		Witness: "scorumwitness1",
		Approve: true,
	})
	require.NotNil(t, err)

	perr, ok := err.(*rpc.RPCError)
	require.True(t, ok)
	require.Equal(t, "assert_exception", perr.Data.Name)
	require.Equal(t, int(10), perr.Data.Code)
}

func TestClient_Broadcast_Transfer(t *testing.T) {
	client := newHTTPClient()
	amount, _ := types.AssetFromString("0.000009 SCR")

	azucena := "5J7FEcpqc1sZ7ZbKx2kVvBHx2oTjWG2wMU2e2FYX85sGA2qu8KT"
	resp, err := client.Broadcast(sign.TestChain, []string{azucena}, &types.TransferOperation{
		From:   "azucena",
		To:     "leonarda",
		Amount: *amount,
		Memo:   "1",
	})
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
