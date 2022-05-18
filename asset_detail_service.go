package binance

import (
	"context"

	"github.com/crypto-zero/go-binance/v2/common"
)

// GetAssetDetailService fetches all asset detail.
//
// See https://binance-docs.github.io/apidocs/spot/en/#asset-detail-user_data
type GetAssetDetailService struct {
	c     *Client
	asset *string
}

// Asset sets the asset parameter.
func (s *GetAssetDetailService) Asset(asset string) *GetAssetDetailService {
	s.asset = &asset
	return s
}

// Do sends the Request.
func (s *GetAssetDetailService) Do(ctx context.Context) (res map[string]AssetDetail, err error) {
	r := &common.Request{
		Method:   "GET",
		Endpoint: "/sapi/v1/asset/assetDetail",
		SecType:  common.SecTypeSigned,
	}
	if s.asset != nil {
		r.SetQuery("asset", *s.asset)
	}

	res = make(map[string]AssetDetail)
	if err = s.c.callAPI(ctx, r, &res); err != nil {
		return
	}
	return res, nil
}

// AssetDetail represents the detail of an asset
type AssetDetail struct {
	MinWithdrawAmount float64 `json:"minWithdrawAmount"`
	DepositStatus     bool    `json:"depositStatus"`
	WithdrawFee       float64 `json:"withdrawFee"`
	WithdrawStatus    bool    `json:"withdrawStatus"`
	DepositTip        string  `json:"depositTip"`
}
