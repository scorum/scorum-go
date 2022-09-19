package rpc

import (
	gohttp "net/http"

	gorilla "github.com/gorilla/websocket"

	"github.com/scorum/scorum-go/rpc/internal/http"
	"github.com/scorum/scorum-go/rpc/internal/websocket"
)

func NewWebSocketTransport(url string, dialer *gorilla.Dialer) *websocket.Transport {
	return websocket.NewTransport(websocket.NewConnector(url, dialer))
}

func NewHTTPTransport(url string, options ...func(*http.Transport)) *http.Transport {
	return http.NewTransport(url, options...)
}

func WithHttpClient(client *gohttp.Client) func(*http.Transport) {
	return http.WithHttpClient(client)
}
