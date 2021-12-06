package types

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/scorum/scorum-go/encoding/transaction"
	"github.com/shopspring/decimal"
)

const Symbol = "SCR"

type Asset struct {
	d decimal.Decimal
}

func AssetFromDecimal(d decimal.Decimal) *Asset{
	return &Asset{d: d}
}

func AssetFromFloat(value float64) *Asset {
	return &Asset{d: decimal.NewFromFloat(value)}
}

func AssetFromString(value string) (*Asset, error) {
	index := strings.Index(value, Symbol)
	if index != -1 {
		if len(value) == len(Symbol)+index {
			value = value[0 : index-1]
		} else {
			return nil, errors.New(fmt.Sprintf("can't convert %s to asset", value))
		}
	}
	d, err := decimal.NewFromString(value)
	if err != nil {
		return nil, err
	}
	return &Asset{d: d}, nil
}

func (as Asset) String() string {
	return fmt.Sprintf("%s %s", as.d.StringFixed(9), Symbol)
}

func (as Asset) Decimal() decimal.Decimal {
	return as.d
}

func (as Asset) MarshalJSON() ([]byte, error) {
	return []byte(`"` + as.String() + `"`), nil
}

func (as *Asset) UnmarshalJSON(data []byte) (err error) {
	str := string(data)
	index := strings.Index(str, Symbol)
	if index == -1 {
		return errors.New("asset does not contain " + Symbol)
	}

	val := str[1 : index-1]
	as.d, err = decimal.NewFromString(val)
	return
}

func (as Asset) MarshalTransaction(encoder *transaction.Encoder) error {
	return encoder.Encode(as.String())
}
