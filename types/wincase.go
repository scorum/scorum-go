package types

import "github.com/scorum/scorum-go/encoding/transaction"

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
	WincaseTotalGoalsHomeOver
	WincaseTotalGoalsHomeUnder
	WincaseTotalGoalsAwayOver
	WincaseTotalGoalsAwayUnder
)

type Wincase interface {
	transaction.TransactionMarshaller
}

type WincaseID int8

type OverUnderWincase struct {
	ID WincaseID

	Threshold int16
}

func (op *OverUnderWincase) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)
	enc.Encode(int8(op.ID))
	enc.Encode(op.Threshold)
	return enc.Err()
}

type ScoreYesNoWincase struct {
	ID WincaseID

	Home uint16
	Away uint16
}

func (op *ScoreYesNoWincase) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)
	enc.Encode(int8(op.ID))
	enc.Encode(op.Home)
	enc.Encode(op.Away)
	return enc.Err()
}

type YesNoWincase struct {
	ID WincaseID
}

func (op *YesNoWincase) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)
	enc.Encode(int8(op.ID))
	return enc.Err()
}
