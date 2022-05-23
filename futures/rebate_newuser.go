package futures

import (
	"context"

	"github.com/crypto-zero/go-binance/v2/common"
)

// GetRebateNewUserService
type GetRebateNewUserService struct {
	c           *Client
	brokerageID string
	type_future int
}

// BrokerageID setting
func (s *GetRebateNewUserService) BrokerageID(brokerageID string) *GetRebateNewUserService {
	s.brokerageID = brokerageID
	return s
}

// Type future setting
func (s *GetRebateNewUserService) Type(type_future int) *GetRebateNewUserService {
	s.type_future = type_future
	return s
}

// Do send request
func (s *GetRebateNewUserService) Do(ctx context.Context, opts ...common.RequestOption) (res *RebateNewUser, err error) {
	r := common.NewGetRequestSigned("/fapi/v1/apiReferral/ifNewUser")

	if s.brokerageID != "" {
		r.SetQuery("brokerId", s.brokerageID)
	}
	if s.type_future != 0 {
		r.SetQuery("type", s.type_future)
	}

	res = &RebateNewUser{}
	if err = s.c.CallAPI(ctx, r, res, opts...); err != nil {
		return nil, err
	}
	return res, nil
}

// PositionRisk define position risk info
type RebateNewUser struct {
	BrokerId      string `json:"brokerId"`
	RebateWorking bool   `json:"rebateWorking"`
	IfNewUser     bool   `json:"ifNewUser"`
}
