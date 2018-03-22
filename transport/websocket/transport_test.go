package websocket

import (
	"sync"
	"testing"

	"github.com/scorum/scorum-go/transport"
	"github.com/stretchr/testify/require"
)

const (
	nodeWS = "ws://blockchain.scorum.com:8003"
)

func TestUnknownAPIID(t *testing.T) {
	caller, err := NewTransport(nodeWS)
	require.NoError(t, err)
	defer caller.Close()

	var reply interface{}
	err = caller.Call("some api", "some method", []interface{}{}, reply)
	require.Error(t, err)

	require.IsType(t, &transport.RPCError{}, err)
	t.Logf("error: %+v", err)
}

func TestUnknownMethod(t *testing.T) {
	caller, err := NewTransport(nodeWS)
	require.NoError(t, err)
	defer caller.Close()

	var reply interface{}
	err = caller.Call("database_api", "some method", []interface{}{}, reply)
	require.Error(t, err)

	require.IsType(t, &transport.RPCError{}, err)
	t.Logf("error: %+v", err)
}

func TestTooFewArgumentsPassedToMethod(t *testing.T) {
	caller, err := NewTransport(nodeWS)
	require.NoError(t, err)
	defer caller.Close()

	var reply interface{}
	err = caller.Call("database_api", "get_block_header", []interface{}{}, reply)
	require.Error(t, err)

	require.IsType(t, &transport.RPCError{}, err)
	t.Logf("error: %+v", err)
}

func TestParallel(t *testing.T) {
	caller, err := NewTransport(nodeWS)
	require.NoError(t, err)
	defer caller.Close()

	const parallel = 20

	wg := sync.WaitGroup{}
	wg.Add(parallel)

	for i := 0; i < 20; i++ {
		go func(num int) {
			var resp interface{}
			err := caller.Call("database_api", "get_block_header", []interface{}{num}, &resp)
			require.NoError(t, err)
			wg.Done()
		}(i)
	}

	wg.Wait()
}
