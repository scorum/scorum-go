package types

import (
	"bytes"
	"encoding/hex"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/scorum/scorum-go/encoding/transaction"
	"github.com/stretchr/testify/require"
)

func TestCreateGameOperation_SerializationWithoutMarkets(t *testing.T) {
	time, err := time.Parse(Layout, `"2018-08-03T10:12:43"`)
	require.NoError(t, err)

	uuid := uuid.UUID{}
	require.NoError(t, uuid.UnmarshalText([]byte("e629f9aa-6b2c-46aa-8fa8-36770e7a7a5f")))

	op := CreateGameOperation{
		UUID:                uuid,
		Moderator:           "admin",
		Name:                "game name",
		StartTime:           Time{&time},
		AutoResolveDelaySec: 33,
	}

	var b bytes.Buffer
	encoder := transaction.NewEncoder(&b)
	require.NoError(t, op.MarshalTransaction(encoder))

	require.EqualValues(t, "23e629f9aa6b2c46aa8fa836770e7a7a5f0561646d696e0967616d65206e616d659b2a645b210000000000", hex.EncodeToString(b.Bytes()))
}

func TestCreateGameOperation_SerializationWithTotalMarket(t *testing.T) {
	time, err := time.Parse(Layout, `"2018-08-03T10:12:43"`)
	require.NoError(t, err)

	uuid := uuid.UUID{}
	require.NoError(t, uuid.UnmarshalText([]byte("e629f9aa-6b2c-46aa-8fa8-36770e7a7a5f")))

	op := CreateGameOperation{
		UUID:                uuid,
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
	require.EqualValues(t, "23e629f9aa6b2c46aa8fa836770e7a7a5f0561646d696e0967616d65206e616d659b2a645b2100000000010ce803", hex.EncodeToString(b.Bytes()))
}

func TestCreateGameOperation_SerializationManyMarkets(t *testing.T) {
	time := time.Unix(1461605400, 0)

	uuid := uuid.UUID{}
	require.NoError(t, uuid.UnmarshalText([]byte("e629f9aa-6b2c-46aa-8fa8-36770e7a7a5f")))

	op := CreateGameOperation{
		UUID:                uuid,
		Moderator:           "moderator_name",
		Name:                "game_name",
		StartTime:           Time{&time},
		AutoResolveDelaySec: 33,
		Markets: []Market{
			&YesNoMarket{
				ID: MarketResultHome,
			},
			&YesNoMarket{
				ID: MarketResultDraw,
			},
			&YesNoMarket{
				ID: MarketResultAway,
			},
			&YesNoMarket{
				ID: MarketRoundHome,
			},
			&OverUnderMarket{
				ID:        MarketHandicap,
				Threshold: -500,
			},
			&OverUnderMarket{
				ID:        MarketHandicap,
				Threshold: 0,
			},
			&OverUnderMarket{
				ID:        MarketHandicap,
				Threshold: 1000,
			},
			&YesNoMarket{
				ID: MarketCorrectScoreHome,
			},
			&YesNoMarket{
				ID: MarketCorrectScoreDraw,
			},
			&YesNoMarket{
				ID: MarketCorrectScoreAway,
			},
			&ScoreYesNoMarket{
				ID:   MarketCorrectScore,
				Home: 1,
				Away: 0,
			},
			&ScoreYesNoMarket{
				ID:   MarketCorrectScore,
				Home: 1,
				Away: 1,
			},
			&YesNoMarket{
				ID: MarketGoalHome,
			},
			&YesNoMarket{
				ID: MarketGoalBoth,
			},
			&YesNoMarket{
				ID: MarketGoalAway,
			},
			&OverUnderMarket{
				ID:        MarketTotal,
				Threshold: 0,
			},
			&OverUnderMarket{
				ID:        MarketTotal,
				Threshold: 500,
			},
			&OverUnderMarket{
				ID:        MarketTotal,
				Threshold: 1000,
			},
		},
	}

	var b bytes.Buffer
	encoder := transaction.NewEncoder(&b)
	require.NoError(t, op.MarshalTransaction(encoder))
	require.EqualValues(
		t,
		"23e629f9aa6b2c46aa8fa836770e7a7a5f0e6d6f64657261746f725f6e616d650967616d655f6e616d6518541e5721000000001200010203040cfe04000004e80305060708010000000801000100090a0b0c00000cf4010ce803",
		hex.EncodeToString(b.Bytes()))
}

func TestPostGameResultsOperation_Serialization(t *testing.T) {
	uuid := uuid.UUID{}
	require.NoError(t, uuid.UnmarshalText([]byte("e629f9aa-6b2c-46aa-8fa8-36770e7a7a5f")))

	op := PostGameResultsOperation{
		UUID:      uuid,
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
				ID: WincaseRoundHomeNo,
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
		"27e629f9aa6b2c46aa8fa836770e7a7a5f05686f6d6572110003040708e803090cfe0900000a0d0f1001000200110300020012151618000019e803",
		hex.EncodeToString(b.Bytes()))
}
