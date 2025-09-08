# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a cryptocurrency trading bot built in Go with a clean architecture. It features:
- Multiple trading strategies (Moving Average, RSI)
- Binance integration with testnet support
- Dry run mode for testing
- Portfolio management
- Interface-driven design for extensibility

## Development Commands

```bash
# Run the bot
go run ./cmd/bot

# Build the bot
go build ./cmd/bot

# Run tests
go test ./...

# Format code
go fmt ./...

# Check for issues
go vet ./...
```

## Architecture

The project follows a clean architecture with clear separation of concerns:

- **`cmd/bot/`**: Application entry point with main.go
- **`internal/bot/`**: Core bot logic and configuration management
- **`internal/exchange/`**: Exchange API abstraction layer (currently Binance)
- **`internal/strategy/`**: Trading strategies (Moving Average, RSI)
- **`internal/portfolio/`**: Portfolio and transaction management
- **`internal/market/`**: Market data handling and mock data for dry runs

### Key Interfaces

- **Exchange Interface** (`internal/exchange/exchange.go`): Defines contract for exchange implementations
- **Strategy Interface** (`internal/strategy/strategy.go`): Defines contract for trading strategies

## Configuration

The bot uses `configs/config.json` with three main sections:
- **binance**: API credentials and testnet flag
- **trading**: Symbol, balance, strategy selection, risk parameters
- **bot**: Interval, dry_run mode, logging

Configuration is validated at startup. Missing config file triggers automatic creation of default config.

## Key Components

### Bot Lifecycle
- Main entry point: `cmd/bot/main.go:main()`
- Bot creation: `internal/bot/bot.go:NewTradingBot()`
- Trading loop: `internal/bot/bot.go:processTick()`

### Strategy Selection
Strategy selection happens in `internal/bot/bot.go:NewTradingBot()` based on config:
- "rsi" → RSI strategy (14-period)
- default → Moving Average strategy (20/50 periods)

### Exchange Integration
- Exchange interface allows easy addition of new exchanges
- Binance client in `internal/exchange/binance.go`
- Mock data used in dry run mode via `internal/market/data.go`

## Extension Points

### Adding New Exchanges
1. Implement Exchange interface in `internal/exchange/`
2. Update exchange creation logic in `internal/bot/bot.go:NewTradingBot()`

### Adding New Strategies
1. Implement Strategy interface in `internal/strategy/`
2. Update strategy selection in `internal/bot/bot.go:NewTradingBot()`

## Safety Features

- Testnet enabled by default
- Dry run mode simulates trades without real money
- Configuration validation before startup
- Graceful shutdown handling (Ctrl+C)
- Error handling continues operation on API failures