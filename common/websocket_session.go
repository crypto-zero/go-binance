package common

import (
	"nhooyr.io/websocket"
)

type websocketSession struct {
	client  WebsocketClient
	handler WebsocketSessionHandler
}

type WebsocketSession interface{}

type WebsocketSessionHandler interface {
	OnMessage(data []byte)
	OnClose(err error)
}

func NewWebsocketSession(client WebsocketClient, handler WebsocketSessionHandler) WebsocketSession {
	return &websocketSession{client: client, handler: handler}
}

func (ws *websocketSession) onMessage(messageType websocket.MessageType, data []byte) (err error) {
	return
}

func (ws *websocketSession) Loop() {
	if err := ws.client.Loop(ws.onMessage); err != nil {
		ws.handler.OnClose(err)
	}
}
