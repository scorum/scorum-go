package sign

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/scorum/scorum-go/encoding/transaction"
	"github.com/scorum/scorum-go/key"
	"github.com/scorum/scorum-go/types"
)

type SignedTransaction struct {
	*types.Transaction
}

func NewSignedTransaction(tx *types.Transaction) *SignedTransaction {
	return &SignedTransaction{tx}
}

func (tx *SignedTransaction) Serialize() ([]byte, error) {
	var b bytes.Buffer
	encoder := transaction.NewEncoder(&b)

	if err := encoder.Encode(tx.Transaction); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (tx *SignedTransaction) Digest(chainID []byte) ([]byte, error) {
	var msgBuffer bytes.Buffer

	if _, err := msgBuffer.Write(chainID); err != nil {
		return nil, fmt.Errorf("failed to write chain ID: %w", err)
	}

	// Write the serialized transaction.
	rawTx, err := tx.Serialize()
	if err != nil {
		return nil, err
	}

	if _, err := msgBuffer.Write(rawTx); err != nil {
		return nil, fmt.Errorf("failed to write serialized transaction: %w", err)
	}

	// Compute the digest.
	digest := sha256.Sum256(msgBuffer.Bytes())
	return digest[:], nil
}

func (tx *SignedTransaction) Sign(chainID []byte, keys ...*key.PrivateKey) error {
	digest, err := tx.Digest(chainID)
	if err != nil {
		return err
	}

	sigsHex := make([]string, len(keys))
	for i, k := range keys {
		sig := k.Sign(digest)
		sigsHex[i] = hex.EncodeToString(sig)
	}

	tx.Transaction.Signatures = sigsHex
	return nil
}

func (tx *SignedTransaction) Verify(chainID []byte, keys ...*key.PublicKey) error {
	dig, err := tx.Digest(chainID)
	if err != nil {
		return fmt.Errorf("failed to get digest: %w", err)
	}

	for _, signature := range tx.Signatures {
		sig, err := hex.DecodeString(signature)
		if err != nil {
			return fmt.Errorf("failed to decode signature: %w", err)
		}

		for _, k := range keys {
			if err := k.Verify(dig, sig); err != nil {
				return fmt.Errorf("verify signature: %w", err)
			}
		}
	}

	return nil
}
