package types

import (
	"bytes"
	"crypto/sha256"
	"errors"

	"github.com/scorum/scorum-go/encoding/transaction"
)

type Transaction struct {
	RefBlockNum    uint16          `json:"ref_block_num"`
	RefBlockPrefix uint32          `json:"ref_block_prefix"`
	Expiration     *Time           `json:"expiration"`
	Operations     OperationsArray `json:"operations"`
	Signatures     []string        `json:"signatures"`
}

func (tx *Transaction) ID() ([]byte, error) {
	var b bytes.Buffer
	encoder := transaction.NewEncoder(&b)
	if err := tx.MarshalTransaction(encoder); err != nil {
		return nil, err
	}
	h := sha256.Sum256(b.Bytes())
	return h[:20], nil
}

// MarshalTransaction implements transaction.Marshaller interface.
func (tx *Transaction) MarshalTransaction(encoder *transaction.Encoder) error {
	if len(tx.Operations) == 0 {
		return errors.New("no operation specified")
	}

	enc := transaction.NewRollingEncoder(encoder)

	enc.Encode(tx.RefBlockNum)
	enc.Encode(tx.RefBlockPrefix)
	enc.Encode(tx.Expiration)

	enc.EncodeUVarint(uint64(len(tx.Operations)))
	for _, op := range tx.Operations {
		enc.Encode(op)
	}

	// Extensions are not supported yet.
	enc.EncodeUVarint(0)

	return enc.Err()
}

// PushOperation can be used to add an operation into the transaction.
func (tx *Transaction) PushOperation(op Operation) {
	tx.Operations = append(tx.Operations, op)
}
