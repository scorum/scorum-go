package http

import (
	"testing"

	"github.com/scorum/scorum-go/transport"
	"github.com/stretchr/testify/require"
)

const (
	nodeWS = "http://blockchain.scorum.com:8003"
)

func TestNodeIsDown(t *testing.T) {
	caller := NewTransport("http://node_is_down.scorum.com")
	defer caller.Close()

	var reply interface{}
	err := caller.Call("some api", "some method", []interface{}{}, reply)
	require.Error(t, err)
}

func TestUnknownAPIID(t *testing.T) {
	caller := NewTransport(nodeWS)
	defer caller.Close()

	var reply interface{}
	err := caller.Call("some api", "some method", []interface{}{}, reply)
	require.Error(t, err)

	require.IsType(t, &transport.RPCError{}, err)
	t.Logf("error: %+v", err)
}
