package binance

import (
	"context"
	"encoding/json"

	"github.com/crypto-zero/go-binance/v2/common"
)

// ListBookTickersService list best price/qty on the order book for a symbol or symbols
type ListBookTickersService struct {
	c      *Client
	symbol *string
}

// Symbol set symbol
func (s *ListBookTickersService) Symbol(symbol string) *ListBookTickersService {
	s.symbol = &symbol
	return s
}

// Do send Request
func (s *ListBookTickersService) Do(ctx context.Context, opts ...common.RequestOption) (res []*BookTicker, err error) {
	r := common.NewGetRequestPublic("/api/v3/ticker/bookTicker")
	if s.symbol != nil {
		r.SetQuery("symbol", *s.symbol)
	}

	f := func(data []byte) error {
		res = make([]*BookTicker, 0)
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

// BookTicker define book ticker info
type BookTicker struct {
	Symbol      string `json:"symbol"`
	BidPrice    string `json:"bidPrice"`
	BidQuantity string `json:"bidQty"`
	AskPrice    string `json:"askPrice"`
	AskQuantity string `json:"askQty"`
}

// ListPricesService list latest price for a symbol or symbols
type ListPricesService struct {
	c      *Client
	symbol *string
}

// Symbol set symbol
func (s *ListPricesService) Symbol(symbol string) *ListPricesService {
	s.symbol = &symbol
	return s
}

// Do send Request
func (s *ListPricesService) Do(ctx context.Context, opts ...common.RequestOption) (res []*SymbolPrice, err error) {
	r := common.NewGetRequestPublic("/api/v3/ticker/price")
	if s.symbol != nil {
		r.SetQuery("symbol", *s.symbol)
	}

	f := func(data []byte) error {
		res = make([]*SymbolPrice, 0)
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

// SymbolPrice define symbol and price pair
type SymbolPrice struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

// ListPriceChangeStatsService show stats of price change in last 24 hours for all symbols
type ListPriceChangeStatsService struct {
	c      *Client
	symbol *string
}

// Symbol set symbol
func (s *ListPriceChangeStatsService) Symbol(symbol string) *ListPriceChangeStatsService {
	s.symbol = &symbol
	return s
}

// Do send Request
func (s *ListPriceChangeStatsService) Do(ctx context.Context, opts ...common.RequestOption) (res []*PriceChangeStats, err error) {
	r := common.NewGetRequestPublic("/api/v3/ticker/24hr")
	if s.symbol != nil {
		r.SetQuery("symbol", *s.symbol)
	}

	f := func(data []byte) error {
		res = make([]*PriceChangeStats, 0)
		data = common.ToJSONList(data)
		if err = json.Unmarshal(data, &res); err != nil {
			return err
		}
		return nil
	}
	if err = s.c.CallAPI(ctx, r, f, opts...); err != nil {
		return res, err
	}
	return res, nil
}

// PriceChangeStats define price change stats
type PriceChangeStats struct {
	Symbol             string `json:"symbol"`
	PriceChange        string `json:"priceChange"`
	PriceChangePercent string `json:"priceChangePercent"`
	WeightedAvgPrice   string `json:"weightedAvgPrice"`
	PrevClosePrice     string `json:"prevClosePrice"`
	LastPrice          string `json:"lastPrice"`
	LastQty            string `json:"lastQty"`
	BidPrice           string `json:"bidPrice"`
	AskPrice           string `json:"askPrice"`
	OpenPrice          string `json:"openPrice"`
	HighPrice          string `json:"highPrice"`
	LowPrice           string `json:"lowPrice"`
	Volume             string `json:"volume"`
	QuoteVolume        string `json:"quoteVolume"`
	OpenTime           int64  `json:"openTime"`
	CloseTime          int64  `json:"closeTime"`
	FristID            int64  `json:"firstId"`
	LastID             int64  `json:"lastId"`
	Count              int64  `json:"count"`
}

// AveragePriceService show current average price for a symbol
type AveragePriceService struct {
	c      *Client
	symbol string
}

// Symbol set symbol
func (s *AveragePriceService) Symbol(symbol string) *AveragePriceService {
	s.symbol = symbol
	return s
}

// Do send Request
func (s *AveragePriceService) Do(ctx context.Context, opts ...common.RequestOption) (res *AvgPrice, err error) {
	r := common.NewGetRequestPublic("/api/v3/avgPrice")
	r.SetQuery("symbol", s.symbol)

	res = new(AvgPrice)
	if err = s.c.CallAPI(ctx, r, res, opts...); err != nil {
		return res, err
	}
	return res, nil
}

// AvgPrice define average price
type AvgPrice struct {
	Mins  int64  `json:"mins"`
	Price string `json:"price"`
}
