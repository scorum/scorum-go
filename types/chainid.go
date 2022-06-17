package types

import "encoding/hex"

type ChainID []byte

func (p *ChainID) MarshalText() (text []byte, err error) {
	return *p, nil
}

func (p *ChainID) UnmarshalText(text []byte) error {
	chainID, err := hex.DecodeString(string(text))
	if err != nil {
		return err
	}
	*p = chainID
	return nil
}
