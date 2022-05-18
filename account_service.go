package binance

import (
	"context"
)

// GetAccountService get account info
type GetAccountService struct {
	c *Client
}

// Do send Request
func (s *GetAccountService) Do(ctx context.Context, opts ...RequestOption) (res *Account, err error) {
	r := &Request{
		Method:   "GET",
		Endpoint: "/api/v3/account",
		SecType:  SecTypeSigned,
	}
	res = new(Account)
	if err = s.c.callAPI(ctx, r, res, opts...); err != nil {
		return nil, err
	}
	return res, nil
}

// Account define account info
type Account struct {
	MakerCommission  int64     `json:"makerCommission"`
	TakerCommission  int64     `json:"takerCommission"`
	BuyerCommission  int64     `json:"buyerCommission"`
	SellerCommission int64     `json:"sellerCommission"`
	CanTrade         bool      `json:"canTrade"`
	CanWithdraw      bool      `json:"canWithdraw"`
	CanDeposit       bool      `json:"canDeposit"`
	UpdateTime       uint64    `json:"updateTime"`
	AccountType      string    `json:"accountType"`
	Balances         []Balance `json:"balances"`
	Permissions      []string  `json:"permissions"`
}

// Balance define user balance of your account
type Balance struct {
	Asset  string  `json:"asset"`
	Free   float64 `json:"free,string"`
	Locked float64 `json:"locked,string"`
}

// GetAccountSnapshotService all account orders; active, canceled, or filled
type GetAccountSnapshotService struct {
	c           *Client
	accountType string
	startTime   *int64
	endTime     *int64
	limit       *int
}

// Type set account type ("SPOT", "MARGIN", "FUTURES")
func (s *GetAccountSnapshotService) Type(accountType string) *GetAccountSnapshotService {
	s.accountType = accountType
	return s
}

// StartTime set starttime
func (s *GetAccountSnapshotService) StartTime(startTime int64) *GetAccountSnapshotService {
	s.startTime = &startTime
	return s
}

// EndTime set endtime
func (s *GetAccountSnapshotService) EndTime(endTime int64) *GetAccountSnapshotService {
	s.endTime = &endTime
	return s
}

// Limit set limit
func (s *GetAccountSnapshotService) Limit(limit int) *GetAccountSnapshotService {
	s.limit = &limit
	return s
}

// Do send Request
func (s *GetAccountSnapshotService) Do(ctx context.Context, opts ...RequestOption) (res *Snapshot, err error) {
	r := &Request{
		Method:   "GET",
		Endpoint: "/sapi/v1/accountSnapshot",
		SecType:  SecTypeSigned,
	}
	r.SetQuery("type", s.accountType)

	if s.startTime != nil {
		r.SetQuery("startTime", *s.startTime)
	}
	if s.endTime != nil {
		r.SetQuery("endTime", *s.endTime)
	}
	if s.limit != nil {
		r.SetQuery("limit", *s.limit)
	}
	res = new(Snapshot)
	if err = s.c.callAPI(ctx, r, res, opts...); err != nil {
		return &Snapshot{}, err
	}
	return res, nil
}

// Snapshot define snapshot
type Snapshot struct {
	Code     int            `json:"code"`
	Msg      string         `json:"msg"`
	Snapshot []*SnapshotVos `json:"snapshotVos"`
}

// SnapshotVos define content of a snapshot
type SnapshotVos struct {
	Data       *SnapshotData `json:"data"`
	Type       string        `json:"type"`
	UpdateTime int64         `json:"updateTime"`
}

// SnapshotData define content of a snapshot
type SnapshotData struct {
	MarginLevel         string `json:"marginLevel"`
	TotalAssetOfBtc     string `json:"totalAssetOfBtc"`
	TotalLiabilityOfBtc string `json:"totalLiabilityOfBtc"`
	TotalNetAssetOfBtc  string `json:"totalNetAssetOfBtc"`

	Balances   []*SnapshotBalances   `json:"balances"`
	UserAssets []*SnapshotUserAssets `json:"userAssets"`
	Assets     []*SnapshotAssets     `json:"assets"`
	Positions  []*SnapshotPositions  `json:"position"`
}

// SnapshotBalances define snapshot balances
type SnapshotBalances struct {
	Asset  string `json:"asset"`
	Free   string `json:"free"`
	Locked string `json:"locked"`
}

// SnapshotUserAssets define snapshot user assets
type SnapshotUserAssets struct {
	Asset    string `json:"asset"`
	Borrowed string `json:"borrowed"`
	Free     string `json:"free"`
	Interest string `json:"interest"`
	Locked   string `json:"locked"`
	NetAsset string `json:"netAsset"`
}

// SnapshotAssets define snapshot assets
type SnapshotAssets struct {
	Asset         string `json:"asset"`
	MarginBalance string `json:"marginBalance"`
	WalletBalance string `json:"walletBalance"`
}

// SnapshotPositions define snapshot positions
type SnapshotPositions struct {
	EntryPrice       string `json:"entryPrice"`
	MarkPrice        string `json:"markPrice"`
	PositionAmt      string `json:"positionAmt"`
	Symbol           string `json:"symbol"`
	UnRealizedProfit string `json:"unRealizedProfit"`
}

// APIRestrictionService query permission from binance.
type APIRestrictionService struct {
	c *Client
}

// Do send Request
func (s *APIRestrictionService) Do(ctx context.Context, opts ...RequestOption) (
	res *APIRestriction, err error,
) {
	r := &Request{
		Method:   "GET",
		Endpoint: "/sapi/v1/account/apiRestrictions",
		SecType:  SecTypeSigned,
	}
	res = new(APIRestriction)
	if err = s.c.callAPI(ctx, r, res, opts...); err != nil {
		return nil, err
	}
	return res, nil
}

type APIRestriction struct {
	IPRestrict                     bool   `json:"ipRestrict"`
	CreateTime                     uint64 `json:"createTime"`
	EnableWithdrawals              bool   `json:"enableWithdrawals"`
	EnableInternalTransfer         bool   `json:"enableInternalTransfer"`
	PermitsUniversalTransfer       bool   `json:"permitsUniversalTransfer"`
	EnableVanillaOptions           bool   `json:"enableVanillaOptions"`
	EnableReading                  bool   `json:"enableReading"`
	EnableFutures                  bool   `json:"enableFutures"`
	EnableMargin                   bool   `json:"enableMargin"`
	EnableSpotAndMarginTrading     bool   `json:"enableSpotAndMarginTrading"`
	TradingAuthorityExpirationTime uint64 `json:"tradingAuthorityExpirationTime"`
}
