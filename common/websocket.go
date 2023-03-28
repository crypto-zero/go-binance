package common

import (
	"context"
	"net/url"
	"time"
)

type WebsocketMessageCallback func(data []byte) error

type WebsocketClient interface {
	Loop(f WebsocketMessageCallback) error
	Delay() time.Duration
	Ping()
	Write(data []byte)
}

var DefaultWebsocketProvider func(ctx context.Context, url string, proxyURL *url.URL) (
	WebsocketClient, error)
