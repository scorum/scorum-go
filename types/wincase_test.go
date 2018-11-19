package types

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	testCorrectScoreWincase = `["correct_score::yes",{"home":1,"away":2}]`
	testHandicapWincase     = `["handicap::over",{"threshold":1000}]`
	testResultWincase       = `["result_draw::no",{}]`
	testInvalidWincase      = `["doka_trade::no_thx",{}]`
)

func TestWincase_CorrectScore_MarshalJSON(t *testing.T) {
	wincase := Wincase{
		WincaseInterface: &ScoreYesNoWincase{
			ID:   WincaseCorrectScoreYes,
			Home: 1,
			Away: 2,
		},
	}

	j, err := json.Marshal(wincase)
	require.NoError(t, err)
	require.EqualValues(t, string(j), testCorrectScoreWincase)
}

func TestWincase_Handicap_MarshalJSON(t *testing.T) {
	wincase := Wincase{
		WincaseInterface: &OverUnderWincase{
			ID:        WincaseHandicapOver,
			Threshold: 1000,
		},
	}

	j, err := json.Marshal(wincase)
	require.NoError(t, err)
	require.EqualValues(t, string(j), testHandicapWincase)
}

func TestWincase_Result_MarshalJSON(t *testing.T) {
	wincase := Wincase{
		WincaseInterface: &YesNoWincase{
			ID: WincaseResultDrawNo,
		},
	}

	j, err := json.Marshal(wincase)
	require.NoError(t, err)
	require.EqualValues(t, string(j), testResultWincase)
}

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

func TestOverUnderWincase_GetMeta(t *testing.T) {
	m := OverUnderWincase{Threshold: 500}
	meta, err := m.GetMeta()
	require.NoError(t, err)
	require.EqualValues(t, testThresholdMeta, meta)
}

func TestScoreYesNoWincase_GetMeta(t *testing.T) {
	m := ScoreYesNoWincase{Home: 1, Away: 2}
	meta, err := m.GetMeta()
	require.NoError(t, err)
	require.EqualValues(t, testScoreMeta, meta)
}

func TestYesNoWincase_GetMeta(t *testing.T) {
	m := YesNoWincase{}
	meta, err := m.GetMeta()
	require.NoError(t, err)
	require.EqualValues(t, testYesNoMeta, meta)
}
