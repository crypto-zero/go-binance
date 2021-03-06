package binance

import (
	"context"
	"fmt"

	"github.com/crypto-zero/go-binance/v2/common"
)

// KlinesService list klines
type KlinesService struct {
	c         *Client
	symbol    string
	interval  string
	limit     *int
	startTime *int64
	endTime   *int64
}

// Symbol set symbol
func (s *KlinesService) Symbol(symbol string) *KlinesService {
	s.symbol = symbol
	return s
}

// Interval set interval
func (s *KlinesService) Interval(interval string) *KlinesService {
	s.interval = interval
	return s
}

// Limit set limit
func (s *KlinesService) Limit(limit int) *KlinesService {
	s.limit = &limit
	return s
}

// StartTime set startTime
func (s *KlinesService) StartTime(startTime int64) *KlinesService {
	s.startTime = &startTime
	return s
}

// EndTime set endTime
func (s *KlinesService) EndTime(endTime int64) *KlinesService {
	s.endTime = &endTime
	return s
}

// Do send Request
func (s *KlinesService) Do(ctx context.Context, opts ...common.RequestOption) (res []*Kline, err error) {
	r := common.NewGetRequestPublic("/api/v3/klines")
	r.SetQuery("symbol", s.symbol)
	r.SetQuery("interval", s.interval)
	if s.limit != nil {
		r.SetQuery("limit", *s.limit)
	}
	if s.startTime != nil {
		r.SetQuery("startTime", *s.startTime)
	}
	if s.endTime != nil {
		r.SetQuery("endTime", *s.endTime)
	}

	f := func(data []byte) error {
		j, err := newJSON(data)
		if err != nil {
			return err
		}
		num := len(j.MustArray())
		res = make([]*Kline, num)
		for i := 0; i < num; i++ {
			item := j.GetIndex(i)
			if len(item.MustArray()) < 11 {
				return fmt.Errorf("invalid kline response")
			}
			res[i] = &Kline{
				OpenTime:                 item.GetIndex(0).MustInt64(),
				Open:                     item.GetIndex(1).MustString(),
				High:                     item.GetIndex(2).MustString(),
				Low:                      item.GetIndex(3).MustString(),
				Close:                    item.GetIndex(4).MustString(),
				Volume:                   item.GetIndex(5).MustString(),
				CloseTime:                item.GetIndex(6).MustInt64(),
				QuoteAssetVolume:         item.GetIndex(7).MustString(),
				TradeNum:                 item.GetIndex(8).MustInt64(),
				TakerBuyBaseAssetVolume:  item.GetIndex(9).MustString(),
				TakerBuyQuoteAssetVolume: item.GetIndex(10).MustString(),
			}
		}
		return nil
	}
	if err = s.c.CallAPI(ctx, r, f, opts...); err != nil {
		return nil, err
	}
	return res, nil
}

// Kline define kline info
type Kline struct {
	OpenTime                 int64  `json:"openTime"`
	Open                     string `json:"open"`
	High                     string `json:"high"`
	Low                      string `json:"low"`
	Close                    string `json:"close"`
	Volume                   string `json:"volume"`
	CloseTime                int64  `json:"closeTime"`
	QuoteAssetVolume         string `json:"quoteAssetVolume"`
	TradeNum                 int64  `json:"tradeNum"`
	TakerBuyBaseAssetVolume  string `json:"takerBuyBaseAssetVolume"`
	TakerBuyQuoteAssetVolume string `json:"takerBuyQuoteAssetVolume"`
}
