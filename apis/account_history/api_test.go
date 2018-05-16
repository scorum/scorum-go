package account_history

import (
	"testing"

	"github.com/scorum/scorum-go/transport/http"
	"github.com/stretchr/testify/require"
)

const nodeHTTPS = "https://testnet.scorum.com"

func TestGetAccountScrToScrTransfers(t *testing.T) {
	transport := http.NewTransport(nodeHTTPS)
	api := NewAPI(transport)

	history, err := api.GetAccountScrToScrTransfers("sheldon", -1, 3)
	require.NoError(t, err)
	require.True(t, len(history) > 0)
}
