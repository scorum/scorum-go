package account_history

import "github.com/scorum/scorum-go/caller"

const APIID = "account_history_api"

type API struct {
	caller caller.Caller
}

func NewAPI(caller caller.Caller) *API {
	return &API{caller}
}

func (api *API) call(method string, args []interface{}, reply interface{}) error {
	return api.caller.Call(APIID, method, args, reply)
}

// Get accounts by the provided names
// Account operations have sequence numbers from 0 to N where N is the most recent operation. This method
// returns operations in the range [from-limit, from]
// from - the absolute sequence number, -1 means most recent, limit is the number of operations before from.
// limit - the maximum number of items that can be queried (0 to 1000], must be less than from
func (api *API) GetAccountHistory(name string, from, limit int32) (AccountHistory, error) {
	resp := make(AccountHistory, 0)
	err := api.call("get_account_history", []interface{}{name, from, limit}, &resp)
	return resp, err
}
