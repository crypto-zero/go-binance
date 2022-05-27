package futures

import (
	"context"
	"os"
	"testing"
)

type testSessionHandler struct {
	*testing.T
	done chan struct{}

	aggTrade, markPrice, kline bool

	markPriceCount int
}

func (t *testSessionHandler) OnUnknownMessage(bytes []byte, i interface{}) error {
	t.Logf("got unknown message: %v\n", i)
	return nil
}

func (t *testSessionHandler) OnClose(err error) {
	t.Logf("got on close err: %v\n", err)
}

func (t *testSessionHandler) OnAggTrade(event *WsAggTradeEvent) {
	t.Logf("got agg trade event: %#v\n", event)
	t.aggTrade = true
	t.triggerDone()
}

func (t *testSessionHandler) OnMarkPrice(event *WsMarkPriceEvent) {
	t.Logf("got mark price event: %#v\n", event)
	t.markPrice = true
	t.markPriceCount++
	t.triggerDone()
}

func (t *testSessionHandler) OnKline(line *WsKlineEvent) {
	t.Logf("got kline event: %#v\n", line)
	t.kline = true
	t.triggerDone()
}

func (t *testSessionHandler) triggerDone() {
	if !t.aggTrade || !t.markPrice || !t.kline ||
		t.markPriceCount < 10 || t.done == nil {
		return
	}
	close(t.done)
	t.done = nil
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
	if err = session.SubscribeMarkPrice(ctx, "BTCUSDT", "BNBUSDT"); err != nil {
		t.Fatal(err)
	}
	if err = session.SubscribeAllMarkPrice(ctx); err != nil {
		t.Fatal(err)
	}
	if err = session.SubscribeKline(ctx, "BTCUSDT", KlineInterval1Minute); err != nil {
		t.Fatal(err)
	}

	<-handler.done

	cancel()
	if err = <-errC; err != nil {
		t.Fatal(err)
	}
}
