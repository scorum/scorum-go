package websocket

import (
	"context"
	"encoding/json"
	"sync"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/scorum/scorum-go/rpc/protocol"
)

const (
	nodeWSS = "wss://testnet.scorum.work"
)

func TestUnknownAPIID(t *testing.T) {
	caller := NewTransport(NewConnector(nodeWSS, websocket.DefaultDialer))
	require.NoError(t, caller.Dial(context.Background()))
	defer func() {
		require.NoError(t, caller.Close())
	}()

	var reply interface{}
	err := caller.Call(context.Background(), "some api", "some method", []interface{}{}, reply)
	require.Error(t, err)

	require.IsType(t, &protocol.RPCError{}, err)
	t.Logf("error: %+v", err)
}

func TestUnknownMethod(t *testing.T) {
	caller := NewTransport(NewConnector(nodeWSS, websocket.DefaultDialer))
	require.NoError(t, caller.Dial(context.Background()))
	defer func() {
		require.NoError(t, caller.Close())
	}()

	var reply interface{}
	err := caller.Call(context.Background(), "database_api", "some method", []interface{}{}, reply)
	require.Error(t, err)

	require.IsType(t, &protocol.RPCError{}, err)
	t.Logf("error: %+v", err)
}

func TestTooFewArgumentsPassedToMethod(t *testing.T) {
	caller := NewTransport(NewConnector(nodeWSS, websocket.DefaultDialer))
	require.NoError(t, caller.Dial(context.Background()))
	defer func() {
		require.NoError(t, caller.Close())
	}()

	var reply interface{}
	err := caller.Call(context.Background(), "database_api", "get_block_header", []interface{}{}, reply)
	require.Error(t, err)

	require.IsType(t, &protocol.RPCError{}, err)
	t.Logf("error: %+v", err)
}

func TestSingleCall(t *testing.T) {
	caller := NewTransport(NewConnector(nodeWSS, websocket.DefaultDialer))
	require.NoError(t, caller.Dial(context.Background()))
	defer func() {
		require.NoError(t, caller.Close())
	}()

	var resp interface{}
	err := caller.Call(context.Background(), "blockchain_history_api", "get_block_header", []interface{}{10}, &resp)
	require.NoError(t, err)

	data, err := json.Marshal(resp)
	require.NoError(t, err)
	require.Equal(t, "{\"extensions\":[],\"previous\":\"00000009300eb6ecf852348bc8f38254f4f616c3\",\"timestamp\":\"2018-12-05T15:17:06\",\"transaction_merkle_root\":\"0000000000000000000000000000000000000000\",\"witness\":\"scorumwitness1\"}", string(data))
}

func TestParallel(t *testing.T) {
	caller := NewTransport(NewConnector(nodeWSS, websocket.DefaultDialer))
	require.NoError(t, caller.Dial(context.Background()))
	defer func() {
		require.NoError(t, caller.Close())
	}()

	const parallel = 20

	wg := sync.WaitGroup{}
	wg.Add(parallel)

	for i := 0; i < parallel; i++ {
		go func(num int) {
			var resp interface{}
			err := caller.Call(context.Background(), "blockchain_history_api", "get_block_header", []interface{}{num}, &resp)
			assert.NoError(t, err)
			wg.Done()
		}(i)
	}

	wg.Wait()
}
