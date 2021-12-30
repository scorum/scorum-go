package chain

import (
	"context"
	"testing"

	"github.com/scorum/scorum-go/transport/http"
	"github.com/stretchr/testify/require"
)

const nodeHTTPS = "https://testnet.scorum.com"

func TestGetChainProperties(t *testing.T) {
	transport := http.NewTransport(nodeHTTPS)
	api := NewAPI(transport)

	props, err := api.GetChainProperties(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, props.ChainID)
	require.True(t, props.HeadBlockNumber > 0)
	require.True(t, props.LastIrreversibleBlockNumber > 0)
}
