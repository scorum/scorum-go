package chain

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/scorum/scorum-go/rpc"
)

const nodeHTTPS = "https://testnet.scorum.work"

func TestGetChainProperties(t *testing.T) {
	transport := rpc.NewHTTPTransport(nodeHTTPS)
	api := NewAPI(transport)

	props, err := api.GetChainProperties(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, props.ChainID)
	require.True(t, props.HeadBlockNumber > 0)
	require.True(t, props.LastIrreversibleBlockNumber > 0)
}
