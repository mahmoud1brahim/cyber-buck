package strategy

import (
	"trading-bot/internal/market"
)

type Action string

const (
	ActionBuy  Action = "BUY"
	ActionSell Action = "SELL"
	ActionHold Action = "HOLD"
)

type Signal struct {
	Action Action
	Symbol string
	Amount float64
}

type Strategy interface {
	Analyze(data *market.Data) Signal
	Name() string
}
