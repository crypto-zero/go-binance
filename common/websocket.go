package common

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"sync"
	"time"

	"nhooyr.io/websocket"
)

type WebsocketMessageCallback func(messageType websocket.MessageType, data []byte) error

type websocketClient struct {
	*websocket.Conn
	ctx          context.Context
	delay        time.Duration
	writerBuffer chan []byte
	pingBuffer   chan struct{}
}

func (wc *websocketClient) Ping() {
	select {
	case wc.pingBuffer <- struct{}{}:
	default:
	}
}

func (wc *websocketClient) Delay() time.Duration {
	return wc.delay
}

func (wc *websocketClient) Write(data []byte) {
	wc.writerBuffer <- data
}

func (wc *websocketClient) Loop(f WebsocketMessageCallback) error {
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

func (wc *websocketClient) writeLoop(ctx context.Context) (err error) {
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

func (wc *websocketClient) readLoop(ctx context.Context, f WebsocketMessageCallback) (err error) {
	for {
		mt, data, err := wc.Read(ctx)
		if err != nil {
			return err
		}
		if err = f(mt, data); err != nil {
			return err
		}
	}
}

func (wc *websocketClient) pingLoop(ctx context.Context) (err error) {
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

type WebsocketClient interface {
	Loop(f WebsocketMessageCallback) error
	Delay() time.Duration
	Ping()
	Write(data []byte)
}

func WebsocketDial(ctx context.Context, url string, httpClient *http.Client) (
	out WebsocketClient, err error,
) {
	opts := &websocket.DialOptions{HTTPClient: httpClient}
	conn, _, err := websocket.Dial(ctx, url, opts)
	if err != nil {
		return nil, err
	}
	cli := &websocketClient{
		Conn:         conn,
		ctx:          ctx,
		writerBuffer: make(chan []byte, 100),
		pingBuffer:   make(chan struct{}, 0),
	}
	return cli, nil
}

func WebsocketDialProxy(ctx context.Context, url string, proxyURL *url.URL) (
	out WebsocketClient, err error,
) {
	hc := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}
	return WebsocketDial(ctx, url, hc)
}
