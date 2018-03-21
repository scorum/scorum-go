package caller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"sync"
)

type HttpCaller struct {
	Url string

	requestID uint64
	reqMutex  sync.Mutex
}

func NewHttpCaller(url string) *HttpCaller {
	return &HttpCaller{
		Url: url,
	}
}

func (caller *HttpCaller) Call(api string, method string, args []interface{}, reply interface{}) error {
	caller.reqMutex.Lock()
	defer caller.reqMutex.Unlock()

	// increase request id
	if caller.requestID == math.MaxUint64 {
		caller.requestID = 0
	}
	caller.requestID++

	request := RPCRequest{
		Method: "call",
		ID:     caller.requestID,
		Params: []interface{}{api, method, args},
	}

	reqBody, err := json.Marshal(request)
	if err != nil {
		return err
	}

	resp, err := http.Post(caller.Url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Println(string(respBody))

	var rpcResponse RPCResponse
	if err = json.Unmarshal(respBody, &rpcResponse); err != nil {
		return err
	}

	if rpcResponse.Error != nil {
		return rpcResponse.Error
	}

	if rpcResponse.Result != nil {
		if err := json.Unmarshal(*rpcResponse.Result, reply); err != nil {
			return err
		}
	}

	return nil
}

func (caller *HttpCaller) SetCallback(api string, method string, notice func(args json.RawMessage)) error {
	panic("not supported")
}

func (caller *HttpCaller) Close() error {
	return nil
}
