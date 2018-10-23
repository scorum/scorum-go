package types

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	testCorrectScoreWincase = `["correct_score::yes", { "home": 1, "away": 2 }]`
	testHandicapWincase     = `[ "handicap::over", { "threshold": 1000 } ]`
	testResultWincase       = `[ "result_draw::no", {} ]`
	testInvalidWincase      = `[ "doka_trade::no_thx", {} ]`
)

func TestWincase_CorrectScore_UnmarshalJSON(t *testing.T) {
	var wincase Wincase
	require.NoError(t, json.Unmarshal([]byte(testCorrectScoreWincase), &wincase))
	require.IsType(t, &ScoreYesNoWincase{}, wincase.WincaseInterface)
	require.EqualValues(t, 1, wincase.WincaseInterface.(*ScoreYesNoWincase).Home)
	require.EqualValues(t, 2, wincase.WincaseInterface.(*ScoreYesNoWincase).Away)
	require.EqualValues(t, 16, wincase.WincaseInterface.(*ScoreYesNoWincase).ID)
}

func TestWincase_Handicap_UnmarshalJSON(t *testing.T) {
	var wincase Wincase
	require.NoError(t, json.Unmarshal([]byte(testHandicapWincase), &wincase))
	require.IsType(t, &OverUnderWincase{}, wincase.WincaseInterface)
	require.EqualValues(t, 1000, wincase.WincaseInterface.(*OverUnderWincase).Threshold)
	require.EqualValues(t, 8, wincase.WincaseInterface.(*OverUnderWincase).ID)
}

func TestWincase_Result_UnmarshalJSON(t *testing.T) {
	var wincase Wincase
	require.NoError(t, json.Unmarshal([]byte(testResultWincase), &wincase))
	require.IsType(t, &YesNoWincase{}, wincase.WincaseInterface)
	require.EqualValues(t, 3, wincase.WincaseInterface.(*YesNoWincase).ID)
}

func TestWincase_Invalid_UnmarshalJSON(t *testing.T) {
	var wincase Wincase
	require.Error(t, json.Unmarshal([]byte(testInvalidWincase), &wincase))
}
