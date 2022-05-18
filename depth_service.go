package binance

import (
	"context"

	"github.com/crypto-zero/go-binance/v2/common"
)

// DepthService show depth info
type DepthService struct {
	c      *Client
	symbol string
	limit  *int
}

// Symbol set symbol
func (s *DepthService) Symbol(symbol string) *DepthService {
	s.symbol = symbol
	return s
}

// Limit set limit
func (s *DepthService) Limit(limit int) *DepthService {
	s.limit = &limit
	return s
}

// Do send Request
func (s *DepthService) Do(ctx context.Context, opts ...common.RequestOption) (res *DepthResponse, err error) {
	r := &common.Request{
		Method:   "GET",
		Endpoint: "/api/v3/depth",
	}
	r.SetQuery("symbol", s.symbol)
	if s.limit != nil {
		r.SetQuery("limit", *s.limit)
	}

	res = new(DepthResponse)
	f := func(data []byte) error {
		j, err := newJSON(data)
		if err != nil {
			return err
		}
		res.LastUpdateID = j.Get("lastUpdateId").MustInt64()
		bidsLen := len(j.Get("bids").MustArray())
		res.Bids = make([]Bid, bidsLen)
		for i := 0; i < bidsLen; i++ {
			item := j.Get("bids").GetIndex(i)
			res.Bids[i] = Bid{
				Price:    item.GetIndex(0).MustString(),
				Quantity: item.GetIndex(1).MustString(),
			}
		}
		asksLen := len(j.Get("asks").MustArray())
		res.Asks = make([]Ask, asksLen)
		for i := 0; i < asksLen; i++ {
			item := j.Get("asks").GetIndex(i)
			res.Asks[i] = Ask{
				Price:    item.GetIndex(0).MustString(),
				Quantity: item.GetIndex(1).MustString(),
			}
		}
		return nil
	}
	if err := s.c.callAPI(ctx, r, f, opts...); err != nil {
		return nil, err
	}
	return res, nil
}

// DepthResponse define depth info with bids and asks
type DepthResponse struct {
	LastUpdateID int64 `json:"lastUpdateId"`
	Bids         []Bid `json:"bids"`
	Asks         []Ask `json:"asks"`
}

// Ask is a type alias for PriceLevel.
type Ask = common.PriceLevel

// Bid is a type alias for PriceLevel.
type Bid = common.PriceLevel
