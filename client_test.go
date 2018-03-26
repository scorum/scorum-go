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
	nodeWS   = "ws://blockchain.scorum.com:8003"
	nodeHTTP = "http://blockchain.scorum.com:8003"
	nodeLive = "https://rpc.scorum.com"
)

func newWebsocketClient(t *testing.T) *Client {
	transport, err := websocket.NewTransport(nodeWS)
	require.NoError(t, err)
	client := NewClient(transport)
	return client
}

func newHTTPClient() *Client {
	transport := http.NewTransport(nodeHTTP)
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
	require.Equal(t, "scorumwitness2", block.Witness)
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

	accounts, err := client.Database.GetAccounts("leonarda", "andrewww")
	require.NoError(t, err)

	require.Len(t, accounts, 2)
	require.Equal(t, "leonarda", accounts[0].Name)
	require.Equal(t, "andrewww", accounts[1].Name)
}

func TestGetAccountHistory(t *testing.T) {
	client := newWebsocketClient(t)
	defer client.Close()

	history, err := client.AccountHistory.GetAccountHistory("leonarda", -1, 1000)
	require.NoError(t, err)
	require.True(t, len(history) > 0)

	t.Logf("history: %+v", history)
}

func TestClient_Broadcast(t *testing.T) {
	transport := http.NewTransport(nodeLive)
	client := NewClient(transport)

	megaherz := "5KHK69Be8P8NQLy46KXugJWyNkxw8Nw3Mzue4wD8ygx48emMugd"
	_, err := client.Broadcast(sign.ScorumChain, []string{megaherz}, &types.AccountWitnessVoteOperation{
		Account: "megaherz",
		Witness: "andrewww",
		Approve: false,
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
