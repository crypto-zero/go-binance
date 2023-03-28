package common

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var _ WebsocketClient = (*gorillaWebsocketClient)(nil)

type gorillaWebsocketClient struct {
	*websocket.Conn
	ctx          context.Context
	delay        time.Duration
	writerBuffer chan []byte
	pingBuffer   chan struct{}
}

func (g *gorillaWebsocketClient) Loop(f WebsocketMessageCallback) (err error) {
	ctx, cancel := context.WithCancel(g.ctx)

	var wg sync.WaitGroup
	c := make(chan error, 1)

	run := func(fn func() error) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := fn(); err != nil {
				select {
				case c <- err:
				default:
				}
			}
		}()
	}

	run(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-g.pingBuffer:
				if err := g.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
					return err
				}
			case data := <-g.writerBuffer:
				if err := g.WriteMessage(websocket.TextMessage, data); err != nil {
					return err
				}
			}
		}
	})
	run(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
			_, data, err := g.Conn.ReadMessage()
			if err != nil {
				return err
			}
			if data != nil {
				if err := f(data); err != nil {
					return err
				}
			}
		}
	})

	err = <-c
	cancel()
	wg.Wait()

	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return nil
	}
	return err
}

func (g *gorillaWebsocketClient) Delay() time.Duration {
	return g.delay
}

func (g *gorillaWebsocketClient) Ping() {
	select {
	case g.pingBuffer <- struct{}{}:
	default:
	}
}

func (g *gorillaWebsocketClient) Write(data []byte) {
	g.writerBuffer <- data
}

func GorillaWebsocketDialProxy(ctx context.Context, url string, proxyURL *url.URL) (
	out WebsocketClient, err error,
) {
	dialer := websocket.Dialer{}
	if proxyURL != nil {
		dialer.Proxy = http.ProxyURL(proxyURL)
	}
	c, _, err := dialer.DialContext(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return &gorillaWebsocketClient{
		Conn:         c,
		ctx:          ctx,
		writerBuffer: make(chan []byte, 100),
		pingBuffer:   make(chan struct{}, 0),
	}, nil
}

func init() {
	DefaultWebsocketProvider = GorillaWebsocketDialProxy
}
