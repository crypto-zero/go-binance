package delivery

import (
	"context"
	"encoding/json"

	"github.com/crypto-zero/go-binance/v2/common"
)

// GetBalanceService get account balance
type GetBalanceService struct {
	c *Client
}

// Do send request
func (s *GetBalanceService) Do(ctx context.Context, opts ...common.RequestOption) (res []*Balance, err error) {
	r := common.NewGetRequestSigned("/dapi/v1/balance")
	data, err := s.c.CallAPIBytes(ctx, r, opts...)
	if err != nil {
		return []*Balance{}, err
	}
	res = make([]*Balance, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return []*Balance{}, err
	}
	return res, nil
}

// Balance define user balance of your account
type Balance struct {
	AccountAlias       string `json:"accountAlias"`
	Asset              string `json:"asset"`
	Balance            string `json:"balance"`
	WithdrawAvailable  string `json:"withdrawAvailable"`
	CrossWalletBalance string `json:"crossWalletBalance"`
	CrossUnPnl         string `json:"crossUnPnl"`
	AvailableBalance   string `json:"availableBalance"`
	UpdateTime         int64  `json:"updateTime"`
}

// GetAccountService get account info
type GetAccountService struct {
	c *Client
}

// Do send request
func (s *GetAccountService) Do(ctx context.Context, opts ...common.RequestOption) (res *Account, err error) {
	r := common.NewGetRequestSigned("/dapi/v1/account")
	data, err := s.c.CallAPIBytes(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(Account)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Account define account info
type Account struct {
	Assets      []*AccountAsset    `json:"assets"`
	CanDeposit  bool               `json:"canDeposit"`
	CanTrade    bool               `json:"canTrade"`
	CanWithdraw bool               `json:"canWithdraw"`
	FeeTier     int                `json:"feeTier"`
	Positions   []*AccountPosition `json:"positions"`
	UpdateTime  int64              `json:"updateType"`
}

// AccountAsset define account asset
type AccountAsset struct {
	Asset                  string `json:"asset"`
	WalletBalance          string `json:"walletBalance"`
	UnrealizedProfit       string `json:"unrealizedProfit"`
	MarginBalance          string `json:"marginBalance"`
	MaintMargin            string `json:"maintMargin"`
	InitialMargin          string `json:"initialMargin"`
	PositionInitialMargin  string `json:"positionInitialMargin"`
	OpenOrderInitialMargin string `json:"openOrderInitialMargin"`
	MaxWithdrawAmount      string `json:"maxWithdrawAmount"`
	CrossWalletBalance     string `json:"crossWalletBalance"`
	CrossUnPnl             string `json:"crossUnPnl"`
	AvailableBalance       string `json:"availableBalance"`
}

// AccountPosition define accoutn position
type AccountPosition struct {
	Symbol                 string `json:"symbol"`
	InitialMargin          string `json:"initialMargin"`
	MaintMargin            string `json:"maintMargin"`
	UnrealizedProfit       string `json:"unrealizedProfit"`
	PositionInitialMargin  string `json:"positionInitialMargin"`
	OpenOrderInitialMargin string `json:"openOrderInitialMargin"`
	Leverage               string `json:"leverage"`
	Isolated               bool   `json:"isolated"`
	PositionSide           string `json:"positionSide"`
	EntryPrice             string `json:"entryPrice"`
	MaxQty                 string `json:"maxQty"`
}
