package binance

import (
	"context"
	"strings"

	"github.com/crypto-zero/go-binance/v2/common"
)

// MarginTransferService transfer between spot account and margin account
type MarginTransferService struct {
	c            *Client
	asset        string
	amount       string
	transferType int
}

// Asset set asset being transferred, e.g., BTC
func (s *MarginTransferService) Asset(asset string) *MarginTransferService {
	s.asset = asset
	return s
}

// Amount the amount to be transferred
func (s *MarginTransferService) Amount(amount string) *MarginTransferService {
	s.amount = amount
	return s
}

// Type 1: transfer from main account to margin account 2: transfer from margin account to main account
func (s *MarginTransferService) Type(transferType MarginTransferType) *MarginTransferService {
	s.transferType = int(transferType)
	return s
}

// Do send Request
func (s *MarginTransferService) Do(ctx context.Context, opts ...common.RequestOption) (res *TransactionResponse, err error) {
	r := &common.Request{
		Method:   "POST",
		Endpoint: "/sapi/v1/margin/transfer",
		SecType:  common.SecTypeSigned,
	}
	m := common.Params{
		"asset":  s.asset,
		"amount": s.amount,
		"type":   s.transferType,
	}
	r.SetFormParams(m)

	res = new(TransactionResponse)
	if err = s.c.callAPI(ctx, r, res, opts...); err != nil {
		return nil, err
	}
	return res, nil
}

// TransactionResponse define transaction response
type TransactionResponse struct {
	TranID int64 `json:"tranId"`
}

// MarginLoanService apply for a loan
type MarginLoanService struct {
	c              *Client
	asset          string
	amount         string
	isolatedSymbol string
}

// Asset set asset being transferred, e.g., BTC
func (s *MarginLoanService) Asset(asset string) *MarginLoanService {
	s.asset = asset
	return s
}

// Amount the amount to be transferred
func (s *MarginLoanService) Amount(amount string) *MarginLoanService {
	s.amount = amount
	return s
}

// IsolatedSymbol set IsolatedSymbol
func (s *MarginLoanService) IsolatedSymbol(isolatedSymbol string) *MarginLoanService {
	s.isolatedSymbol = isolatedSymbol
	return s
}

// Do send Request
func (s *MarginLoanService) Do(ctx context.Context, opts ...common.RequestOption) (res *TransactionResponse, err error) {
	r := &common.Request{
		Method:   "POST",
		Endpoint: "/sapi/v1/margin/loan",
		SecType:  common.SecTypeSigned,
	}
	m := common.Params{
		"asset":  s.asset,
		"amount": s.amount,
	}
	r.SetFormParams(m)
	if s.isolatedSymbol != "" {
		r.SetQuery("isolatedSymbol", s.isolatedSymbol)
	}

	res = new(TransactionResponse)
	if err = s.c.callAPI(ctx, r, res, opts...); err != nil {
		return nil, err
	}
	return res, nil
}

// MarginRepayService repay loan for margin account
type MarginRepayService struct {
	c              *Client
	asset          string
	amount         string
	isolatedSymbol string
}

// Asset set asset being transferred, e.g., BTC
func (s *MarginRepayService) Asset(asset string) *MarginRepayService {
	s.asset = asset
	return s
}

// Amount the amount to be transferred
func (s *MarginRepayService) Amount(amount string) *MarginRepayService {
	s.amount = amount
	return s
}

// IsolatedSymbol set IsolatedSymbol
func (s *MarginRepayService) IsolatedSymbol(isolatedSymbol string) *MarginRepayService {
	s.isolatedSymbol = isolatedSymbol
	return s
}

// Do send Request
func (s *MarginRepayService) Do(ctx context.Context, opts ...common.RequestOption) (res *TransactionResponse, err error) {
	r := &common.Request{
		Method:   "POST",
		Endpoint: "/sapi/v1/margin/repay",
		SecType:  common.SecTypeSigned,
	}
	m := common.Params{
		"asset":  s.asset,
		"amount": s.amount,
	}
	r.SetFormParams(m)
	if s.isolatedSymbol != "" {
		r.SetQuery("isolatedSymbol", s.isolatedSymbol)
	}

	res = new(TransactionResponse)
	if err = s.c.callAPI(ctx, r, res, opts...); err != nil {
		return nil, err
	}
	return res, nil
}

