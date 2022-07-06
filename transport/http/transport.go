package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/scorum/scorum-go/transport"
)

type requestIDGenerator interface {
	Generate() uint64
}

type Transport struct {
	Url    string
	client *http.Client

	requestID requestIDGenerator
}

func NewTransport(url string, options ...func(*Transport)) *Transport {
	t := Transport{
		Url:       url,
		requestID: transport.NewSequenceGenerator(0),
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
	request := transport.RPCRequest{
		Method: "call",
		ID:     caller.requestID.Generate(),
		Params: []interface{}{api, method, args},
	}

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

	var rpcResponse transport.RPCResponse
	if err := decodeJSON(resp.Body, &rpcResponse); err != nil {
		return fmt.Errorf("failed decode rpc response: %w", err)
	}

	if rpcResponse.Error != nil {
		return rpcResponse.Error
	}

	if rpcResponse.Result != nil {
		if err := decodeJSON(bytes.NewReader(*rpcResponse.Result), reply); err != nil {
			return fmt.Errorf("failed decode rpc result: %w", err)
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

func decodeJSON(reader io.Reader, out interface{}) error {
	buf := new(bytes.Buffer)
	tr := io.TeeReader(reader, buf)

	data, err := ioutil.ReadAll(tr)
	if err != nil {
		return fmt.Errorf("failed to read body: %w", err)
	}
	if !json.Valid(data) {
		return fmt.Errorf("invalid json: %q", buf.String())
	}

	if err := json.NewDecoder(buf).Decode(out); err != nil {
		return fmt.Errorf("failed to decode json: %w", err)
	}
	return nil
}
