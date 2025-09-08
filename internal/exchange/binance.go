package exchange

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"trading-bot/internal/market"
)

type BinanceClient struct {
	APIKey     string
	SecretKey  string
	BaseURL    string
	HTTPClient *http.Client
}

type BinanceTicker struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

func NewBinanceClient(apiKey, secretKey string, testNet bool) (*BinanceClient, error) {
	baseURL := "https://api.binance.com"
	if testNet {
		baseURL = "https://testnet.binance.vision"
	}

	return &BinanceClient{
		APIKey:     apiKey,
		SecretKey:  secretKey,
		BaseURL:    baseURL,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}, nil
}

func (bc *BinanceClient) GetMarketData(symbol string) (*market.Data, error) {
	endpoint := "/api/v3/ticker/price"
	params := url.Values{}
	params.Add("symbol", symbol)

	url := fmt.Sprintf("%s%s?%s", bc.BaseURL, endpoint, params.Encode())

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := bc.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ticker BinanceTicker
	if err := json.Unmarshal(body, &ticker); err != nil {
		return nil, err
	}

	price, err := strconv.ParseFloat(ticker.Price, 64)
	if err != nil {
		return nil, err
	}

	return &market.Data{
		Symbol:    symbol,
		Price:     price,
		Volume:    0,
		Timestamp: time.Now(),
	}, nil
}

func (bc *BinanceClient) PlaceOrder(symbol string, side Side, orderType OrderType, quantity, price float64) (*OrderResponse, error) {
	endpoint := "/api/v3/order"

	params := url.Values{}
	params.Add("symbol", symbol)
	params.Add("side", string(side))
	params.Add("type", string(orderType))
	params.Add("quantity", fmt.Sprintf("%.8f", quantity))

	if orderType == TypeLimit {
		params.Add("price", fmt.Sprintf("%.8f", price))
		params.Add("timeInForce", "GTC")
	}

	params.Add("timestamp", strconv.FormatInt(time.Now().Unix()*1000, 10))

	signature := bc.generateSignature(params.Encode())
	params.Add("signature", signature)

	url := fmt.Sprintf("%s%s", bc.BaseURL, endpoint)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-MBX-APIKEY", bc.APIKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Body = io.NopCloser(strings.NewReader(params.Encode()))

	resp, err := bc.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("binance API error: %s", string(body))
	}

	var orderResp OrderResponse
	if err := json.Unmarshal(body, &orderResp); err != nil {
		return nil, err
	}

	return &orderResp, nil
}

func (bc *BinanceClient) TestConnection() error {
	endpoint := "/api/v3/ping"
	url := fmt.Sprintf("%s%s", bc.BaseURL, endpoint)

	resp, err := bc.HTTPClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to ping Binance API: status %d", resp.StatusCode)
	}

	return nil
}

func (bc *BinanceClient) generateSignature(queryString string) string {
	h := hmac.New(sha256.New, []byte(bc.SecretKey))
	h.Write([]byte(queryString))
	return hex.EncodeToString(h.Sum(nil))
}
