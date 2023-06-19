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

func (b *BinanceClient) GetRate(currency string, coin string) (float64, time.Time, error) {
	path := fmt.Sprintf("/ticker/price?symbol=%s%s", coin, currency)
	url := b.baseURL.String() + path

	resp, err := b.client.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return 0.0, time.Time{}, &domain.EndpointInaccessibleError{Message: "Couldn't access the Binance endpoint"}
	}

	timestamp := time.Now()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0.0, time.Time{}, err
	}

	var result PriceResponse
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return 0.0, time.Time{}, err
	}

	price, err := strconv.ParseFloat(result.Price, 64)
	if err != nil {
		return 0.0, time.Time{}, err
	}

	return price, timestamp, nil
}

type PriceResponse struct {
	Symbol string
	Price  string
}
