package common

import (
	"context"
	"fmt"
	"os"
	"testing"
)

type testWebsocketSessionHandler struct {
	t            *testing.T
	messageCount int
	done         chan struct{}
}

func (t *testWebsocketSessionHandler) OnUnknownMessage(data []byte, m interface{}) error {
	t.t.Log(m)
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

	cli, err := DefaultWebsocketProvider(ctx, urlAddress, nil)
	if err != nil {
		t.Fatal(err)
		return
	}

	handler := &testWebsocketSessionHandler{t, 0, make(chan struct{})}
	session := NewWebsocketSession(cli, handler)
	go session.Loop()

	symbols := []string{
		"bnbusdt", "btcusdt", "etcusdt", "lunausdt", "ethusdt", "linkusdt",
		"eosusdt", "ltcusdt", "bchusdt", "dashusdt", "ontusdt", "neousdt",
	}

	var streams []string
	for _, s := range symbols {
		streams = append(streams, fmt.Sprintf("%s@markPrice", s))
	}

	reply, err := session.Subscribe(ctx, streams...)
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
