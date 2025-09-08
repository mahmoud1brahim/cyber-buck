package strategy

import (
	"math"
	"trading-bot/internal/market"
)

type RSIStrategy struct {
	period       int
	priceHistory []float64
	maxHistory   int
	overbought   float64
	oversold     float64
}

func NewRSIStrategy(period int) Strategy {
	return &RSIStrategy{
		period:       period,
		priceHistory: make([]float64, 0),
		maxHistory:   period + 10,
		overbought:   70.0,
		oversold:     30.0,
	}
}

func (rsi *RSIStrategy) Name() string {
	return "RSI"
}

func (rsi *RSIStrategy) Analyze(data *market.Data) Signal {
	rsi.priceHistory = append(rsi.priceHistory, data.Price)

	if len(rsi.priceHistory) > rsi.maxHistory {
		rsi.priceHistory = rsi.priceHistory[1:]
	}

	if len(rsi.priceHistory) < rsi.period+1 {
		return Signal{Action: ActionHold, Symbol: data.Symbol, Amount: 0}
	}

	rsiValue := rsi.calculateRSI()

	if rsiValue < rsi.oversold {
		return Signal{
			Action: ActionBuy,
			Symbol: "BTC",
			Amount: 500.0,
		}
	} else if rsiValue > rsi.overbought {
		return Signal{
			Action: ActionSell,
			Symbol: "BTC",
			Amount: 0.3,
		}
	}

	return Signal{Action: ActionHold, Symbol: data.Symbol, Amount: 0}
}

func (rsi *RSIStrategy) calculateRSI() float64 {
	gains := make([]float64, 0)
	losses := make([]float64, 0)

	for i := 1; i < len(rsi.priceHistory); i++ {
		change := rsi.priceHistory[i] - rsi.priceHistory[i-1]
		if change > 0 {
			gains = append(gains, change)
			losses = append(losses, 0)
		} else {
			gains = append(gains, 0)
			losses = append(losses, math.Abs(change))
		}
	}

	if len(gains) < rsi.period {
		return 50.0
	}

	avgGain := 0.0
	avgLoss := 0.0

	start := len(gains) - rsi.period
	for i := start; i < len(gains); i++ {
		avgGain += gains[i]
		avgLoss += losses[i]
	}

	avgGain /= float64(rsi.period)
	avgLoss /= float64(rsi.period)

	if avgLoss == 0 {
		return 100.0
	}

	rs := avgGain / avgLoss
	return 100.0 - (100.0 / (1.0 + rs))
}
