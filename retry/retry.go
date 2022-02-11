package retry

import (
	"time"
)

type Options struct {
	Timeout    time.Duration
	RetryLimit int
}

type Func func() error

func Do(f Func, opt Options) error {
	var err error
	var attempt int

	for {
		err = f()
		if err == nil {
			return nil
		}

		attempt++
		if attempt > opt.RetryLimit {
			return err
		}

		time.Sleep(opt.Timeout)
	}
}
