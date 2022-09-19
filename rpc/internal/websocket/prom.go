package websocket

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	isShutdown = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "ws_is_shutdown",
		Help: "The status of websocket connection",
	})
)
