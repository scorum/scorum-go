package account_history

import (
	"encoding/json"
	"errors"

	"github.com/scorum/scorum-go/types"
)

type AccountHistory map[uint32]*types.OperationObject

func (ah *AccountHistory) UnmarshalJSON(b []byte) (err error) {
	// unmarshal array
	var o []json.RawMessage
	if err := json.Unmarshal(b, &o); err != nil {
		return err
	}

	// foreach operation
	for _, op := range o {
		var kv []json.RawMessage
		if err := json.Unmarshal(op, &kv); err != nil {
			return err
		}

		if len(kv) != 2 {
			return errors.New("invalid operation format: should be sequence number, value")
		}

		var key uint32
		if err := json.Unmarshal(kv[0], &key); err != nil {
			return err
		}

		var ops types.OperationObject
		if err := json.Unmarshal(kv[1], &ops); err != nil {
			return err
		}

		(*ah)[key] = &ops
	}

	return nil
}
