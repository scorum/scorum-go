package types

import "github.com/scorum/scorum-go/encoding/transaction"

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

type Market interface {
	transaction.TransactionMarshaller
}

type MarketID int8

type OverUnderMarket struct {
	ID MarketID

	Threshold int16
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

func (op *YesNoMarket) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)
	enc.Encode(int8(op.ID))
	return enc.Err()
}
