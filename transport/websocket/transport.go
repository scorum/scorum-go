package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/scorum/scorum-go/transport"
)

var (
	ErrWaitResponseTimeout = errors.New("wait response timeout")
)

type connection interface {
	WriteJSON(v interface{}) error
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, data []byte) error
	Close() error
}

type Transport struct {
	conn connection

	reqMutex sync.Mutex
	pending  map[uint64]*callRequest

	callbackMutex sync.Mutex
	callbackID    uint64
	callbacks     map[uint64]func(args json.RawMessage)

	closing  bool // user has called Close
	shutdown bool // server has told us to stop

	waitResponseTimeout time.Duration

	mutex sync.Mutex
}

// Represent an async call
type callRequest struct {
	Error error            // after completion, the error status.
	Done  chan bool        // strobes when call is complete.
	Reply *json.RawMessage // reply message
}

func NewTransport(conn connection) *Transport {
	tr := Transport{
		conn:                conn,
		pending:             make(map[uint64]*callRequest),
		callbacks:           make(map[uint64]func(args json.RawMessage)),
		waitResponseTimeout: 10 * time.Second,
	}

	go tr.readPump()

	return &tr
}

func (tr *Transport) Call(ctx context.Context, api string, method string, args []interface{}, reply interface{}) error {
	var (
		requestID = uint64(rand.Uint32())
		call      = callRequest{Done: make(chan bool, 1)}
	)

	tr.mutex.Lock()
	if tr.closing || tr.shutdown {
		tr.mutex.Unlock()
		return transport.ErrShutdown
	}
	tr.pending[requestID] = &call
	tr.mutex.Unlock()

	send := func(v interface{}) error {
		tr.reqMutex.Lock()
		defer tr.reqMutex.Unlock()

		return tr.conn.WriteJSON(v)
	}

	r := transport.RPCRequest{
		Method: "call",
		ID:     requestID,
		Params: []interface{}{api, method, args},
	}

	if err := send(&r); err != nil {
		tr.mutex.Lock()
		delete(tr.pending, requestID)
		tr.mutex.Unlock()

		return fmt.Errorf("send: %w", err)
	}

	select {
	case <-time.After(tr.waitResponseTimeout):
		return ErrWaitResponseTimeout
	case <-ctx.Done():
		return ctx.Err()
	case <-call.Done:
		// wait for the call to complete
	}

	if call.Error != nil {
		return call.Error
	}

	if call.Reply != nil {
		if err := json.Unmarshal(*call.Reply, reply); err != nil {
			return fmt.Errorf("json unmarshall: %w", err)
		}
	}
	return nil
}

// readPump pumps messages from the websocket connection and dispatches them.
func (tr *Transport) readPump() {
	for {
		_, message, err := tr.conn.ReadMessage()
		if err != nil {
			tr.stop(err)
			return
		}

		var response transport.RPCResponse
		if err := json.Unmarshal(message, &response); err != nil {
			tr.stop(fmt.Errorf("json unmarshal: %w", err))
			return
		}

		tr.mutex.Lock()
		call, ok := tr.pending[response.ID]
		tr.mutex.Unlock()

		if ok {
			tr.onCallResponse(response, call)
		} else {
			// the message is not a pending call, but probably a callback notice
			var incoming transport.RPCIncoming
			if err := json.Unmarshal(message, &incoming); err != nil {
				tr.stop(fmt.Errorf("json unmarshall: %w", err))
				return
			}
			if incoming.Method == "notice" {
				if err := tr.onNotice(incoming); err != nil {
					tr.stop(err)
					return
				}
			} else {
				log.Printf("protocol error: unknown message received: %+v\n", incoming)
			}
		}
	}
}

// Return pending clients and shutdown the client
func (tr *Transport) stop(err error) {
	tr.mutex.Lock()
	defer tr.mutex.Unlock()

	tr.shutdown = true
	for _, call := range tr.pending {
		call.Error = err
		call.Done <- true
	}
}

// Call response handler
func (tr *Transport) onCallResponse(response transport.RPCResponse, call *callRequest) {
	tr.mutex.Lock()
	defer tr.mutex.Unlock()

	delete(tr.pending, response.ID)
	if response.Error != nil {
		call.Error = response.Error
	}
	call.Reply = response.Result
	call.Done <- true
}

// Incoming notice handler
func (tr *Transport) onNotice(incoming transport.RPCIncoming) error {
	length := len(incoming.Params)

	if length == 0 {
		return nil
	}

	if length == 1 {
		return fmt.Errorf("invalid notice(%+v) message with odd number of params", incoming)
	}

	for i := 0; i < length; i += 2 {
		callbackID, err := strconv.ParseUint(string(incoming.Params[i]), 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse callbackID: %w", err)
		}

		notice := tr.callbacks[callbackID]
		if notice == nil {
			return fmt.Errorf("callback %d is not registered", callbackID)
		}

		// invoke callback
		notice(incoming.Params[i+1])
	}

	return nil
}

func (tr *Transport) SetCallback(api string, method string, notice func(args json.RawMessage)) error {
	// increase callback id
	tr.callbackMutex.Lock()
	if tr.callbackID == math.MaxUint64 {
		tr.callbackID = 0
	}
	tr.callbackID++
	tr.callbacks[tr.callbackID] = notice
	tr.callbackMutex.Unlock()

	return tr.Call(context.Background(), api, method, []interface{}{tr.callbackID}, nil)
}

// Close calls the underlying web socket Close method. If the connection is already
// shutting down, ErrShutdown is returned.
func (tr *Transport) Close() error {
	tr.mutex.Lock()
	if tr.closing {
		tr.mutex.Unlock()
		return transport.ErrShutdown
	}
	tr.closing = true
	tr.mutex.Unlock()

	tr.reqMutex.Lock()
	defer tr.reqMutex.Unlock()

	err := tr.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("write close:", err)
	}

	if err := tr.conn.Close(); err != nil {
		return fmt.Errorf("conn close: %w", err)
	}

	return nil
}
