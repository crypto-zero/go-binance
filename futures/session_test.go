package futures

import (
	"context"
	"os"
	"testing"
)

type testSessionHandler struct {
	*testing.T
	done  chan struct{}
	count int
}

func (t *testSessionHandler) OnUnknownMessage(bytes []byte, i interface{}) error {
	t.Logf("got unknown message: %v\n", i)
	return nil
}

func (t *testSessionHandler) OnClose(err error) {
	t.Logf("got on close err: %v\n", err)
}

func (t *testSessionHandler) OnAggTradeEvent(event *WsAggTradeEvent) {
	t.Logf("got agg trade event: %#v\n", event)
	t.count++
	if t.done != nil && t.count > 5 {
		x := t.done
		t.done = nil
		close(x)
	}
}

func TestSession(t *testing.T) {
	if value := os.Getenv("TEST_FUTURES_WS_SESSION"); value == "" {
		t.Skip("skip futures websocket session tests")
		return
	}

	handler := &testSessionHandler{T: t, done: make(chan struct{})}

	ctx, cancel := context.WithCancel(context.Background())
	session, err := NewSession(ctx, false, "", nil, handler)
	if err != nil {
		t.Fatal(err)
		return
	}
	errC := session.RunLoop()

	if err = session.SubscribeAggTrade(ctx, "BTCUSDT", "BNBUSDT"); err != nil {
		t.Fatal(err)
	}

	<-handler.done

	cancel()
	if err = <-errC; err != nil {
		t.Fatal(err)
	}
}
