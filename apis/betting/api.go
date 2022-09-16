package betting

import (
	"context"

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

func (api *API) call(ctx context.Context, method string, args []interface{}, reply interface{}) error {
	return api.caller.Call(ctx, APIID, method, args, reply)
}

func (api *API) GetGameWinners(ctx context.Context, gameID uuid.UUID) ([]Winner, error) {
	var resp []Winner
	err := api.call(ctx, "get_game_winners", []interface{}{gameID.String()}, &resp)
	return resp, err
}
