package betting

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/scorum/scorum-go/rpc"
)

const nodeHTTPS = "https://testnet.scorum.work"

func TestGetGameWinners(t *testing.T) {
	t.Skip("need to start and finish game to get results")
	api := NewAPI(rpc.NewHTTPTransport(nodeHTTPS))

	gameUUID, err := uuid.Parse("3bd3fb0a-4c3c-4103-b736-61849157062a")
	require.NoError(t, err)

	winners, err := api.GetGameWinners(context.Background(), gameUUID)
	require.NoError(t, err)
	require.NotEmpty(t, winners)
}