// ListMarginLoansService list loan record
type ListMarginLoansService struct {
	c         *Client
	asset     string
	txID      *int64
	startTime *int64
	endTime   *int64
	current   *int64
	size      *int64
}

// Asset set asset
func (s *ListMarginLoansService) Asset(asset string) *ListMarginLoansService {
	s.asset = asset
	return s
}

// TxID set transaction id
func (s *ListMarginLoansService) TxID(txID int64) *ListMarginLoansService {
	s.txID = &txID
	return s
}

// StartTime set start time
func (s *ListMarginLoansService) StartTime(startTime int64) *ListMarginLoansService {
	s.startTime = &startTime
	return s
}

// EndTime set end time
func (s *ListMarginLoansService) EndTime(endTime int64) *ListMarginLoansService {
	s.endTime = &endTime
	return s
}

// Current currently querying page. Start from 1. Default:1
func (s *ListMarginLoansService) Current(current int64) *ListMarginLoansService {
	s.current = &current
	return s
}

// Size default:10 max:100
func (s *ListMarginLoansService) Size(size int64) *ListMarginLoansService {
	s.size = &size
	return s
}

// Do send Request
func (s *ListMarginLoansService) Do(ctx context.Context, opts ...common.RequestOption) (res *MarginLoanResponse, err error) {
	r := &common.Request{
		Method:   "GET",
		Endpoint: "/sapi/v1/margin/loan",
		SecType:  common.SecTypeSigned,
	}
	r.SetQuery("asset", s.asset)
	if s.txID != nil {
		r.SetQuery("txId", *s.txID)
	}
	if s.startTime != nil {
		r.SetQuery("startTime", *s.startTime)
	}
	if s.endTime != nil {
		r.SetQuery("endTime", *s.endTime)
	}
	if s.current != nil {
		r.SetQuery("current", *s.current)
	}
	if s.size != nil {
		r.SetQuery("size", *s.size)
	}

	res = new(MarginLoanResponse)
	if err = s.c.callAPI(ctx, r, res, opts...); err != nil {
		return nil, err
	}
	return res, nil
}

// MarginLoanResponse define margin loan response
type MarginLoanResponse struct {
	Rows  []MarginLoan `json:"rows"`
	Total int64        `json:"total"`
}

// MarginLoan define margin loan
type MarginLoan struct {
	Asset     string               `json:"asset"`
	Principal string               `json:"principal"`
	Timestamp int64                `json:"timestamp"`
	Status    MarginLoanStatusType `json:"status"`
}

// ListMarginRepaysService list repay record
type ListMarginRepaysService struct {
	c         *Client
	asset     string
	txID      *int64
	startTime *int64
	endTime   *int64
	current   *int64
	size      *int64
}

// Asset set asset
func (s *ListMarginRepaysService) Asset(asset string) *ListMarginRepaysService {
	s.asset = asset
	return s
}

// TxID set transaction id
func (s *ListMarginRepaysService) TxID(txID int64) *ListMarginRepaysService {
	s.txID = &txID
	return s
}

// StartTime set start time
func (s *ListMarginRepaysService) StartTime(startTime int64) *ListMarginRepaysService {
	s.startTime = &startTime
	return s
}

// EndTime set end time
func (s *ListMarginRepaysService) EndTime(endTime int64) *ListMarginRepaysService {
	s.endTime = &endTime
	return s
}

// Current currently querying page. Start from 1. Default:1
func (s *ListMarginRepaysService) Current(current int64) *ListMarginRepaysService {
	s.current = &current
	return s
}

// Size default:10 max:100
func (s *ListMarginRepaysService) Size(size int64) *ListMarginRepaysService {
	s.size = &size
	return s
}

