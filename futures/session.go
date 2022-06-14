package futures

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/crypto-zero/go-binance/v2/common"
)

type Session struct {
	common.WebsocketSession
	handler SessionHandler
}

type SessionHandler interface {
	common.WebsocketSessionHandler
	OnAggTrade(*WsAggTradeEvent)
	OnMarkPrice(*WsMarkPriceEvent)
	OnKline(*WsKlineEvent)
	OnContinuousKline(*WsContinuousKlineEvent)
	OnMiniMarketTicker(*WsMiniMarketTickerEvent)
	OnMarketTicker(*WsMarketTickerEvent)
	OnBookTicker(*WsBookTickerEvent)
	OnWsLiquidationOrder(*WsLiquidationOrderEvent)
	OnDepth(*WsDepthEvent)
	OnCompositeIndex(*WsCompositeIndexEvent)
	OnUserData(*WsUserDataEvent)
}

func (s *Session) SubscribeAggTrade(ctx context.Context, symbol ...string) (err error) {
	var streams []string
	for _, s := range symbol {
		streams = append(streams, fmt.Sprintf("%s@aggTrade", strings.ToLower(s)))
	}
	return s.SubscribeNoReply(ctx, streams...)
}

func (s *Session) SubscribeMarkPrice(ctx context.Context, symbol ...string) (err error) {
	var streams []string
	for _, s := range symbol {
		streams = append(streams, fmt.Sprintf("%s@markPrice", strings.ToLower(s)))
	}
	return s.SubscribeNoReply(ctx, streams...)
}

func (s *Session) SubscribeAllMarkPrice(ctx context.Context) error {
	return s.SubscribeNoReply(ctx, "!markPrice@arr@1s")
}

func (s *Session) KlineStreamName(symbol string, interval KlineInterval) string {
	return fmt.Sprintf("%s@kline_%s", strings.ToLower(symbol), strings.ToLower(string(interval)))
}

func (s *Session) SubscribeKline(ctx context.Context, symbol string, interval KlineInterval) error {
	return s.SubscribeNoReply(ctx, s.KlineStreamName(symbol, interval))
}

func (s *Session) SubscribeContinuousKline(ctx context.Context, symbol string,
	contractType ContractType, interval KlineInterval,
) error {
	return s.SubscribeNoReply(
		ctx,
		fmt.Sprintf(
			"%s_%s@continuousKline_%s",
			strings.ToLower(symbol), strings.ToLower(string(contractType)), interval,
		),
	)
}

func (s *Session) SubscribeMiniMarketTicker(ctx context.Context, symbol ...string) error {
	var streams []string
	for _, s := range symbol {
		streams = append(streams, fmt.Sprintf("%s@miniTicker", strings.ToLower(s)))
	}
	return s.SubscribeNoReply(ctx, streams...)
}

func (s *Session) SubscribeAllMiniMarketTicker(ctx context.Context) error {
	return s.SubscribeNoReply(ctx, "!miniTicker@arr")
}

func (s *Session) SubscribeMarketTicker(ctx context.Context, symbol ...string) error {
	var streams []string
	for _, s := range symbol {
		streams = append(streams, fmt.Sprintf("%s@ticker", strings.ToLower(s)))
	}
	return s.SubscribeNoReply(ctx, streams...)
}

func (s *Session) SubscribeAllMarketTicker(ctx context.Context) error {
	return s.SubscribeNoReply(ctx, "!ticker@arr")
}

func (s *Session) SubscribeBookTicker(ctx context.Context, symbol ...string) error {
	var streams []string
	for _, s := range symbol {
		streams = append(streams, fmt.Sprintf("%s@bookTicker", strings.ToLower(s)))
	}
	return s.SubscribeNoReply(ctx, streams...)
}

func (s *Session) SubscribeAllBookTicker(ctx context.Context) error {
	return s.SubscribeNoReply(ctx, "!bookTicker")
}

func (s *Session) SubscribeLiquidationOrder(ctx context.Context, symbol ...string) error {
	var streams []string
	for _, s := range symbol {
		streams = append(streams, fmt.Sprintf("%s@forceOrder", strings.ToLower(s)))
	}
	return s.SubscribeNoReply(ctx, streams...)
}

func (s *Session) SubscribeAllLiquidationOrder(ctx context.Context) error {
	return s.SubscribeNoReply(ctx, "!forceOrder@arr")
}

func (s *Session) SubscribeDepth(ctx context.Context, symbol string, level int,
	interval time.Duration,
) error {
	stream := fmt.Sprintf("%s@depth", strings.ToLower(symbol))
	if level > 0 {
		stream = fmt.Sprintf("%s%d", stream, level)
	}
	if interval > 0 {
		stream = fmt.Sprintf("%s@%s", stream, interval.String())
	}
	return s.SubscribeNoReply(ctx, stream)
}

func (s *Session) SubscribeCompositeIndex(ctx context.Context, symbol ...string) error {
	var streams []string
	for _, s := range symbol {
		streams = append(streams, fmt.Sprintf("%s@compositeIndex", strings.ToLower(s)))
	}
	return s.SubscribeNoReply(ctx, streams...)
}

