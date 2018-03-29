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
		val, err := unmarshalOperation(key, o[1])
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
	TransferOpType:                    reflect.TypeOf(TransferOperation{}),
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
