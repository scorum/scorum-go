package types

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	testValidScoreMarket     = `["correct_score",{"home":1,"away":2}]`
	testValidThresholdMarket = `["correct_score",{"threshold":500}]`
	testValidYesNoMarket     = `["correct_score",{}]`
	testInvalidMarket        = `["doka_market", {}]`
)

func TestMarket_validScoreMarket_UnmarshalJSON(t *testing.T) {
	var market Market
	require.NoError(t, json.Unmarshal([]byte(testValidScoreMarket), &market))
	require.IsType(t, &ScoreYesNoMarket{}, market.MarketInterface)
	require.EqualValues(t, 2, market.MarketInterface.(*ScoreYesNoMarket).Away)
	require.EqualValues(t, 1, market.MarketInterface.(*ScoreYesNoMarket).Home)
	require.EqualValues(t, 8, market.MarketInterface.(*ScoreYesNoMarket).ID)
}

func TestMarket_validThresholdMarket_UnmarshalJSON(t *testing.T) {
	var market Market
	require.NoError(t, json.Unmarshal([]byte(testValidThresholdMarket), &market))
	require.IsType(t, &OverUnderMarket{}, market.MarketInterface)
	require.EqualValues(t, 500, market.MarketInterface.(*OverUnderMarket).Threshold)
	require.EqualValues(t, 8, market.MarketInterface.(*OverUnderMarket).ID)
}

func TestMarket_yesNoMarket_UnmarshalJSON(t *testing.T) {
	var market Market
	require.NoError(t, json.Unmarshal([]byte(testValidYesNoMarket), &market))
	require.IsType(t, &YesNoMarket{}, market.MarketInterface)
	require.EqualValues(t, 8, market.MarketInterface.(*YesNoMarket).ID)
}

func TestMarket_invalidMarket_UnmarshalJSON(t *testing.T) {
	var market Market
	require.Error(t, json.Unmarshal([]byte(testInvalidMarket), &market))
}
