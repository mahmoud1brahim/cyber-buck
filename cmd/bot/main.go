package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"trading-bot/internal/bot"
)

func main() {
	configFile := "configs/config.json"

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		log.Printf("Config file not found, creating default config: %s", configFile)

		if err := os.MkdirAll("configs", 0755); err != nil {
			log.Fatalf("Failed to create configs directory: %v", err)
		}

		if err := bot.CreateDefaultConfig(configFile); err != nil {
			log.Fatalf("Failed to create default config: %v", err)
		}
		log.Printf("Please edit %s with your Binance API credentials", configFile)
		return
	}

	config, err := bot.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := config.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	tradingBot, err := bot.NewTradingBot(config)
	if err != nil {
		log.Fatalf("Failed to create trading bot: %v", err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Shutdown signal received...")
		tradingBot.Stop()
	}()

	fmt.Printf("Starting Trading Bot for %s...\n", config.Trading.Symbol)
	fmt.Printf("Strategy: %s\n", config.Trading.Strategy)
	fmt.Printf("Dry Run: %v\n", config.Bot.DryRun)

	if err := tradingBot.Start(); err != nil {
		log.Fatalf("Trading bot error: %v", err)
	}
}
