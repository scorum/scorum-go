package blockchain_history

import (
	"testing"

	"github.com/scorum/scorum-go/transport/http"
	"github.com/stretchr/testify/require"
	"math"
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

func TestGetOperationsInBlock(t *testing.T) {
	transport := http.NewTransport(nodeHTTPS)
	api := NewAPI(transport)

	ops, err := api.GetOperationsInBlock(uint32(127), AllOp)
	require.NoError(t, err)
	require.NotEmpty(t, ops)

	for _, op := range ops {
		require.True(t, len(op.Operations) > 0)
	}
}

func TestGetBlocksHistory(t *testing.T) {
	transport := http.NewTransport(nodeHTTPS)
	api := NewAPI(transport)

	t.Run("from beginning", func(t *testing.T) {
		blocks, err := api.GetBlocksHistory(100, 100)
		require.NoError(t, err)
		require.Len(t, blocks, 100)
	})

	t.Run("from end", func(t *testing.T) {
		blocks, err := api.GetBlocksHistory(math.MaxUint32, 100)
		require.NoError(t, err)
		require.True(t, len(blocks) > 0)
	})

	t.Run("exceeded limit", func(t *testing.T) {
		_, err := api.GetBlocksHistory(math.MaxUint32, 2000)
		require.Error(t, err)
	})

}
