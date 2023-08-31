package types

import (
	"encoding/json"
	"errors"
	"reflect"
)

const (
	SoccerGameType = 0
	HockeyGameType = 1
)

var (
	errUnsupportedGameType = errors.New("unsupported game type")
)

type GameType uint8

var GameTypeNames = map[GameType]string{
	SoccerGameType: "soccer_game",
	HockeyGameType: "hockey_game",
}

func (g GameType) MarshalJSON() ([]byte, error) {
	var err error

	a := make([]json.RawMessage, 2)
	a[0], err = json.Marshal(GameTypeNames[g])
	if err != nil {
		return nil, err
	}
	a[1] = json.RawMessage("{}")

	return json.Marshal(a)
}

func (g *GameType) UnmarshalJSON(b []byte) error {
	var gt []interface{}
	if err := json.Unmarshal(b, &gt); err != nil {
		return err
	}

	gName := gt[0].(string)
	for k, v := range GameTypeNames {
		if v == gName {
			reflect.Indirect(reflect.ValueOf(g)).SetUint(uint64(k))
			return nil
		}
	}

	return errUnsupportedGameType
}
