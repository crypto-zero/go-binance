package futures

import (
	"context"
	"encoding/json"

	"github.com/crypto-zero/go-binance/v2/common"
)

// PremiumIndexService get premium index
type PremiumIndexService struct {
	c      *Client
	symbol *string
}

// Symbol set symbol
func (s *PremiumIndexService) Symbol(symbol string) *PremiumIndexService {
	s.symbol = &symbol
	return s
}

// Do send request
func (s *PremiumIndexService) Do(ctx context.Context, opts ...common.RequestOption) (res []*PremiumIndex, err error) {
	r := common.NewGetRequestPublic("/fapi/v1/premiumIndex")
	if s.symbol != nil {
		r.SetQuery("symbol", *s.symbol)
	}

	f := func(data []byte) error {
		res = make([]*PremiumIndex, 0)
		data = common.ToJSONList(data)
		if err = json.Unmarshal(data, &res); err != nil {
			return err
		}
		return nil
	}
	if err = s.c.CallAPI(ctx, r, f, opts...); err != nil {
		return nil, err
	}
	return res, nil
}

// PremiumIndex define premium index of mark price
type PremiumIndex struct {
	Symbol          string `json:"symbol"`
	MarkPrice       string `json:"markPrice"`
	LastFundingRate string `json:"lastFundingRate"`
	NextFundingTime int64  `json:"nextFundingTime"`
	Time            int64  `json:"time"`
}

// FundingRateService get funding rate
type FundingRateService struct {
	c         *Client
	symbol    string
	startTime *int64
	endTime   *int64
	limit     *int
}

// Symbol set symbol
func (s *FundingRateService) Symbol(symbol string) *FundingRateService {
	s.symbol = symbol
	return s
}

// StartTime set startTime
func (s *FundingRateService) StartTime(startTime int64) *FundingRateService {
	s.startTime = &startTime
	return s
}

// EndTime set startTime
func (s *FundingRateService) EndTime(endTime int64) *FundingRateService {
	s.endTime = &endTime
	return s
}

// Limit set limit
func (s *FundingRateService) Limit(limit int) *FundingRateService {
	s.limit = &limit
	return s
}

// Do send request
func (s *FundingRateService) Do(ctx context.Context, opts ...common.RequestOption) (res []*FundingRate, err error) {
	r := common.NewGetRequestPublic("/fapi/v1/fundingRate")
	r.SetQuery("symbol", s.symbol)
	if s.startTime != nil {
		r.SetQuery("startTime", *s.startTime)
	}
	if s.endTime != nil {
		r.SetQuery("endTime", *s.endTime)
	}
	if s.limit != nil {
		r.SetQuery("limit", *s.limit)
	}

	res = make([]*FundingRate, 0)
	if err = s.c.CallAPI(ctx, r, &res, opts...); err != nil {
		return nil, err
	}
	return res, nil
}

// FundingRate define funding rate of mark price
type FundingRate struct {
	Symbol      string `json:"symbol"`
	FundingRate string `json:"fundingRate"`
	FundingTime int64  `json:"fundingTime"`
	Time        int64  `json:"time"`
}

// GetLeverageBracketService get funding rate
type GetLeverageBracketService struct {
	c      *Client
	symbol string
}

// Symbol set symbol
func (s *GetLeverageBracketService) Symbol(symbol string) *GetLeverageBracketService {
	s.symbol = symbol
	return s
}

// Do send request
func (s *GetLeverageBracketService) Do(ctx context.Context, opts ...common.RequestOption) (res []*LeverageBracket, err error) {
	r := common.NewGetRequestSigned("/fapi/v1/leverageBracket")
	r.SetQuery("symbol", s.symbol)
	if s.symbol != "" {
		r.SetQuery("symbol", s.symbol)
	}

	res = make([]*LeverageBracket, 0)
	f := func(data []byte) error {
		if s.symbol != "" {
			data = common.ToJSONList(data)
		}
		if err = json.Unmarshal(data, &res); err != nil {
			return err
		}
		return nil
	}
	if err = s.c.CallAPI(ctx, r, f, opts...); err != nil {
		return nil, err
	}
	return res, nil
}

// LeverageBracket define the leverage bracket
type LeverageBracket struct {
	Symbol   string    `json:"symbol"`
	Brackets []Bracket `json:"brackets"`
}

// Bracket define the bracket
type Bracket struct {
	Bracket          int     `json:"bracket"`
	InitialLeverage  int     `json:"initialLeverage"`
	NotionalCap      float64 `json:"notionalCap"`
	NotionalFloor    float64 `json:"notionalFloor"`
	MaintMarginRatio float64 `json:"maintMarginRatio"`
	Cum              float64 `json:"cum"`
}
