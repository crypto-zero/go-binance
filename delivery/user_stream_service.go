package delivery

import (
	"context"

	"github.com/crypto-zero/go-binance/v2/common"
)

// StartUserStreamService create listen key for user stream service
type StartUserStreamService struct {
	c *Client
}

// Do send request
func (s *StartUserStreamService) Do(ctx context.Context, opts ...common.RequestOption) (listenKey string, err error) {
	r := common.NewPostRequestSigned("/dapi/v1/listenKey")
	data, err := s.c.CallAPIBytes(ctx, r, opts...)
	if err != nil {
		return "", err
	}
	j, err := newJSON(data)
	if err != nil {
		return "", err
	}
	listenKey = j.Get("listenKey").MustString()
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

// Do send request
func (s *KeepaliveUserStreamService) Do(ctx context.Context, opts ...common.RequestOption) (err error) {
	r := common.NewPutRequestSigned("/dapi/v1/listenKey")
	r.SetForm("listenKey", s.listenKey)
	_, err = s.c.CallAPIBytes(ctx, r, opts...)
	return err
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

// Do send request
func (s *CloseUserStreamService) Do(ctx context.Context, opts ...common.RequestOption) (err error) {
	r := common.NewDeleteRequestSigned("/dapi/v1/listenKey")
	r.SetForm("listenKey", s.listenKey)
	_, err = s.c.CallAPIBytes(ctx, r, opts...)
	return err
}