func (s *Session) registerHandler(handler SessionHandler) {
	s.RegisterMessageHandler(
		common.WebsocketSessionMessageFactoryBuild[WsAggTradeEvent](),
		common.WebsocketSessionMessageHandlerBuild(handler.OnAggTrade),
		s.RequireMapKeyValue("e", "aggTrade"),
	)
	s.RegisterMessageHandler(
		common.WebsocketSessionMessageFactoryBuild[WsMarkPriceEvent](),
		common.WebsocketSessionMessageHandlerBuild(handler.OnMarkPrice),
		s.RequireMapKeyValue("e", "markPriceUpdate"),
	)
	s.RegisterMessageHandler(
		common.WebsocketSessionMessageFactoryBuild[WsKlineEvent](),
		common.WebsocketSessionMessageHandlerBuild(handler.OnKline),
		s.RequireMapKeyValue("e", "kline"),
	)
	s.RegisterMessageHandler(
		common.WebsocketSessionMessageFactoryBuild[WsContinuousKlineEvent](),
		common.WebsocketSessionMessageHandlerBuild(handler.OnContinuousKline),
		s.RequireMapKeyValue("e", "continuous_kline"),
	)
	s.RegisterMessageHandler(
		common.WebsocketSessionMessageFactoryBuild[WsMiniMarketTickerEvent](),
		common.WebsocketSessionMessageHandlerBuild(handler.OnMiniMarketTicker),
		s.RequireMapKeyValue("e", "24hrMiniTicker"),
	)
	s.RegisterMessageHandler(
		common.WebsocketSessionMessageFactoryBuild[WsMarketTickerEvent](),
		common.WebsocketSessionMessageHandlerBuild(handler.OnMarketTicker),
		s.RequireMapKeyValue("e", "24hrTicker"),
	)
	s.RegisterMessageHandler(
		common.WebsocketSessionMessageFactoryBuild[WsBookTickerEvent](),
		common.WebsocketSessionMessageHandlerBuild(handler.OnBookTicker),
		s.RequireMapKeyValue("e", "bookTicker"),
	)
	s.RegisterMessageHandler(
		common.WebsocketSessionMessageFactoryBuild[WsLiquidationOrderEvent](),
		common.WebsocketSessionMessageHandlerBuild(handler.OnWsLiquidationOrder),
		s.RequireMapKeyValue("e", "forceOrder"),
	)
	s.RegisterMessageHandler(
		common.WebsocketSessionMessageFactoryBuild[WsDepthEvent](),
		common.WebsocketSessionMessageHandlerBuild(handler.OnDepth),
		s.RequireMapKeyValue("e", "depthUpdate"),
	)
	s.RegisterMessageHandler(
		common.WebsocketSessionMessageFactoryBuild[WsCompositeIndexEvent](),
		common.WebsocketSessionMessageHandlerBuild(handler.OnCompositeIndex),
		s.RequireMapKeyValue("e", "compositeIndex"),
	)

	s.RegisterMessageHandler(
		common.WebsocketSessionMessageFactoryBuild[WsUserDataEvent](),
		common.WebsocketSessionMessageHandlerBuild(handler.OnUserData),
		s.RequireMapKeyValue("e", string(UserDataEventTypeListenKeyExpired)),
	)
	s.RegisterMessageHandler(
		common.WebsocketSessionMessageFactoryBuild[WsUserDataEvent](),
		common.WebsocketSessionMessageHandlerBuild(handler.OnUserData),
		s.RequireMapKeyValue("e", string(UserDataEventTypeMarginCall)),
	)
	s.RegisterMessageHandler(
		common.WebsocketSessionMessageFactoryBuild[WsUserDataEvent](),
		common.WebsocketSessionMessageHandlerBuild(handler.OnUserData),
		s.RequireMapKeyValue("e", string(UserDataEventTypeAccountUpdate)),
	)
	s.RegisterMessageHandler(
		common.WebsocketSessionMessageFactoryBuild[WsUserDataEvent](),
		common.WebsocketSessionMessageHandlerBuild(handler.OnUserData),
		s.RequireMapKeyValue("e", string(UserDataEventTypeOrderTradeUpdate)),
	)
	s.RegisterMessageHandler(
		common.WebsocketSessionMessageFactoryBuild[WsUserDataEvent](),
		common.WebsocketSessionMessageHandlerBuild(handler.OnUserData),
		s.RequireMapKeyValue("e", string(UserDataEventTypeAccountConfigUpdate)),
	)
}

func newMockSession(handler SessionHandler) (*Session, common.MockWebsocketSession) {
	wss := common.NewMockWebsocketSession(handler)
	session := new(Session)
	session.WebsocketSession = wss
	session.handler = handler
	session.registerHandler(handler)
	return session, wss
}

func NewSession(ctx context.Context, testnet bool, listenKey string, proxyURL *url.URL,
	handler SessionHandler,
) (session *Session, err error) {
	address := baseWsMainUrl
	if testnet {
		address = baseWsTestnetUrl
	}
	if listenKey != "" {
		address = fmt.Sprintf("%s/%s", address, listenKey)
	}

	var cli common.WebsocketClient
	if proxyURL == nil {
		cli, err = common.WebsocketDial(ctx, address, nil)
	} else {
		cli, err = common.WebsocketDialProxy(ctx, address, proxyURL)
	}
	if err != nil {
		return nil, err
	}

	session = new(Session)
	session.WebsocketSession = common.NewWebsocketSession(cli, handler)
	session.handler = handler
	session.registerHandler(handler)
	return session, nil
}
