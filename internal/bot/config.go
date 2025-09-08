package bot

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Binance struct {
		APIKey    string `json:"api_key"`
		SecretKey string `json:"secret_key"`
		TestNet   bool   `json:"testnet"`
	} `json:"binance"`

	Trading struct {
		Symbol         string  `json:"symbol"`
		InitialBalance float64 `json:"initial_balance"`
		Strategy       string  `json:"strategy"`
		MaxRisk        float64 `json:"max_risk"`
		StopLoss       float64 `json:"stop_loss"`
	} `json:"trading"`

	Bot struct {
		IntervalSeconds int    `json:"interval_seconds"`
		DryRun          bool   `json:"dry_run"`
		LogLevel        string `json:"log_level"`
	} `json:"bot"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	return &config, nil
}

func CreateDefaultConfig(filename string) error {
	defaultConfig := Config{}

	defaultConfig.Binance.APIKey = "your_binance_api_key"
	defaultConfig.Binance.SecretKey = "your_binance_secret_key"
	defaultConfig.Binance.TestNet = true

	defaultConfig.Trading.Symbol = "BTCUSDT"
	defaultConfig.Trading.InitialBalance = 10000.0
	defaultConfig.Trading.Strategy = "moving_average"
	defaultConfig.Trading.MaxRisk = 0.02
	defaultConfig.Trading.StopLoss = 0.05

	defaultConfig.Bot.IntervalSeconds = 10
	defaultConfig.Bot.DryRun = true
	defaultConfig.Bot.LogLevel = "info"

	data, err := json.MarshalIndent(defaultConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling default config: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	return nil
}

func (c *Config) Validate() error {
	if !c.Bot.DryRun {
		if c.Binance.APIKey == "" || c.Binance.APIKey == "your_binance_api_key" {
			return fmt.Errorf("binance API key not configured")
		}
		if c.Binance.SecretKey == "" || c.Binance.SecretKey == "your_binance_secret_key" {
			return fmt.Errorf("binance secret key not configured")
		}
	}

	if c.Trading.InitialBalance <= 0 {
		return fmt.Errorf("initial balance must be positive")
	}

	if c.Trading.MaxRisk <= 0 || c.Trading.MaxRisk >= 1 {
		return fmt.Errorf("max risk must be between 0 and 1")
	}

	if c.Bot.IntervalSeconds <= 0 {
		return fmt.Errorf("interval seconds must be positive")
	}

	return nil
}
