package bot

import (
	"fmt"
	"log"
	"time"

	"trading-bot/internal/exchange"
	"trading-bot/internal/market"
	"trading-bot/internal/portfolio"
	"trading-bot/internal/strategy"
)

type TradingBot struct {
	strategy  strategy.Strategy
	portfolio *portfolio.Portfolio
	exchange  exchange.Exchange
	config    *Config
	running   bool
}

func NewTradingBot(config *Config) (*TradingBot, error) {
	var strat strategy.Strategy
	switch config.Trading.Strategy {
	case "rsi":
		strat = strategy.NewRSIStrategy(14)
	default:
		strat = strategy.NewMovingAverageStrategy(20, 50)
	}

	exch, err := exchange.NewBinanceClient(
		config.Binance.APIKey,
		config.Binance.SecretKey,
		config.Binance.TestNet,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create exchange client: %w", err)
	}

	return &TradingBot{
		strategy:  strat,
		portfolio: portfolio.NewPortfolio(config.Trading.InitialBalance),
		exchange:  exch,
		config:    config,
		running:   false,
	}, nil
}

func (bot *TradingBot) Start() error {
	bot.running = true
	log.Printf("Trading bot started (DryRun: %v)", bot.config.Bot.DryRun)

	if !bot.config.Bot.DryRun {
		if err := bot.exchange.TestConnection(); err != nil {
			log.Printf("Failed to connect to exchange: %v", err)
			return err
		}
		log.Println("Connected to exchange API")
	}

	ticker := time.NewTicker(time.Duration(bot.config.Bot.IntervalSeconds) * time.Second)
	defer ticker.Stop()

	for bot.running {
		select {
		case <-ticker.C:
			if err := bot.processTick(); err != nil {
				log.Printf("Error processing tick: %v", err)
			}
		}
	}
	return nil
}

func (bot *TradingBot) processTick() error {
	var marketData *market.Data
	var err error

	if bot.config.Bot.DryRun {
		marketData, err = market.FetchMockData(bot.config.Trading.Symbol)
	} else {
		marketData, err = bot.exchange.GetMarketData(bot.config.Trading.Symbol)
	}

	if err != nil {
		return fmt.Errorf("error fetching market data: %w", err)
	}

	signal := bot.strategy.Analyze(marketData)

	switch signal.Action {
	case strategy.ActionBuy:
		if bot.portfolio.GetBalance() >= signal.Amount {
			if bot.config.Bot.DryRun {
				bot.portfolio.Buy(signal.Symbol, signal.Amount, marketData.Price)
			} else {
				quantity := signal.Amount / marketData.Price
				_, err := bot.exchange.PlaceOrder(bot.config.Trading.Symbol, exchange.SideBuy, exchange.TypeMarket, quantity, 0)
				if err != nil {
					log.Printf("Failed to place BUY order: %v", err)
				} else {
					bot.portfolio.Buy(signal.Symbol, signal.Amount, marketData.Price)
				}
			}
		}
	case strategy.ActionSell:
		position := bot.portfolio.GetPosition(signal.Symbol)
		if position >= signal.Amount {
			if bot.config.Bot.DryRun {
				bot.portfolio.Sell(signal.Symbol, signal.Amount, marketData.Price)
			} else {
				_, err := bot.exchange.PlaceOrder(bot.config.Trading.Symbol, exchange.SideSell, exchange.TypeMarket, signal.Amount, 0)
				if err != nil {
					log.Printf("Failed to place SELL order: %v", err)
				} else {
					bot.portfolio.Sell(signal.Symbol, signal.Amount, marketData.Price)
				}
			}
		}
	}

	currentPrices := map[string]float64{
		signal.Symbol: marketData.Price,
	}

	if len(bot.portfolio.GetHistory())%10 == 0 {
		bot.portfolio.PrintSummary(currentPrices)
	}

	return nil
}

func (bot *TradingBot) Stop() {
	bot.running = false
	log.Println("Trading bot stopped")
}
