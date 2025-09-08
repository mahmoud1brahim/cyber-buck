package market

import (
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"time"
)

type Data struct {
	Symbol    string    `json:"symbol"`
	Price     float64   `json:"price"`
	Volume    float64   `json:"volume"`
	Timestamp time.Time `json:"timestamp"`
}

type CoinGeckoResponse struct {
	Bitcoin struct {
		USD float64 `json:"usd"`
	} `json:"bitcoin"`
}

func FetchMockData(symbol string) (*Data, error) {
	if symbol == "BTC/USD" || symbol == "BTCUSDT" {
		return fetchBTCPrice()
	}
	return generateMockData(symbol), nil
}

func fetchBTCPrice() (*Data, error) {
	url := "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd"

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return generateMockData("BTC/USD"), nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return generateMockData("BTC/USD"), nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return generateMockData("BTC/USD"), nil
	}

	var cgResp CoinGeckoResponse
	if err := json.Unmarshal(body, &cgResp); err != nil {
		return generateMockData("BTC/USD"), nil
	}

	return &Data{
		Symbol:    "BTC/USD",
		Price:     cgResp.Bitcoin.USD,
		Volume:    rand.Float64() * 1000000,
		Timestamp: time.Now(),
	}, nil
}

func generateMockData(symbol string) *Data {
	basePrice := 45000.0
	if symbol == "BTC/USD" || symbol == "BTCUSDT" {
		variation := (rand.Float64() - 0.5) * 2000
		return &Data{
			Symbol:    symbol,
			Price:     basePrice + variation,
			Volume:    rand.Float64() * 1000000,
			Timestamp: time.Now(),
		}
	}

	return &Data{
		Symbol:    symbol,
		Price:     1000.0 + (rand.Float64()-0.5)*100,
		Volume:    rand.Float64() * 100000,
		Timestamp: time.Now(),
	}
}
