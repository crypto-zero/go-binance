package binance

import (
	"context"

	"github.com/crypto-zero/go-binance/v2/common"
)

// PingService ping server
type PingService struct {
	c *Client
}

// Do send Request
func (s *PingService) Do(ctx context.Context, opts ...common.RequestOption) (err error) {
	r := &common.Request{
		Method:   "GET",
		Endpoint: "/api/v3/ping",
	}
	return s.c.callAPI(ctx, r, nil, opts...)
}

// ServerTimeService get server time
type ServerTimeService struct {
	c *Client
}

// Do send Request
func (s *ServerTimeService) Do(ctx context.Context, opts ...common.RequestOption) (serverTime int64, err error) {
	r := &common.Request{
		Method:   "GET",
		Endpoint: "/api/v3/time",
	}

	f := func(data []byte) error {
		j, err := newJSON(data)
		if err != nil {
			return err
		}
		serverTime = j.Get("serverTime").MustInt64()
		return nil
	}
	if err = s.c.callAPI(ctx, r, f, opts...); err != nil {
		return 0, err
	}
	return serverTime, nil
}

// SetServerTimeService set server time
type SetServerTimeService struct {
	c *Client
}

// Do send Request
func (s *SetServerTimeService) Do(ctx context.Context, opts ...common.RequestOption) (timeOffset int64, err error) {
	serverTime, err := s.c.NewServerTimeService().Do(ctx)
	if err != nil {
		return 0, err
	}
	timeOffset = currentTimestamp() - serverTime
	s.c.TimeOffset = timeOffset
	return timeOffset, nil
}
