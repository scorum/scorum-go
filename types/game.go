package types

import (
	"encoding/json"
)

const (
	SoccerGameType = 0
	HockeyGameType = 1
)

type GameType int8

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
