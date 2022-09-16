package rpc

import (
	"github.com/scorum/scorum-go/rpc/http"
	"github.com/scorum/scorum-go/rpc/websocket"
)

func NewWebSocketTransport(conn websocket.Connection) *websocket.Transport {
	return websocket.NewTransport(conn)
}

func NewHTTPTransport(url string, options ...func(*http.Transport)) *http.Transport {
	return http.NewTransport(url, options...)
}
