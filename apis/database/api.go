package database

import (
	"context"
	"encoding/json"

	"github.com/scorum/scorum-go/caller"
	"github.com/scorum/scorum-go/types"
)

const APIID = "database_api"

type API struct {
	caller caller.Caller
}

func NewAPI(caller caller.Caller) *API {
	return &API{caller}
}

func (api *API) call(ctx context.Context, method string, args []interface{}, reply interface{}) error {
	return api.caller.Call(ctx, APIID, method, args, reply)
}

func (api *API) setCallback(method string, callback func(raw json.RawMessage)) error {
	return api.caller.SetCallback(APIID, method, callback)
}

func (api *API) GetConfig(ctx context.Context) (*Config, error) {
	var config Config
	err := api.call(ctx, "get_config", caller.EmptyParams, &config)
	return &config, err
}

func (api *API) GetDynamicGlobalProperties(ctx context.Context) (*DynamicGlobalProperties, error) {
	var resp DynamicGlobalProperties
	err := api.call(ctx, "get_dynamic_global_properties", caller.EmptyParams, &resp)
	return &resp, err
}

// Get accounts by the provided names
func (api *API) GetAccounts(ctx context.Context, names ...string) ([]*Account, error) {
	var resp []*Account
	err := api.call(ctx, "get_accounts", []interface{}{names}, &resp)
	return resp, err
}

// GetAccountsCount returns account count
func (api *API) GetAccountsCount(ctx context.Context) (int, error) {
	var resp int
	err := api.call(ctx, "get_account_count", caller.EmptyParams, &resp)
	return resp, err
}

// LookupAccounts get names and IDs for registered accounts.
// lowerBoundName Lower bound of the first name to return.
// limit Maximum number of results to return -- must not exceed 1000
func (api *API) LookupAccounts(ctx context.Context, lowerBoundName string, limit uint16) ([]string, error) {
	var resp []string
	err := api.call(ctx, "lookup_accounts", []interface{}{lowerBoundName, limit}, &resp)
	return resp, err
}

// Set callback to invoke as soon as a new block is applied
func (api *API) SetBlockAppliedCallback(notice func(header *types.BlockHeader, error error)) (err error) {
	err = api.setCallback("set_block_applied_callback", func(raw json.RawMessage) {
		var header []types.BlockHeader
		if err := json.Unmarshal(raw, &header); err != nil {
			notice(nil, err)
		}
		for _, b := range header {
			notice(&b, nil)
		}
	})
	return
}
