package blockchain_history

import (
	"github.com/scorum/scorum-go/caller"
)

const APIID = "blockchain_history_api"

type API struct {
	caller caller.Caller
}

func NewAPI(caller caller.Caller) *API {
	return &API{caller}
}

func (api *API) call(method string, args []interface{}, reply interface{}) error {
	return api.caller.Call(APIID, method, args, reply)
}

type AppliedOperationType int32

const (
	AllOp AppliedOperationType = iota
	NotVirtualOp
	VirtualOp
	MarketOp
)

// Get sequence of operations included/generated within a particular block
func (api *API) GetOperationsInBlock(blockNum uint32, opType AppliedOperationType) (History, error) {
	var resp History
	err := api.call("get_ops_in_block", []interface{}{blockNum, opType}, &resp)
	return resp, err
}