// Do send Request
func (s *ListMarginRepaysService) Do(ctx context.Context, opts ...common.RequestOption) (res *MarginRepayResponse, err error) {
	r := &common.Request{
		Method:   "GET",
		Endpoint: "/sapi/v1/margin/repay",
		SecType:  common.SecTypeSigned,
	}
	r.SetQuery("asset", s.asset)
	if s.txID != nil {
		r.SetQuery("txId", *s.txID)
	}
	if s.startTime != nil {
		r.SetQuery("startTime", *s.startTime)
	}
	if s.endTime != nil {
		r.SetQuery("endTime", *s.endTime)
	}
	if s.current != nil {
		r.SetQuery("current", *s.current)
	}
	if s.size != nil {
		r.SetQuery("size", *s.size)
	}

	res = new(MarginRepayResponse)
	if err = s.c.callAPI(ctx, r, res, opts...); err != nil {
		return nil, err
	}
	return res, nil
}

// MarginRepayResponse define margin repay response
type MarginRepayResponse struct {
	Rows  []MarginRepay `json:"rows"`
	Total int64         `json:"total"`
}

// MarginRepay define margin repay
type MarginRepay struct {
	Asset     string                `json:"asset"`
	Amount    string                `json:"amount"`
	Interest  string                `json:"interest"`
	Principal string                `json:"principal"`
	Timestamp int64                 `json:"timestamp"`
	Status    MarginRepayStatusType `json:"status"`
	TxID      int64                 `json:"txId"`
}

// GetIsolatedMarginAccountService gets isolated margin account info
type GetIsolatedMarginAccountService struct {
	c *Client

	symbols []string
}

// Symbols set symbols to the isolated margin account
func (s *GetIsolatedMarginAccountService) Symbols(symbols ...string) *GetIsolatedMarginAccountService {
	s.symbols = symbols
	return s
}

// Do send Request
func (s *GetIsolatedMarginAccountService) Do(ctx context.Context, opts ...common.RequestOption) (res *IsolatedMarginAccount, err error) {
	r := &common.Request{
		Method:   "GET",
		Endpoint: "/sapi/v1/margin/isolated/account",
		SecType:  common.SecTypeSigned,
	}

	if len(s.symbols) > 0 {
		r.SetQuery("symbols", strings.Join(s.symbols, ","))
	}

	res = new(IsolatedMarginAccount)
	if err = s.c.callAPI(ctx, r, res, opts...); err != nil {
		return nil, err
	}
	return res, nil
}

// IsolatedMarginAccount defines isolated user assets of margin account
type IsolatedMarginAccount struct {
	TotalAssetOfBTC     string                `json:"totalAssetOfBtc"`
	TotalLiabilityOfBTC string                `json:"totalLiabilityOfBtc"`
	TotalNetAssetOfBTC  string                `json:"totalNetAssetOfBtc"`
	Assets              []IsolatedMarginAsset `json:"assets"`
}

// IsolatedMarginAsset defines isolated margin asset information, like margin level, liquidation price... etc
type IsolatedMarginAsset struct {
	Symbol     string            `json:"symbol"`
	QuoteAsset IsolatedUserAsset `json:"quoteAsset"`
	BaseAsset  IsolatedUserAsset `json:"baseAsset"`

	IsolatedCreated   bool   `json:"isolatedCreated"`
	MarginLevel       string `json:"marginLevel"`
	MarginLevelStatus string `json:"marginLevelStatus"`
	MarginRatio       string `json:"marginRatio"`
	IndexPrice        string `json:"indexPrice"`
	LiquidatePrice    string `json:"liquidatePrice"`
	LiquidateRate     string `json:"liquidateRate"`
	TradeEnabled      bool   `json:"tradeEnabled"`
}

