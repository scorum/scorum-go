package blockchain_history

import (
	"testing"

	"github.com/scorum/scorum-go/transport/http"
	"github.com/stretchr/testify/require"
)

const nodeHTTPS = "https://testnet.scorum.com"

func TestGetBlockHeader(t *testing.T) {
	transport := http.NewTransport(nodeHTTPS)
	api := NewAPI(transport)

	block, err := api.GetBlockHeader(24)
	require.NoError(t, err)

	require.NotEmpty(t, block.Previous)
	require.NotEmpty(t, block.Witness)
}

func TestGetBlock(t *testing.T) {
	transport := http.NewTransport(nodeHTTPS)
	api := NewAPI(transport)

	block, err := api.GetBlock(uint32(50))
	require.NoError(t, err)

	require.NotEmpty(t, block.Previous)
	require.NotEmpty(t, block.TransactionMerkleRoot)
	require.NotEmpty(t, "00000032cfc128aff54138d97d183c416a352ec7", block.BlockID)
	require.Equal(t, "scorumwitness1", block.Witness)
}
