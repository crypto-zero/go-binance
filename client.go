package binance

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync/atomic"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/crypto-zero/go-binance/v2/common"
	"github.com/crypto-zero/go-binance/v2/delivery"
	"github.com/crypto-zero/go-binance/v2/futures"
)

// SideType define side type of order
type SideType string

// OrderType define order type
type OrderType string

// TimeInForceType define time in force type of order
type TimeInForceType string

// NewOrderRespType define response JSON verbosity
type NewOrderRespType string

// OrderStatusType define order status type
type OrderStatusType string

// SymbolType define symbol type
type SymbolType string

// SymbolStatusType define symbol status type
type SymbolStatusType string

// SymbolFilterType define symbol filter type
type SymbolFilterType string

// MarginTransferType define margin transfer type
type MarginTransferType int

// MarginLoanStatusType define margin loan status type
type MarginLoanStatusType string

// MarginRepayStatusType define margin repay status type
type MarginRepayStatusType string

// FuturesTransferStatusType define futures transfer status type
type FuturesTransferStatusType string

// SideEffectType define side effect type for orders
type SideEffectType string

// FuturesTransferType define futures transfer type
type FuturesTransferType int

// Endpoints
const (
	baseAPIMainURL    = "https://api.binance.com"
	baseAPITestnetURL = "https://testnet.binance.vision"
)

// UseTestnet switch all the API endpoints from production to the testnet
var UseTestnet = false

// Global enums
const (
	SideTypeBuy  SideType = "BUY"
	SideTypeSell SideType = "SELL"

	OrderTypeLimit           OrderType = "LIMIT"
	OrderTypeMarket          OrderType = "MARKET"
	OrderTypeLimitMaker      OrderType = "LIMIT_MAKER"
	OrderTypeStopLoss        OrderType = "STOP_LOSS"
	OrderTypeStopLossLimit   OrderType = "STOP_LOSS_LIMIT"
	OrderTypeTakeProfit      OrderType = "TAKE_PROFIT"
	OrderTypeTakeProfitLimit OrderType = "TAKE_PROFIT_LIMIT"

	TimeInForceTypeGTC TimeInForceType = "GTC"
	TimeInForceTypeIOC TimeInForceType = "IOC"
	TimeInForceTypeFOK TimeInForceType = "FOK"

	NewOrderRespTypeACK    NewOrderRespType = "ACK"
	NewOrderRespTypeRESULT NewOrderRespType = "RESULT"
	NewOrderRespTypeFULL   NewOrderRespType = "FULL"

	OrderStatusTypeNew             OrderStatusType = "NEW"
	OrderStatusTypePartiallyFilled OrderStatusType = "PARTIALLY_FILLED"
	OrderStatusTypeFilled          OrderStatusType = "FILLED"
	OrderStatusTypeCanceled        OrderStatusType = "CANCELED"
	OrderStatusTypePendingCancel   OrderStatusType = "PENDING_CANCEL"
	OrderStatusTypeRejected        OrderStatusType = "REJECTED"
	OrderStatusTypeExpired         OrderStatusType = "EXPIRED"

	SymbolTypeSpot SymbolType = "SPOT"

	SymbolStatusTypePreTrading   SymbolStatusType = "PRE_TRADING"
	SymbolStatusTypeTrading      SymbolStatusType = "TRADING"
	SymbolStatusTypePostTrading  SymbolStatusType = "POST_TRADING"
	SymbolStatusTypeEndOfDay     SymbolStatusType = "END_OF_DAY"
	SymbolStatusTypeHalt         SymbolStatusType = "HALT"
	SymbolStatusTypeAuctionMatch SymbolStatusType = "AUCTION_MATCH"
	SymbolStatusTypeBreak        SymbolStatusType = "BREAK"

	SymbolFilterTypeLotSize          SymbolFilterType = "LOT_SIZE"
	SymbolFilterTypePriceFilter      SymbolFilterType = "PRICE_FILTER"
	SymbolFilterTypePercentPrice     SymbolFilterType = "PERCENT_PRICE"
	SymbolFilterTypeMinNotional      SymbolFilterType = "MIN_NOTIONAL"
	SymbolFilterTypeIcebergParts     SymbolFilterType = "ICEBERG_PARTS"
	SymbolFilterTypeMarketLotSize    SymbolFilterType = "MARKET_LOT_SIZE"
	SymbolFilterTypeMaxNumAlgoOrders SymbolFilterType = "MAX_NUM_ALGO_ORDERS"

	MarginTransferTypeToMargin MarginTransferType = 1
	MarginTransferTypeToMain   MarginTransferType = 2

	FuturesTransferTypeSpotToFutures  FuturesTransferType = 1
	FuturesTransferTypeFuturesToSpot  FuturesTransferType = 2
	FuturesTransferTypeSpotToFuturesM FuturesTransferType = 3
	FuturesTransferTypeFuturesMToSpot FuturesTransferType = 4

	MarginLoanStatusTypePending   MarginLoanStatusType = "PENDING"
	MarginLoanStatusTypeConfirmed MarginLoanStatusType = "CONFIRMED"
	MarginLoanStatusTypeFailed    MarginLoanStatusType = "FAILED"

	MarginRepayStatusTypePending   MarginRepayStatusType = "PENDING"
	MarginRepayStatusTypeConfirmed MarginRepayStatusType = "CONFIRMED"
	MarginRepayStatusTypeFailed    MarginRepayStatusType = "FAILED"

	FuturesTransferStatusTypePending   FuturesTransferStatusType = "PENDING"
	FuturesTransferStatusTypeConfirmed FuturesTransferStatusType = "CONFIRMED"
	FuturesTransferStatusTypeFailed    FuturesTransferStatusType = "FAILED"

	SideEffectTypeNoSideEffect SideEffectType = "NO_SIDE_EFFECT"
	SideEffectTypeMarginBuy    SideEffectType = "MARGIN_BUY"
	SideEffectTypeAutoRepay    SideEffectType = "AUTO_REPAY"

	timestampKey  = "timestamp"
	signatureKey  = "signature"
	recvWindowKey = "recvWindow"
)