// IsolatedUserAsset defines isolated user assets of the margin account
type IsolatedUserAsset struct {
	Asset         string `json:"asset"`
	Borrowed      string `json:"borrowed"`
	Free          string `json:"free"`
	Interest      string `json:"interest"`
	Locked        string `json:"locked"`
	NetAsset      string `json:"netAsset"`
	NetAssetOfBtc string `json:"netAssetOfBtc"`

	BorrowEnabled bool   `json:"borrowEnabled"`
	RepayEnabled  bool   `json:"repayEnabled"`
	TotalAsset    string `json:"totalAsset"`
}

// GetMarginAccountService get margin account info
type GetMarginAccountService struct {
	c *Client
}

// Do send Request
func (s *GetMarginAccountService) Do(ctx context.Context, opts ...common.RequestOption) (res *MarginAccount, err error) {
	r := &common.Request{
		Method:   "GET",
		Endpoint: "/sapi/v1/margin/account",
		SecType:  common.SecTypeSigned,
	}

	res = new(MarginAccount)
	if err = s.c.callAPI(ctx, r, res, opts...); err != nil {
		return nil, err
	}
	return res, nil
}

// MarginAccount define margin account info
type MarginAccount struct {
	BorrowEnabled       bool        `json:"borrowEnabled"`
	MarginLevel         string      `json:"marginLevel"`
	TotalAssetOfBTC     string      `json:"totalAssetOfBtc"`
	TotalLiabilityOfBTC string      `json:"totalLiabilityOfBtc"`
	TotalNetAssetOfBTC  string      `json:"totalNetAssetOfBtc"`
	TradeEnabled        bool        `json:"tradeEnabled"`
	TransferEnabled     bool        `json:"transferEnabled"`
	UserAssets          []UserAsset `json:"userAssets"`
}

// UserAsset define user assets of margin account
type UserAsset struct {
	Asset    string `json:"asset"`
	Borrowed string `json:"borrowed"`
	Free     string `json:"free"`
	Interest string `json:"interest"`
	Locked   string `json:"locked"`
	NetAsset string `json:"netAsset"`
}

// GetMarginAssetService get margin asset info
type GetMarginAssetService struct {
	c     *Client
	asset string
}

// Asset set asset
func (s *GetMarginAssetService) Asset(asset string) *GetMarginAssetService {
	s.asset = asset
	return s
}

// Do send Request
func (s *GetMarginAssetService) Do(ctx context.Context, opts ...common.RequestOption) (res *MarginAsset, err error) {
	r := &common.Request{
		Method:   "GET",
		Endpoint: "/sapi/v1/margin/asset",
		SecType:  common.SecTypeAPIKey,
	}
	r.SetQuery("asset", s.asset)

	res = new(MarginAsset)
	if err = s.c.callAPI(ctx, r, res, opts...); err != nil {
		return nil, err
	}
	return res, nil
}

// MarginAsset define margin asset info
type MarginAsset struct {
	FullName      string `json:"assetFullName"`
	Name          string `json:"assetName"`
	Borrowable    bool   `json:"isBorrowable"`
	Mortgageable  bool   `json:"isMortgageable"`
	UserMinBorrow string `json:"userMinBorrow"`
	UserMinRepay  string `json:"userMinRepay"`
}

// GetMarginPairService get margin pair info
type GetMarginPairService struct {
	c      *Client
	symbol string
}

// Symbol set symbol
func (s *GetMarginPairService) Symbol(symbol string) *GetMarginPairService {
	s.symbol = symbol
	return s
}

// Do send Request
func (s *GetMarginPairService) Do(ctx context.Context, opts ...common.RequestOption) (res *MarginPair, err error) {
	r := &common.Request{
		Method:   "GET",
		Endpoint: "/sapi/v1/margin/pair",
		SecType:  common.SecTypeAPIKey,
	}
	r.SetQuery("symbol", s.symbol)

	res = new(MarginPair)
	if err = s.c.callAPI(ctx, r, res, opts...); err != nil {
		return nil, err
	}
	return res, nil
}

