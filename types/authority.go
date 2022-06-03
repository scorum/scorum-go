package types

import (
	"encoding/json"
	"errors"

	"github.com/elliotchance/orderedmap"
	"github.com/scorum/scorum-go/encoding/transaction"
)

type Authority struct {
	WeightThreshold uint32               `json:"weight_threshold"`
	AccountAuths    *AccountAuthorityMap `json:"account_auths"`
	KeyAuths        *KeyAuthorityMap     `json:"key_auths"`
}

func (m Authority) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)
	enc.Encode(m.WeightThreshold)
	enc.Encode(m.AccountAuths)
	enc.Encode(m.KeyAuths)
	return enc.Err()
}

type KeyAuthority struct {
	Key    PublicKey
	Weight uint16
}

type KeyAuthorityMap struct {
	*orderedmap.OrderedMap
}

func NewKeyAuthorityMap(items ...KeyAuthority) *KeyAuthorityMap {
	orderedMap := orderedmap.NewOrderedMap()

	for _, v := range items {
		orderedMap.Set(v.Key, v.Weight)
	}

	return &KeyAuthorityMap{
		OrderedMap: orderedMap,
	}
}

func (m *KeyAuthorityMap) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)
	enc.Encode(uint8(m.Len()))
	for el := m.Front(); el != nil; el = el.Next() {
		enc.Encode(el.Key)
		enc.Encode(el.Value)
	}
	return enc.Err()
}

func (m KeyAuthorityMap) MarshalJSON() ([]byte, error) {
	xs := make([]interface{}, 0, m.Len())
	for el := m.Front(); el != nil; el = el.Next() {
		xs = append(xs, []interface{}{el.Key, el.Value})
	}

	return json.Marshal(xs)
}

func (m *KeyAuthorityMap) UnmarshalJSON(data []byte) error {
	var xs [][]interface{}
	if err := json.Unmarshal(data, &xs); err != nil {
		return err
	}

	orderedMap := orderedmap.NewOrderedMap()

	var invalid bool
	for _, kv := range xs {
		if len(kv) != 2 {
			invalid = true
			break
		}

		k, ok := kv[0].(string)
		if !ok {
			invalid = true
			break
		}

		var v uint16
		switch t := kv[1].(type) {
		case uint16:
			v = t
		case float64:
			v = uint16(t)
		default:
			invalid = true
			break
		}

		orderedMap.Set(PublicKey(k), v)
	}
	if invalid {
		return errors.New("invalid map encoding")
	}

	m.OrderedMap = orderedMap

	return nil
}

type AccountAuthority struct {
	AccountName string
	Weight      uint16
}

type AccountAuthorityMap struct {
	*orderedmap.OrderedMap
}

func NewAccountAuthorityMap(items ...AccountAuthority) *AccountAuthorityMap {
	orderedMap := orderedmap.NewOrderedMap()

	for _, v := range items {
		orderedMap.Set(v.AccountName, v.Weight)
	}

	return &AccountAuthorityMap{
		OrderedMap: orderedMap,
	}
}

func (m *AccountAuthorityMap) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)
	enc.Encode(uint8(m.Len()))
	for el := m.Front(); el != nil; el = el.Next() {
		enc.Encode(el.Key)
		enc.Encode(el.Value)
	}
	return enc.Err()
}

func (m AccountAuthorityMap) MarshalJSON() ([]byte, error) {
	xs := make([]interface{}, 0, m.Len())
	for el := m.Front(); el != nil; el = el.Next() {
		xs = append(xs, []interface{}{el.Key, el.Value})
	}

	return json.Marshal(xs)
}

func (m *AccountAuthorityMap) UnmarshalJSON(data []byte) error {
	var xs [][]interface{}
	if err := json.Unmarshal(data, &xs); err != nil {
		return err
	}

	var invalid bool
	orderedMap := orderedmap.NewOrderedMap()
	for _, kv := range xs {
		if len(kv) != 2 {
			invalid = true
			break
		}

		k, ok := kv[0].(string)
		if !ok {
			invalid = true
			break
		}

		var v uint16
		switch t := kv[1].(type) {
		case uint16:
			v = t
		case float64:
			v = uint16(t)
		default:
			invalid = true
			break
		}

		orderedMap.Set(PublicKey(k), v)
	}
	if invalid {
		return errors.New("invalid map encoding")
	}

	m.OrderedMap = orderedMap

	return nil
}
