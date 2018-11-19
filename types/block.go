package types

import (
	"encoding/json"
	"fmt"

	"github.com/bitly/go-simplejson"
)

type BlockHeader struct {
	TransactionMerkleRoot string            `json:"transaction_merkle_root"`
	Previous              string            `json:"previous"`
	Timestamp             Time              `json:"timestamp"`
	Witness               string            `json:"witness"`
	Extensions            []json.RawMessage `json:"extensions"`
}

type Block struct {
	Previous              string            `json:"previous"`
	BlockID               string            `json:"block_id"`
	WitnessSignature      string            `json:"witness_signature"`
	SigningKey            string            `json:"signing_key"`
	TransactionIDs        []string          `json:"transaction_ids"`
	Timestamp             string            `json:"timestamp"`
	Witness               string            `json:"witness"`
	TransactionMerkleRoot string            `json:"transaction_merkle_root"`
	Transactions          []Transaction     `json:"transactions"`
	Extensions            []json.RawMessage `json:"extensions"`
	Signatures            []string          `json:"signatures"`
}

type OperationsBlock struct {
	BlockNum              uint32            `json:"block_num"`
	Previous              string            `json:"previous"`
	WitnessSignature      string            `json:"witness_signature"`
	Timestamp             string            `json:"timestamp"`
	Witness               string            `json:"witness"`
	TransactionMerkleRoot string            `json:"transaction_merkle_root"`
	Operations            []OperationInfo   `json:"operations"`
	Extensions            []json.RawMessage `json:"extensions"`
}

type OperationInfo struct {
	Operation     Operation
	Timestamp     Time
	TransactionID string
}

func (oi *OperationInfo) UnmarshalJSON(b []byte) error {
	j, err := simplejson.NewJson(b)
	if err != nil {
		return err
	}

	tb, err := j.Get("timestamp").String()
	if err != nil {
		return err
	}

	if err := oi.Timestamp.UnmarshalJSON([]byte(fmt.Sprintf(`"%s"`, tb))); err != nil {
		return err
	}

	oi.TransactionID, err = j.Get("trx_id").String()
	if err != nil {
		return err
	}

	tuple := j.GetPath("op")

	opType, err := tuple.GetIndex(0).String()
	if err != nil {
		return err
	}

	data, err := tuple.GetIndex(1).Encode()
	if err != nil {
		return err
	}

	oi.Operation, err = unmarshalOperation(opType, data)
	if err != nil {
		return err
	}

	return nil
}
