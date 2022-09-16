package websocket

import (
	"context"
	"sync"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"

	"github.com/scorum/scorum-go/rpc/protocol"
)

const (
	nodeWSS = "wss://testnet.scorum.work"
)

func TestUnknownAPIID(t *testing.T) {
	ws, _, err := websocket.DefaultDialer.Dial(nodeWSS, nil)
	require.NoError(t, err)

	caller := NewTransport(ws)
	require.NoError(t, err)
	defer caller.Close()

	var reply interface{}
	err = caller.Call(context.Background(), "some api", "some method", []interface{}{}, reply)
	require.Error(t, err)

	require.IsType(t, &protocol.RPCError{}, err)
	t.Logf("error: %+v", err)
}

func TestUnknownMethod(t *testing.T) {
	ws, _, err := websocket.DefaultDialer.Dial(nodeWSS, nil)
	require.NoError(t, err)

	caller := NewTransport(ws)
	require.NoError(t, err)
	defer caller.Close()

	var reply interface{}
	err = caller.Call(context.Background(), "database_api", "some method", []interface{}{}, reply)
	require.Error(t, err)

	require.IsType(t, &protocol.RPCError{}, err)
	t.Logf("error: %+v", err)
}

func TestTooFewArgumentsPassedToMethod(t *testing.T) {
	ws, _, err := websocket.DefaultDialer.Dial(nodeWSS, nil)
	require.NoError(t, err)

	caller := NewTransport(ws)
	require.NoError(t, err)
	defer caller.Close()

	var reply interface{}
	err = caller.Call(context.Background(), "database_api", "get_block_header", []interface{}{}, reply)
	require.Error(t, err)

	require.IsType(t, &protocol.RPCError{}, err)
	t.Logf("error: %+v", err)
}

func TestParallel(t *testing.T) {
	ws, _, err := websocket.DefaultDialer.Dial(nodeWSS, nil)
	require.NoError(t, err)

	caller := NewTransport(ws)
	require.NoError(t, err)
	defer caller.Close()

	const parallel = 20

	wg := sync.WaitGroup{}
	wg.Add(parallel)

	for i := 0; i < 20; i++ {
		go func(num int) {
			var resp interface{}
			err := caller.Call(context.Background(), "blockchain_history_api", "get_block_header", []interface{}{num}, &resp)
			require.NoError(t, err)
			wg.Done()
		}(i)
	}

	wg.Wait()
}
