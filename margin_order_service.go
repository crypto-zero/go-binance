package binance

import (
	"context"

	"github.com/crypto-zero/go-binance/v2/common"
)

// CreateMarginOrderService create order
type CreateMarginOrderService struct {
	c                *Client
	symbol           string
	side             SideType
	orderType        OrderType
	quantity         *string
	quoteOrderQty    *string
	price            *string
	stopPrice        *string
	newClientOrderID *string
	icebergQuantity  *string
	newOrderRespType *NewOrderRespType
	sideEffectType   *SideEffectType
	timeInForce      *TimeInForceType
	isIsolated       *bool
}

// Symbol set symbol
func (s *CreateMarginOrderService) Symbol(symbol string) *CreateMarginOrderService {
	s.symbol = symbol
	return s
}

// IsIsolated sets the order to isolated margin
func (s *CreateMarginOrderService) IsIsolated(isIsolated bool) *CreateMarginOrderService {
	s.isIsolated = &isIsolated
	return s
}

// Side set side
func (s *CreateMarginOrderService) Side(side SideType) *CreateMarginOrderService {
	s.side = side
	return s
}

// Type set type
func (s *CreateMarginOrderService) Type(orderType OrderType) *CreateMarginOrderService {
	s.orderType = orderType
	return s
}

// TimeInForce set timeInForce
func (s *CreateMarginOrderService) TimeInForce(timeInForce TimeInForceType) *CreateMarginOrderService {
	s.timeInForce = &timeInForce
	return s
}

// Quantity set quantity
func (s *CreateMarginOrderService) Quantity(quantity string) *CreateMarginOrderService {
	s.quantity = &quantity
	return s
}

// QuoteOrderQty set quoteOrderQty
func (s *CreateMarginOrderService) QuoteOrderQty(quoteOrderQty string) *CreateMarginOrderService {
	s.quoteOrderQty = &quoteOrderQty
	return s
}

// Price set price
func (s *CreateMarginOrderService) Price(price string) *CreateMarginOrderService {
	s.price = &price
	return s
}

// NewClientOrderID set newClientOrderID
func (s *CreateMarginOrderService) NewClientOrderID(newClientOrderID string) *CreateMarginOrderService {
	s.newClientOrderID = &newClientOrderID
	return s
}

// StopPrice set stopPrice
func (s *CreateMarginOrderService) StopPrice(stopPrice string) *CreateMarginOrderService {
	s.stopPrice = &stopPrice
	return s
}

// IcebergQuantity set icebergQuantity
func (s *CreateMarginOrderService) IcebergQuantity(icebergQuantity string) *CreateMarginOrderService {
	s.icebergQuantity = &icebergQuantity
	return s
}

// NewOrderRespType set icebergQuantity
func (s *CreateMarginOrderService) NewOrderRespType(newOrderRespType NewOrderRespType) *CreateMarginOrderService {
	s.newOrderRespType = &newOrderRespType
	return s
}

// SideEffectType set sideEffectType
func (s *CreateMarginOrderService) SideEffectType(sideEffectType SideEffectType) *CreateMarginOrderService {
	s.sideEffectType = &sideEffectType
	return s
}

// Do send Request
func (s *CreateMarginOrderService) Do(ctx context.Context, opts ...common.RequestOption) (res *CreateOrderResponse, err error) {
	r := common.NewPostRequestSigned("/sapi/v1/margin/order")
	m := common.Params{
		"symbol": s.symbol,
		"side":   s.side,
		"type":   s.orderType,
	}
	if s.quantity != nil {
		m["quantity"] = *s.quantity
	}
	if s.quoteOrderQty != nil {
		m["quoteOrderQty"] = *s.quoteOrderQty
	}
	if s.isIsolated != nil {
		if *s.isIsolated {
			m["isIsolated"] = "TRUE"
		} else {
			m["isIsolated"] = "FALSE"
		}
	}
	if s.timeInForce != nil {
		m["timeInForce"] = *s.timeInForce
	}
	if s.price != nil {
		m["price"] = *s.price
	}
	if s.newClientOrderID != nil {
		m["newClientOrderId"] = *s.newClientOrderID
	}
	if s.stopPrice != nil {
		m["stopPrice"] = *s.stopPrice
	}
	if s.icebergQuantity != nil {
		m["icebergQty"] = *s.icebergQuantity
	}
	if s.newOrderRespType != nil {
		m["newOrderRespType"] = *s.newOrderRespType
	}
	if s.sideEffectType != nil {
		m["sideEffectType"] = *s.sideEffectType
	}
	r.SetFormParams(m)

	res = new(CreateOrderResponse)
	if err = s.c.CallAPI(ctx, r, res, opts...); err != nil {
		return nil, err
	}
	return res, nil
}

