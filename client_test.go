package scorumgo

import (
	"testing"
	"time"

	"github.com/scorum/scorum-go/apis/database"
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

func TestGetConfigViaHttp(t *testing.T) {
	client := newHTTPClient()
	defer client.Close()

	config, err := client.Database.GetConfig()
	require.NoError(t, err)
	require.Equal(t, "SCR", config.ScorumAddressPrefix)
}

func TestGetConfigViaWS(t *testing.T) {
	client := newWebsocketClient(t)
	defer client.Close()

	config, err := client.Database.GetConfig()
	require.NoError(t, err)
	require.Equal(t, "SCR", config.ScorumAddressPrefix)
}

func TestGetChainProperties(t *testing.T) {
	client := newWebsocketClient(t)
	defer client.Close()

	props, err := client.Database.GetChainProperties()
	require.NoError(t, err)
	require.True(t, props.MaximumBlockSize != 0)

	t.Logf("%+v", props)
}

func TestGetDynamicGlobalProperties(t *testing.T) {
	client := newWebsocketClient(t)
	defer client.Close()

	config, err := client.Database.GetDynamicGlobalProperties()
	require.NoError(t, err)
	t.Logf("dynamic properties: %+v", config)
}

func TestGetBlockHeader(t *testing.T) {
	client := newWebsocketClient(t)
	defer client.Close()

	block, err := client.Database.GetBlockHeader(24)
	require.NoError(t, err)

	require.NotEmpty(t, block.Previous)
	require.NotEmpty(t, block.Witness)
}

func TestGetBlock(t *testing.T) {
	client := newWebsocketClient(t)
	defer client.Close()

	block, err := client.Database.GetBlock(uint32(50))
	require.NoError(t, err)

	require.NotEmpty(t, block.Previous)
	require.NotEmpty(t, "00000032cfc128aff54138d97d183c416a352ec7", block.BlockID)
	require.Equal(t, "scorumwitness14", block.Witness)
	t.Logf("block: %+v", block)
}

func TestGetOpsInBlock(t *testing.T) {
	client := newWebsocketClient(t)
	defer client.Close()

	ops, err := client.Database.GetOpsInBlock(int32(686), false)
	require.NoError(t, err)
	require.Len(t, ops, 1)
	require.Len(t, ops[0].Operations, 1)
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
	transport := http.NewTransport(nodeHTTPS)
	client := NewClient(transport)

	history, err := client.AccountHistory.GetAccountHistory("roselle", -1, 1000)
	require.NoError(t, err)
	require.True(t, len(history) > 0)

	t.Logf("history: %+v", history)
}

func TestClient_Broadcast(t *testing.T) {
	transport := http.NewTransport(nodeHTTPS)
	client := NewClient(transport)

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

func TestSetBlockAppliedCallback(t *testing.T) {
	client := newWebsocketClient(t)
	defer client.Close()

	var called bool
	err := client.Database.SetBlockAppliedCallback(func(block *database.BlockHeader, err error) {
		t.Log("block:", block, "error", err)
		called = true
	})
	require.NoError(t, err)
	time.Sleep(10 * time.Second)
	require.True(t, called)
}
