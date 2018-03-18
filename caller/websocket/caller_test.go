package websocket

import (
	"sync"
	"testing"

	"github.com/scorum/scorum-go/apis/database"
	"github.com/stretchr/testify/require"
)

const (
	node = "ws://localhost:8090"
)

func TestUnknownAPIID(t *testing.T) {
	caller, err := NewCaller(node)
	require.NoError(t, err)
	defer caller.Close()

	var reply interface{}
	err = caller.Call("some api", "some method", []interface{}{}, reply)
	require.Error(t, err)

	require.IsType(t, &RPCError{}, err)
	t.Logf("error: %+v", err)
}

func TestUnknownMethod(t *testing.T) {
	caller, err := NewCaller(node)
	require.NoError(t, err)
	defer caller.Close()

	var reply interface{}
	err = caller.Call(database.APIID, "some method", []interface{}{}, reply)
	require.Error(t, err)

	require.IsType(t, &RPCError{}, err)
	t.Logf("error: %+v", err)
}

func TestTooFewArgumentsPassedToMethod(t *testing.T) {
	caller, err := NewCaller(node)
	require.NoError(t, err)
	defer caller.Close()

	var reply interface{}
	err = caller.Call(database.APIID, "get_block_header", []interface{}{}, reply)
	require.Error(t, err)

	require.IsType(t, &RPCError{}, err)
	t.Logf("error: %+v", err)
}

func TestParallel(t *testing.T) {
	caller, err := NewCaller(node)
	require.NoError(t, err)
	defer caller.Close()

	const parallel = 20

	wg := sync.WaitGroup{}
	wg.Add(parallel)

	for i := 0; i < 20; i++ {
		go func(num int) {
			var resp interface{}
			err := caller.Call(database.APIID, "get_block_header", []interface{}{num}, &resp)
			require.NoError(t, err)
			wg.Done()
		}(i)
	}

	wg.Wait()
}