// CancelMarginOrderService cancel an order
type CancelMarginOrderService struct {
	c                 *Client
	symbol            string
	orderID           *int64
	origClientOrderID *string
	newClientOrderID  *string
	isIsolated        bool
}

// Symbol set symbol
func (s *CancelMarginOrderService) Symbol(symbol string) *CancelMarginOrderService {
	s.symbol = symbol
	return s
}

// IsIsolated set isIsolated
func (s *CancelMarginOrderService) IsIsolated(isIsolated bool) *CancelMarginOrderService {
	s.isIsolated = isIsolated
	return s
}

// OrderID set orderID
func (s *CancelMarginOrderService) OrderID(orderID int64) *CancelMarginOrderService {
	s.orderID = &orderID
	return s
}

// OrigClientOrderID set origClientOrderID
func (s *CancelMarginOrderService) OrigClientOrderID(origClientOrderID string) *CancelMarginOrderService {
	s.origClientOrderID = &origClientOrderID
	return s
}

// NewClientOrderID set newClientOrderID
func (s *CancelMarginOrderService) NewClientOrderID(newClientOrderID string) *CancelMarginOrderService {
	s.newClientOrderID = &newClientOrderID
	return s
}

// Do send Request
func (s *CancelMarginOrderService) Do(ctx context.Context, opts ...common.RequestOption) (res *CancelMarginOrderResponse, err error) {
	r := common.NewDeleteRequestSigned("/sapi/v1/margin/order")
	r.SetForm("symbol", s.symbol)
	if s.orderID != nil {
		r.SetForm("orderId", *s.orderID)
	}
	if s.origClientOrderID != nil {
		r.SetForm("origClientOrderId", *s.origClientOrderID)
	}
	if s.newClientOrderID != nil {
		r.SetForm("newClientOrderId", *s.newClientOrderID)
	}
	if s.isIsolated {
		r.SetForm("isIsolated", "TRUE")
	}

	res = new(CancelMarginOrderResponse)
	if err = s.c.CallAPI(ctx, r, res, opts...); err != nil {
		return nil, err
	}
	return res, nil
}

// GetMarginOrderService get an order
type GetMarginOrderService struct {
	c                 *Client
	symbol            string
	orderID           *int64
	origClientOrderID *string
	isIsolated        bool
}

// IsIsolated set isIsolated
func (s *GetMarginOrderService) IsIsolated(isIsolated bool) *GetMarginOrderService {
	s.isIsolated = isIsolated
	return s
}

// Symbol set symbol
func (s *GetMarginOrderService) Symbol(symbol string) *GetMarginOrderService {
	s.symbol = symbol
	return s
}

// OrderID set orderID
func (s *GetMarginOrderService) OrderID(orderID int64) *GetMarginOrderService {
	s.orderID = &orderID
	return s
}

// OrigClientOrderID set origClientOrderID
func (s *GetMarginOrderService) OrigClientOrderID(origClientOrderID string) *GetMarginOrderService {
	s.origClientOrderID = &origClientOrderID
	return s
}

// Do send Request
func (s *GetMarginOrderService) Do(ctx context.Context, opts ...common.RequestOption) (res *Order, err error) {
	r := common.NewGetRequestSigned("/sapi/v1/margin/order")
	r.SetQuery("symbol", s.symbol)
	if s.orderID != nil {
		r.SetQuery("orderId", *s.orderID)
	}
	if s.origClientOrderID != nil {
		r.SetQuery("origClientOrderId", *s.origClientOrderID)
	}
	if s.isIsolated {
		r.SetQuery("isIsolated", "TRUE")
	}

	res = new(Order)
	if err = s.c.CallAPI(ctx, r, res, opts...); err != nil {
		return nil, err
	}
	return res, nil
}

