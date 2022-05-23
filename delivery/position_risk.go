package delivery

import (
	"context"
	"encoding/json"

	"github.com/crypto-zero/go-binance/v2/common"
)

// GetPositionRiskService get account balance
type GetPositionRiskService struct {
	c           *Client
	pair        *string
	marginAsset *string
}

// MarginAsset set margin asset
func (s *GetPositionRiskService) MarginAsset(marginAsset string) *GetPositionRiskService {
	s.marginAsset = &marginAsset
	return s
}

// Pair set pair
func (s *GetPositionRiskService) Pair(pair string) *GetPositionRiskService {
	s.pair = &pair
	return s
}

// Do send request
func (s *GetPositionRiskService) Do(ctx context.Context, opts ...common.RequestOption) (res []*PositionRisk, err error) {
	r := common.NewGetRequestSigned("/dapi/v1/positionRisk")
	if s.marginAsset != nil {
		r.SetQuery("marginAsset", *s.marginAsset)
	}
	if s.pair != nil {
		r.SetQuery("pair", *s.pair)
	}
	data, err := s.c.CallAPIBytes(ctx, r, opts...)
	if err != nil {
		return []*PositionRisk{}, err
	}
	res = make([]*PositionRisk, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return []*PositionRisk{}, err
	}
	return res, nil
}

// PositionRisk define position risk info
type PositionRisk struct {
	Symbol           string `json:"symbol"`
	PositionAmt      string `json:"positionAmt"`
	EntryPrice       string `json:"entryPrice"`
	MarkPrice        string `json:"markPrice"`
	UnRealizedProfit string `json:"unRealizedProfit"`
	LiquidationPrice string `json:"liquidationPrice"`
	Leverage         string `json:"leverage"`
	MaxQuantity      string `json:"maxQty"`
	MarginType       string `json:"marginType"`
	IsolatedMargin   string `json:"isolatedMargin"`
	IsAutoAddMargin  string `json:"isAutoAddMargin"`
	PositionSide     string `json:"positionSide"`
}
