package types

import (
	"encoding/json"
	"github.com/scorum/scorum-go/encoding/transaction"
)

const (
	MarketResultHome MarketID = iota
	MarketResultDraw
	MarketResultAway
	MarketRoundHome
	MarketHandicap
	MarketCorrectScoreHome
	MarketCorrectScoreDraw
	MarketCorrectScoreAway
	MarketCorrectScore
	MarketGoalHome
	MarketGoalBoth
	MarketGoalAway
	MarketTotal
	MarketTotalGoalsHome
	MarketTotalGoalsAway
)

var MarketNames = map[MarketID]string{
	MarketResultHome:       "result_home",
	MarketResultDraw:       "result_draw",
	MarketResultAway:       "result_away",
	MarketRoundHome:        "round_home",
	MarketHandicap:         "handicap",
	MarketCorrectScoreHome: "correct_score_home",
	MarketCorrectScoreDraw: "correct_score_draw",
	MarketCorrectScoreAway: "correct_score_away",
	MarketCorrectScore:     "correct_score",
	MarketGoalHome:         "goal_home",
	MarketGoalBoth:         "goal_both",
	MarketGoalAway:         "goal_away",
	MarketTotal:            "total",
	MarketTotalGoalsHome:   "total_goals_home",
	MarketTotalGoalsAway:   "total_goals_away",
}

type Market interface {
	transaction.TransactionMarshaller
}

type MarketID int8

type OverUnderMarket struct {
	ID MarketID

	Threshold int16
}

func (m OverUnderMarket) MarshalJSON() ([]byte, error) {
	var err error

	a := make([]json.RawMessage, 2)
	a[0], err = json.Marshal(MarketNames[m.ID])
	if err != nil {
		return nil, err
	}

	a[1], err = json.Marshal(struct {
		Threshold int16 `json:"threshold"`
	}{m.Threshold})
	if err != nil {
		return nil, err
	}

	return json.Marshal(a)
}

func (op *OverUnderMarket) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)
	enc.Encode(int8(op.ID))
	enc.Encode(op.Threshold)
	return enc.Err()
}

type ScoreYesNoMarket struct {
	ID MarketID

	Home uint16
	Away uint16
}

func (m ScoreYesNoMarket) MarshalJSON() ([]byte, error) {
	var err error

	a := make([]json.RawMessage, 2)
	a[0], err = json.Marshal(MarketNames[m.ID])
	if err != nil {
		return nil, err
	}

	a[1], err = json.Marshal(struct {
		Home uint16 `json:"home"`
		Away uint16 `json:"away"`
	}{Home: m.Home, Away: m.Away})
	if err != nil {
		return nil, err
	}

	return json.Marshal(a)
}

func (op *ScoreYesNoMarket) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)
	enc.Encode(int8(op.ID))
	enc.Encode(op.Home)
	enc.Encode(op.Away)
	return enc.Err()
}

type YesNoMarket struct {
	ID MarketID
}

func (m YesNoMarket) MarshalJSON() ([]byte, error) {
	var err error

	a := make([]json.RawMessage, 2)
	a[0], err = json.Marshal(MarketNames[m.ID])
	if err != nil {
		return nil, err
	}
	a[1] = json.RawMessage("{}")

	return json.Marshal(a)
}

func (op *YesNoMarket) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)
	enc.Encode(int8(op.ID))
	return enc.Err()
}
