package blockchain_history

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/scorum/scorum-go/types"
)

// History, the key is operation sequence number
type History map[uint32]*types.OperationObject

func (h *History) UnmarshalJSON(b []byte) (err error) {
	var o []json.RawMessage
	if err := json.Unmarshal(b, &o); err != nil {
		return err
	}

	ops := make(map[uint32]*types.OperationObject, len(o))

	// foreach operation objects
	for _, op := range o {
		var kv []json.RawMessage
		if err := json.Unmarshal(op, &kv); err != nil {
			return err
		}

		if len(kv) != 2 {
			return errors.New("invalid history encoding: should be id, operation object")
		}

		key, err := strconv.ParseUint(string(kv[0]), 10, 64)
		if err != nil {
			return err
		}

		var val *types.OperationObject
		if err := json.Unmarshal(kv[1], &val); err != nil {
			println(string(kv[1]))
			return err
		}

		ops[uint32(key)] = val
	}

	*h = ops
	return nil
}
