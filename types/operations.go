package types

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/scorum/scorum-go/encoding/transaction"
)

type Operation interface {
	Type() OpType
}

// []Operations coming from the Api in the following form: [["op1", {}], ["op2", {}], ...]
type OperationsArray []Operation

func (ops OperationsArray) MarshalJSON() ([]byte, error) {
	tuples := make([]*operationTuple, 0, len(ops))
	for _, op := range ops {
		tuples = append(tuples, &operationTuple{
			Type: op.Type(),
			Data: op,
		})
	}
	return json.Marshal(tuples)
}

type operationTuple struct {
	Type OpType
	Data Operation
}

func (op *operationTuple) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{
		op.Type,
		op.Data,
	})
}

type OperationObject struct {
	BlockNumber             uint32         `json:"block"`
	TransactionID           string         `json:"trx_id"`
	TransactionsInBlock     uint32         `json:"trx_in_block"`
	OperationsInTransaction uint32         `json:"op_in_trx"`
	VirtualOperations       uint32         `json:"virtual_op"`
	Timestamp               Time           `json:"timestamp"`
	Operations              OperationsFlat `json:"op"`
}

func (ops *OperationsArray) UnmarshalJSON(b []byte) (err error) {
	// unmarshal array
	var o []json.RawMessage
	if err := json.Unmarshal(b, &o); err != nil {
		return err
	}

	// foreach operation
	for _, op := range o {
		var kv []json.RawMessage
		if err := json.Unmarshal(op, &kv); err != nil {
			return err
		}

		if len(kv) != 2 {
			return errors.New("invalid operation format: should be name, value")
		}

		var key string
		if err := json.Unmarshal(kv[0], &key); err != nil {
			return err
		}

		val, err := unmarshalOperation(key, kv[1])
		if err != nil {
			return err
		}

		*ops = append(*ops, val)
	}

	return nil
}

// []Operations coming from the Api in the following form: ["op1", {}, "op2", {}, ...]
type OperationsFlat []Operation

func (t *OperationsFlat) UnmarshalJSON(b []byte) (err error) {

	// unmarshal array
	var o []json.RawMessage
	if err := json.Unmarshal(b, &o); err != nil {
		return err
	}

	for i := 0; i < len(o); i += 2 {
		var key string
		if err := json.Unmarshal(o[i], &key); err != nil {
			return err
		}
		val, err := unmarshalOperation(key, o[i+1])
		if err != nil {
			return err
		}

		*t = append(*t, val)
	}
	return nil
}

func unmarshalOperation(key string, obj json.RawMessage) (Operation, error) {
	opType, ok := knownOperations[OpType(key)]
	if !ok {
		// operation is unknown wrap it as a general operation
		val := UnknownOperation{
			kind: OpType(key),
			Data: obj,
		}
		return &val, nil
	} else {
		val := reflect.New(opType).Interface()
		if err := json.Unmarshal(obj, val); err != nil {
			return nil, err
		}
		return val.(Operation), nil
	}
}

