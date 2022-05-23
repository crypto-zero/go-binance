package futures

import (
	"context"

	"github.com/crypto-zero/go-binance/v2/common"
)

// GetPositionMarginHistoryService get position margin history service
type GetPositionMarginHistoryService struct {
	c         *Client
	symbol    string
	_type     *int
	startTime *int64
	endTime   *int64
	limit     *int64
}

// Symbol set symbol
func (s *GetPositionMarginHistoryService) Symbol(symbol string) *GetPositionMarginHistoryService {
	s.symbol = symbol
	return s
}

// Type set type
func (s *GetPositionMarginHistoryService) Type(_type int) *GetPositionMarginHistoryService {
	s._type = &_type
	return s
}

// StartTime set startTime
func (s *GetPositionMarginHistoryService) StartTime(startTime int64) *GetPositionMarginHistoryService {
	s.startTime = &startTime
	return s
}

// EndTime set endTime
func (s *GetPositionMarginHistoryService) EndTime(endTime int64) *GetPositionMarginHistoryService {
	s.endTime = &endTime
	return s
}

// Limit set limit
func (s *GetPositionMarginHistoryService) Limit(limit int64) *GetPositionMarginHistoryService {
	s.limit = &limit
	return s
}

// Do send request
func (s *GetPositionMarginHistoryService) Do(ctx context.Context, opts ...common.RequestOption) (res []*PositionMarginHistory, err error) {
	r := common.NewGetRequestSigned("/fapi/v1/positionMargin/history")
	r.SetQuery("symbol", s.symbol)
	if s._type != nil {
		r.SetQuery("type", *s._type)
	}
	if s.startTime != nil {
		r.SetQuery("startTime", *s.startTime)
	}
	if s.endTime != nil {
		r.SetQuery("endTime", *s.endTime)
	}
	if s.limit != nil {
		r.SetQuery("limit", *s.limit)
	}

	res = make([]*PositionMarginHistory, 0)
	if err = s.c.CallAPI(ctx, r, &res, opts...); err != nil {
		return nil, err
	}
	return res, nil
}

// PositionMarginHistory define position margin history info
type PositionMarginHistory struct {
	Amount       string `json:"amount"`
	Asset        string `json:"asset"`
	Symbol       string `json:"symbol"`
	Time         int64  `json:"time"`
	Type         int    `json:"type"`
	PositionSide string `json:"positionSide"`
}
