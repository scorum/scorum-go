package types

import (
	"github.com/scorum/scorum-go/encoding/transaction"
	"github.com/scorum/scorum-go/key"
)

type PublicKey string

func (k PublicKey) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)

	pubKey, err := key.NewPublicKey(string(k))
	if err != nil {
		return err
	}
	enc.Encode(pubKey.Serialize())
	return enc.Err()
}
