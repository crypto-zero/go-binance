package binance

import (
	"context"
)

// CreateWithdrawService submits a withdraw Request.
//
// See https://binance-docs.github.io/apidocs/spot/en/#withdraw
type CreateWithdrawService struct {
	c                  *Client
	coin               string
	withdrawOrderID    *string
	network            *string
	address            string
	addressTag         *string
	amount             string
	transactionFeeFlag *bool
	name               *string
}

// Coin sets the coin parameter (MANDATORY).
func (s *CreateWithdrawService) Coin(v string) *CreateWithdrawService {
	s.coin = v
	return s
}

// WithdrawOrderID sets the withdrawOrderID parameter.
func (s *CreateWithdrawService) WithdrawOrderID(v string) *CreateWithdrawService {
	s.withdrawOrderID = &v
	return s
}

// Network sets the network parameter.
func (s *CreateWithdrawService) Network(v string) *CreateWithdrawService {
	s.network = &v
	return s
}

// Address sets the address parameter (MANDATORY).
func (s *CreateWithdrawService) Address(v string) *CreateWithdrawService {
	s.address = v
	return s
}

// AddressTag sets the addressTag parameter.
func (s *CreateWithdrawService) AddressTag(v string) *CreateWithdrawService {
	s.addressTag = &v
	return s
}

// Amount sets the amount parameter (MANDATORY).
func (s *CreateWithdrawService) Amount(v string) *CreateWithdrawService {
	s.amount = v
	return s
}

// TransactionFeeFlag sets the transactionFeeFlag parameter.
func (s *CreateWithdrawService) TransactionFeeFlag(v bool) *CreateWithdrawService {
	s.transactionFeeFlag = &v
	return s
}

// Name sets the name parameter.
func (s *CreateWithdrawService) Name(v string) *CreateWithdrawService {
	s.name = &v
	return s
}

// Do sends the Request.
func (s *CreateWithdrawService) Do(ctx context.Context) (res *CreateWithdrawResponse, err error) {
	r := &Request{
		Method:   "POST",
		Endpoint: "/sapi/v1/capital/withdraw/apply",
		SecType:  SecTypeSigned,
	}
	r.SetQuery("coin", s.coin)
	r.SetQuery("address", s.address)
	r.SetQuery("amount", s.amount)
	if v := s.withdrawOrderID; v != nil {
		r.SetQuery("withdrawOrderId", *v)
	}
	if v := s.network; v != nil {
		r.SetQuery("network", *v)
	}
	if v := s.addressTag; v != nil {
		r.SetQuery("addressTag", *v)
	}
	if v := s.transactionFeeFlag; v != nil {
		r.SetQuery("transactionFeeFlag", *v)
	}
	if v := s.name; v != nil {
		r.SetQuery("name", *v)
	}

	res = &CreateWithdrawResponse{}
	if err = s.c.callAPI(ctx, r, res); err != nil {
		return nil, err
	}
	return res, nil
}

// CreateWithdrawResponse represents a response from CreateWithdrawService.
type CreateWithdrawResponse struct {
	ID string `json:"id"`
}

// ListWithdrawsService fetches withdraw history.
//
// See https://binance-docs.github.io/apidocs/spot/en/#withdraw-history-supporting-network-user_data
type ListWithdrawsService struct {
	c         *Client
	coin      *string
	status    *int
	startTime *int64
	endTime   *int64
	offset    *int
	limit     *int
}

// Coin sets the coin parameter.
func (s *ListWithdrawsService) Coin(coin string) *ListWithdrawsService {
	s.coin = &coin
	return s
}

// Status sets the status parameter.
func (s *ListWithdrawsService) Status(status int) *ListWithdrawsService {
	s.status = &status
	return s
}

// StartTime sets the startTime parameter.
// If present, EndTime MUST be specified. The difference between EndTime - StartTime MUST be between 0-90 days.
func (s *ListWithdrawsService) StartTime(startTime int64) *ListWithdrawsService {
	s.startTime = &startTime
	return s
}

// EndTime sets the endTime parameter.
// If present, StartTime MUST be specified. The difference between EndTime - StartTime MUST be between 0-90 days.
func (s *ListWithdrawsService) EndTime(endTime int64) *ListWithdrawsService {
	s.endTime = &endTime
	return s
}

// Offset set offset
func (s *ListWithdrawsService) Offset(offset int) *ListWithdrawsService {
	s.offset = &offset
	return s
}

// Limit set limit
func (s *ListWithdrawsService) Limit(limit int) *ListWithdrawsService {
	s.limit = &limit
	return s
}

// Do sends the Request.
func (s *ListWithdrawsService) Do(ctx context.Context) (res []*Withdraw, err error) {
	r := &Request{
		Method:   "GET",
		Endpoint: "/sapi/v1/capital/withdraw/history",
		SecType:  SecTypeSigned,
	}
	if s.coin != nil {
		r.SetQuery("coin", *s.coin)
	}
	if s.status != nil {
		r.SetQuery("status", *s.status)
	}
	if s.startTime != nil {
		r.SetQuery("startTime", *s.startTime)
	}
	if s.endTime != nil {
		r.SetQuery("endTime", *s.endTime)
	}
	if s.offset != nil {
		r.SetQuery("offset", *s.offset)
	}
	if s.limit != nil {
		r.SetQuery("limit", *s.limit)
	}

	res = make([]*Withdraw, 0)
	if err = s.c.callAPI(ctx, r, &res); err != nil {
		return
	}
	return res, nil
}

// Withdraw represents a single withdraw entry.
type Withdraw struct {
	Address         string `json:"address"`
	Amount          string `json:"amount"`
	ApplyTime       string `json:"applyTime"`
	Coin            string `json:"coin"`
	ID              string `json:"id"`
	WithdrawOrderID string `json:"withdrawOrderID"`
	Network         string `json:"network"`
	TransferType    int    `json:"transferType"`
	Status          int    `json:"status"`
	TransactionFee  string `json:"transactionFee"`
	TxID            string `json:"txId"`
}
