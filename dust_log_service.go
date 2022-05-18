/*******************************************************************************
** @Author:					Thomas Bouder <Tbouder>
** @Email:					Tbouder@protonmail.com
** @Date:					Friday 19 June 2020 - 18:51:45
** @Filename:				0_dust_service.go
**
** @Last modified by:		Tbouder
*******************************************************************************/

package binance

import (
	"context"

	"github.com/crypto-zero/go-binance/v2/common"
)

// ListDustLogService fetch small amounts of assets exchanged versus BNB
// See https://binance-docs.github.io/apidocs/spot/en/#dustlog-user_data
type ListDustLogService struct {
	c         *Client
	startTime *int64
	endTime   *int64
}

// StartTime sets the startTime parameter.
// If present, EndTime MUST be specified. The difference between EndTime - StartTime MUST be between 0-90 days.
func (s *ListDustLogService) StartTime(startTime int64) *ListDustLogService {
	s.startTime = &startTime
	return s
}

// EndTime sets the endTime parameter.
// If present, StartTime MUST be specified. The difference between EndTime - StartTime MUST be between 0-90 days.
func (s *ListDustLogService) EndTime(endTime int64) *ListDustLogService {
	s.endTime = &endTime
	return s
}

// Do sends the Request.
func (s *ListDustLogService) Do(ctx context.Context) (res *DustResult, err error) {
	r := common.NewGetRequestSigned("/sapi/v1/asset/dribblet")
	if s.startTime != nil {
		r.SetQuery("startTime", *s.startTime)
	}
	if s.endTime != nil {
		r.SetQuery("endTime", *s.endTime)
	}
	res = new(DustResult)
	if err = s.c.callAPI(ctx, r, res); err != nil {
		return
	}
	return res, nil
}

// DustResult represents the result of a DustLog API Call.
type DustResult struct {
	Total              uint8               `json:"total"` // Total counts of exchange
	UserAssetDribblets []UserAssetDribblet `json:"userAssetDribblets"`
}

// UserAssetDribblet represents one dust log row
type UserAssetDribblet struct {
	OperateTime              int64                     `json:"operateTime"`
	TotalTransferedAmount    string                    `json:"totalTransferedAmount"`    // Total transfered BNB amount for this exchange.
	TotalServiceChargeAmount string                    `json:"totalServiceChargeAmount"` // Total service charge amount for this exchange.
	TransID                  int64                     `json:"transId"`
	UserAssetDribbletDetails []UserAssetDribbletDetail `json:"userAssetDribbletDetails"` // Details of this exchange.
}

// DustLog represents one dust log informations
type UserAssetDribbletDetail struct {
	TransID             int    `json:"transId"`
	ServiceChargeAmount string `json:"serviceChargeAmount"`
	Amount              string `json:"amount"`
	OperateTime         int64  `json:"operateTime"` // The time of this exchange.
	TransferedAmount    string `json:"transferedAmount"`
	FromAsset           string `json:"fromAsset"`
}

// DustTransferService convert dust assets to BNB.
// See https://binance-docs.github.io/apidocs/spot/en/#dust-transfer-user_data
type DustTransferService struct {
	c     *Client
	asset []string
}

// Asset set asset.
func (s *DustTransferService) Asset(asset []string) *DustTransferService {
	s.asset = asset
	return s
}

// Do send the Request.
func (s *DustTransferService) Do(ctx context.Context) (withdraws *DustTransferResponse, err error) {
	r := common.NewPostRequestSigned("/sapi/v1/asset/dust")
	for _, a := range s.asset {
		r.AddQuery("asset", a)
	}
	res := new(DustTransferResponse)
	if err = s.c.callAPI(ctx, r, res); err != nil {
		return
	}
	return res, nil
}

// DustTransferResponse represents the response from DustTransferService.
type DustTransferResponse struct {
	TotalServiceCharge string                `json:"totalServiceCharge"`
	TotalTransfered    string                `json:"totalTransfered"`
	TransferResult     []*DustTransferResult `json:"transferResult"`
}

// DustTransferResult represents the result of a dust transfer.
type DustTransferResult struct {
	Amount              string `json:"amount"`
	FromAsset           string `json:"fromAsset"`
	OperateTime         int64  `json:"operateTime"`
	ServiceChargeAmount string `json:"serviceChargeAmount"`
	TranID              int64  `json:"tranId"`
	TransferedAmount    string `json:"transferedAmount"`
}