// MarginPair define margin pair info
type MarginPair struct {
	ID            int64  `json:"id"`
	Symbol        string `json:"symbol"`
	Base          string `json:"base"`
	Quote         string `json:"quote"`
	IsMarginTrade bool   `json:"isMarginTrade"`
	IsBuyAllowed  bool   `json:"isBuyAllowed"`
	IsSellAllowed bool   `json:"isSellAllowed"`
}

// GetMarginAllPairsService get margin pair info
type GetMarginAllPairsService struct {
	c *Client
}

// Do send Request
func (s *GetMarginAllPairsService) Do(ctx context.Context, opts ...common.RequestOption) (res []*MarginAllPair, err error) {
	r := &common.Request{
		Method:   "GET",
		Endpoint: "/sapi/v1/margin/allPairs",
		SecType:  common.SecTypeAPIKey,
	}

	res = make([]*MarginAllPair, 0)
	if err = s.c.callAPI(ctx, r, &res, opts...); err != nil {
		return nil, err
	}
	return res, nil
}

// MarginAllPair define margin pair info
type MarginAllPair struct {
	ID            int64  `json:"id"`
	Symbol        string `json:"symbol"`
	Base          string `json:"base"`
	Quote         string `json:"quote"`
	IsMarginTrade bool   `json:"isMarginTrade"`
	IsBuyAllowed  bool   `json:"isBuyAllowed"`
	IsSellAllowed bool   `json:"isSellAllowed"`
}

// GetMarginPriceIndexService get margin price index
type GetMarginPriceIndexService struct {
	c      *Client
	symbol string
}

// Symbol set symbol
func (s *GetMarginPriceIndexService) Symbol(symbol string) *GetMarginPriceIndexService {
	s.symbol = symbol
	return s
}

// Do send Request
func (s *GetMarginPriceIndexService) Do(ctx context.Context, opts ...common.RequestOption) (res *MarginPriceIndex, err error) {
	r := &common.Request{
		Method:   "GET",
		Endpoint: "/sapi/v1/margin/priceIndex",
		SecType:  common.SecTypeAPIKey,
	}
	r.SetQuery("symbol", s.symbol)

	res = new(MarginPriceIndex)
	if err = s.c.callAPI(ctx, r, res, opts...); err != nil {
		return nil, err
	}
	return res, nil
}

// MarginPriceIndex define margin price index
type MarginPriceIndex struct {
	CalcTime int64  `json:"calcTime"`
	Price    string `json:"price"`
	Symbol   string `json:"symbol"`
}

// ListMarginTradesService list trades
type ListMarginTradesService struct {
	c          *Client
	symbol     string
	startTime  *int64
	endTime    *int64
	limit      *int
	fromID     *int64
	isIsolated bool
}

// Symbol set symbol
func (s *ListMarginTradesService) Symbol(symbol string) *ListMarginTradesService {
	s.symbol = symbol
	return s
}

// IsIsolated set isIsolated
func (s *ListMarginTradesService) IsIsolated(isIsolated bool) *ListMarginTradesService {
	s.isIsolated = isIsolated
	return s
}

// StartTime set starttime
func (s *ListMarginTradesService) StartTime(startTime int64) *ListMarginTradesService {
	s.startTime = &startTime
	return s
}

// EndTime set endtime
func (s *ListMarginTradesService) EndTime(endTime int64) *ListMarginTradesService {
	s.endTime = &endTime
	return s
}

// Limit set limit
func (s *ListMarginTradesService) Limit(limit int) *ListMarginTradesService {
	s.limit = &limit
	return s
}

// FromID set fromID
func (s *ListMarginTradesService) FromID(fromID int64) *ListMarginTradesService {
	s.fromID = &fromID
	return s
}

