package rpc

import (
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/scorum/scorum-go/apis/database"
	"github.com/scorum/scorum-go/caller/websocket"
	"github.com/stretchr/testify/require"
)

const (
	node = "ws://localhost:8090"
)

func newClient(t *testing.T) *Client {
	caller, err := websocket.NewCaller(node)
	require.NoError(t, err)
	client := NewClient(caller)
	return client
}

func TestGetConfig(t *testing.T) {
	client := newClient(t)
	defer client.Close()

	config, err := client.Database.GetConfig()
	require.NoError(t, err)
	require.Equal(t, "SCR", config.ScorumAddressPrefix)
	t.Logf("config: %+v", config)
}

func TestGetChainProperties(t *testing.T) {
	client := newClient(t)
	defer client.Close()

	props, err := client.Database.GetChainProperties()
	require.NoError(t, err)
	require.True(t, props.MaximumBlockSize != 0)

	t.Logf("%+v", props)
}

func TestGetDynamicGlobalProperties(t *testing.T) {
	client := newClient(t)
	defer client.Close()

	config, err := client.Database.GetDynamicGlobalProperties()
	require.NoError(t, err)
	t.Logf("dynamic properties: %+v", config)
}

func TestGetBlockHeader(t *testing.T) {
	client := newClient(t)
	defer client.Close()

	block, err := client.Database.GetBlockHeader(24)
	require.NoError(t, err)

	require.NotEmpty(t, block.Previous)
	require.NotEmpty(t, block.Witness)
	require.Equal(t, block.Timestamp.Time, time.Date(2018, 1, 30, 12, 27, 6, 0, time.UTC))
	t.Logf("block header: %+v", block)
}

func TestGetBlock(t *testing.T) {
	client := newClient(t)
	defer client.Close()

	block, err := client.Database.GetBlock(int32(686))
	require.NoError(t, err)

	require.NotEmpty(t, block.Previous)
	require.NotEmpty(t, block.Witness)
	require.True(t, len(block.TransactionIDs) != 0)
	t.Logf("block: %+v", block)
}

func TestGetOpsInBlock(t *testing.T) {
	client := newClient(t)
	defer client.Close()

	ops, err := client.Database.GetOpsInBlock(int32(686), false)
	require.NoError(t, err)
	require.Len(t, ops, 1)
	require.Len(t, ops[0].Operations, 1)
}

func TestGetAccounts(t *testing.T) {
	client := newClient(t)
	defer client.Close()

	accounts, err := client.Database.GetAccounts("user20", "bob")
	require.NoError(t, err)

	require.Len(t, accounts, 2)
	require.Equal(t, "user20", accounts[0].Name)
	require.Equal(t, "bob", accounts[1].Name)
}

func TestGetAccountHistory(t *testing.T) {
	client := newClient(t)
	defer client.Close()

	history, err := client.Database.GetAccountHistory("initdelegate", -1, 1000)
	require.NoError(t, err)
	require.True(t, len(history) > 0)

	t.Logf("history: %+v", history)
	spew.Dump(history)
}

func TestSetBlockAppliedCallback(t *testing.T) {
	client := newClient(t)
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
