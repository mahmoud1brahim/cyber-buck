package strategy

import (
	"trading-bot/internal/market"
)

type MovingAverageStrategy struct {
	shortPeriod  int
	longPeriod   int
	priceHistory []float64
	maxHistory   int
}

func NewMovingAverageStrategy(shortPeriod, longPeriod int) Strategy {
	return &MovingAverageStrategy{
		shortPeriod:  shortPeriod,
		longPeriod:   longPeriod,
		priceHistory: make([]float64, 0),
		maxHistory:   longPeriod + 10,
	}
}

func (mas *MovingAverageStrategy) Name() string {
	return "MovingAverage"
}

func (mas *MovingAverageStrategy) Analyze(data *market.Data) Signal {
	mas.priceHistory = append(mas.priceHistory, data.Price)

	if len(mas.priceHistory) > mas.maxHistory {
		mas.priceHistory = mas.priceHistory[1:]
	}

	if len(mas.priceHistory) < mas.longPeriod {
		return Signal{Action: ActionHold, Symbol: data.Symbol, Amount: 0}
	}

	shortMA := mas.calculateMA(mas.shortPeriod)
	longMA := mas.calculateMA(mas.longPeriod)

	if shortMA > longMA {
		return Signal{
			Action: ActionBuy,
			Symbol: "BTC",
			Amount: 1000.0,
		}
	} else if shortMA < longMA {
		return Signal{
			Action: ActionSell,
			Symbol: "BTC",
			Amount: 0.5,
		}
	}

	return Signal{Action: ActionHold, Symbol: data.Symbol, Amount: 0}
}

func (mas *MovingAverageStrategy) calculateMA(period int) float64 {
	if len(mas.priceHistory) < period {
		return 0
	}

	sum := 0.0
	start := len(mas.priceHistory) - period
	for i := start; i < len(mas.priceHistory); i++ {
		sum += mas.priceHistory[i]
	}

	return sum / float64(period)
}
