package blockchain_history

import (
	"context"
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/scorum/scorum-go/rpc"
)

const nodeHTTPS = "https://testnet.scorum.work"

func TestGetBlockHeader(t *testing.T) {
	transport := rpc.NewHTTPTransport(nodeHTTPS)
	api := NewAPI(transport)

	block, err := api.GetBlockHeader(context.Background(), 24)
	require.NoError(t, err)

	require.NotEmpty(t, block.Previous)
	require.NotEmpty(t, block.Witness)
}

func TestGetBlock(t *testing.T) {
	transport := rpc.NewHTTPTransport(nodeHTTPS)
	api := NewAPI(transport)

	block, err := api.GetBlock(context.Background(), uint32(50))
	require.NoError(t, err)

	require.NotEmpty(t, block.Previous)
	require.NotEmpty(t, block.TransactionMerkleRoot)
	require.NotEmpty(t, "00000032cfc128aff54138d97d183c416a352ec7", block.BlockID)
	require.Equal(t, "scorumwitness2", block.Witness)
}

func TestGetOperationsInBlock(t *testing.T) {
	transport := rpc.NewHTTPTransport(nodeHTTPS)
	api := NewAPI(transport)

	ops, err := api.GetOperationsInBlock(context.Background(), uint32(127), AllOp)
	require.NoError(t, err)
	require.NotEmpty(t, ops)

	for _, op := range ops {
		require.True(t, len(op.Operations) > 0)
	}
}

func TestGetBlocksHistory(t *testing.T) {
	transport := rpc.NewHTTPTransport(nodeHTTPS)
	api := NewAPI(transport)

	t.Run("from beginning", func(t *testing.T) {
		blocks, err := api.GetBlocksHistory(context.Background(), 100, 100)
		require.NoError(t, err)
		require.Len(t, blocks, 100)
	})

	t.Run("from end", func(t *testing.T) {
		blocks, err := api.GetBlocksHistory(context.Background(), math.MaxUint32, 100)
		require.NoError(t, err)
		require.True(t, len(blocks) > 0)
	})

	t.Run("exceeded limit", func(t *testing.T) {
		_, err := api.GetBlocksHistory(context.Background(), math.MaxUint32, 2000)
		require.Error(t, err)
	})

}

func TestGetBlocks(t *testing.T) {
	transport := rpc.NewHTTPTransport(nodeHTTPS)
	api := NewAPI(transport)

	t.Run("from beginning", func(t *testing.T) {
		blocks, err := api.GetBlocks(context.Background(), 100, 100)
		require.NoError(t, err)
		require.Len(t, blocks, 100)
		for _, v := range blocks {
			require.NotEmpty(t, v.Operations)
		}
	})

	t.Run("from end", func(t *testing.T) {
		blocks, err := api.GetBlocks(context.Background(), math.MaxUint32, 100)
		require.NoError(t, err)
		require.NotEmpty(t, blocks)
	})

	t.Run("exceeded limit", func(t *testing.T) {
		_, err := api.GetBlocks(context.Background(), math.MaxUint32, 2000)
		require.Error(t, err)
	})
}
