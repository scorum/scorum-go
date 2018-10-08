package types

import (
	"bytes"
	"encoding/hex"
	"testing"
	"time"

	"github.com/scorum/scorum-go/encoding/transaction"
	"github.com/stretchr/testify/require"
)

func TestCreateGameOperation_SerializationWithoutMarkets(t *testing.T) {
	time, err := time.Parse(Layout, `"2018-08-03T10:12:43"`)
	require.NoError(t, err)

	op := CreateGameOperation{
		Moderator:           "admin",
		Name:                "game name",
		StartTime:           Time{&time},
		AutoResolveDelaySec: 33,
	}

	var b bytes.Buffer
	encoder := transaction.NewEncoder(&b)
	require.NoError(t, op.MarshalTransaction(encoder))

	require.EqualValues(t, "230561646d696e0967616d65206e616d659b2a645b210000000000", hex.EncodeToString(b.Bytes()))
}

func TestCreateGameOperation_SerializationWithTotalMarket(t *testing.T) {
	time, err := time.Parse(Layout, `"2018-08-03T10:12:43"`)
	require.NoError(t, err)

	op := CreateGameOperation{
		Moderator:           "admin",
		Name:                "game name",
		StartTime:           Time{&time},
		AutoResolveDelaySec: 33,
		Markets: []Market{&OverUnderMarket{
			ID:        MarketTotal,
			Threshold: 1000,
		}},
	}

	var b bytes.Buffer
	encoder := transaction.NewEncoder(&b)
	require.NoError(t, op.MarshalTransaction(encoder))
	require.EqualValues(t, "230561646d696e0967616d65206e616d659b2a645b2100000000010ce803", hex.EncodeToString(b.Bytes()))
}

func TestPostGameResultsOperation_Serialization(t *testing.T) {
	op := PostGameResultsOperation{
		GameID:    42,
		Moderator: "homer",
		Wincases: []Wincase{
			&YesNoWincase{
				ID: WincaseResultHomeYes,
			},
			&YesNoWincase{
				ID: WincaseResultDrawNo,
			},
			&YesNoWincase{
				ID: WincaseResultAwayYes,
			},
			&YesNoWincase{
				ID: WincaseResultHomeNo,
			},
			&OverUnderWincase{
				ID:        WincaseHandicapOver,
				Threshold: 1000,
			},
			&OverUnderWincase{
				ID:        WincaseHandicapUnder,
				Threshold: -500,
			},
			&OverUnderWincase{
				ID:        WincaseHandicapUnder,
				Threshold: 0,
			},
			&YesNoWincase{
				ID: WincaseCorrectScoreHomeYes,
			},
			&YesNoWincase{
				ID: WincaseCorrectScoreDrawNo,
			},
			&YesNoWincase{
				ID: WincaseCorrectScoreAwayNo,
			},
			&ScoreYesNoWincase{
				ID:   WincaseCorrectScoreYes,
				Home: 1,
				Away: 2,
			},
			&ScoreYesNoWincase{
				ID:   WincaseCorrectScoreNo,
				Home: 3,
				Away: 2,
			},
			&YesNoWincase{
				ID: WincaseGoalHomeYes,
			},
			&YesNoWincase{
				ID: WincaseGoalBothNo,
			},
			&YesNoWincase{
				ID: WincaseGoalAwayYes,
			},
			&OverUnderWincase{
				ID:        WincaseTotalOver,
				Threshold: 0,
			},
			&OverUnderWincase{
				ID:        WincaseTotalUnder,
				Threshold: 1000,
			},
		},
	}
	var b bytes.Buffer
	encoder := transaction.NewEncoder(&b)

	require.NoError(t, op.MarshalTransaction(encoder))

	require.EqualValues(
		t,
		"2705686f6d65722a00000000000000110003040108e803090cfe0900000a0d0f1001000200110300020012151618000019e803",
		hex.EncodeToString(b.Bytes()))
}
