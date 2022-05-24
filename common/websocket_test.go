package common

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"nhooyr.io/websocket"
)

func TestWebsocketConnection(t *testing.T) {
	if value := os.Getenv("TEST_WS"); value == "" {
		t.Skip("skip websocket client tests")
		return
	}

	ctx := context.Background()
	url := "wss://demo.piesocket.com/v3/channel_1?api_key=VCXCEuvhGcBDP7XhiJJUDvR1e1D3eiVjgZ9VRiaV&notify_self"

	cli, err := WebsocketDial(ctx, url, nil)
	if err != nil {
		panic(err)
	}

	done := make(chan struct{}, 2)

	f := func(mt websocket.MessageType, data []byte) error {
		t.Logf("ws client got message: %s\n", string(data))
		done <- struct{}{}
		return nil
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err = cli.Loop(ctx, f); err != nil {
			panic(err)
		} else {
			done <- struct{}{}
		}
	}()

	cli.Ping()
	cli.Ping()
	cli.Ping()
	cli.Ping()
	cli.Ping()

	time.Sleep(time.Minute)

	cli.Write([]byte("hello world"))

	<-done
	if err = cli.Close(); err != nil {
		panic(err)
	}
	wg.Wait()
}
