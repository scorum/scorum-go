package caller

import (
	"encoding/json"
	"errors"
	"strconv"
)

var ErrShutdown = errors.New("connection is shut down")

type (
	RPCRequest struct {
		Method string      `json:"method"`
		Params interface{} `json:"params,omitempty"`
		ID     uint64      `json:"id"`
	}

	RPCResponse struct {
		Result *json.RawMessage `json:"result,omitempty"`
		Error  *RPCError        `json:"error,omitempty"`
		ID     uint64           `json:"id"`
	}

	RPCError struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			Code    int    `json:"code"`
			Name    string `json:"name"`
			Message string `json:"message"`
			Stack   []struct {
				Context struct {
					Level      string `json:"level"`
					File       string `json:"file"`
					Line       int    `json:"line"`
					Method     string `json:"method"`
					Hostname   string `json:"hostname"`
					ThreadName string `json:"thread_name"`
					Timestamp  string `json:"timestamp"`
				} `json:"context"`
				Format string      `json:"format"`
				Data   interface{} `json:"data"`
			} `json:"stack"`
		} `json:"data"`
	}

	rpcIncoming struct {
		Method string            `json:"method"`
		Params []json.RawMessage `json:"params"`
	}
)

func (e *RPCError) Error() string {
	return strconv.Itoa(e.Code) + ": " + e.Message
}
