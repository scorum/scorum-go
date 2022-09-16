package blockchain_history

import (
	"context"

	"github.com/scorum/scorum-go/caller"
	"github.com/scorum/scorum-go/types"
)

const APIID = "blockchain_history_api"

type API struct {
	caller caller.Caller
}

func NewAPI(caller caller.Caller) *API {
	return &API{caller}
}

func (api *API) call(ctx context.Context, method string, args []interface{}, reply interface{}) error {
	return api.caller.Call(ctx, APIID, method, args, reply)
}

type AppliedOperationType int32

const (
	AllOp AppliedOperationType = iota
	NotVirtualOp
	VirtualOp
	MarketOp
)

// Get a full signed block by the given block number
func (api *API) GetBlock(ctx context.Context, blockNum uint32) (*types.Block, error) {
	var resp types.Block
	err := api.call(ctx, "get_block", []interface{}{blockNum}, &resp)
	return &resp, err
}

// Get block header by the given block number
func (api *API) GetBlockHeader(ctx context.Context, blockNum int32) (*types.BlockHeader, error) {
	var resp types.BlockHeader
	err := api.call(ctx, "get_block_header", []interface{}{blockNum}, &resp)
	return &resp, err
}

// Get sequence of operations included/generated within a particular block
func (api *API) GetOperationsInBlock(ctx context.Context, blockNum uint32, opType AppliedOperationType) (History, error) {
	var resp History
	err := api.call(ctx, "get_ops_in_block", []interface{}{blockNum, opType}, &resp)
	return resp, err
}

// Get sequence of 'limit' blocks with transactions from offset
func (api *API) GetBlocksHistory(ctx context.Context, blockNum uint32, limit uint32) (BlockHistory, error) {
	var resp BlockHistory
	err := api.call(ctx, "get_blocks_history", []interface{}{blockNum, limit}, &resp)
	return resp, err
}

// Get sequence of 'limit' blocks with operations from offset
func (api *API) GetBlocks(ctx context.Context, blockNum uint32, limit uint32) (Blocks, error) {
	var resp Blocks
	err := api.call(ctx, "get_blocks", []interface{}{blockNum, limit}, &resp)
	return resp, err
}
