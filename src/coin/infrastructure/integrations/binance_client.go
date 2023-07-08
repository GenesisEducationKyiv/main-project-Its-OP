package integrations

import (
	"btcRate/coin/application"
	"btcRate/coin/domain"
	"btcRate/coin/infrastructure"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type BinanceClient struct {
	client  infrastructure.IHttpClient
	baseURL *url.URL
}

func NewBinanceClient(client infrastructure.IHttpClient) *BinanceClient {
	baseUrl := &url.URL{Scheme: "https", Host: "api.binance.com", Path: "/api/v3"}
	return &BinanceClient{client: client, baseURL: baseUrl}
}

func (b *BinanceClient) GetRate(currency string, coin string) (*application.SpotPrice, error) {
	path := fmt.Sprintf("/ticker/price?symbol=%s%s", coin, currency)
	url := b.baseURL.String() + path

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := b.client.SendRequest(req)
	if err != nil || resp.Code != http.StatusOK {
		return nil, &domain.EndpointInaccessibleError{Message: endpointInaccessibleErrorMessage}
	}

	timestamp := time.Now()

	var result binanceResponse
	err = json.Unmarshal(resp.Body, &result)
	if err != nil {
		return nil, err
	}

	price, err := strconv.ParseFloat(result.Price, 64)
	if err != nil {
		return nil, err
	}

	return &application.SpotPrice{Amount: price, Timestamp: timestamp}, nil
}

type binanceResponse struct {
	Symbol string
	Price  string
}
