package key

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
)

type PrivateKey struct {
	wif *btcutil.WIF
}

func NewPrivateKey() (*PrivateKey, error) {
	raw, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return nil, fmt.Errorf("new private key: %w", err)
	}

	wif, err := btcutil.NewWIF(raw, &chaincfg.Params{PrivateKeyID: 128}, false)
	if err != nil {
		return nil, fmt.Errorf("new wif: %w", err)
	}

	return &PrivateKey{wif: wif}, nil
}

func PrivateKeyFromString(wif string) (*PrivateKey, error) {
	w, err := btcutil.DecodeWIF(wif)
	if err != nil {
		return nil, fmt.Errorf("decode wif: %w", err)
	}

	return &PrivateKey{wif: w}, nil
}

func PrivateKeyFromBytes(privKey []byte) (*PrivateKey, error) {
	raw, _ := btcec.PrivKeyFromBytes(btcec.S256(), privKey)
	wif, err := btcutil.NewWIF(raw, &chaincfg.Params{PrivateKeyID: 128}, false)
	if err != nil {
		return nil, fmt.Errorf("new wif: %w", err)
	}

	return &PrivateKey{wif: wif}, nil
}

func (p *PrivateKey) Serialize() []byte {
	return p.wif.PrivKey.Serialize()
}

func (p *PrivateKey) String() string {
	return p.wif.String()
}

func (p *PrivateKey) PublicKey() *PublicKey {
	return &PublicKey{
		raw: p.wif.PrivKey.PubKey(),
	}
}

func (p *PrivateKey) Sign(hash []byte) []byte {
	return SignBufferSha256(hash, p.wif.PrivKey.ToECDSA())
}

func (p *PrivateKey) MarshalText() (text []byte, err error) {
	return []byte(p.wif.String()), nil
}

func (p *PrivateKey) UnmarshalText(text []byte) error {
	w, err := btcutil.DecodeWIF(string(text))
	if err != nil {
		return fmt.Errorf("decode wif: %w", err)
	}
	p.wif = w
	return nil
}
