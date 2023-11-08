package futures

import (
	"context"

	"github.com/crypto-zero/go-binance/v2/common"
)

// GetBalanceService get account balance
type GetBalanceService struct {
	c *Client
}

// Do send request
func (s *GetBalanceService) Do(ctx context.Context, opts ...common.RequestOption) (res []*Balance, err error) {
	r := common.NewGetRequestSigned("/fapi/v2/balance")

	res = make([]*Balance, 0)
	if err = s.c.CallAPI(ctx, r, &res, opts...); err != nil {
		return nil, err
	}
	return res, nil
}

// Balance define user balance of your account
type Balance struct {
	AccountAlias       string `json:"accountAlias"`
	Asset              string `json:"asset"`
	Balance            string `json:"balance"`
	CrossWalletBalance string `json:"crossWalletBalance"`
	CrossUnPnl         string `json:"crossUnPnl"`
	AvailableBalance   string `json:"availableBalance"`
	MaxWithdrawAmount  string `json:"maxWithdrawAmount"`
}

// GetAccountService get account info
type GetAccountService struct {
	c *Client
}

// Do send request
func (s *GetAccountService) Do(ctx context.Context, opts ...common.RequestOption) (res *Account, err error) {
	r := common.NewGetRequestSigned("/fapi/v2/account")

	res = new(Account)
	if err = s.c.CallAPI(ctx, r, res, opts...); err != nil {
		return nil, err
	}
	return res, nil
}

// Account define account info
type Account struct {
	Assets                      []*AccountAsset    `json:"assets"`
	CanDeposit                  bool               `json:"canDeposit"`
	CanTrade                    bool               `json:"canTrade"`
	CanWithdraw                 bool               `json:"canWithdraw"`
	FeeTier                     int                `json:"feeTier"`
	MultiAssetsMargin           bool               `json:"multiAssetsMargin"`
	TradeGroupID                int64              `json:"tradeGroupId"`
	MaxWithdrawAmount           string             `json:"maxWithdrawAmount"`
	Positions                   []*AccountPosition `json:"positions"`
	TotalInitialMargin          string             `json:"totalInitialMargin"`
	TotalMaintMargin            string             `json:"totalMaintMargin"`
	TotalMarginBalance          string             `json:"totalMarginBalance"`
	TotalOpenOrderInitialMargin string             `json:"totalOpenOrderInitialMargin"`
	TotalPositionInitialMargin  string             `json:"totalPositionInitialMargin"`
	TotalUnrealizedProfit       string             `json:"totalUnrealizedProfit"`
	TotalWalletBalance          string             `json:"totalWalletBalance"`
	TotalCrossWalletBalance     string             `json:"totalCrossWalletBalance"` // 全仓账户余额, 仅计算 USDT 资产
	TotalCrossUnrealizedProfit  string             `json:"totalCrossUnPnl"`         // 全仓持仓未实现盈亏总额, 仅计算 USDT 资产
	AvailableBalance            string             `json:"availableBalance"`        // 可用余额, 仅计算 USDT 资产
	UpdateTime                  int64              `json:"updateTime"`
}

// AccountAsset define account asset
type AccountAsset struct {
	Asset                  string `json:"asset"`
	InitialMargin          string `json:"initialMargin"`
	MaintMargin            string `json:"maintMargin"`
	MarginBalance          string `json:"marginBalance"`
	MaxWithdrawAmount      string `json:"maxWithdrawAmount"`
	OpenOrderInitialMargin string `json:"openOrderInitialMargin"`
	PositionInitialMargin  string `json:"positionInitialMargin"`
	UnrealizedProfit       string `json:"unrealizedProfit"`
	WalletBalance          string `json:"walletBalance"`
	CrossWalletBalance     string `json:"crossWalletBalance"` // 全仓账户余额
	CrossUnrealizedProfit  string `json:"crossUnPnl"`         // 全仓持仓未实现盈亏
	AvailableBalance       string `json:"availableBalance"`   // 可用余额
	MarginAvailable        bool   `json:"marginAvailable"`    // 是否可用作联合保证金
	UpdateTime             int64  `json:"updateTime"`         // 更新时间
}

// AccountPosition define account position
type AccountPosition struct {
	Isolated               bool             `json:"isolated"`
	Leverage               string           `json:"leverage"`
	InitialMargin          string           `json:"initialMargin"`
	MaintMargin            string           `json:"maintMargin"`
	OpenOrderInitialMargin string           `json:"openOrderInitialMargin"`
	PositionInitialMargin  string           `json:"positionInitialMargin"`
	Symbol                 string           `json:"symbol"`
	UnrealizedProfit       string           `json:"unrealizedProfit"`
	EntryPrice             string           `json:"entryPrice"`
	MaxNotional            string           `json:"maxNotional"`
	BidNotional            string           `json:"bidNotional"` // 买单净值，忽略
	AskNotional            string           `json:"askNotional"` // 卖单净值，忽略
	PositionSide           PositionSideType `json:"positionSide"`
	PositionAmt            string           `json:"positionAmt"`
	Notional               string           `json:"notional"`
	IsolatedWallet         string           `json:"isolatedWallet"`
	UpdateTime             int64            `json:"updateTime"` // 更新时间
}
