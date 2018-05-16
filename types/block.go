package types

import "encoding/json"

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
