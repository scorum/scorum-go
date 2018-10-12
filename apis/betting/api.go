package betting

import (
	"github.com/google/uuid"
	"github.com/scorum/scorum-go/caller"
)

const APIID = "betting_api"

type API struct {
	caller caller.Caller
}

func NewAPI(caller caller.Caller) *API {
	return &API{caller}
}

func (api *API) call(method string, args []interface{}, reply interface{}) error {
	return api.caller.Call(APIID, method, args, reply)
}

func (api *API) GetGameWinners(gameID uuid.UUID) ([]Winner, error) {
	var resp []Winner
	err := api.call("get_game_winners", []interface{}{gameID.String()}, &resp)
	return resp, err
}
