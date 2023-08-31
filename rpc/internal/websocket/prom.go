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

	pingPongCounter = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "ping_pong_count",
		Help: "The balance of pings sent (inc) to pongs received (dec)",
	})
)
