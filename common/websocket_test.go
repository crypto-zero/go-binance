package common

import (
	"context"
	"os"
	"sync"
	"testing"

	"nhooyr.io/websocket"
)

func TestWebsocketConnection(t *testing.T) {
	if value := os.Getenv("TEST_WS"); value == "" {
		t.Skip("skip websocket client tests")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	url := "wss://demo.piesocket.com/v3/channel_1?api_key=VCXCEuvhGcBDP7XhiJJUDvR1e1D3eiVjgZ9VRiaV&notify_self"

	cli, err := WebsocketDial(ctx, url, nil)
	if err != nil {
		t.Fatal(err)
		return
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
		err := cli.Loop(f)
		done <- struct{}{}
		if err != nil {
			t.Error(err)
		}
	}()

	cli.Ping()
	cli.Write([]byte("hello world"))

	<-done

	// request client to close
	cancel()

	wg.Wait()
	t.Log(cli.Delay())
}