// Do send Request
func (s *ListMarginTradesService) Do(ctx context.Context, opts ...common.RequestOption) (res []*TradeV3, err error) {
	r := &common.Request{
		Method:   "GET",
		Endpoint: "/sapi/v1/margin/myTrades",
		SecType:  common.SecTypeSigned,
	}
	r.SetQuery("symbol", s.symbol)
	if s.limit != nil {
		r.SetQuery("limit", *s.limit)
	}
	if s.startTime != nil {
		r.SetQuery("startTime", *s.startTime)
	}
	if s.endTime != nil {
		r.SetQuery("endTime", *s.endTime)
	}
	if s.fromID != nil {
		r.SetQuery("fromId", *s.fromID)
	}
	if s.isIsolated {
		r.SetQuery("isIsolated", "TRUE")
	}

	res = make([]*TradeV3, 0)
	if err = s.c.callAPI(ctx, r, &res, opts...); err != nil {
		return nil, err
	}
	return res, nil
}

// GetMaxBorrowableService get max borrowable of asset
type GetMaxBorrowableService struct {
	c     *Client
	asset string
}

// Asset set asset
func (s *GetMaxBorrowableService) Asset(asset string) *GetMaxBorrowableService {
	s.asset = asset
	return s
}

// Do send Request
func (s *GetMaxBorrowableService) Do(ctx context.Context, opts ...common.RequestOption) (res *MaxBorrowable, err error) {
	r := &common.Request{
		Method:   "GET",
		Endpoint: "/sapi/v1/margin/maxBorrowable",
		SecType:  common.SecTypeSigned,
	}
	r.SetQuery("asset", s.asset)

	res = new(MaxBorrowable)
	if err = s.c.callAPI(ctx, r, &res, opts...); err != nil {
		return nil, err
	}
	return res, nil
}

// MaxBorrowable define max borrowable response
type MaxBorrowable struct {
	Amount string `json:"amount"`
}

// GetMaxTransferableService get max transferable of asset
type GetMaxTransferableService struct {
	c     *Client
	asset string
}

// Asset set asset
func (s *GetMaxTransferableService) Asset(asset string) *GetMaxTransferableService {
	s.asset = asset
	return s
}

// Do send Request
func (s *GetMaxTransferableService) Do(ctx context.Context, opts ...common.RequestOption) (res *MaxTransferable, err error) {
	r := &common.Request{
		Method:   "GET",
		Endpoint: "/sapi/v1/margin/maxTransferable",
		SecType:  common.SecTypeSigned,
	}
	r.SetQuery("asset", s.asset)

	res = new(MaxTransferable)
	if err = s.c.callAPI(ctx, r, res, opts...); err != nil {
		return nil, err
	}
	return res, nil
}

// MaxTransferable define max transferable response
type MaxTransferable struct {
	Amount string `json:"amount"`
}

// StartIsolatedMarginUserStreamService create listen key for margin user stream service
type StartIsolatedMarginUserStreamService struct {
	c      *Client
	symbol string
}

// Symbol sets the user stream to isolated margin user stream
func (s *StartIsolatedMarginUserStreamService) Symbol(symbol string) *StartIsolatedMarginUserStreamService {
	s.symbol = symbol
	return s
}

// Do send Request
func (s *StartIsolatedMarginUserStreamService) Do(ctx context.Context, opts ...common.RequestOption) (listenKey string, err error) {
	r := &common.Request{
		Method:   "POST",
		Endpoint: "/sapi/v1/userDataStream/isolated",
		SecType:  common.SecTypeAPIKey,
	}

	r.SetForm("symbol", s.symbol)

	f := func(data []byte) error {
		j, err := newJSON(data)
		if err != nil {
			return err
		}
		listenKey = j.Get("listenKey").MustString()
		return nil
	}
	if err = s.c.callAPI(ctx, r, f, opts...); err != nil {
		return "", err
	}
	return listenKey, nil
}

// KeepaliveIsolatedMarginUserStreamService updates listen key for isolated margin user data stream
type KeepaliveIsolatedMarginUserStreamService struct {
	c         *Client
	listenKey string
	symbol    string
}

// Symbol set symbol to the isolated margin keepalive Request
func (s *KeepaliveIsolatedMarginUserStreamService) Symbol(symbol string) *KeepaliveIsolatedMarginUserStreamService {
	s.symbol = symbol
	return s
}

