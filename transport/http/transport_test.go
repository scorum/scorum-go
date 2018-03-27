package http

import (
	"testing"

	"github.com/scorum/scorum-go/transport"
	"github.com/stretchr/testify/require"
)

const (
	nodeHTTPS = "https://testnet.scorum.com"
)

func TestNodeIsDown(t *testing.T) {
	caller := NewTransport("http://nonode.scorum.com")
	defer caller.Close()

	var reply interface{}
	err := caller.Call("some api", "some method", []interface{}{}, reply)
	require.Error(t, err)
}

func TestUnknownAPIID(t *testing.T) {
	caller := NewTransport(nodeHTTPS)
	defer caller.Close()

	var reply interface{}
	err := caller.Call("some api", "some method", []interface{}{}, reply)
	require.Error(t, err)

	require.IsType(t, &transport.RPCError{}, err)
	t.Logf("error: %+v", err)
}
