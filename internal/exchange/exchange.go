package exchange

import (
	"trading-bot/internal/market"
)

type Side string
type OrderType string

const (
	SideBuy  Side = "BUY"
	SideSell Side = "SELL"
)

const (
	TypeMarket OrderType = "MARKET"
	TypeLimit  OrderType = "LIMIT"
)

type OrderResponse struct {
	Symbol        string `json:"symbol"`
	OrderID       int64  `json:"orderId"`
	ClientOrderID string `json:"clientOrderId"`
	Status        string `json:"status"`
	Type          string `json:"type"`
	Side          string `json:"side"`
	Quantity      string `json:"origQty"`
	Price         string `json:"price"`
}

type Exchange interface {
	GetMarketData(symbol string) (*market.Data, error)
	PlaceOrder(symbol string, side Side, orderType OrderType, quantity, price float64) (*OrderResponse, error)
	TestConnection() error
}
