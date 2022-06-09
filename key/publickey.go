package key

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"

	"github.com/scorum/scorum-go/encoding/transaction"
)

const (
	publicKeyPrefix = "SCR"
	checkSumLen     = 4
	keyLen          = 53
)

var (
	ErrWrongPrefix   = errors.New("wrong prefix")
	ErrWrongChecksum = errors.New("wrong check sum")
	ErrWrongKeyLen   = errors.New("wrong key len")
	ErrKeyMismatch   = errors.New("key mismatch")
)

type PublicKey struct {
	raw *btcec.PublicKey
}

func NewPublicKey(pubKey string) (*PublicKey, error) {
	if len(pubKey) != keyLen {
		return nil, ErrWrongKeyLen
	}

	if !strings.HasPrefix(pubKey, publicKeyPrefix) {
		return nil, ErrWrongPrefix
	}

	keyWithChecksum := base58.Decode(pubKey[len(publicKeyPrefix):])

	key := keyWithChecksum[:len(keyWithChecksum)-checkSumLen]

	h := ripemd160.New()
	h.Write(key)
	hash := h.Sum(nil)

	if !bytes.Equal(hash[:checkSumLen], keyWithChecksum[len(keyWithChecksum)-checkSumLen:]) {
		return nil, ErrWrongChecksum
	}

	return PublicKeyFromBytes(key)
}

func PublicKeyFromBytes(key []byte) (*PublicKey, error) {
	p, err := btcec.ParsePubKey(key, btcec.S256())
	if err != nil {
		return nil, err
	}

	return &PublicKey{raw: p}, nil
}

func (p *PublicKey) String() string {
	h := ripemd160.New()
	h.Write(p.raw.SerializeCompressed())
	hash := h.Sum(nil)

	b := bytes.Buffer{}
	b.Write(p.raw.SerializeCompressed())
	b.Write(hash[:checkSumLen])

	return publicKeyPrefix + base58.Encode(b.Bytes())
}

func (p *PublicKey) Serialize() []byte {
	return p.raw.SerializeCompressed()
}

func (p *PublicKey) Verify(digest []byte, signature []byte) error {
	pub, _, err := btcec.RecoverCompact(btcec.S256(), signature, digest)
	if err != nil {
		return fmt.Errorf("recover compact: %w", err)
	}

	if !p.raw.IsEqual(pub) {
		return ErrKeyMismatch
	}

	return nil
}

func (p *PublicKey) MarshalTransaction(encoder *transaction.Encoder) error {
	if p.raw == nil {
		return errors.New("public key is nil")
	}

	enc := transaction.NewRollingEncoder(encoder)
	enc.Encode(p.raw.SerializeCompressed())
	return enc.Err()
}

func (p *PublicKey) MarshalText() (text []byte, err error) {
	return []byte(p.String()), nil
}

func (p *PublicKey) UnmarshalText(text []byte) error {
	pk, err := NewPublicKey(string(text))
	if err != nil {
		return err
	}

	p.raw = pk.raw

	return nil
}
