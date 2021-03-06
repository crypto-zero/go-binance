package futures

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

// Do send request
func (s *ListBookTickersService) Do(ctx context.Context, opts ...common.RequestOption) (res []*BookTicker, err error) {
	r := common.NewGetRequestPublic("/fapi/v1/ticker/bookTicker")
	if s.symbol != nil {
		r.SetQuery("symbol", *s.symbol)
	}

	res = make([]*BookTicker, 0)
	f := func(data []byte) error {
		data = common.ToJSONList(data)
		err = json.Unmarshal(data, &res)
		if err != nil {
			return err
		}
		return nil
	}
	if err = s.c.CallAPI(ctx, r, f, opts...); err != nil {
		return []*BookTicker{}, err
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

// Do send request
func (s *ListPricesService) Do(ctx context.Context, opts ...common.RequestOption) (res []*SymbolPrice, err error) {
	r := common.NewGetRequestPublic("/fapi/v1/ticker/price")
	if s.symbol != nil {
		r.SetQuery("symbol", *s.symbol)
	}

	f := func(data []byte) error {
		if err = json.Unmarshal(common.ToJSONList(data), &res); err != nil {
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

// Do send request
func (s *ListPriceChangeStatsService) Do(ctx context.Context, opts ...common.RequestOption) (res []*PriceChangeStats, err error) {
	r := common.NewGetRequestPublic("/fapi/v1/ticker/24hr")
	if s.symbol != nil {
		r.SetQuery("symbol", *s.symbol)
	}

	f := func(data []byte) error {
		if err = json.Unmarshal(common.ToJSONList(data), &res); err != nil {
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
	LastQuantity       string `json:"lastQty"`
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
