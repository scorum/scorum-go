package retry

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_DoNotRetry_When_AttemptSucceed(t *testing.T) {
	var counter int

	operation := func() error {
		counter++
		return nil
	}

	opts := Options{
		Timeout:    0,
		RetryLimit: 5,
	}

	require.NoError(t, Do(operation, opts))
	require.Equal(t, 1, counter)
}

func Test_DoNotRetry_When_SecondAttemptSucceed(t *testing.T) {
	var counter int

	operation := func() error {
		counter++
		if counter == 2 {
			return nil
		}
		return errors.New("some error")
	}

	opts := Options{
		Timeout:    0,
		RetryLimit: 10,
	}

	require.NoError(t, Do(operation, opts))
	require.Equal(t, 2, counter)
}

func Test_RetryMaxAttempts_WhenAttemptReturnError(t *testing.T) {
	var counter int

	operation := func() error {
		counter++
		return errors.New("some error")
	}

	opts := Options{
		Timeout:    0,
		RetryLimit: 4,
	}

	require.Error(t, Do(operation, opts), "some error")
	require.Equal(t, opts.RetryLimit+1, counter)
}

func Test_RetryWithTimeout(t *testing.T) {
	var counter int

	operation := func() error {
		counter++
		return errors.New("some error")
	}

	opts := Options{
		Timeout:    time.Millisecond * 100,
		RetryLimit: 4,
	}

	expected := time.Now().Add(opts.Timeout * time.Duration(opts.RetryLimit))

	require.Error(t, Do(operation, opts), "some error")
	require.Equal(t, opts.RetryLimit+1, counter)

	require.WithinDuration(t, expected, time.Now(), 2*time.Millisecond)
}
