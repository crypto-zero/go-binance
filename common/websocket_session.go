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
	Loop()
	Subscribe(ctx context.Context, streams ...string) (reply *WebsocketReply, err error)
}

type WebsocketSessionHandler interface {
	OnMessage(interface{}) error
	OnClose(err error)
}

func NewWebsocketSession(client WebsocketClient, handler WebsocketSessionHandler) WebsocketSession {
	return &websocketSession{
		client:          client,
		handler:         handler,
		pendingRequests: make(map[uint64]*websocketSessionRequest),
	}
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

	switch result := unpackResult.(type) {
	case map[string]interface{}:
		if MapHasKeys(result, "id", "method", "code") {
			return ws.onRequestReply(data)
		}
	}
	return ws.handler.OnMessage(unpackResult)
}

func (ws *websocketSession) onRequestReply(data []byte) (err error) {
	var reply WebsocketReply
	if err = json.Unmarshal(data, &reply); err != nil {
		return err
	}

	ws.requestLock.Lock()
	request, ok := ws.pendingRequests[reply.ID]
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

func (ws *websocketSession) Loop() {
	err := ws.client.Loop(ws.onMessage)
	if err == nil {
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
