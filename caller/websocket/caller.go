package websocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"sync"

	"golang.org/x/net/websocket"
)

var ErrShutdown = errors.New("connection is shut down")

type Caller struct {
	conn *websocket.Conn

	reqMutex  sync.Mutex
	requestID uint64
	pending   map[uint64]*callRequest

	callbackMutex sync.Mutex
	callbackID    uint64
	callbacks     map[uint64]func(args json.RawMessage)

	closing  bool // user has called Close
	shutdown bool // server has told us to stop

	mutex sync.Mutex
}

// Represent an async call
type callRequest struct {
	Error error            // after completion, the error status.
	Done  chan bool        // strobes when call is complete.
	Reply *json.RawMessage // reply message
}

func NewCaller(url string) (*Caller, error) {
	ws, err := websocket.Dial(url, "", "http://localhost")
	if err != nil {
		return nil, err
	}

	client := &Caller{
		conn:      ws,
		pending:   make(map[uint64]*callRequest),
		callbacks: make(map[uint64]func(args json.RawMessage)),
	}

	go client.input()
	return client, nil
}

func (client *Caller) Call(api string, method string, args []interface{}, reply interface{}) error {
	client.reqMutex.Lock()
	defer client.reqMutex.Unlock()

	client.mutex.Lock()
	if client.closing || client.shutdown {
		client.mutex.Unlock()
		return ErrShutdown
	}

	// increase request id
	if client.requestID == math.MaxUint64 {
		client.requestID = 0
	}
	client.requestID++
	seq := client.requestID

	c := &callRequest{
		Done: make(chan bool, 1),
	}
	client.pending[seq] = c
	client.mutex.Unlock()

	// send Json Rcp request
	err := websocket.JSON.Send(client.conn, RPCRequest{
		Method: "call",
		ID:     client.requestID,
		Params: []interface{}{api, method, args},
	})
	if err != nil {
		client.mutex.Lock()
		delete(client.pending, seq)
		client.mutex.Unlock()
		return err
	}

	// wait for the call to complete
	<-c.Done
	if c.Error != nil {
		return c.Error
	}

	if c.Reply != nil {
		if err := json.Unmarshal(*c.Reply, reply); err != nil {
			return err
		}
	}
	return nil
}

func (client *Caller) input() {
	for {
		var message string
		if err := websocket.Message.Receive(client.conn, &message); err != nil {
			client.stop(err)
			return
		}

		var response RPCResponse
		if err := json.Unmarshal([]byte(message), &response); err != nil {
			client.stop(err)
			return
		} else {
			if call, ok := client.pending[response.ID]; ok {
				client.onCallResponse(response, call)
			} else {
				//the message is not a pending call, but probably a callback notice
				var incoming rpcIncoming
				if err := json.Unmarshal([]byte(message), &incoming); err != nil {
					client.stop(err)
					return
				}
				if incoming.Method == "notice" {
					if err := client.onNotice(incoming); err != nil {
						client.stop(err)
						return
					}
				} else {
					log.Printf("protocol error: unknown message received: %+v\n", incoming)
				}
			}
		}
	}
}

// Return pending clients and shutdown the client
func (client *Caller) stop(err error) {
	client.reqMutex.Lock()
	client.shutdown = true
	for _, call := range client.pending {
		call.Error = err
		call.Done <- true
	}
	client.reqMutex.Unlock()
}

// Call response handler
func (client *Caller) onCallResponse(response RPCResponse, call *callRequest) {
	client.mutex.Lock()
	delete(client.pending, response.ID)
	if response.Error != nil {
		call.Error = response.Error
	}
	call.Reply = response.Result
	call.Done <- true
	client.mutex.Unlock()
}

// Incoming notice handler
func (client *Caller) onNotice(incoming rpcIncoming) error {
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
			return fmt.Errorf("failed to parse %s as callbackID in notice(%+v): %s", incoming.Params[i], incoming, err)
		}

		notice := client.callbacks[callbackID]
		if notice == nil {
			return fmt.Errorf("callback %d is not registered", callbackID)
		}

		// invoke callback
		notice(incoming.Params[i+1])
	}

	return nil
}

func (client *Caller) SetCallback(api string, method string, notice func(args json.RawMessage)) error {
	// increase callback id
	client.callbackMutex.Lock()
	if client.callbackID == math.MaxUint64 {
		client.callbackID = 0
	}
	client.callbackID++
	client.callbacks[client.callbackID] = notice
	client.callbackMutex.Unlock()

	return client.Call(api, method, []interface{}{client.callbackID}, nil)
}

// Close calls the underlying web socket Close method. If the connection is already
// shutting down, ErrShutdown is returned.
func (client *Caller) Close() error {
	client.mutex.Lock()
	if client.closing {
		client.mutex.Unlock()
		return ErrShutdown
	}
	client.closing = true
	client.mutex.Unlock()
	return client.conn.Close()
}
