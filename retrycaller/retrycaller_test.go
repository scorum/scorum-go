package retrycaller

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/scorum/scorum-go/caller"
)

var (
	ctx    = context.Background()
	api    = "database_api"
	method = "get_account"
	args   = []interface{}{"leo"}
	reply  = ""
)

func TestRetryCaller_Call_Once(t *testing.T) {
	ctrl := gomock.NewController(t)
	transport := caller.NewMockCallCloser(ctrl)
	transport.EXPECT().
		Call(ctx, api, method, args, gomock.Any()).
		Return(errors.New("some error"))

	client := NewRetryCaller(transport)

	require.Error(t, client.Call(ctx, api, method, args, &reply), "some error")
}

func TestRetryCaller_Call_ReturnReply(t *testing.T) {
	ctrl := gomock.NewController(t)
	transport := caller.NewMockCallCloser(ctrl)
	transport.EXPECT().
		Call(ctx, api, method, args, gomock.Any()).
		DoAndReturn(func(ctx context.Context, api string, method string, args []interface{}, reply interface{}) error {
			r, _ := reply.(*string)
			*r = "result"
			return nil
		})

	client := NewRetryCaller(transport)

	require.NoError(t, client.Call(ctx, api, method, args, &reply))
	require.Equal(t, "result", reply)
}

func TestRetryCaller_Call_WithRetry_MaxAttempts(t *testing.T) {
	ctrl := gomock.NewController(t)
	transport := caller.NewMockCallCloser(ctrl)
	transport.EXPECT().
		Call(ctx, api, method, args, gomock.Any()).
		Return(errors.New("some error")).
		Times(5)

	client := NewRetryCaller(transport, WithRetry(api, method, RetryOptions{
		Timeout:    0,
		RetryLimit: 4,
	}))

	require.Error(t, client.Call(ctx, api, method, args, &reply), "some error")
}

func TestRetryCaller_Call_WithRetry_NotMatch_Once(t *testing.T) {
	ctrl := gomock.NewController(t)
	transport := caller.NewMockCallCloser(ctrl)
	transport.EXPECT().
		Call(ctx, api, method, args, gomock.Any()).
		Return(errors.New("some error"))

	client := NewRetryCaller(transport, WithRetry(api, "get_block", RetryOptions{
		Timeout:    0,
		RetryLimit: 5,
	}))

	require.Error(t, client.Call(ctx, api, method, args, &reply), "some error")
}

func TestRetryCaller_Call_WithDefaultRetry_MaxAttempts(t *testing.T) {
	ctrl := gomock.NewController(t)
	transport := caller.NewMockCallCloser(ctrl)
	transport.EXPECT().
		Call(ctx, api, method, args, gomock.Any()).
		Return(errors.New("some error")).
		Times(5)

	client := NewRetryCaller(transport, WithDefaultRetry(RetryOptions{
		Timeout:    0,
		RetryLimit: 4,
	}))

	require.Error(t, client.Call(ctx, api, method, args, &reply), "some error")
}
