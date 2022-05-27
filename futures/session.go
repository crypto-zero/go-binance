package futures

import (
	"context"
	"fmt"
	"net/url"
	"strings"

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

	session.RegisterMessageHandler(
		common.WebsocketSessionMessageFactoryBuild[WsAggTradeEvent](),
		common.WebsocketSessionMessageHandlerBuild(handler.OnAggTrade),
		session.RequireMapKeyValue("e", "aggTrade"),
	)

	session.RegisterMessageHandler(
		common.WebsocketSessionMessageFactoryBuild[WsMarkPriceEvent](),
		common.WebsocketSessionMessageHandlerBuild(handler.OnMarkPrice),
		session.RequireMapKeyValue("e", "markPriceUpdate"),
	)
	return session, nil
}
