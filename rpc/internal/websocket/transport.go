package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/scorum/scorum-go/rpc/protocol"
)

type Transport struct {
	conn *Connector

	mutex   sync.Mutex
	pending map[uint64]*callRequest

	callbackMutex sync.Mutex
	callbackID    uint64
	callbacks     map[uint64]func(args json.RawMessage)

	waitResponseTimeout time.Duration
}

// Represent an async call
type callRequest struct {
	Error error            // after completion, the error status.
	Done  chan bool        // strobes when call is complete.
	Reply *json.RawMessage // reply message
}

func NewTransport(conn *Connector) *Transport {
	tr := Transport{
		conn:                conn,
		pending:             make(map[uint64]*callRequest),
		callbacks:           make(map[uint64]func(args json.RawMessage)),
		waitResponseTimeout: 10 * time.Second,
	}

	return &tr
}

func (tr *Transport) Dial(ctx context.Context) error {
	return tr.conn.Dial(ctx, tr.OnMessage, tr.OnReconnect)
}

func (tr *Transport) Close() error {
	return tr.conn.Close()
}

func (tr *Transport) Call(ctx context.Context, api string, method string, args []interface{}, reply interface{}) error {
	var (
		requestID = uint64(rand.Uint32())
		call      = callRequest{Done: make(chan bool, 1)}
	)

	tr.mutex.Lock()
	tr.pending[requestID] = &call

	tr.mutex.Unlock()

	r := protocol.RPCRequest{
		Method: "call",
		ID:     requestID,
		Params: []interface{}{api, method, args},
	}

	if err := tr.conn.WriteJSON(&r); err != nil {
		tr.finishPending(requestID)

		return fmt.Errorf("send: %w", err)
	}

	select {
	case <-time.After(tr.waitResponseTimeout):
		tr.finishPending(requestID)
		return protocol.ErrWaitResponseTimeout

	case <-ctx.Done():
		tr.finishPending(requestID)

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

func (tr *Transport) OnReconnect() {
	tr.stopAllPending(protocol.ErrShutdown)
}

func (tr *Transport) OnMessage(message []byte) {
	log := logrus.WithField("message", string(message))

	var response protocol.RPCResponse
	if err := json.Unmarshal(message, &response); err != nil {
		log.WithError(err).Error("json unmarshall rpc response")
		return
	}

	tr.mutex.Lock()
	call, ok := tr.pending[response.ID]
	tr.mutex.Unlock()

	if ok {
		tr.onCallResponse(response, call)
		return
	}

	// the message is not a pending call, but probably a callback notice
	var incoming protocol.RPCIncoming
	if err := json.Unmarshal(message, &incoming); err != nil {
		log.WithError(err).Error("json unmarshall rpc incoming")
		return
	}

	if incoming.Method != "notice" {
		log.Debugf("protocol error: unknown message received: %+v\n", incoming)
		return
	}

	if err := tr.onNotice(incoming); err != nil {
		log.WithError(err).Error("on notice")
		return
	}
}

func (tr *Transport) ConnectionAliveAt() time.Time {
	return tr.conn.GetAliveAt()
}

// readPump pumps messages from the websocket Connection and dispatches them.
// func (tr *Transport) readPump() {
// 	for {
// 		_, message, err := tr.conn.ReadMessage()
// 		if err != nil {
// 			tr.stopAllPending(fmt.Errorf("conn read message: %w", err))
// 			return
// 		}
//
// 		var response protocol.RPCResponse
// 		if err := json.Unmarshal(message, &response); err != nil {
// 			logrus.WithError(err).WithFields(logrus.Fields{
// 				"message": string(message),
// 			}).Error("json unmarshall rpc response")
//
// 			continue
// 			// tr.stopAllPending(fmt.Errorf("json unmarshal rpc response: %w", err))
// 			// return
// 		}
//
// 		tr.mutex.Lock()
// 		call, ok := tr.pending[response.ID]
// 		tr.mutex.Unlock()
//
// 		if ok {
// 			tr.onCallResponse(response, call)
// 			continue
// 		}
//
// 		// the message is not a pending call, but probably a callback notice
// 		var incoming protocol.RPCIncoming
// 		if err := json.Unmarshal(message, &incoming); err != nil {
// 			tr.stopAllPending(fmt.Errorf("json unmarshall rpc incoming: %w", err))
// 			return
// 		}
//
// 		if incoming.Method != "notice" {
// 			log.Printf("protocol error: unknown message received: %+v\n", incoming)
// 			continue
// 		}
//
// 		if err := tr.onNotice(incoming); err != nil {
// 			tr.stopAllPending(fmt.Errorf("on notice: %w", err))
// 			return
// 		}
// 	}
// }

func (tr *Transport) finishPending(requestID uint64) {
	tr.mutex.Lock()
	defer tr.mutex.Unlock()

	delete(tr.pending, requestID)
}

// Return pending clients and shutdown the client
func (tr *Transport) stopAllPending(err error) {
	tr.mutex.Lock()
	defer tr.mutex.Unlock()

	for _, call := range tr.pending {
		call.Error = err
		call.Done <- true
	}
}

// Call response handler
func (tr *Transport) onCallResponse(response protocol.RPCResponse, call *callRequest) {
	tr.finishPending(response.ID)

	if response.Error != nil {
		call.Error = response.Error
	}

	call.Reply = response.Result
	call.Done <- true
}

// Incoming notice handler
func (tr *Transport) onNotice(incoming protocol.RPCIncoming) error {
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
