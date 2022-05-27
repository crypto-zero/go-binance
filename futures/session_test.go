package futures

import (
	"context"
	"os"
	"testing"
	"time"
)

type testSessionHandler struct {
	*testing.T
	done chan struct{}

	aggTrade, markPrice, kline, continuousKline, miniMarketTicker,
	marketTicker, bookTicker bool

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
	// t.Logf("got agg trade event: %#v\n", event)
	t.aggTrade = true
	t.triggerDone()
}

func (t *testSessionHandler) OnMarkPrice(event *WsMarkPriceEvent) {
	// t.Logf("got mark price event: %#v\n", event)
	t.markPrice = true
	t.markPriceCount++
	t.triggerDone()
}

func (t *testSessionHandler) OnKline(line *WsKlineEvent) {
	// t.Logf("got kline event: %#v\n", line)
	t.kline = true
	t.triggerDone()
}

func (t *testSessionHandler) OnContinuousKline(line *WsContinuousKlineEvent) {
	// t.Logf("got continuous kline event: %#v\n", line)
	t.continuousKline = true
	t.triggerDone()
}

func (t *testSessionHandler) OnMiniMarketTicker(ticker *WsMiniMarketTickerEvent) {
	// t.Logf("got mini market ticker event: %#v\n", ticker)
	t.miniMarketTicker = true
	t.triggerDone()
}

func (t *testSessionHandler) OnMarketTicker(ticker *WsMarketTickerEvent) {
	// t.Logf("got market ticker event: %#v\n", ticker)
	t.marketTicker = true
	t.triggerDone()
}

func (t *testSessionHandler) OnBookTicker(ticker *WsBookTickerEvent) {
	// t.Logf("got book ticker event: %#v\n", ticker)
	t.bookTicker = true
	t.triggerDone()
}

func (t *testSessionHandler) triggerDone() {
	if !t.aggTrade || !t.markPrice || !t.kline || !t.continuousKline || !t.miniMarketTicker ||
		!t.marketTicker || !t.bookTicker || t.markPriceCount < 10 || t.done == nil {
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
	if err = session.SubscribeContinuousKline(ctx, "BTCUSDT", ContractTypePerpetual,
		KlineInterval1Minute); err != nil {
		t.Fatal(err)
	}
	if err = session.SubscribeMiniMarketTicker(ctx, "BTCUSDT", "BNBUSDT"); err != nil {
		t.Fatal(err)
	}
	if err = session.SubscribeAllMiniMarketTicker(ctx); err != nil {
		t.Fatal(err)
	}
	if err = session.SubscribeMarketTicker(ctx, "BTCUSDT", "BNBUSDT"); err != nil {
		t.Fatal(err)
	}
	if err = session.SubscribeAllMarketTicker(ctx); err != nil {
		t.Fatal(err)
	}
	if err = session.SubscribeBookTicker(ctx, "BTCUSDT", "BNBUSDT"); err != nil {
		t.Fatal(err)
	}

	// sleep for a while
	time.Sleep(time.Second)

	if err = session.SubscribeAllBookTicker(ctx); err != nil {
		t.Fatal(err)
	}

	if handler.done != nil {
		<-handler.done
	}

	cancel()
	if err = <-errC; err != nil {
		t.Fatal(err)
	}
}
