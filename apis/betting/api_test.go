package betting

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/scorum/scorum-go/transport/http"
	"github.com/stretchr/testify/require"
)

const nodeHTTPS = "https://testnet.scorum.com"

func TestGetGameWinners(t *testing.T) {
	t.Skip("need to start and finish game to get results")
	transport := http.NewTransport(nodeHTTPS)
	api := NewAPI(transport)

	uuid, err := uuid.Parse("3bd3fb0a-4c3c-4103-b736-61849157062a")
	require.NoError(t, err)

	winners, err := api.GetGameWinners(context.Background(), uuid)
	require.NoError(t, err)
	require.NotEmpty(t, winners)
}
