package futures

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"
)

type testSessionHandler struct {
	*testing.T
	done chan struct{}

	aggTrade, markPrice, kline, continuousKline, miniMarketTicker,
	marketTicker, bookTicker, forceOrder, depth, compositeIndex, userData bool
	userDataMarginCall, userDataAccountUpdate, userDataOrderUpdate,
	userDataConfigUpdated, userDataLicenseKeyExpired bool
	doneWhenUserData bool

	markPriceCount int
}

func newTestSessionHandler(t *testing.T) *testSessionHandler {
	return &testSessionHandler{T: t, done: make(chan struct{})}
}

func (t *testSessionHandler) OnUnknownMessage(bytes []byte, i interface{}) error {
	t.Logf("got unknown message: %v\n", i)
	return nil
}

func (t *testSessionHandler) OnClose(err error) {
	t.Logf("got on close err: %v\n", err)
}

func (t *testSessionHandler) OnAggTrade(event *WsAggTradeEvent) {
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

func (t *testSessionHandler) OnWsLiquidationOrder(event *WsLiquidationOrderEvent) {
	// t.Logf("force order event: %#v\n", event)
	t.forceOrder = true
	t.triggerDone()
}

func (t *testSessionHandler) OnDepth(event *WsDepthEvent) {
	// t.Logf("depth event: %#v\n", event)
	t.depth = true
	t.triggerDone()
}

func (t *testSessionHandler) OnCompositeIndex(event *WsCompositeIndexEvent) {
	// t.Logf("composite index event: %#v\n", event)
	t.compositeIndex = true
	t.triggerDone()
}

func (t *testSessionHandler) OnUserData(event *WsUserDataEvent) {
	t.Logf("user data event: %#v\n", event)
	t.userData = true
	switch event.Event {
	case UserDataEventTypeListenKeyExpired:
		t.userDataLicenseKeyExpired = true
	case UserDataEventTypeAccountUpdate:
		t.userDataAccountUpdate = true
	case UserDataEventTypeAccountConfigUpdate:
		t.userDataConfigUpdated = true
	case UserDataEventTypeOrderTradeUpdate:
		t.userDataOrderUpdate = true
	case UserDataEventTypeMarginCall:
		t.userDataMarginCall = true
	}
	t.triggerDone()
}

func (t *testSessionHandler) triggerDone() {
	if t.done == nil {
		return
	}
	if !t.doneWhenUserData {
		if !t.aggTrade || !t.markPrice || !t.kline || !t.continuousKline ||
			!t.miniMarketTicker || !t.marketTicker || !t.bookTicker || !t.depth ||
			!t.compositeIndex || t.markPriceCount < 10 {
			return
		}
	} else {
		if !t.userData || !t.userDataOrderUpdate || !t.userDataConfigUpdated {
			return
		}
	}
	close(t.done)
	t.done = nil
	t.Log("handler ok.")
}

func TestSession(t *testing.T) {
	if value := os.Getenv("TEST_FUTURES_WS_SESSION"); value == "" {
		t.Skip("skip futures websocket session tests")
		return
	}

	handler := newTestSessionHandler(t)

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

	// sleep for a while
	time.Sleep(time.Second)

	if err = session.SubscribeBookTicker(ctx, "BTCUSDT", "BNBUSDT"); err != nil {
		t.Fatal(err)
	}
	if err = session.SubscribeAllBookTicker(ctx); err != nil {
		t.Fatal(err)
	}
	if err = session.SubscribeLiquidationOrder(ctx, "BTCUSDT"); err != nil {
		t.Fatal(err)
	}
	if err = session.SubscribeAllLiquidationOrder(ctx); err != nil {
		t.Fatal(err)
	}
	if err = session.SubscribeDepth(ctx, "BTCUSDT", 20, 100*time.Millisecond); err != nil {
		t.Fatal(err)
	}
	if err = session.SubscribeCompositeIndex(ctx, "DEFIUSDT"); err != nil {
		t.Fatal(err)
	}

	t.Log("waiting ..")
	if handler.done != nil {
		<-handler.done
	}
	t.Log("wait ok ..")

	cancel()
	if err = <-errC; err != nil {
		t.Fatal(err)
	}
}

func TestMockSession(t *testing.T) {
	handler := newTestSessionHandler(t)
	session, wss := newMockSession(handler)

	messages := []string{
		`{"e":"listenKeyExpired","E":1576653824250}`,
		`{"e":"MARGIN_CALL","E":1587727187525,"cw":"3.16812045","p":[{"s":"ETHUSDT","ps":"LONG","pa":"1.327","mt":"CROSSED","iw":"0","mp":"187.17127","up":"-1.166074","mm":"1.614445"}]}`,
		`{"e":"ACCOUNT_UPDATE","E":1564745798939,"T":1564745798938,"a":{"m":"ORDER","B":[{"a":"USDT","wb":"122624.12345678","cw":"100.12345678","bc":"50.12345678"},{"a":"BUSD","wb":"1.00000000","cw":"0.00000000","bc":"-49.12345678"}],"P":[{"s":"BTCUSDT","pa":"0","ep":"0.00000","cr":"200","up":"0","mt":"isolated","iw":"0.00000000","ps":"BOTH"},{"s":"BTCUSDT","pa":"20","ep":"6563.66500","cr":"0","up":"2850.21200","mt":"isolated","iw":"13200.70726908","ps":"LONG"},{"s":"BTCUSDT","pa":"-10","ep":"6563.86000","cr":"-45.04000000","up":"-1423.15600","mt":"isolated","iw":"6570.42511771","ps":"SHORT"}]}}`,
		`{"e":"ORDER_TRADE_UPDATE","E":1568879465651,"T":1568879465650,"o":{"s":"BTCUSDT","c":"TEST","S":"SELL","o":"TRAILING_STOP_MARKET","f":"GTC","q":"0.001","p":"0","ap":"0","sp":"7103.04","x":"NEW","X":"NEW","i":8886774,"l":"0","z":"0","L":"0","N":"USDT","n":"0","T":1568879465651,"t":0,"b":"0","a":"9.91","m":false,"R":false,"wt":"CONTRACT_PRICE","ot":"TRAILING_STOP_MARKET","ps":"LONG","cp":false,"AP":"7476.89","cr":"5.0","pP":false,"si":0,"ss":0,"rp":"0"}}
`,
		`{"e":"ACCOUNT_CONFIG_UPDATE","E":1611646737479,"T":1611646737476,"ac":{"s":"BTCUSDT","l":25}}`,
		`{"e":"ACCOUNT_CONFIG_UPDATE","E":1611646737479,"T":1611646737476,"ai":{"j":true}}`,
	}

	for _, msg := range messages {
		if err := wss.MockProcessMessage([]byte(msg)); err != nil {
			t.Fatal(err, msg)
		}
	}
	t.Log(session)
	if !handler.userData || !handler.userDataLicenseKeyExpired || !handler.userDataAccountUpdate ||
		!handler.userDataOrderUpdate || !handler.userDataConfigUpdated || !handler.userDataMarginCall {
		t.Fatal("handler did not get user events")
	}
}

func TestUserData(t *testing.T) {
	key, secret := os.Getenv("TEST_SESSION_KEY"), os.Getenv("TEST_SESSION_SECRET")
	if key == "" || secret == "" {
		t.Skip("skip test user data")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	handler := newTestSessionHandler(t)
	handler.doneWhenUserData = true

	c := NewClient(key, secret, false)
	s := c.NewStartUserStreamService()
	rsp, err := s.Do(ctx)
	if err != nil {
		t.Fatal(err)
	}

	session, err := NewSession(ctx, false, rsp, nil, handler)
	if err != nil {
		t.Fatal(err)
	}

	loopC := session.RunLoop()

	symbol := "EOSUSDT"
	leverage := 0

	accountReply, err := c.NewGetAccountService().Do(ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range accountReply.Positions {
		if p.Symbol == symbol {
			if leverage, err = strconv.Atoi(p.Leverage); err != nil {
				t.Fatal(err)
			}
		}
	}
	if leverage >= 50 {
		leverage = 30
	} else if leverage >= 10 && leverage < 50 {
		leverage++
	}
	t.Log(leverage)

	// do operations
	if _, err = c.NewChangeLeverageService().Symbol(symbol).Leverage(leverage).Do(ctx); err != nil {
		t.Fatal(err)
	}

	reply, err := c.NewCreateOrderService().Symbol(symbol).Side(SideTypeBuy).Price("0.5").
		Quantity("10").Type(OrderTypeLimit).TimeInForce(TimeInForceTypeGTC).Do(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(reply.OrderID)

	if _, err = c.NewCancelOrderService().Symbol(symbol).OrderID(reply.OrderID).Do(ctx); err != nil {
		t.Log(err)
	}

	if handler.done != nil {
		<-handler.done
	}
	cancel()
	<-loopC
}
