package chain

import (
	"context"

	"github.com/scorum/scorum-go/caller"
)

const APIID = "chain_api"

type API struct {
	caller caller.Caller
}

func NewAPI(caller caller.Caller) *API {
	return &API{caller}
}

func (api *API) GetChainProperties(ctx context.Context) (*ChainProperties, error) {
	var resp ChainProperties
	err := api.caller.Call(ctx, APIID, "get_chain_properties", caller.EmptyParams, &resp)
	return &resp, err
}
