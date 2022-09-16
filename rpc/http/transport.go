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

	"github.com/scorum/scorum-go/rpc/protocol"
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

	// increase request id
	if caller.requestID == math.MaxUint64 {
		caller.requestID = 0
	}
	caller.requestID++

	request := protocol.RPCRequest{
		Method: "call",
		ID:     caller.requestID,
		Params: []interface{}{api, method, args},
	}

	caller.reqMutex.Unlock()

	reqBody, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("json marshall: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", caller.Url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("http new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := caller.client.Do(req)
	if err != nil {
		return fmt.Errorf("http client do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %w", err)
	}

	var rpcResponse protocol.RPCResponse
	if err = json.Unmarshal(respBody, &rpcResponse); err != nil {
		return fmt.Errorf("json unmarshall rpc reponse: %w: %+v", err, string(respBody))
	}

	if rpcResponse.Error != nil {
		return rpcResponse.Error
	}

	if rpcResponse.Result != nil {
		if err := json.Unmarshal(*rpcResponse.Result, reply); err != nil {
			return fmt.Errorf("json unmarshall rpc result: %w: %+v", err, string(*rpcResponse.Result))
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
