package common

import (
	"context"
	"os"
	"testing"
)

type testWebsocketSessionHandler struct {
	t            *testing.T
	messageCount int
	done         chan struct{}
}

func (t *testWebsocketSessionHandler) OnMessage(i interface{}) error {
	t.t.Log(i)
	t.messageCount++
	if t.done != nil && t.messageCount >= 3 {
		close(t.done)
		t.done = nil
	}
	return nil
}

func (t *testWebsocketSessionHandler) OnClose(err error) {
	if err != nil {
		t.t.Logf("handler on close: %v\n", err)
	}
}

func TestWebsocketSession(t *testing.T) {
	if value := os.Getenv("TEST_WS_SESSION"); value == "" {
		t.Skip("skip websocket session tests")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	urlAddress := "wss://fstream.binance.com/ws"

	cli, err := WebsocketDial(ctx, urlAddress, nil)
	if err != nil {
		t.Fatal(err)
		return
	}

	handler := &testWebsocketSessionHandler{t, 0, make(chan struct{})}
	session := NewWebsocketSession(cli, handler)
	go session.Loop()

	reply, err := session.Subscribe(ctx, "btcusdt@markPrice")
	if err != nil {
		t.Fatal(err)
		return
	}
	if err = reply.OK(); err != nil {
		t.Fatal(err)
		return
	}
	<-handler.done
	cancel()
}
