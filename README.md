# Trading Bot

A cryptocurrency trading bot built in Go with Binance integration and clean architecture.

## Features

- **Multiple Trading Strategies**: Moving Average and RSI strategies
- **Binance Integration**: Real-time price data and order execution
- **Dry Run Mode**: Test strategies without real trades
- **Portfolio Management**: Track balance and positions
- **Configurable**: JSON-based configuration
- **Clean Architecture**: Well-organized, maintainable code structure

## Project Structure

```
trading-bot/
├── cmd/bot/                    # Application entry point
│   └── main.go
├── internal/                   # Private application code
│   ├── bot/                    # Core bot logic
│   │   ├── bot.go
│   │   └── config.go
│   ├── exchange/               # Exchange interfaces and implementations
│   │   ├── exchange.go
│   │   └── binance.go
│   ├── strategy/               # Trading strategies
│   │   ├── strategy.go
│   │   ├── moving_average.go
│   │   └── rsi.go
│   ├── portfolio/              # Portfolio management
│   │   └── portfolio.go
│   └── market/                 # Market data handling
│       └── data.go
├── configs/                    # Configuration files
│   └── config.json
├── go.mod                      # Go module definition
└── README.md                   # This file
```

## Quick Start

1. **Build and run the bot**:
   ```bash
   go run ./cmd/bot
   ```

2. **Configure API keys** (edit `configs/config.json`):
   ```json
   {
     "binance": {
       "api_key": "your_actual_api_key",
       "secret_key": "your_actual_secret_key",
       "testnet": true
     }
   }
   ```

3. **Set dry_run to false** for live trading:
   ```json
   {
     "bot": {
       "dry_run": false
     }
   }
   ```

## Configuration

The bot uses `configs/config.json` for configuration:

- **binance**: API credentials and testnet settings
- **trading**: Symbol, balance, strategy, and risk parameters  
- **bot**: Interval, dry run mode, and logging

## Strategies

### Moving Average Strategy
- **File**: `internal/strategy/moving_average.go`
- Uses short (20) and long (50) period moving averages
- Buy when short MA crosses above long MA
- Sell when short MA crosses below long MA

### RSI Strategy
- **File**: `internal/strategy/rsi.go`
- Uses 14-period RSI indicator
- Buy when RSI < 30 (oversold)
- Sell when RSI > 70 (overbought)

## Architecture Benefits

### Clean Separation of Concerns
- **`internal/bot/`**: Core trading logic and configuration
- **`internal/exchange/`**: Exchange API abstraction (easy to add new exchanges)
- **`internal/strategy/`**: Trading strategies (easy to add new strategies)
- **`internal/portfolio/`**: Portfolio and transaction management
- **`internal/market/`**: Market data fetching and processing

### Interface-Driven Design
- **Exchange Interface**: Easy to add Coinbase, Kraken, etc.
- **Strategy Interface**: Easy to add ML strategies, custom indicators
- **Testable**: All components can be easily mocked for testing

### Scalability
- Add new exchanges by implementing the `Exchange` interface
- Add new strategies by implementing the `Strategy` interface
- Clean module boundaries make the code maintainable

## Safety Features

- **Testnet by default**: Prevents accidental live trading
- **Dry run mode**: Simulates trades without real money
- **Input validation**: Validates configuration before starting
- **Error handling**: Continues operation on API errors
- **Graceful shutdown**: Handles Ctrl+C properly

## Building

```bash
# Build the application
go build ./cmd/bot

# Run tests (when added)
go test ./...

# Format code
go fmt ./...

# Check for issues
go vet ./...
```

## Adding New Features

### Adding a New Exchange
1. Create `internal/exchange/newexchange.go`
2. Implement the `Exchange` interface
3. Update bot creation logic in `internal/bot/bot.go`

### Adding a New Strategy
1. Create `internal/strategy/newstrategy.go`
2. Implement the `Strategy` interface
3. Update strategy selection in `internal/bot/bot.go`

This architecture makes the trading bot highly maintainable and extensible!