var knownOperations = map[OpType]reflect.Type{
	AccountCreateOpType:               reflect.TypeOf(AccountCreateOperation{}),
	TransferToScorumpowerOpType:       reflect.TypeOf(TransferToScorumpowerOperation{}),
	AccountWitnessVoteOpType:          reflect.TypeOf(AccountWitnessVoteOperation{}),
	WitnessUpdateOpType:               reflect.TypeOf(WitnessUpdateOperation{}),
	AccountCreateByCommitteeOpType:    reflect.TypeOf(AccountCreateByCommitteeOperation{}),
	AccountCreateWithDelegationOpType: reflect.TypeOf(AccountCreateWithDelegationOperation{}),
	AccountUpdateOpType:               reflect.TypeOf(AccountUpdateOperation{}),
	TransferOpType:                    reflect.TypeOf(TransferOperation{}),
	ProducerRewardOpType:              reflect.TypeOf(ProducerRewardOperation{}),
	CommentOptionsOpType:              reflect.TypeOf(CommentOptionsOperation{}),
	CommentOpType:                     reflect.TypeOf(CommentOperation{}),
	DeleteCommentOpType:               reflect.TypeOf(DeleteCommentOperation{}),
	VoteOpType:                        reflect.TypeOf(VoteOperation{}),
	WithdrawScorumpowerOpType:         reflect.TypeOf(WithdrawScorumpowerOperation{}),
	DelegateScorumpower:               reflect.TypeOf(DelegateScorumpowerOperation{}),
	CreateGame:                        reflect.TypeOf(CreateGameOperation{}),
	CancelGame:                        reflect.TypeOf(CancelGameOperation{}),
	UpdateGameStartTime:               reflect.TypeOf(UpdateGameStartTimeOperation{}),
	PostGameResults:                   reflect.TypeOf(PostGameResultsOperation{}),
	PostBet:                           reflect.TypeOf(PostBetOperation{}),
	CancelPendingBets:                 reflect.TypeOf(CancelPendingBetsOperation{}),
	BetsMatched:                       reflect.TypeOf(BetsMatchedVirtualOperation{}),
}

// UnknownOperation
type UnknownOperation struct {
	kind OpType
	Data json.RawMessage
}

func (op *UnknownOperation) Type() OpType { return op.kind }

// AccountCreateWithDelegationOperation
type AccountCreateWithDelegationOperation struct {
	Fee            string            `json:"fee"`
	Creator        string            `json:"creator"`
	NewAccountName string            `json:"new_account_name"`
	Owner          Authority         `json:"owner"`
	Active         Authority         `json:"active"`
	Posting        Authority         `json:"posting"`
	MemoKey        string            `json:"memo_key"`
	JsonMetadata   string            `json:"json_metadata"`
	Extensions     []json.RawMessage `json:"extensions"`
}

func (op *AccountCreateWithDelegationOperation) Type() OpType {
	return AccountCreateWithDelegationOpType
}

// AccountCreateByCommitteeOperation
type AccountCreateByCommitteeOperation struct {
	Creator        string    `json:"creator"`
	NewAccountName string    `json:"new_account_name"`
	Owner          Authority `json:"owner"`
	Active         Authority `json:"active"`
	Posting        Authority `json:"posting"`
	MemoKey        string    `json:"memo_key"`
	JsonMetadata   string    `json:"json_metadata"`
}

func (op *AccountCreateByCommitteeOperation) Type() OpType { return AccountCreateByCommitteeOpType }

// TransferToScorumpowerOperation
type TransferToScorumpowerOperation struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount string `json:"amount"`
}

func (op *TransferToScorumpowerOperation) Type() OpType { return TransferToScorumpowerOpType }

// AccountCreateOperation
type AccountCreateOperation struct {
	Fee            string    `json:"fee"`
	Creator        string    `json:"creator"`
	NewAccountName string    `json:"new_account_name"`
	Owner          Authority `json:"owner"`
	Active         Authority `json:"active"`
	Posting        Authority `json:"posting"`
	MemoKey        string    `json:"memo_key"`
	JsonMetadata   string    `json:"json_metadata"`
}

func (op *AccountCreateOperation) Type() OpType { return AccountCreateOpType }

// AccountWitnessVoteOperation
type AccountWitnessVoteOperation struct {
	Account string `json:"account"`
	Witness string `json:"witness"`
	Approve bool   `json:"approve"`
}

func (op *AccountWitnessVoteOperation) Type() OpType { return AccountWitnessVoteOpType }

func (op *AccountWitnessVoteOperation) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)
	enc.EncodeUVarint(uint64(op.Type().Code()))
	enc.Encode(op.Account)
	enc.Encode(op.Witness)
	enc.EncodeBool(op.Approve)
	return enc.Err()
}

// WitnessUpdateOperation
type WitnessUpdateOperation struct {
	Owner           string                      `json:"owner"`
	Url             string                      `json:"url"`
	BlockSigningKey string                      `json:"block_signing_key"`
	Props           WitnessUpdateOperationProps `json:"props"`
	Fee             string                      `json:"fee"`
}

