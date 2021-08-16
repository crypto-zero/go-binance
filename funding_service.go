package binance

import (
	"context"
	"encoding/json"
)

// GetFundingAssetService fetches all asset detail.
//
// See https://binance-docs.github.io/apidocs/spot/en/#user_data-16
type GetFundingAssetService struct {
	c                *Client
	asset            *string
	needBtcValuation *bool
}

// Asset set the asset parameter.
func (s *GetFundingAssetService) Asset(asset string) *GetFundingAssetService {
	s.asset = &asset
	return s
}

// NeedBTCValuation set the needBtcValuation parameter.
func (s *GetFundingAssetService) NeedBTCValuation(needBtcValuation bool) *GetFundingAssetService {
	s.needBtcValuation = &needBtcValuation
	return s
}

// Do send the request.
func (s *GetFundingAssetService) Do(ctx context.Context) (out map[string]FundingAsset, err error) {
	r := &request{
		method:   "POST",
		endpoint: "/sapi/v1/asset/get-funding-asset",
		secType:  secTypeSigned,
	}
	if s.asset != nil {
		r.setFormParam("asset", *s.asset)
	}
	if s.needBtcValuation != nil {
		val := "true"
		if !*s.needBtcValuation {
			val = "false"
		}
		r.setFormParam("needBtcValuation", val)
	}
	data, err := s.c.callAPI(ctx, r)
	if err != nil {
		return
	}
	out = make(map[string]FundingAsset)
	var rsp []FundingAsset
	if err = json.Unmarshal(data, &rsp); err != nil {
		return
	}
	for _, x := range rsp {
		out[x.Asset] = x
	}
	return out, nil
}

// FundingAsset represents the detail of an asset
type FundingAsset struct {
	Asset        string  `json:"asset"`
	Free         float64 `json:"free,string"`
	Locked       float64 `json:"locked,string"`
	Freeze       float64 `json:"freeze,string"`
	Withdrawing  float64 `json:"withdrawing,string"`
	BTCValuation float64 `json:"btcValuation,string"`
}
