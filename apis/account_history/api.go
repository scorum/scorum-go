package account_history

import (
	"context"

	"github.com/scorum/scorum-go/caller"
)

const APIID = "account_history_api"

type API struct {
	caller caller.Caller
}

func NewAPI(caller caller.Caller) *API {
	return &API{caller}
}

func (api *API) call(ctx context.Context, method string, args []interface{}, reply interface{}) error {
	return api.caller.Call(ctx, APIID, method, args, reply)
}

// GetAccountHistory returns operations history for the given account
// Owner operations have sequence numbers from 0 to N where N is the most recent operation. This method
// returns operations in the range [from-limit, from]
// from - the absolute sequence number, -1 means most recent, limit is the number of operations before from.
// limit - the maximum number of items that can be queried [1 to 1000], must be less than from
func (api *API) GetAccountHistory(ctx context.Context, account string, from, limit int32) (AccountHistory, error) {
	resp := make(AccountHistory, 0)
	err := api.call(ctx, "get_account_history", []interface{}{account, from, limit}, &resp)
	return resp, err
}

// GetAccountScrToScrTransfers returns transactions history for the given account
// from - the absolute sequence number, -1 means most recent, limit is the number of operations before from.
// limit - the maximum number of items that can be queried [1 to 1000], must be less than from
func (api *API) GetAccountScrToScrTransfers(ctx context.Context, account string, from, limit int32) (AccountHistory, error) {
	resp := make(AccountHistory, 0)
	err := api.call(ctx, "get_account_scr_to_scr_transfers", []interface{}{account, from, limit}, &resp)
	return resp, err
}