// ListenKey set listen key
func (s *KeepaliveIsolatedMarginUserStreamService) ListenKey(listenKey string) *KeepaliveIsolatedMarginUserStreamService {
	s.listenKey = listenKey
	return s
}

// Do send Request
func (s *KeepaliveIsolatedMarginUserStreamService) Do(ctx context.Context, opts ...common.RequestOption) (err error) {
	r := &common.Request{
		Method:   "PUT",
		Endpoint: "/sapi/v1/userDataStream/isolated",
		SecType:  common.SecTypeAPIKey,
	}
	r.SetForm("listenKey", s.listenKey)
	r.SetForm("symbol", s.symbol)

	return s.c.callAPI(ctx, r, nil, opts...)
}

// CloseIsolatedMarginUserStreamService delete listen key
type CloseIsolatedMarginUserStreamService struct {
	c         *Client
	listenKey string

	symbol string
}

// ListenKey set listen key
func (s *CloseIsolatedMarginUserStreamService) ListenKey(listenKey string) *CloseIsolatedMarginUserStreamService {
	s.listenKey = listenKey
	return s
}

// Symbol set symbol to the isolated margin user stream close Request
func (s *CloseIsolatedMarginUserStreamService) Symbol(symbol string) *CloseIsolatedMarginUserStreamService {
	s.symbol = symbol
	return s
}

// Do send Request
func (s *CloseIsolatedMarginUserStreamService) Do(ctx context.Context, opts ...common.RequestOption) (err error) {
	r := &common.Request{
		Method:   "DELETE",
		Endpoint: "/sapi/v1/userDataStream/isolated",
		SecType:  common.SecTypeAPIKey,
	}

	r.SetForm("listenKey", s.listenKey)
	r.SetForm("symbol", s.symbol)

	return s.c.callAPI(ctx, r, nil, opts...)
}

// StartMarginUserStreamService create listen key for margin user stream service
type StartMarginUserStreamService struct {
	c *Client
}

// Do send Request
func (s *StartMarginUserStreamService) Do(ctx context.Context, opts ...common.RequestOption) (listenKey string, err error) {
	r := &common.Request{
		Method:   "POST",
		Endpoint: "/sapi/v1/userDataStream",
		SecType:  common.SecTypeAPIKey,
	}

	f := func(data []byte) error {
		j, err := newJSON(data)
		if err != nil {
			return err
		}
		listenKey = j.Get("listenKey").MustString()
		return nil
	}
	if err = s.c.callAPI(ctx, r, f, opts...); err != nil {
		return "", err
	}
	return listenKey, nil
}

// KeepaliveMarginUserStreamService update listen key
type KeepaliveMarginUserStreamService struct {
	c         *Client
	listenKey string
}

// ListenKey set listen key
func (s *KeepaliveMarginUserStreamService) ListenKey(listenKey string) *KeepaliveMarginUserStreamService {
	s.listenKey = listenKey
	return s
}

// Do send Request
func (s *KeepaliveMarginUserStreamService) Do(ctx context.Context, opts ...common.RequestOption) (err error) {
	r := &common.Request{
		Method:   "PUT",
		Endpoint: "/sapi/v1/userDataStream",
		SecType:  common.SecTypeAPIKey,
	}
	r.SetForm("listenKey", s.listenKey)
	return s.c.callAPI(ctx, r, nil, opts...)
}

// CloseMarginUserStreamService delete listen key
type CloseMarginUserStreamService struct {
	c         *Client
	listenKey string
}

// ListenKey set listen key
func (s *CloseMarginUserStreamService) ListenKey(listenKey string) *CloseMarginUserStreamService {
	s.listenKey = listenKey
	return s
}

// Do send Request
func (s *CloseMarginUserStreamService) Do(ctx context.Context, opts ...common.RequestOption) (err error) {
	r := &common.Request{
		Method:   "DELETE",
		Endpoint: "/sapi/v1/userDataStream",
		SecType:  common.SecTypeAPIKey,
	}

	r.SetForm("listenKey", s.listenKey)
	return s.c.callAPI(ctx, r, nil, opts...)
}
