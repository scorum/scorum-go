package chain

import (
	"github.com/scorum/scorum-go/caller"
)

const APIID = "chain_api"

type API struct {
	caller caller.Caller
}

func NewAPI(caller caller.Caller) *API {
	return &API{caller}
}

func (api *API) call(method string, args []interface{}, reply interface{}) error {
	return api.caller.Call(APIID, method, args, reply)
}

// Get chain properties
func (api *API) GetChainProperties() (*ChainProperties, error) {
	var resp ChainProperties
	err := api.call("get_chain_properties", caller.EmptyParams, &resp)
	return &resp, err
}
