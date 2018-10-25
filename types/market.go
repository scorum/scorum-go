package types

import (
	"encoding/json"
	"github.com/bitly/go-simplejson"

	"errors"
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

var (
	errUnknownMarket = errors.New("unknown market id")
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

type Market struct {
	MarketInterface
}

type MarketInterface interface {
	transaction.TransactionMarshaller

	GetName() string
	GetID() int8
	GetMeta() (json.RawMessage, error)
}

func (m Market) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.MarketInterface)
}

func (m *Market) UnmarshalJSON(b []byte) error {
	json, err := simplejson.NewJson(b)
	if err != nil {
		return err
	}

	marketName, err := json.GetIndex(0).String()
	if err != nil {
		return err
	}

	marketID := MarketID(-1)
	for k, v := range MarketNames {
		if v == marketName {
			marketID = k
		}
	}
	if marketID == -1 {
		return errUnknownMarket
	}

	marketObj := json.GetIndex(1)

	threshold, err := marketObj.Get("threshold").Int()
	if err == nil {
		market := OverUnderMarket{
			ID:        marketID,
			Threshold: int16(threshold),
		}
		m.MarketInterface = MarketInterface(&market)
		return nil
	}

	home, err := marketObj.Get("home").Int()
	if err == nil {
		away, err := marketObj.Get("away").Int()
		if err == nil {
			market := ScoreYesNoMarket{
				ID:   marketID,
				Home: uint16(home),
				Away: uint16(away),
			}
			m.MarketInterface = MarketInterface(&market)
			return nil
		}
	}

	market := YesNoMarket{
		ID: marketID,
	}
	m.MarketInterface = MarketInterface(&market)

	return nil
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

	a[1], err = m.GetMeta()
	if err != nil {
		return nil, err
	}

	return json.Marshal(a)
}

func (op *OverUnderMarket) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)
	enc.Encode(uint8(op.ID))
	enc.Encode(op.Threshold)
	return enc.Err()
}

func (op *OverUnderMarket) GetName() string {
	return MarketNames[op.ID]
}

func (op *OverUnderMarket) GetID() int8 {
	return int8(op.ID)
}

func (op *OverUnderMarket) GetMeta() (json.RawMessage, error) {
	s := struct {
		Threshold int16 `json:"threshold"`
	}{
		Threshold: op.Threshold,
	}

	return json.Marshal(s)
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

	a[1], err = m.GetMeta()
	if err != nil {
		return nil, err
	}

	return json.Marshal(a)
}

func (op *ScoreYesNoMarket) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)
	enc.Encode(uint8(op.ID))
	enc.Encode(op.Home)
	enc.Encode(op.Away)
	return enc.Err()
}

func (op *ScoreYesNoMarket) GetName() string {
	return MarketNames[op.ID]
}

func (op *ScoreYesNoMarket) GetID() int8 {
	return int8(op.ID)
}

func (op *ScoreYesNoMarket) GetMeta() (json.RawMessage, error) {
	s := struct {
		Home uint16 `json:"home"`
		Away uint16 `json:"away"`
	}{
		Home: op.Home,
		Away: op.Away,
	}

	return json.Marshal(s)
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
	a[1], err = m.GetMeta()
	if err != nil {
		return nil, err
	}

	return json.Marshal(a)
}

func (op *YesNoMarket) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)
	enc.Encode(uint8(op.ID))
	return enc.Err()
}

func (op *YesNoMarket) GetName() string {
	return MarketNames[op.ID]
}

func (op *YesNoMarket) GetID() int8 {
	return int8(op.ID)
}

func (op *YesNoMarket) GetMeta() (json.RawMessage, error) {
	return json.Marshal(struct{}{})
}
