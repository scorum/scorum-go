package types

import (
	"fmt"
	"strings"

	"github.com/scorum/scorum-go/encoding/transaction"
	"github.com/shopspring/decimal"
)

const Symbol = "SCR"

type Asset struct {
	d decimal.Decimal
}

func AssetFromDecimal(d decimal.Decimal) *Asset {
	return &Asset{d: d}
}

func AssetFromFloat(value float64) *Asset {
	return &Asset{d: decimal.NewFromFloat(value)}
}

func AssetFromString(value string) (*Asset, error) {
	var a Asset
	if err := a.UnmarshalText([]byte(value)); err != nil {
		return nil, err
	}
	return &a, nil
}

func (as Asset) String() string {
	return fmt.Sprintf("%s %s", as.d.StringFixed(9), Symbol)
}

func (as Asset) Decimal() decimal.Decimal {
	return as.d
}

func (as *Asset) MarshalText() (text []byte, err error) {
	return []byte(as.String()), nil
}

func (as *Asset) UnmarshalText(data []byte) error {
	value := string(data)
	index := strings.Index(value, Symbol)
	if index != -1 {
		if len(value) == len(Symbol)+index {
			value = value[0 : index-1]
		} else {
			return fmt.Errorf("can't convert %s to asset", value)
		}
	}
	d, err := decimal.NewFromString(value)
	if err != nil {
		return err
	}
	as.d = d
	return nil
}

func (as Asset) MarshalTransaction(encoder *transaction.Encoder) error {
	return encoder.Encode(as.String())
}
