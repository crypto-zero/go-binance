package futures

import (
	"context"
	"fmt"
	"net/url"

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
