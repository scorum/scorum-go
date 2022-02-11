package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/scorum/scorum-go/transport"
)

type Transport struct {
	Url    string
	client *http.Client

	requestID uint64
	reqMutex  sync.Mutex
}

func NewTransport(url string, options ...func(*Transport)) *Transport {
	t := Transport{
		Url: url,
	}

	for _, o := range options {
		o(&t)
	}

	if t.client == nil {
		t.client = &http.Client{
			Timeout: 5 * time.Second,
		}
	}

	return &t
}

func WithHttpClient(client *http.Client) func(*Transport) {
	return func(t *Transport) {
		t.client = client
	}
}

func (caller *Transport) Call(ctx context.Context, api string, method string, args []interface{}, reply interface{}) error {
	caller.reqMutex.Lock()
	defer caller.reqMutex.Unlock()

	// increase request id
	if caller.requestID == math.MaxUint64 {
		caller.requestID = 0
	}
	caller.requestID++

	request := transport.RPCRequest{
		Method: "call",
		ID:     caller.requestID,
		Params: []interface{}{api, method, args},
	}

	reqBody, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", caller.Url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := caller.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %w", err)
	}

	var rpcResponse transport.RPCResponse
	if err = json.Unmarshal(respBody, &rpcResponse); err != nil {
		return fmt.Errorf("failed to unmarshal response: %+v: %w", string(respBody), err)
	}

	if rpcResponse.Error != nil {
		return rpcResponse.Error
	}

	if rpcResponse.Result != nil {
		if err := json.Unmarshal(*rpcResponse.Result, reply); err != nil {
			return fmt.Errorf("failed to unmarshal rpc result: %+v: %w", string(*rpcResponse.Result), err)
		}
	}

	return nil
}

func (caller *Transport) SetCallback(api string, method string, notice func(args json.RawMessage)) error {
	panic("not supported")
}

func (caller *Transport) Close() error {
	return nil
}
