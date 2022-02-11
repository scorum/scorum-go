package retrycaller

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/scorum/scorum-go/caller"
	"github.com/scorum/scorum-go/retry"
)

type RetryOptions struct {
	Timeout    time.Duration
	RetryLimit int
}

type Option func(r *retryCaller)

func WithRetry(api, method string, options RetryOptions) func(r *retryCaller) {
	return func(r *retryCaller) {
		r.retryOptions[fmt.Sprintf("%s.%s", api, method)] = options
	}
}

func WithDefaultRetry(options RetryOptions) func(r *retryCaller) {
	return func(r *retryCaller) {
		r.defaultRetryOptions = options
	}
}

type retryCaller struct {
	transport           caller.Caller
	retryOptions        map[string]RetryOptions
	defaultRetryOptions RetryOptions
}

func NewRetryCaller(transport caller.CallCloser, options ...Option) caller.CallCloser {
	client := retryCaller{
		transport:    transport,
		retryOptions: map[string]RetryOptions{},
		defaultRetryOptions: RetryOptions{
			Timeout:    0,
			RetryLimit: 0,
		},
	}

	for _, option := range options {
		option(&client)
	}

	return &client
}

func (r retryCaller) Call(ctx context.Context, api string, method string, args []interface{}, reply interface{}) error {
	retryOptions := r.getRetryOptions(api, method)

	options := retry.Options{
		Timeout:    retryOptions.Timeout,
		RetryLimit: retryOptions.RetryLimit,
	}

	operation := func() error {
		return r.transport.Call(ctx, api, method, args, reply)
	}

	return retry.Do(operation, options)
}

func (r retryCaller) SetCallback(api string, method string, callback func(raw json.RawMessage)) error {
	return r.transport.SetCallback(api, method, callback)
}

func (r retryCaller) Close() error {
	return nil
}

func (r retryCaller) getRetryOptions(api string, method string) RetryOptions {
	opt, ok := r.retryOptions[fmt.Sprintf("%s.%s", api, method)]
	if !ok {
		return r.defaultRetryOptions
	}
	return opt
}
