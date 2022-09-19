package scorumgo

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/scorum/scorum-go/caller"
	"github.com/scorum/scorum-go/rpc/protocol"
)

var (
	callsProcessed = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "call_pending",
		Help: "The number of pending calls at the moment",
	}, []string{"api", "method"})

	callDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "call_duration_millisecond",
		Help: "",
	}, []string{"status", "api", "method"})
)

type PrometheusInterceptor struct {
	caller caller.CallCloser
}

func NewPrometheusInterceptor(caller caller.CallCloser) *PrometheusInterceptor {
	return &PrometheusInterceptor{caller: caller}
}

func (c *PrometheusInterceptor) Call(ctx context.Context, api string, method string, args []interface{}, reply interface{}) error {
	start := time.Now()

	callsProcessed.WithLabelValues(api, method).Inc()
	defer callsProcessed.WithLabelValues(api, method).Dec()

	err := c.caller.Call(ctx, api, method, args, reply)
	if err != nil {
		var status = "error"
		if errors.Is(err, protocol.ErrWaitResponseTimeout) {
			status = "timeout"
		}

		callDuration.WithLabelValues(status, api, method).
			Observe(float64(time.Since(start).Nanoseconds()) / 1000000)

		return err
	}

	callDuration.WithLabelValues("ok", api, method).
		Observe(float64(time.Since(start).Nanoseconds()) / 1000000)

	return nil
}

func (c *PrometheusInterceptor) SetCallback(api string, method string, callback func(raw json.RawMessage)) error {
	err := c.caller.SetCallback(api, method, callback)
	if err != nil {
		return err
	}

	return nil
}

func (c *PrometheusInterceptor) Close() error {
	return c.caller.Close()
}