func (op *WitnessUpdateOperation) Type() OpType { return WitnessUpdateOpType }

type WitnessUpdateOperationProps struct {
	AccountCreationFee string `json:"account_creation_fee"`
	MaximumBlockSize   int32  `json:"maximum_block_size"`
}

// TransferOperation
type TransferOperation struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount Asset  `json:"amount"`
	Memo   string `json:"memo"`
}

func (op *TransferOperation) Type() OpType { return TransferOpType }

func (op *TransferOperation) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)
	enc.EncodeUVarint(uint64(op.Type().Code()))
	enc.Encode(op.From)
	enc.Encode(op.To)
	enc.EncodeMoney(op.Amount.String())
	enc.Encode(op.Memo)
	return enc.Err()
}

// Equal returns whether the numbers represented by d and d2 are equal.
func (op TransferOperation) Equals(t2 TransferOperation) bool {
	return op.To == t2.To &&
		op.From == t2.From &&
		op.Memo == t2.Memo &&
		op.Amount.Decimal().Equals(t2.Amount.Decimal())
}

// VoteOperation
type VoteOperation struct {
	Voter    string `json:"voter"`
	Author   string `json:"author"`
	Permlink string `json:"permlink"`
	Weight   int16  `json:"weight"`
}

func (op *VoteOperation) Type() OpType { return VoteOpType }

func (op *VoteOperation) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)
	enc.EncodeUVarint(uint64(op.Type().Code()))
	enc.Encode(op.Voter)
	enc.Encode(op.Author)
	enc.Encode(op.Permlink)
	enc.Encode(op.Weight)
	return enc.Err()
}

// CommentOperation represents either a new post or a comment.
// In case Title is filled in and ParentAuthor is empty, it is a new post.
// The post category can be read from ParentPermlink.
type CommentOperation struct {
	ParentAuthor   string `json:"parent_author"`
	ParentPermlink string `json:"parent_permlink"`
	Author         string `json:"author"`
	Permlink       string `json:"permlink"`
	Title          string `json:"title"`
	Body           string `json:"body"`
	JsonMetadata   string `json:"json_metadata"`
}

func (op *CommentOperation) Type() OpType {
	return CommentOpType
}

type DeleteCommentOperation struct {
	Author   string `json:"author"`
	Permlink string `json:"permlink"`
}

func (op *DeleteCommentOperation) Type() OpType {
	return DeleteCommentOpType
}

// Authors of posts may not want all of the benefits that come from creating a post. This
// operation allows authors to update properties associated with their post.
//
// The max_accepted_payout may be decreased, but never increased.
// The percent_scrs may be decreased, but never increased
type CommentOptionsOperation struct {
	Author               string        `json:"author"`
	Permlink             string        `json:"permlink"`
	MaxAcceptedPayout    string        `json:"max_accepted_payout"`
	PercentSCRs          uint16        `json:"percent_scrs"`
	AllowVotes           bool          `json:"allow_votes"`
	AllowCurationRewards bool          `json:"allow_curation_rewards"`
	Extensions           []interface{} `json:"extensions"`
}

func (op *CommentOptionsOperation) Type() OpType {
	return CommentOptionsOpType
}

type ProducerRewardOperation struct {
	Producer    string `json:"producer"`
	Scorumpower string `json:"scorumpower"`
}

func (op *ProducerRewardOperation) Type() OpType {
	return ProducerRewardOpType
}

type AccountUpdateOperation struct {
	Account      string    `json:"account"`
	Owner        Authority `json:"owner"`
	Active       Authority `json:"active"`
	Posting      Authority `json:"posting"`
	MemoKey      string    `json:"memo_key"`
	JsonMetadata string    `json:"json_metadata"`
}

func (op *AccountUpdateOperation) Type() OpType {
	return AccountUpdateOpType
}

type WithdrawScorumpowerOperation struct {
	Account     string `json:"account"`
	Scorumpower string `json:"scorumpower"`
}

