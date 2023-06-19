package infrastructure

import (
	"btcRate/domain"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type BinanceClient struct {
	client  *http.Client
	baseURL *url.URL
}

func NewBinanceClient() *BinanceClient {
	baseUrl := &url.URL{Scheme: "https", Host: "api.binance.com", Path: "/api/v3"}
	return &BinanceClient{client: http.DefaultClient, baseURL: baseUrl}
}

func (binanceClient *BinanceClient) GetRate(currency string, coin string) (float64, time.Time, error) {
	path := fmt.Sprintf("/ticker/price?symbol=%s%s", coin, currency)
	url := binanceClient.baseURL.String() + path

	resp, err := binanceClient.client.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return 0.0, time.Time{}, &domain.EndpointInaccessibleError{Message: "Couldn't access the Binance endpoint"}
	}

	timestamp := time.Now()

	respBody, _ := io.ReadAll(resp.Body)

	var result PriceResponse
	_ = json.Unmarshal(respBody, &result)

	price, _ := strconv.ParseFloat(result.Price, 64)

	return price, timestamp, nil
}

type PriceResponse struct {
	Symbol string
	Price  string
}
