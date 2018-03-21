package types

import (
	"encoding/json"
	"errors"
	"reflect"
)

type Operation interface {
	GetName() string
}

type OperationsArray []Operation

type OperationObject struct {
	BlockNumber uint32         `json:"block"`
	TrxID       string         `json:"trx_id"`
	TrxInBlock  uint32         `json:"trx_in_block"`
	OpInTrx     uint32         `json:"op_in_trx"`
	VirtualOp   uint32         `json:"virtual_op"`
	Timestamp   Time           `json:"timestamp"`
	Operations  OperationsFlat `json:"op"`
}

func (t *OperationsArray) UnmarshalJSON(b []byte) (err error) {
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

		*t = append(*t, val)
	}

	return nil
}

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
	opType, ok := knownOperations[key]
	if !ok {
		// operation is unknown wrap it as a general operation
		val := GeneralOperation{
			Name: key,
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

var knownOperations = map[string]reflect.Type{
	"account_create":                 reflect.TypeOf(AccountCreateOperation{}),
	"transfer_to_vesting":            reflect.TypeOf(TransferToVestingOperation{}),
	"account_witness_vote":           reflect.TypeOf(AccountWitnessVoteOperation{}),
	"witness_update":                 reflect.TypeOf(WitnessUpdateOperation{}),
	"account_create_by_committee":    reflect.TypeOf(AccountCreateByCommitteeOperation{}),
	"account_create_with_delegation": reflect.TypeOf(AccountCreateWithDelegationOperation{}),
	"transfer_operation":             reflect.TypeOf(TransferOperation{}),
}

// GeneralOperation
type GeneralOperation struct {
	Name string
	Data json.RawMessage
}

func (op *GeneralOperation) GetName() string { return op.Name }

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

func (op *AccountCreateWithDelegationOperation) GetName() string {
	return "account_create_with_delegation"
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

func (op *AccountCreateByCommitteeOperation) GetName() string { return "account_create_by_committee" }

// TransferToVestingOperation
type TransferToVestingOperation struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount string `json:"amount"`
}

func (op *TransferToVestingOperation) GetName() string { return "transfer_to_vesting" }

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

func (op *AccountCreateOperation) GetName() string { return "account_create" }

// AccountWitnessVoteOperation
type AccountWitnessVoteOperation struct {
	Account string `json:"account"`
	Witness string `json:"witness"`
	Approve bool   `json:"approve"`
}

func (op *AccountWitnessVoteOperation) GetName() string { return "account_witness_vote" }

// WitnessUpdateOperation
type WitnessUpdateOperation struct {
	Owner           string                      `json:"owner"`
	Url             string                      `json:"url"`
	BlockSigningKey string                      `json:"block_signing_key"`
	Props           WitnessUpdateOperationProps `json:"props"`
	Fee             string                      `json:"fee"`
}

func (op *WitnessUpdateOperation) GetName() string { return "witness_update" }

type WitnessUpdateOperationProps struct {
	AccountCreationFee string `json:"account_creation_fee"`
	MaximumBlockSize   int32  `json:"maximum_block_size"`
}

// TransferOperation
type TransferOperation struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount string `json:"amount"`
	Memo   string `json:"memo"`
}

func (op *TransferOperation) GetName() string { return "transfer_operation" }
