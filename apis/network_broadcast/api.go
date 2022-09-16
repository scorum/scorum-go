package network_broadcast

import (
	"context"

	"github.com/scorum/scorum-go/caller"
	"github.com/scorum/scorum-go/types"
)

const APIID = "network_broadcast_api"

type API struct {
	caller caller.Caller
}

func NewAPI(caller caller.Caller) *API {
	return &API{caller}
}

func (api *API) call(ctx context.Context, method string, args []interface{}, reply interface{}) error {
	return api.caller.Call(ctx, APIID, method, args, reply)
}

func (api *API) BroadcastTransaction(ctx context.Context, tx *types.Transaction) error {
	return api.call(ctx, "broadcast_transaction", []interface{}{tx}, nil)
}

func (api *API) BroadcastTransactionSynchronous(ctx context.Context, tx *types.Transaction) (*BroadcastResponse, error) {
	var resp BroadcastResponse
	err := api.call(ctx, "broadcast_transaction_synchronous", []interface{}{tx}, &resp)
	return &resp, err
}
