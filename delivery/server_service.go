package delivery

import (
	"context"

	"github.com/crypto-zero/go-binance/v2/common"
)

// PingService ping server
type PingService struct {
	c *Client
}

// Do send request
func (s *PingService) Do(ctx context.Context, opts ...common.RequestOption) (err error) {
	r := common.NewGetRequestPublic("/dapi/v1/ping")
	_, err = s.c.CallAPIBytes(ctx, r, opts...)
	return err
}

// ServerTimeService get server time
type ServerTimeService struct {
	c *Client
}

// Do send request
func (s *ServerTimeService) Do(ctx context.Context, opts ...common.RequestOption) (serverTime int64, err error) {
	r := common.NewGetRequestPublic("/dapi/v1/time")
	data, err := s.c.CallAPIBytes(ctx, r, opts...)
	if err != nil {
		return 0, err
	}
	j, err := newJSON(data)
	if err != nil {
		return 0, err
	}
	serverTime = j.Get("serverTime").MustInt64()
	return serverTime, nil
}

// SetServerTimeService set server time
type SetServerTimeService struct {
	c *Client
}

// Do send request
func (s *SetServerTimeService) Do(ctx context.Context, opts ...common.RequestOption) (timeOffset int64, err error) {
	serverTime, err := s.c.NewServerTimeService().Do(ctx, opts...)
	if err != nil {
		return 0, err
	}
	timeOffset = currentTimestamp() - serverTime
	s.c.UpdateTimeOffset(timeOffset)
	return timeOffset, nil
}
