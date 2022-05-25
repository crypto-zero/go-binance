package common

import (
	"context"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

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
