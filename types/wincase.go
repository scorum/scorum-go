package types

import (
	"encoding/json"
	"errors"
	"github.com/bitly/go-simplejson"
	"github.com/scorum/scorum-go/encoding/transaction"
)

const (
	WincaseResultHomeYes WincaseID = iota
	WincaseResultHomeNo
	WincaseResultDrawYes
	WincaseResultDrawNo
	WincaseResultAwayYes
	WincaseResultAwayNo
	WincaseRoundHomeYes
	WincaseRoundHomeNo
	WincaseHandicapOver
	WincaseHandicapUnder
	WincaseCorrectScoreHomeYes
	WincaseCorrectScoreHomeNo
	WincaseCorrectScoreDrawYes
	WincaseCorrectScoreDrawNo
	WincaseCorrectScoreAwayYes
	WincaseCorrectScoreAwayNo
	WincaseCorrectScoreYes
	WincaseCorrectScoreNo
	WincaseGoalHomeYes
	WincaseGoalHomeNo
	WincaseGoalBothYes
	WincaseGoalBothNo
	WincaseGoalAwayYes
	WincaseGoalAwayNo
	WincaseTotalOver
	WincaseTotalUnder
)

var (
	errUnknownWincase = errors.New("unknown wincase id")
)

var WincaseNames = map[WincaseID]string{
	WincaseResultHomeYes:       "result_home::yes",
	WincaseResultHomeNo:        "result_home::no",
	WincaseResultDrawYes:       "result_draw::yes",
	WincaseResultDrawNo:        "result_draw::no",
	WincaseResultAwayYes:       "result_away::yes",
	WincaseResultAwayNo:        "result_away::no",
	WincaseRoundHomeYes:        "round_home::yes",
	WincaseRoundHomeNo:         "round_home::no",
	WincaseHandicapOver:        "handicap::over",
	WincaseHandicapUnder:       "handicap::under",
	WincaseCorrectScoreHomeYes: "correct_score_home::yes",
	WincaseCorrectScoreHomeNo:  "correct_score_home::no",
	WincaseCorrectScoreDrawYes: "correct_score_draw::yes",
	WincaseCorrectScoreDrawNo:  "correct_score_draw::no",
	WincaseCorrectScoreAwayYes: "correct_score_away::yes",
	WincaseCorrectScoreAwayNo:  "correct_score_away::no",
	WincaseCorrectScoreYes:     "correct_score::yes",
	WincaseCorrectScoreNo:      "correct_score::no",
	WincaseGoalHomeYes:         "goal_home::yes",
	WincaseGoalHomeNo:          "goal_home::no",
	WincaseGoalBothYes:         "goal_both::yes",
	WincaseGoalBothNo:          "goal_both::no",
	WincaseGoalAwayYes:         "goal_away::yes",
	WincaseGoalAwayNo:          "goal_away::no",
	WincaseTotalOver:           "total::over",
	WincaseTotalUnder:          "total::under",
}

type Wincase struct {
	WincaseInterface
}

type WincaseInterface interface {
	transaction.TransactionMarshaller

	GetName() string
	GetID() int8
	GetMeta() (json.RawMessage, error)
}

func (w Wincase) MarshalJSON() ([]byte, error) {
	return json.Marshal(w.WincaseInterface)
}

func (w *Wincase) UnmarshalJSON(b []byte) error {
	json, err := simplejson.NewJson(b)
	if err != nil {
		return err
	}

	wincaseName, err := json.GetIndex(0).String()
	if err != nil {
		return err
	}

	wincaseID := WincaseID(-1)
	for k, v := range WincaseNames {
		if v == wincaseName {
			wincaseID = k
		}
	}
	if wincaseID == -1 {
		return errUnknownWincase
	}

	wincaseObj := json.GetIndex(1)

	threshold, err := wincaseObj.Get("threshold").Int()
	if err == nil {
		wincase := OverUnderWincase{
			ID:        wincaseID,
			Threshold: int16(threshold),
		}
		w.WincaseInterface = WincaseInterface(&wincase)
		return nil
	}

	home, err := wincaseObj.Get("home").Int()
	if err == nil {
		away, err := wincaseObj.Get("away").Int()
		if err == nil {
			wincase := ScoreYesNoWincase{
				ID:   wincaseID,
				Home: uint16(home),
				Away: uint16(away),
			}
			w.WincaseInterface = WincaseInterface(&wincase)
			return nil
		}
	}

	wincase := YesNoWincase{
		ID: wincaseID,
	}
	w.WincaseInterface = WincaseInterface(&wincase)

	return nil
}

type WincaseID int8

type OverUnderWincase struct {
	ID WincaseID

	Threshold int16
}

func (w OverUnderWincase) MarshalJSON() ([]byte, error) {
	var err error

	a := make([]json.RawMessage, 2)
	a[0], err = json.Marshal(WincaseNames[w.ID])
	if err != nil {
		return nil, err
	}

	a[1], err = w.GetMeta()
	if err != nil {
		return nil, err
	}

	return json.Marshal(a)
}

func (op *OverUnderWincase) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)
	enc.Encode(int8(op.ID))
	enc.Encode(op.Threshold)
	return enc.Err()
}

func (op *OverUnderWincase) GetName() string {
	return WincaseNames[op.ID]
}

func (op *OverUnderWincase) GetID() int8 {
	return int8(op.ID)
}

func (op *OverUnderWincase) GetMeta() (json.RawMessage, error) {
	s := struct {
		Threshold int16 `json:"threshold"`
	}{
		Threshold: op.Threshold,
	}

	return json.Marshal(s)
}

type ScoreYesNoWincase struct {
	ID WincaseID

	Home uint16
	Away uint16
}

func (w ScoreYesNoWincase) MarshalJSON() ([]byte, error) {
	var err error

	a := make([]json.RawMessage, 2)
	a[0], err = json.Marshal(WincaseNames[w.ID])
	if err != nil {
		return nil, err
	}

	a[1], err = w.GetMeta()
	if err != nil {
		return nil, err
	}

	return json.Marshal(a)
}

func (op *ScoreYesNoWincase) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)
	enc.Encode(int8(op.ID))
	enc.Encode(op.Home)
	enc.Encode(op.Away)
	return enc.Err()
}

func (op *ScoreYesNoWincase) GetName() string {
	return WincaseNames[op.ID]
}

func (op *ScoreYesNoWincase) GetID() int8 {
	return int8(op.ID)
}

func (op *ScoreYesNoWincase) GetMeta() (json.RawMessage, error) {
	s := struct {
		Home uint16 `json:"home"`
		Away uint16 `json:"away"`
	}{
		Home: op.Home,
		Away: op.Away,
	}

	return json.Marshal(s)
}

type YesNoWincase struct {
	ID WincaseID
}

func (w YesNoWincase) MarshalJSON() ([]byte, error) {
	var err error

	a := make([]json.RawMessage, 2)
	a[0], err = json.Marshal(WincaseNames[w.ID])
	if err != nil {
		return nil, err
	}
	a[1], err = w.GetMeta()
	if err != nil {
		return nil, err
	}

	return json.Marshal(a)
}

func (op *YesNoWincase) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)
	enc.Encode(int8(op.ID))
	return enc.Err()
}

func (op *YesNoWincase) GetName() string {
	return WincaseNames[op.ID]
}

func (op *YesNoWincase) GetID() int8 {
	return int8(op.ID)
}

func (op *YesNoWincase) GetMeta() (json.RawMessage, error) {
	return json.Marshal(struct{}{})
}
