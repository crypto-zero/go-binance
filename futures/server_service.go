package futures

import (
	"context"
)

// PingService ping server
type PingService struct {
	c *Client
}

// Do send request
func (s *PingService) Do(ctx context.Context, opts ...RequestOption) (err error) {
	r := &request{
		method:   "GET",
		endpoint: "/fapi/v1/ping",
	}
	return s.c.callAPI(ctx, r, nil, opts...)
}

// ServerTimeService get server time
type ServerTimeService struct {
	c *Client
}

// Do send request
func (s *ServerTimeService) Do(ctx context.Context, opts ...RequestOption) (serverTime int64, err error) {
	r := &request{
		method:   "GET",
		endpoint: "/fapi/v1/time",
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

// Do send request
func (s *SetServerTimeService) Do(ctx context.Context, opts ...RequestOption) (timeOffset int64, err error) {
	serverTime, err := s.c.NewServerTimeService().Do(ctx)
	if err != nil {
		return 0, err
	}
	timeOffset = currentTimestamp() - serverTime
	s.c.TimeOffset = timeOffset
	return timeOffset, nil
}