// ListMarginOpenOrdersService list margin open orders
type ListMarginOpenOrdersService struct {
	c          *Client
	symbol     string
	isIsolated bool
}

// Symbol set symbol
func (s *ListMarginOpenOrdersService) Symbol(symbol string) *ListMarginOpenOrdersService {
	s.symbol = symbol
	return s
}

// IsIsolated set isIsolated
func (s *ListMarginOpenOrdersService) IsIsolated(isIsolated bool) *ListMarginOpenOrdersService {
	s.isIsolated = isIsolated
	return s
}

// Do send Request
func (s *ListMarginOpenOrdersService) Do(ctx context.Context, opts ...common.RequestOption) (res []*Order, err error) {
	r := common.NewGetRequestSigned("/sapi/v1/margin/openOrders")
	if s.symbol != "" {
		r.SetQuery("symbol", s.symbol)
	}
	if s.isIsolated {
		r.SetQuery("isIsolated", "TRUE")
	}

	res = make([]*Order, 0)
	if err = s.c.CallAPI(ctx, r, &res, opts...); err != nil {
		return nil, err
	}
	return res, nil
}

// ListMarginOrdersService all account orders; active, canceled, or filled
type ListMarginOrdersService struct {
	c          *Client
	symbol     string
	orderID    *int64
	startTime  *int64
	endTime    *int64
	limit      *int
	isIsolated bool
}

// Symbol set symbol
func (s *ListMarginOrdersService) Symbol(symbol string) *ListMarginOrdersService {
	s.symbol = symbol
	return s
}

// IsIsolated set isIsolated
func (s *ListMarginOrdersService) IsIsolated(isIsolated bool) *ListMarginOrdersService {
	s.isIsolated = isIsolated
	return s
}

// OrderID set orderID
func (s *ListMarginOrdersService) OrderID(orderID int64) *ListMarginOrdersService {
	s.orderID = &orderID
	return s
}

// StartTime set starttime
func (s *ListMarginOrdersService) StartTime(startTime int64) *ListMarginOrdersService {
	s.startTime = &startTime
	return s
}

// EndTime set endtime
func (s *ListMarginOrdersService) EndTime(endTime int64) *ListMarginOrdersService {
	s.endTime = &endTime
	return s
}

// Limit set limit
func (s *ListMarginOrdersService) Limit(limit int) *ListMarginOrdersService {
	s.limit = &limit
	return s
}

// Do send Request
func (s *ListMarginOrdersService) Do(ctx context.Context, opts ...common.RequestOption) (res []*Order, err error) {
	r := common.NewGetRequestSigned("/sapi/v1/margin/allOrders")
	r.SetQuery("symbol", s.symbol)
	if s.orderID != nil {
		r.SetQuery("orderId", *s.orderID)
	}
	if s.startTime != nil {
		r.SetQuery("startTime", *s.startTime)
	}
	if s.endTime != nil {
		r.SetQuery("endTime", *s.endTime)
	}
	if s.limit != nil {
		r.SetQuery("limit", *s.limit)
	}
	if s.isIsolated {
		r.SetQuery("isIsolated", "TRUE")
	}

	res = make([]*Order, 0)
	if err := s.c.CallAPI(ctx, r, &res, opts...); err != nil {
		return nil, err
	}
	return res, nil
}

// CancelMarginOrderResponse define response of canceling order
type CancelMarginOrderResponse struct {
	Symbol                   string          `json:"symbol"`
	OrigClientOrderID        string          `json:"origClientOrderId"`
	OrderID                  string          `json:"orderId"`
	ClientOrderID            string          `json:"clientOrderId"`
	TransactTime             int64           `json:"transactTime"`
	Price                    string          `json:"price"`
	OrigQuantity             string          `json:"origQty"`
	ExecutedQuantity         string          `json:"executedQty"`
	CummulativeQuoteQuantity string          `json:"cummulativeQuoteQty"`
	Status                   OrderStatusType `json:"status"`
	TimeInForce              TimeInForceType `json:"timeInForce"`
	Type                     OrderType       `json:"type"`
	Side                     SideType        `json:"side"`
}
