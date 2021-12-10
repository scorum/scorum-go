package sign

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
	"github.com/scorum/scorum-go/encoding/transaction"
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

func (tx *SignedTransaction) Digest(chain *Chain) ([]byte, error) {
	var msgBuffer bytes.Buffer

	// Write the chain ID.
	rawChainID, err := hex.DecodeString(chain.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to decode chain ID: %v", chain.ID)
	}

	if _, err := msgBuffer.Write(rawChainID); err != nil {
		return nil, errors.Wrap(err, "failed to write chain ID")
	}

	// Write the serialized transaction.
	rawTx, err := tx.Serialize()
	if err != nil {
		return nil, err
	}

	if _, err := msgBuffer.Write(rawTx); err != nil {
		return nil, errors.Wrap(err, "failed to write serialized transaction")
	}

	// Compute the digest.
	digest := sha256.Sum256(msgBuffer.Bytes())
	return digest[:], nil
}

func (tx *SignedTransaction) Sign(wifs []string, chain *Chain) error {
	digest, err := tx.Digest(chain)
	if err != nil {
		return err
	}

	privKeys := make([]*btcec.PrivateKey, len(wifs))
	for index, wif := range wifs {
		w, err := btcutil.DecodeWIF(wif)
		if err != nil {
			return err
		}
		privKeys[index] = w.PrivKey
	}

	// Set the signature array in the transaction.
	sigsHex := make([]string, len(privKeys))
	for index, privKey := range privKeys {
		sig := SignBufferSha256(digest, privKey.ToECDSA())
		sigsHex[index] = hex.EncodeToString(sig)
	}
	tx.Transaction.Signatures = sigsHex
	return nil
}

func (tx *SignedTransaction) Verify(chain *Chain, keys [][]byte) (bool, error) {
	dig, err := tx.Digest(chain)
	if err != nil {
		return false, fmt.Errorf("failed to get digest: %w", err)
	}

	pubKeysFound := make([]*btcec.PublicKey, 0, len(tx.Signatures))
	for _, signature := range tx.Signatures {
		sig, err := hex.DecodeString(signature)
		if err != nil {
			return false, fmt.Errorf("failed to decode signature: %w", err)
		}

		p, _, err := btcec.RecoverCompact(btcec.S256(), sig, dig)
		if err != nil {
			return false, fmt.Errorf("failed to RecoverCompact: %w", err)
		}

		pubKeysFound = append(pubKeysFound, p)
	}

find:
	for _, pub := range pubKeysFound {
		for _, v := range keys {
			pb, err := btcec.ParsePubKey(v, btcec.S256())
			if err != nil {
				return false, fmt.Errorf("failed to parse pub key: %w", err)
			}

			if pub.IsEqual(pb) {
				continue find
			}
		}
		return false, nil
	}

	return true, nil
}
