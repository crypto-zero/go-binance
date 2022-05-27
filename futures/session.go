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
	OnAggTradeEvent(*WsAggTradeEvent)
}

func (s *Session) SubscribeAggTrade(ctx context.Context, symbol ...string) (err error) {
	var streams []string
	for _, s := range symbol {
		streams = append(streams, fmt.Sprintf("%s@aggTrade", strings.ToLower(s)))
	}
	reply, err := s.Subscribe(ctx, streams...)
	if err != nil {
		return err
	}
	if err = reply.OK(); err != nil {
		return err
	}
	return nil
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

	// for WsAggTradeEvent message handler
	session.RegisterMessageHandler(
		session.RequireMapHasAllKeys("e", "E", "s", "a", "p", "q", "f", "l", "T", "m"),
		common.WebsocketSessionMessageFactoryBuild[WsAggTradeEvent](),
		common.WebsocketSessionMessageHandlerBuild(handler.OnAggTradeEvent),
	)
	return session, nil
}