func (op *WithdrawScorumpowerOperation) Type() OpType {
	return WithdrawScorumpowerOpType
}

type DelegateScorumpowerOperation struct {
	Delegator   string `json:"delegator"`
	Delegatee   string `json:"delegatee"`
	Scorumpower string `json:"scorumpower"`
}

func (op *DelegateScorumpowerOperation) Type() OpType {
	return DelegateScorumpower
}

type CreateGameOperation struct {
	Moderator           string   `json:"moderator"`
	Name                string   `json:"name"`
	GameType            int8     `json:"game_type"`
	StartTime           Time     `json:"start_time"`
	AutoResolveDelaySec uint32   `json:"auto_resolve_delay_sec"`
	Markets             []Market `json:"markets"`
}

func (op *CreateGameOperation) Type() OpType {
	return CreateGame
}

func (op *CreateGameOperation) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)
	enc.EncodeUVarint(uint64(op.Type().Code()))
	enc.Encode(op.Moderator)
	enc.Encode(op.Name)
	op.StartTime.MarshalTransaction(encoder)
	enc.Encode(op.AutoResolveDelaySec)
	enc.Encode(op.GameType)
	enc.Encode(int8(len(op.Markets)))
	for _, m := range op.Markets {
		enc.Encode(m)
	}
	return enc.Err()
}

type CancelGameOperation struct {
	GameID    int64  `json:"game_id"`
	Moderator string `json:"moderator"`
}

func (op *CancelGameOperation) Type() OpType {
	return CancelGame
}

func (op *CancelGameOperation) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)
	enc.EncodeUVarint(uint64(op.Type().Code()))
	enc.Encode(op.Moderator)
	enc.Encode(op.GameID)
	return enc.Err()
}

type UpdateGameStartTimeOperation struct {
	GameID    int64  `json:"game_id"`
	Moderator string `json:"moderator"`
	StartTime Time   `json:"start_time"`
}

func (op *UpdateGameStartTimeOperation) Type() OpType {
	return UpdateGameStartTime
}
func (op *UpdateGameStartTimeOperation) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)
	enc.EncodeUVarint(uint64(op.Type().Code()))
	enc.Encode(op.Moderator)
	enc.Encode(op.GameID)
	op.StartTime.MarshalTransaction(encoder)
	return enc.Err()
}

type PostGameResultsOperation struct {
	GameID    int64     `json:"game_id"`
	Moderator string    `json:"moderator"`
	Wincases  []Wincase `json:"wincases"`
}

func (op *PostGameResultsOperation) Type() OpType {
	return PostGameResults
}

func (op *PostGameResultsOperation) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)
	enc.EncodeUVarint(uint64(op.Type().Code()))
	enc.Encode(op.Moderator)
	enc.Encode(op.GameID)
	enc.Encode(int8(len(op.Wincases)))
	for _, m := range op.Wincases {
		enc.Encode(m)
	}
	return enc.Err()
}

type Odds struct {
	Numerator   int16 `json:"numerator"`
	Denominator int16 `json:"denominator"`
}

type PostBetOperation struct {
	Better  string  `json:"better"`
	GameID  int64   `json:"game_id"`
	Wincase Wincase `json:"wincase"`
	Odds    Odds    `json:"odds"`
	Stake   Asset   `json:"stake"`
	Live    bool    `json:"live"`
}

func (op *PostBetOperation) Type() OpType {
	return PostBet
}

type CancelPendingBetsOperation struct {
	BetIDs []int64 `json:"bet_ids"`
	Better string  `json:"better"`
}

func (op *CancelPendingBetsOperation) Type() OpType {
	return CancelPendingBets
}

type BetsMatchedVirtualOperation struct {
	Better1       string `json:"better1"`
	Better2       string `json:"better2"`
	MatchedStake1 Asset  `json:"matched_stake1"`
	MatchedStake2 Asset  `json:"matched_stake2"`
	MatchedBetID  int64  `json:"matched_bet_id"`
}

func (op *BetsMatchedVirtualOperation) Type() OpType {
	return BetsMatched
}
