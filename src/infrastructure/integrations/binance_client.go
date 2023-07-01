package integrations

import (
	"btcRate/domain"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type BinanceClient struct {
	client  IExtendedHttpClient
	baseURL *url.URL
}

func NewBinanceClient(client IExtendedHttpClient) *BinanceClient {
	baseUrl := &url.URL{Scheme: "https", Host: "api.binance.com", Path: "/api/v3"}
	return &BinanceClient{client: client, baseURL: baseUrl}
}

func (b *BinanceClient) GetRate(currency string, coin string) (float64, time.Time, error) {
	path := fmt.Sprintf("/ticker/price?symbol=%s%s", coin, currency)
	url := b.baseURL.String() + path

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, time.Time{}, err
	}

	respBody, code, err := b.client.SendRequest(req)
	if err != nil || code != http.StatusOK {
		return 0.0, time.Time{}, &domain.EndpointInaccessibleError{Message: "Couldn't access the Binance endpoint"}
	}

	timestamp := time.Now()

	var result binanceResponse
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return 0, time.Time{}, err
	}

	price, err := strconv.ParseFloat(result.Price, 64)
	if err != nil {
		return 0, time.Time{}, err
	}

	return price, timestamp, nil
}

type binanceResponse struct {
	Symbol string
	Price  string
}
