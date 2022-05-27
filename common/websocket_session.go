package common

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"

	"nhooyr.io/websocket"
)

type websocketSession struct {
	client          WebsocketClient
	handler         WebsocketSessionHandler
	globalRequestID uint64
	pendingRequests map[uint64]*websocketSessionRequest
	requestLock     sync.Mutex
	messagePatterns []*websocketSessionMessagePattern
}

type websocketSessionRequest struct {
	WebsocketRequest
	reply WebsocketReply
	done  chan error
}

func newWebsocketSessionRequest() *websocketSessionRequest {
	return &websocketSessionRequest{
		done: make(chan error),
	}
}

type WebsocketSession interface {
	// Loop hold this websocket connection until the connection disconnected or got error
	// from message processor. return error same as SessionHandler.OnClose function.
	Loop() (err error)
	// RunLoop create new go routine and call IOLoop function.
	RunLoop() chan error
	Subscribe(ctx context.Context, streams ...string) (reply *WebsocketReply, err error)
	SubscribeNoReply(ctx context.Context, streams ...string) (err error)
	RegisterMessageHandler(factory WebsocketSessionMessageFactory, callback WebsocketSessionMessageCallback,
		checker ...WebsocketSessionMessageChecker)
	RequireMapHasAllKeys(keys ...string) WebsocketSessionMessageChecker
	RequireMapKeyValue(key, value string) WebsocketSessionMessageChecker
}

type WebsocketSessionHandler interface {
	OnUnknownMessage([]byte, interface{}) error
	OnClose(err error)
}

type (
	WebsocketSessionMessageChecker  func(m interface{}) bool
	WebsocketSessionMessageFactory  func() interface{}
	WebsocketSessionMessageCallback func(m interface{})
)

func WebsocketSessionMessageFactoryBuild[T any]() WebsocketSessionMessageFactory {
	return func() interface{} {
		return new(T)
	}
}

func WebsocketSessionMessageHandlerBuild[T any](f func(*T)) WebsocketSessionMessageCallback {
	return func(m interface{}) {
		x := m.(*T)
		f(x)
	}
}

type websocketSessionMessagePattern struct {
	Check    []WebsocketSessionMessageChecker
	New      WebsocketSessionMessageFactory
	Callback WebsocketSessionMessageCallback
}

func NewWebsocketSession(client WebsocketClient, handler WebsocketSessionHandler) WebsocketSession {
	return &websocketSession{
		client:          client,
		handler:         handler,
		pendingRequests: make(map[uint64]*websocketSessionRequest),
	}
}

func (ws *websocketSession) RequireMapHasAllKeys(keys ...string) WebsocketSessionMessageChecker {
	return func(m interface{}) bool {
		switch result := m.(type) {
		case map[string]interface{}:
			return MapHasAllKeys(result, keys...)
		}
		return false
	}
}

func (ws *websocketSession) RequireMapKeyValue(key, value string) WebsocketSessionMessageChecker {
	return func(m interface{}) bool {
		switch result := m.(type) {
		case map[string]interface{}:
			x, ok := result[key]
			if !ok {
				return false
			}
			stringValue, ok := x.(string)
			if !ok {
				return false
			}
			return stringValue == value
		}
		return false
	}
}

func (ws *websocketSession) RegisterMessageHandler(factory WebsocketSessionMessageFactory,
	callback WebsocketSessionMessageCallback, checker ...WebsocketSessionMessageChecker,
) {
	ws.messagePatterns = append(ws.messagePatterns, &websocketSessionMessagePattern{
		Check:    checker,
		New:      factory,
		Callback: callback,
	})
}

func (ws *websocketSession) onMessage(messageType websocket.MessageType, data []byte) (err error) {
	defer func() {
		if err == nil {
			return
		}
		err = fmt.Errorf("wssession handle message failed: %w", err)
	}()
	var unpackResult interface{}
	if err = json.Unmarshal(data, &unpackResult); err != nil {
		return err
	}
	return ws.processMessage(unpackResult, data)
}

func (ws *websocketSession) processMessage(m interface{}, data []byte) (err error) {
	switch result := m.(type) {
	case map[string]interface{}:
		if MapHasKeys(result, "id", "method", "code") {
			return ws.onRequestReply(data)
		}

	CHECK_LOOP:
		for _, p := range ws.messagePatterns {
			for _, c := range p.Check {
				if !c(result) {
					continue CHECK_LOOP
				}
			}
			x := p.New()
			if err = json.Unmarshal(data, x); err != nil {
				return err
			}
			p.Callback(x)
			return nil
		}
	case []interface{}:
		var list []json.RawMessage
		if err = json.Unmarshal(data, &list); err != nil {
			return err
		}
		for idx, x := range result {
			if err = ws.processMessage(x, list[idx]); err != nil {
				return err
			}
		}
		return nil
	}
	return ws.handler.OnUnknownMessage(data, m)
}

func (ws *websocketSession) onRequestReply(data []byte) (err error) {
	var reply WebsocketReply
	if err = json.Unmarshal(data, &reply); err != nil {
		return err
	}

	ws.requestLock.Lock()
	request, ok := ws.pendingRequests[reply.ID]
	if ok {
		delete(ws.pendingRequests, reply.ID)
	}
	ws.requestLock.Unlock()

	if !ok {
		return nil
	}

	request.reply = reply
	close(request.done)
	return
}

func (ws *websocketSession) request(request *websocketSessionRequest) (err error) {
	request.ID = atomic.AddUint64(&ws.globalRequestID, 1)

	d, err := json.Marshal(request.WebsocketRequest)
	if err != nil {
		return fmt.Errorf("wssession build request failed: %w", err)
	}

	ws.requestLock.Lock()
	ws.pendingRequests[request.ID] = request
	ws.requestLock.Unlock()

	ws.client.Write(d)
	return
}

func (ws *websocketSession) Loop() (err error) {
	if err = ws.client.Loop(ws.onMessage); err == nil {
		return
	}
	ws.handler.OnClose(err)

	var pendingRequest []*websocketSessionRequest
	ws.requestLock.Lock()
	for _, req := range ws.pendingRequests {
		pendingRequest = append(pendingRequest, req)
	}
	ws.pendingRequests = make(map[uint64]*websocketSessionRequest)
	ws.requestLock.Unlock()

	for _, req := range pendingRequest {
		select {
		case req.done <- err:
		default:
		}
	}
	return
}

func (ws *websocketSession) RunLoop() chan error {
	c := make(chan error, 1)
	go func() { c <- ws.Loop() }()
	return c
}

func (ws *websocketSession) Subscribe(ctx context.Context, streams ...string) (
	reply *WebsocketReply, err error,
) {
	request := newWebsocketSessionRequest()
	request.Method = "SUBSCRIBE"
	request.Params = streams

	if err = ws.request(request); err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err = <-request.done:
		return &request.reply, err
	}
}

func (ws *websocketSession) SubscribeNoReply(ctx context.Context, streams ...string) (err error) {
	reply, err := ws.Subscribe(ctx, streams...)
	if err != nil {
		return err
	}
	if err = reply.OK(); err != nil {
		return err
	}
	return nil
}
