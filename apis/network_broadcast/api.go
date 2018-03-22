package network_broadcast

import (
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

func (api *API) call(method string, args []interface{}, reply interface{}) error {
	return api.caller.Call(APIID, method, args, reply)
}

func (api *API) BroadcastTransaction(tx *types.Transaction) error {
	return api.call("broadcast_transaction", []interface{}{tx}, nil)
}

func (api *API) BroadcastTransactionSynchronous(tx *types.Transaction) (*BroadcastResponse, error) {
	var resp BroadcastResponse
	err := api.call("broadcast_transaction_synchronous", []interface{}{tx}, &resp)
	return &resp, err
}
