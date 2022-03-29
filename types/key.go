package types

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"

	// nolint
	"golang.org/x/crypto/ripemd160"

	"github.com/scorum/scorum-go/encoding/transaction"
	"github.com/scorum/scorum-go/internal"
)

const (
	PublicKeyPrefix   = "SCR"
	checkSumLen       = 4
	brainKeyWordCount = 16
)

var (
	ErrWrongPrefix   = errors.New("wrong prefix")
	ErrWrongChecksum = errors.New("wrong check sum")
)

func PublicKeyToString(key *btcec.PublicKey) string {
	h := ripemd160.New()
	h.Write(key.SerializeCompressed())
	hash := h.Sum(nil)

	b := bytes.Buffer{}
	b.Write(key.SerializeCompressed())
	b.Write(hash[:checkSumLen])

	return PublicKeyPrefix + base58.Encode(b.Bytes())
}

func NewPublicKey(pubKey string) (*btcec.PublicKey, error) {
	if !strings.HasPrefix(pubKey, PublicKeyPrefix) {
		return nil, ErrWrongPrefix
	}
	keyWithChecksum := base58.Decode(pubKey[len(PublicKeyPrefix):])
	key := keyWithChecksum[:len(keyWithChecksum)-checkSumLen]

	h := ripemd160.New()
	h.Write(key)
	hash := h.Sum(nil)

	if !bytes.Equal(hash[:checkSumLen], keyWithChecksum[len(keyWithChecksum)-checkSumLen:]) {
		return nil, ErrWrongChecksum
	}

	return btcec.ParsePubKey(key, btcec.S256())
}

func Verify(digest []byte, signature []byte, publicKey *btcec.PublicKey) (bool, error) {
	p, _, err := btcec.RecoverCompact(btcec.S256(), signature, digest)
	if err != nil {
		return false, err
	}
	return p.IsEqual(publicKey), nil
}

func GenerateBrainKey() (res string) {
	max := len(internal.Words)

	r, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
	res = internal.Words[r.Int64()]

	for i := 0; i < brainKeyWordCount-1; i++ {
		r, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
		res += fmt.Sprintf(" %s", internal.Words[r.Int64()])
	}

	return
}

func GenerateWIF() (*btcutil.WIF, error) {
	pk, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return nil, err
	}

	return btcutil.NewWIF(pk, &chaincfg.Params{PrivateKeyID: 128}, false)
}

func WIFFromBytes(b []byte) (*btcutil.WIF, error) {
	pk, _ := btcec.PrivKeyFromBytes(btcec.S256(), b)
	return btcutil.NewWIF(pk, &chaincfg.Params{PrivateKeyID: 128}, false)
}

func WIFFromBrainKey(brainKey string) (*btcutil.WIF, error) {
	hash := sha256.Sum256([]byte(brainKey))
	pk, _ := btcec.PrivKeyFromBytes(btcec.S256(), hash[:])
	return btcutil.NewWIF(pk, &chaincfg.Params{PrivateKeyID: 128}, false)
}

type PublicKey string

func (k PublicKey) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)

	pubKey, err := NewPublicKey(string(k))
	if err != nil {
		return err
	}
	enc.Encode(pubKey.SerializeCompressed())
	return enc.Err()
}