func currentTimestamp() int64 {
	return FormatTimestamp(time.Now())
}

// FormatTimestamp formats a time into Unix timestamp in milliseconds, as requested by Binance.
func FormatTimestamp(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

func newJSON(data []byte) (j *simplejson.Json, err error) {
	j, err = simplejson.NewJson(data)
	if err != nil {
		return nil, err
	}
	return j, nil
}

// getAPIEndpoint return the base endpoint of the Rest API according the testnet flag
func getAPIEndpoint(testnet bool) string {
	if testnet {
		return baseAPITestnetURL
	}
	return baseAPIMainURL
}

// NewClient initialize an API client instance with API key and secret key.
// You should always call this function before using this SDK.
// Services will be created by the form client.NewXXXService().
func NewClient(apiKey, secretKey string, testnet bool) *Client {
	return &Client{
		APIKey:     apiKey,
		SecretKey:  secretKey,
		BaseURL:    getAPIEndpoint(testnet),
		UserAgent:  "Binance/golang",
		HTTPClient: http.DefaultClient,
		Logger: common.NewDefaultLogger(common.LogInfo, log.New(os.Stderr,
			"Binance-golang ", log.LstdFlags)),
	}
}

// NewFuturesClient initialize client for futures API
func NewFuturesClient(apiKey, secretKey string, testnet bool) *futures.Client {
	return futures.NewClient(apiKey, secretKey, testnet)
}

// NewDeliveryClient initialize client for coin-M futures API
func NewDeliveryClient(apiKey, secretKey string, testnet bool) *delivery.Client {
	return delivery.NewClient(apiKey, secretKey, testnet)
}

type doFunc func(req *http.Request) (*http.Response, error)

// Client define API client
type Client struct {
	globalRequestID uint64

	APIKey     string
	SecretKey  string
	BaseURL    string
	UserAgent  string
	HTTPClient *http.Client
	Logger     common.Logger
	TimeOffset int64
	do         doFunc
}

func (c *Client) parseRequest(r *common.Request, opts ...common.RequestOption) (bodyString string, err error) {
	// set Request options from user
	for _, opt := range opts {
		opt(r)
	}
	err = r.Validate()
	if err != nil {
		return "", err
	}

	r.ID = atomic.AddUint64(&c.globalRequestID, 1)

	fullURL := fmt.Sprintf("%s%s", c.BaseURL, r.Endpoint)
	if r.RecvWindow > 0 {
		r.SetQuery(recvWindowKey, r.RecvWindow)
	}
	if r.SecType == common.SecTypeSigned {
		r.SetQuery(timestampKey, currentTimestamp()-c.TimeOffset)
	}
	queryString := r.Query.Encode()
	body := &bytes.Buffer{}
	bodyString = r.Form.Encode()
	header := http.Header{}
	if r.Header != nil {
		header = r.Header.Clone()
	}
	if bodyString != "" {
		header.Set("Content-Type", "application/x-www-form-urlencoded")
		body = bytes.NewBufferString(bodyString)
	}
	if r.SecType == common.SecTypeAPIKey || r.SecType == common.SecTypeSigned {
		header.Set("X-MBX-APIKEY", c.APIKey)
	}

	if r.SecType == common.SecTypeSigned {
		raw := fmt.Sprintf("%s%s", queryString, bodyString)
		mac := hmac.New(sha256.New, []byte(c.SecretKey))
		_, err = mac.Write([]byte(raw))
		if err != nil {
			return "", err
		}
		v := url.Values{}
		v.Set(signatureKey, fmt.Sprintf("%x", mac.Sum(nil)))
		if queryString == "" {
			queryString = v.Encode()
		} else {
			queryString = fmt.Sprintf("%s&%s", queryString, v.Encode())
		}
	}
	if queryString != "" {
		fullURL = fmt.Sprintf("%s?%s", fullURL, queryString)
	}

	r.FullURL = fullURL
	r.Header = header
	r.Body = body
	return bodyString, nil
}

func (c *Client) callAPI(ctx context.Context, r *common.Request, result interface{},
	opts ...common.RequestOption,
) (err error) {
	bodyString, err := c.parseRequest(r, opts...)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(r.Method, r.FullURL, r.Body)
	if err != nil {
		return err
	}

	req = req.WithContext(ctx)
	req.Header = r.Header

	c.Logger.Debugw("call api prepare", "id", r.ID, "url", r.FullURL, "body", bodyString)

	f := c.do
	if f == nil {
		f = c.HTTPClient.Do
	}

	res, err := f(req)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	defer func() {
		// Only overwrite the returned error if the original error was nil and an
		// error occurred while closing the body.
		if cerr := res.Body.Close(); err == nil && cerr != nil {
			err = cerr
		}
	}()

	c.Logger.Debugw("call api reply", "id", r.ID, "status_code", res.StatusCode,
		"response_headers", res.Header, "response_body", string(data))

	if res.StatusCode >= 400 {
		apiErr := &common.APIError{Status: res.StatusCode}
		if e := json.Unmarshal(data, apiErr); e != nil {
			c.Logger.Debugw("call api parse error failed", "id", r.ID, "err", e)
		}
		return apiErr
	}

	if result != nil {
		f, ok := result.(func(data []byte) error)
		if ok {
			if err = f(data); err != nil {
				return err
			}
			return nil
		}
		if err = json.Unmarshal(data, result); err != nil {
			return err
		}
	}
	return nil
}

// NewPingService init ping service
func (c *Client) NewPingService() *PingService {
	return &PingService{c: c}
}

// NewServerTimeService init server time service
func (c *Client) NewServerTimeService() *ServerTimeService {
	return &ServerTimeService{c: c}
}

// NewSetServerTimeService init set server time service
func (c *Client) NewSetServerTimeService() *SetServerTimeService {
	return &SetServerTimeService{c: c}
}

// NewDepthService init depth service
func (c *Client) NewDepthService() *DepthService {
	return &DepthService{c: c}
}

// NewAggTradesService init aggregate trades service
func (c *Client) NewAggTradesService() *AggTradesService {
	return &AggTradesService{c: c}
}

// NewRecentTradesService init recent trades service
func (c *Client) NewRecentTradesService() *RecentTradesService {
	return &RecentTradesService{c: c}
}

// NewKlinesService init klines service
func (c *Client) NewKlinesService() *KlinesService {
	return &KlinesService{c: c}
}

// NewListPriceChangeStatsService init list prices change stats service
func (c *Client) NewListPriceChangeStatsService() *ListPriceChangeStatsService {
	return &ListPriceChangeStatsService{c: c}
}

// NewListPricesService init listing prices service
func (c *Client) NewListPricesService() *ListPricesService {
	return &ListPricesService{c: c}
}

// NewListBookTickersService init listing booking tickers service
func (c *Client) NewListBookTickersService() *ListBookTickersService {
	return &ListBookTickersService{c: c}
}

// NewCreateOrderService init creating order service
func (c *Client) NewCreateOrderService() *CreateOrderService {
	return &CreateOrderService{c: c}
}

// NewCreateOCOService init creating OCO service
func (c *Client) NewCreateOCOService() *CreateOCOService {
	return &CreateOCOService{c: c}
}

// NewCancelOCOService init cancel OCO service
func (c *Client) NewCancelOCOService() *CancelOCOService {
	return &CancelOCOService{c: c}
}

// NewGetOrderService init get order service
func (c *Client) NewGetOrderService() *GetOrderService {
	return &GetOrderService{c: c}
}

// NewCancelOrderService init cancel order service
func (c *Client) NewCancelOrderService() *CancelOrderService {
	return &CancelOrderService{c: c}
}

// NewCancelOpenOrdersService init cancel open orders service
func (c *Client) NewCancelOpenOrdersService() *CancelOpenOrdersService {
	return &CancelOpenOrdersService{c: c}
}

// NewListOpenOrdersService init list open orders service
func (c *Client) NewListOpenOrdersService() *ListOpenOrdersService {
	return &ListOpenOrdersService{c: c}
}

// NewListOrdersService init listing orders service
func (c *Client) NewListOrdersService() *ListOrdersService {
	return &ListOrdersService{c: c}
}

// NewGetAccountService init getting account service
func (c *Client) NewGetAccountService() *GetAccountService {
	return &GetAccountService{c: c}
}

// NewGetAccountSnapshotService init getting account snapshot service
func (c *Client) NewGetAccountSnapshotService() *GetAccountSnapshotService {
	return &GetAccountSnapshotService{c: c}
}

// NewListTradesService init listing trades service
func (c *Client) NewListTradesService() *ListTradesService {
	return &ListTradesService{c: c}
}

// NewHistoricalTradesService init listing trades service
func (c *Client) NewHistoricalTradesService() *HistoricalTradesService {
	return &HistoricalTradesService{c: c}
}

// NewListDepositsService init listing deposits service
func (c *Client) NewListDepositsService() *ListDepositsService {
	return &ListDepositsService{c: c}
}

// NewGetDepositAddressService init getting deposit address service
func (c *Client) NewGetDepositAddressService() *GetDepositsAddressService {
	return &GetDepositsAddressService{c: c}
}

// NewCreateWithdrawService init creating withdraw service
func (c *Client) NewCreateWithdrawService() *CreateWithdrawService {
	return &CreateWithdrawService{c: c}
}

// NewListWithdrawsService init listing withdraw service
func (c *Client) NewListWithdrawsService() *ListWithdrawsService {
	return &ListWithdrawsService{c: c}
}

// NewStartUserStreamService init starting user stream service
func (c *Client) NewStartUserStreamService() *StartUserStreamService {
	return &StartUserStreamService{c: c}
}

// NewKeepaliveUserStreamService init keep alive user stream service
func (c *Client) NewKeepaliveUserStreamService() *KeepaliveUserStreamService {
	return &KeepaliveUserStreamService{c: c}
}

// NewCloseUserStreamService init closing user stream service
func (c *Client) NewCloseUserStreamService() *CloseUserStreamService {
	return &CloseUserStreamService{c: c}
}

// NewExchangeInfoService init exchange info service
func (c *Client) NewExchangeInfoService() *ExchangeInfoService {
	return &ExchangeInfoService{c: c}
}

// NewGetAssetDetailService init get asset detail service
func (c *Client) NewGetAssetDetailService() *GetAssetDetailService {
	return &GetAssetDetailService{c: c}
}

// NewGetFundingAssetService init get asset detail service
func (c *Client) NewGetFundingAssetService() *GetFundingAssetService {
	return &GetFundingAssetService{c: c}
}

// NewAveragePriceService init average price service
func (c *Client) NewAveragePriceService() *AveragePriceService {
	return &AveragePriceService{c: c}
}

// NewMarginTransferService init margin account transfer service
func (c *Client) NewMarginTransferService() *MarginTransferService {
	return &MarginTransferService{c: c}
}

// NewMarginLoanService init margin account loan service
func (c *Client) NewMarginLoanService() *MarginLoanService {
	return &MarginLoanService{c: c}
}

// NewMarginRepayService init margin account repay service
func (c *Client) NewMarginRepayService() *MarginRepayService {
	return &MarginRepayService{c: c}
}

// NewCreateMarginOrderService init creating margin order service
func (c *Client) NewCreateMarginOrderService() *CreateMarginOrderService {
	return &CreateMarginOrderService{c: c}
}

// NewCancelMarginOrderService init cancel order service
func (c *Client) NewCancelMarginOrderService() *CancelMarginOrderService {
	return &CancelMarginOrderService{c: c}
}

// NewGetMarginOrderService init get order service
func (c *Client) NewGetMarginOrderService() *GetMarginOrderService {
	return &GetMarginOrderService{c: c}
}

// NewListMarginLoansService init list margin loan service
func (c *Client) NewListMarginLoansService() *ListMarginLoansService {
	return &ListMarginLoansService{c: c}
}

// NewListMarginRepaysService init list margin repay service
func (c *Client) NewListMarginRepaysService() *ListMarginRepaysService {
	return &ListMarginRepaysService{c: c}
}

// NewGetMarginAccountService init get margin account service
func (c *Client) NewGetMarginAccountService() *GetMarginAccountService {
	return &GetMarginAccountService{c: c}
}

// NewGetIsolatedMarginAccountService init get isolated margin asset service
func (c *Client) NewGetIsolatedMarginAccountService() *GetIsolatedMarginAccountService {
	return &GetIsolatedMarginAccountService{c: c}
}

// NewGetMarginAssetService init get margin asset service
func (c *Client) NewGetMarginAssetService() *GetMarginAssetService {
	return &GetMarginAssetService{c: c}
}

// NewGetMarginPairService init get margin pair service
func (c *Client) NewGetMarginPairService() *GetMarginPairService {
	return &GetMarginPairService{c: c}
}

// NewGetMarginAllPairsService init get margin all pairs service
func (c *Client) NewGetMarginAllPairsService() *GetMarginAllPairsService {
	return &GetMarginAllPairsService{c: c}
}

// NewGetMarginPriceIndexService init get margin price index service
func (c *Client) NewGetMarginPriceIndexService() *GetMarginPriceIndexService {
	return &GetMarginPriceIndexService{c: c}
}

// NewListMarginOpenOrdersService init list margin open orders service
func (c *Client) NewListMarginOpenOrdersService() *ListMarginOpenOrdersService {
	return &ListMarginOpenOrdersService{c: c}
}

// NewListMarginOrdersService init list margin all orders service
func (c *Client) NewListMarginOrdersService() *ListMarginOrdersService {
	return &ListMarginOrdersService{c: c}
}

// NewListMarginTradesService init list margin trades service
func (c *Client) NewListMarginTradesService() *ListMarginTradesService {
	return &ListMarginTradesService{c: c}
}

// NewGetMaxBorrowableService init get max borrowable service
func (c *Client) NewGetMaxBorrowableService() *GetMaxBorrowableService {
	return &GetMaxBorrowableService{c: c}
}

// NewGetMaxTransferableService init get max transferable service
func (c *Client) NewGetMaxTransferableService() *GetMaxTransferableService {
	return &GetMaxTransferableService{c: c}
}

// NewStartMarginUserStreamService init starting margin user stream service
func (c *Client) NewStartMarginUserStreamService() *StartMarginUserStreamService {
	return &StartMarginUserStreamService{c: c}
}

// NewKeepaliveMarginUserStreamService init keep alive margin user stream service
func (c *Client) NewKeepaliveMarginUserStreamService() *KeepaliveMarginUserStreamService {
	return &KeepaliveMarginUserStreamService{c: c}
}

// NewCloseMarginUserStreamService init closing margin user stream service
func (c *Client) NewCloseMarginUserStreamService() *CloseMarginUserStreamService {
	return &CloseMarginUserStreamService{c: c}
}

// NewStartIsolatedMarginUserStreamService init starting margin user stream service
func (c *Client) NewStartIsolatedMarginUserStreamService() *StartIsolatedMarginUserStreamService {
	return &StartIsolatedMarginUserStreamService{c: c}
}

// NewKeepaliveIsolatedMarginUserStreamService init keep alive margin user stream service
func (c *Client) NewKeepaliveIsolatedMarginUserStreamService() *KeepaliveIsolatedMarginUserStreamService {
	return &KeepaliveIsolatedMarginUserStreamService{c: c}
}

// NewCloseIsolatedMarginUserStreamService init closing margin user stream service
func (c *Client) NewCloseIsolatedMarginUserStreamService() *CloseIsolatedMarginUserStreamService {
	return &CloseIsolatedMarginUserStreamService{c: c}
}

// NewFuturesTransferService init futures transfer service
func (c *Client) NewFuturesTransferService() *FuturesTransferService {
	return &FuturesTransferService{c: c}
}

// NewListFuturesTransferService init list futures transfer service
func (c *Client) NewListFuturesTransferService() *ListFuturesTransferService {
	return &ListFuturesTransferService{c: c}
}

// NewListDustLogService init list dust log service
func (c *Client) NewListDustLogService() *ListDustLogService {
	return &ListDustLogService{c: c}
}

// NewDustTransferService init dust transfer service
func (c *Client) NewDustTransferService() *DustTransferService {
	return &DustTransferService{c: c}
}

// NewAPIRestrictionService init api restriction service
func (c *Client) NewAPIRestrictionService() *APIRestrictionService {
	return &APIRestrictionService{c: c}
}
