package account_history

import (
	"context"
	"testing"

	"github.com/scorum/scorum-go/transport/http"
	"github.com/stretchr/testify/require"
)

const nodeHTTPS = "https://testnet.scorum.work"

func TestGetAccountScrToScrTransfers(t *testing.T) {
	transport := http.NewTransport(nodeHTTPS)
	api := NewAPI(transport)

	history, err := api.GetAccountScrToScrTransfers(context.Background(), "sheldon", -1, 3)
	require.NoError(t, err)
	require.True(t, len(history) > 0)
}
