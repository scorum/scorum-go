package types

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	testGameSoccerJSON      = `["soccer_game",{}]`
	testGameHockeyJSON      = `["hockey_game",{}]`
	testGameUnsupportedJSON = `["doka2_trade_game",{}]`
)

func TestGameSocckerUnmarshalJSON(t *testing.T) {
	var game GameType
	require.NoError(t, json.Unmarshal([]byte(testGameSoccerJSON), &game))
	require.EqualValues(t, 0, game)
}

func TestGameHockeyUnmarshalJSON(t *testing.T) {
	var game GameType
	require.NoError(t, json.Unmarshal([]byte(testGameHockeyJSON), &game))
	require.EqualValues(t, 1, game)
}

func TestUnsupportedUnmarshalJSON(t *testing.T) {
	var game GameType
	require.Error(t, json.Unmarshal([]byte(testGameUnsupportedJSON), &game))
}
