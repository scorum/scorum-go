package caller

import (
	"context"
	"encoding/json"
	"io"
)

var EmptyParams = []interface{}{}

type Caller interface {
	Call(ctx context.Context, api string, method string, args []interface{}, reply interface{}) error
	SetCallback(api string, method string, callback func(raw json.RawMessage)) error
}

type CallCloser interface {
	Caller
	io.Closer
}
