package binance

import (
	"context"

	"github.com/crypto-zero/go-binance/v2/common"
)

// StartUserStreamService create listen key for user stream service
type StartUserStreamService struct {
	c *Client
}

// Do send Request
func (s *StartUserStreamService) Do(ctx context.Context, opts ...common.RequestOption) (listenKey string, err error) {
	r := common.NewPostRequestAPIKey("/api/v3/userDataStream")
	f := func(data []byte) error {
		j, err := newJSON(data)
		if err != nil {
			return err
		}
		listenKey = j.Get("listenKey").MustString()
		return nil
	}
	if err = s.c.callAPI(ctx, r, f, opts...); err != nil {
		return "", err
	}
	return listenKey, nil
}

// KeepaliveUserStreamService update listen key
type KeepaliveUserStreamService struct {
	c         *Client
	listenKey string
}

// ListenKey set listen key
func (s *KeepaliveUserStreamService) ListenKey(listenKey string) *KeepaliveUserStreamService {
	s.listenKey = listenKey
	return s
}

// Do send Request
func (s *KeepaliveUserStreamService) Do(ctx context.Context, opts ...common.RequestOption) (err error) {
	r := common.NewPutRequestAPIKey("/api/v3/userDataStream")
	r.SetForm("listenKey", s.listenKey)
	return s.c.callAPI(ctx, r, nil, opts...)
}

// CloseUserStreamService delete listen key
type CloseUserStreamService struct {
	c         *Client
	listenKey string
}

// ListenKey set listen key
func (s *CloseUserStreamService) ListenKey(listenKey string) *CloseUserStreamService {
	s.listenKey = listenKey
	return s
}

// Do send Request
func (s *CloseUserStreamService) Do(ctx context.Context, opts ...common.RequestOption) (err error) {
	r := common.NewDeleteRequestAPIKey("/api/v3/userDataStream")
	r.SetForm("listenKey", s.listenKey)
	return s.c.callAPI(ctx, r, nil, opts...)
}
