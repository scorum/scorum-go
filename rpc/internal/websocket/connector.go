package websocket

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	_ "github.com/prometheus/client_golang/prometheus"
	_ "github.com/prometheus/client_golang/prometheus/promauto"
	_ "github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	"github.com/scorum/scorum-go/rpc/protocol"
)

const (
	reconnectDelay = 2 * time.Second
	writeDeadline  = 10 * time.Second

	pingDuration = 5 * time.Second
	pingTimeout  = 10 * time.Second
)

type Connector struct {
	URL       string
	dialer    *websocket.Dialer
	conn      *websocket.Conn
	mutex     sync.RWMutex
	connMutex sync.Mutex

	isShutdown bool
	isClosing  bool

	messageHandler func(message []byte)
	connectHandler func()

	aliveMutex sync.Mutex
	aliveAt    time.Time
}

func NewConnector(url string, dialer *websocket.Dialer) *Connector {
	return &Connector{
		URL:        url,
		dialer:     dialer,
		isClosing:  true,
		isShutdown: true,
	}
}

func (r *Connector) Dial(ctx context.Context, messageHandler func(message []byte), connectHandler func()) error {
	if err := r.dial(ctx); err != nil {
		return err
	}

	r.messageHandler = messageHandler
	r.connectHandler = connectHandler

	go r.loop(ctx)

	return nil
}

func (r *Connector) dial(ctx context.Context) error {
	r.connMutex.Lock()
	defer r.connMutex.Unlock()

	r.mutex.Lock()
	defer r.mutex.Unlock()

	// already connected
	if !r.isShutdown && !r.isClosing {
		return protocol.ErrShutdown
	}

	conn, _, err := r.dialer.DialContext(ctx, r.URL, nil)
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}

	r.updateAlive()
	conn.SetPongHandler(func(_ string) error {
		pingPongCounter.Dec()
		r.updateAlive()
		return nil
	})

	r.isClosing = false
	r.isShutdown = false
	r.conn = conn

	if r.connectHandler != nil {
		r.connectHandler()
	}

	isShutdown.Inc()

	return nil
}

func (r *Connector) loop(ctx context.Context) {
	pingTicker := time.NewTicker(pingDuration)

	for {
		select {
		case <-ctx.Done():
			return
		case <-pingTicker.C:
			go func() {
				r.connMutex.Lock()
				defer r.connMutex.Unlock()

				err := r.conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(pingTimeout))
				if err != nil {
					pingPongCounter.Inc()
				}
			}()
		default:
			r.mutex.RLock()
			if r.isClosing {
				r.mutex.RUnlock()
				return
			}

			if r.isShutdown {
				r.mutex.RUnlock()
				if err := r.dial(ctx); err != nil {
					logrus.WithError(err).Error("reconnect dial")
					time.Sleep(reconnectDelay)
				}
				continue
			}

			r.mutex.RUnlock()

			_, message, err := r.conn.ReadMessage()
			if err != nil {
				r.shutdown()
				logrus.WithError(err).Error("read message")
				continue
			}

			r.updateAlive()
			if r.messageHandler != nil {
				r.messageHandler(message)
			}
		}
	}
}

func (r *Connector) shutdown() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	isShutdown.Dec()

	r.isShutdown = true
}

func (r *Connector) WriteJSON(v interface{}) error {
	r.mutex.RLock()
	if r.isShutdown || r.isClosing {
		r.mutex.RUnlock()
		return protocol.ErrShutdown
	}
	r.mutex.RUnlock()

	r.connMutex.Lock()
	defer r.connMutex.Unlock()

	_ = r.conn.SetWriteDeadline(time.Now().UTC().Add(writeDeadline))
	if err := r.conn.WriteJSON(v); err != nil {
		r.shutdown()
		logrus.WithError(err).Error("write json")
		return fmt.Errorf("conn write json: %w", err)
	}
	return nil
}

// func (r *Connector) ReadMessage() (messageType int, p []byte, err error) {
// 	r.mutex.RLock()
// 	if r.isShutdown || r.isClosing {
// 		r.mutex.RUnlock()
// 		return 0, nil, protocol.ErrShutdown
// 	}
// 	r.mutex.RUnlock()
//
// 	messageType, p, err = r.conn.ReadMessage()
// 	if err != nil {
// 		r.shutdown()
// 		logrus.WithError(err).Error("read message")
// 	}
// 	return
// }

// Close calls the underlying web socket Close method. If the Connection is already
// shutting down, ErrShutdown is returned.
func (r *Connector) Close() error {
	r.mutex.Lock()
	if r.isShutdown || r.isClosing {
		r.mutex.Unlock()
		return protocol.ErrShutdown
	}

	r.isClosing = true
	r.mutex.Unlock()

	r.connMutex.Lock()
	defer r.connMutex.Unlock()

	msg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
	err := r.conn.WriteControl(websocket.CloseMessage, msg, time.Now().Add(time.Second))
	if err != nil {
		logrus.WithError(err).Error("conn write control close")
	}

	if err := r.conn.Close(); err != nil {
		return fmt.Errorf("conn close: %w", err)
	}

	return nil
}

func (r *Connector) updateAlive() {
	r.aliveMutex.Lock()
	defer r.aliveMutex.Unlock()

	r.aliveAt = time.Now()
}

func (r *Connector) GetAliveAt() time.Time {
	r.aliveMutex.Lock()
	defer r.aliveMutex.Unlock()

	return r.aliveAt
}
