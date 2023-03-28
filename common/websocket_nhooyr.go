package common

import (
	"context"
	"errors"
	"math"
	"net/http"
	"net/url"
	"sync"
	"time"

	"nhooyr.io/websocket"
)

var _ WebsocketClient = (*nhooyrWebsocketClient)(nil)

type nhooyrWebsocketClient struct {
	*websocket.Conn
	ctx          context.Context
	delay        time.Duration
	writerBuffer chan []byte
	pingBuffer   chan struct{}
}

func (wc *nhooyrWebsocketClient) Ping() {
	select {
	case wc.pingBuffer <- struct{}{}:
	default:
	}
}

func (wc *nhooyrWebsocketClient) Delay() time.Duration {
	return wc.delay
}

func (wc *nhooyrWebsocketClient) Write(data []byte) {
	wc.writerBuffer <- data
}

func (wc *nhooyrWebsocketClient) Loop(f WebsocketMessageCallback) error {
	ctx, cancel := context.WithCancel(wc.ctx)

	var wg sync.WaitGroup
	wg.Add(3)

	c := make(chan error, 1)

	run := func(fn func() error) {
		var err error
		go func() {
			defer wg.Done()
			if err = fn(); err == nil {
				return
			}
			select {
			case c <- err:
			default:
			}
		}()
	}

	run(func() error { return wc.writeLoop(ctx) })
	run(func() error { return wc.readLoop(ctx, f) })
	run(func() error { return wc.pingLoop(ctx) })

	err := <-c
	cancel()
	wg.Wait()

	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return nil
	}
	return err
}

func (wc *nhooyrWebsocketClient) writeLoop(ctx context.Context) (err error) {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-wc.writerBuffer:
			if err = wc.Conn.Write(ctx, websocket.MessageText, msg); err != nil {
				return err
			}
		}
	}
}

func (wc *nhooyrWebsocketClient) readLoop(ctx context.Context, f WebsocketMessageCallback) (err error) {
	for {
		_, data, err := wc.Read(ctx)
		if err != nil {
			return err
		}
		if err = f(data); err != nil {
			return err
		}
	}
}

func (wc *nhooyrWebsocketClient) pingLoop(ctx context.Context) (err error) {
	d := 2 * time.Minute
	t := time.NewTimer(d)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-wc.pingBuffer:
			from := time.Now()
			if err = wc.Conn.Ping(ctx); err != nil {
				return err
			}
			wc.delay = time.Since(from)
		case <-t.C:
			from := time.Now()
			if err = wc.Conn.Ping(ctx); err != nil {
				return err
			}
			wc.delay = time.Since(from)
			t.Reset(d)
		}
	}
}

func NhooyrWebsocketDial(ctx context.Context, url string, httpClient *http.Client) (
	out WebsocketClient, err error,
) {
	// disable compression for
	// failed to WebSocket dial: unsupported permessage-deflate parameter: "server_max_window_bits=15"
	// pull request ref: https://github.com/nhooyr/websocket/pull/258
	// issue ref: https://github.com/nhooyr/websocket/issues/351
	opts := &websocket.DialOptions{
		HTTPClient:      httpClient,
		CompressionMode: websocket.CompressionDisabled,
	}
	conn, _, err := websocket.Dial(ctx, url, opts)
	if err != nil {
		return nil, err
	}
	conn.SetReadLimit(math.MaxUint16)
	cli := &nhooyrWebsocketClient{
		Conn:         conn,
		ctx:          ctx,
		writerBuffer: make(chan []byte, 100),
		pingBuffer:   make(chan struct{}, 0),
	}
	return cli, nil
}

func NhooyrWebsocketDialProxy(ctx context.Context, url string, proxyURL *url.URL) (
	out WebsocketClient, err error,
) {
	t := &http.Transport{}
	hc := &http.Client{Transport: t}
	if proxyURL != nil {
		t.Proxy = http.ProxyURL(proxyURL)
	}
	return NhooyrWebsocketDial(ctx, url, hc)
}
