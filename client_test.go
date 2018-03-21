package rpc

import (
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/scorum/scorum-go/apis/database"
	"github.com/scorum/scorum-go/caller"
	"github.com/stretchr/testify/require"
)

const (
	nodeWS   = "ws://blockchain.scorum.com:8003"
	nodeHTTP = "http://blockchain.scorum.com:8003"
)

func newWebsocketClient(t *testing.T) *Client {
	caller, err := caller.NewWebsocketCaller(nodeWS)
	require.NoError(t, err)
	client := NewClient(caller)
	return client
}

func newHttpClient(t *testing.T) *Client {
	caller := caller.NewHttpCaller(nodeHTTP)
	client := NewClient(caller)
	return client
}

func TestGetConfigViaHttp(t *testing.T) {
	client := newHttpClient(t)
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

	block, err := client.Database.GetBlock(int32(50))
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
	client := newHttpClient(t)
	defer client.Close()

	accounts, err := client.Database.GetAccounts("andrewww")
	require.NoError(t, err)

	require.Len(t, accounts, 1)
	require.Equal(t, "andrewww", accounts[0].Name)
}

func TestGetAccountHistory(t *testing.T) {
	client := newWebsocketClient(t)
	defer client.Close()

	history, err := client.AccountHistory.GetAccountHistory("andrewww", -1, 1000)
	require.NoError(t, err)
	require.True(t, len(history) > 0)

	t.Logf("history: %+v", history)
	spew.Dump(history)
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
