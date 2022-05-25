package common

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"nhooyr.io/websocket"
)

type WebsocketMessageCallback func(messageType websocket.MessageType, data []byte) error

type websocketClient struct {
	*websocket.Conn
	writerBuffer chan []byte
	pingBuffer   chan struct{}
}

func (wc *websocketClient) Ping() {
	select {
	case wc.pingBuffer <- struct{}{}:
	default:
	}
}

func (wc *websocketClient) Write(data []byte) {
	wc.writerBuffer <- data
}

func (wc *websocketClient) Loop(ctx context.Context, f WebsocketMessageCallback) error {
	ctx, cancel := context.WithCancel(ctx)

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

	var ce websocket.CloseError
	if errors.As(err, &ce) && ce.Code == websocket.StatusNormalClosure &&
		ce.Reason == "wsc request close" {
		return nil
	}
	return err
}

func (wc *websocketClient) Close() error {
	err := wc.Conn.Close(websocket.StatusNormalClosure, "wsc request close")
	if err != nil {
		return err
	}
	return nil
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
			if err = wc.Conn.Ping(ctx); err != nil {
				return err
			}
		case <-t.C:
			if err = wc.Conn.Ping(ctx); err != nil {
				return err
			}
			t.Reset(d)
		}
	}
}

type WebsocketClient interface {
	Loop(ctx context.Context, f WebsocketMessageCallback) error
	Ping()
	Write(data []byte)
	Close() error
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
		writerBuffer: make(chan []byte, 100),
		pingBuffer:   make(chan struct{}, 0),
	}
	return cli, nil
}

type websocketSubscribe struct {
	WebsocketClient
	limiter                  *rate.Limiter
	pendingSubscribeTopics   []string
	pendingUnsubscribeTopics []string
	topicLock                sync.Mutex
	topics                   map[string]bool
	subscribeRunning         int32
}

func (wss *websocketSubscribe) Run(ctx context.Context) {
	d := time.Second
	t := time.NewTimer(d)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			wss.subscribe()
			t.Reset(d)
		}
	}
}

func (wss *websocketSubscribe) Subscribe(topics ...string) bool {
	wss.topicLock.Lock()
	defer wss.topicLock.Unlock()
	if len(topics)+len(wss.pendingSubscribeTopics)+len(wss.topics) > 200 {
		return false
	}
	wss.pendingSubscribeTopics = append(wss.pendingSubscribeTopics, topics...)
	return true
}

func (wss *websocketSubscribe) subscribe() {
	var subscribeBuffers, unsubscribeBuffers []string

	wss.topicLock.Lock()
	for len(wss.pendingUnsubscribeTopics) > 0 {
		r := wss.limiter.Reserve()
		if !r.OK() {
			return
		}

		topic := wss.pendingUnsubscribeTopics[0]
		for i := 0; i < len(wss.pendingUnsubscribeTopics)-1; i++ {
			wss.pendingUnsubscribeTopics[i] = wss.pendingUnsubscribeTopics[i+1]
		}
		wss.pendingUnsubscribeTopics = wss.pendingUnsubscribeTopics[:len(wss.pendingUnsubscribeTopics)-1]

		if _, ok := wss.topics[topic]; !ok {
			r.Cancel()
			continue
		}

		delete(wss.topics, topic)
		unsubscribeBuffers = append(unsubscribeBuffers, topic)
	}

	for len(wss.pendingSubscribeTopics) > 0 {
		r := wss.limiter.Reserve()
		if !r.OK() {
			return
		}

		topic := wss.pendingSubscribeTopics[0]
		for i := 0; i < len(wss.pendingSubscribeTopics)-1; i++ {
			wss.pendingSubscribeTopics[i] = wss.pendingSubscribeTopics[i+1]
		}
		wss.pendingSubscribeTopics = wss.pendingSubscribeTopics[:len(wss.pendingSubscribeTopics)-1]

		if _, ok := wss.topics[topic]; ok || len(wss.topics) >= 200 {
			r.Cancel()
			continue
		}

		wss.topics[topic] = true
		subscribeBuffers = append(subscribeBuffers, topic)
	}
	wss.topicLock.Unlock()
}

type WebsocketSubscribe interface {
	Run(ctx context.Context)
	Subscribe(topics ...string) bool
}

func NewWebsocketSubscribe(wsc WebsocketClient) WebsocketSubscribe {
	return &websocketSubscribe{
		WebsocketClient: wsc,
		limiter:         rate.NewLimiter(10, 10),
		topics:          make(map[string]bool),
	}
}
