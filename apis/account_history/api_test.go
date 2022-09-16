package account_history

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/scorum/scorum-go/rpc"
)

const nodeHTTPS = "https://testnet.scorum.work"

func TestGetAccountScrToScrTransfers(t *testing.T) {
	transport := rpc.NewHTTPTransport(nodeHTTPS)
	api := NewAPI(transport)

	history, err := api.GetAccountScrToScrTransfers(context.Background(), "sheldon", -1, 3)
	require.NoError(t, err)
	require.True(t, len(history) > 0)
}
