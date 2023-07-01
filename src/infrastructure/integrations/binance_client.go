package integrations

import (
	"btcRate/application"
	"btcRate/domain"
	"btcRate/infrastructure"
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
	next    application.ICoinClient
}

func NewBinanceClient(client infrastructure.IHttpClient) *BinanceClient {
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

	resp, err := b.client.SendRequest(req)
	if err != nil || resp.Code != http.StatusOK {
		if b.next != nil {
			return b.next.GetRate(currency, coin)
		}

		return 0.0, time.Time{}, &domain.EndpointInaccessibleError{Message: endpointInaccessibleErrorMessage}
	}

	timestamp := time.Now()

	var result binanceResponse
	err = json.Unmarshal(resp.Body, &result)
	if err != nil {
		if b.next != nil {
			return b.next.GetRate(currency, coin)
		}

		return 0, time.Time{}, err
	}

	price, err := strconv.ParseFloat(result.Price, 64)
	if err != nil {
		return 0, time.Time{}, err
	}

	return price, timestamp, nil
}

func (b *BinanceClient) SetNext(client application.ICoinClient) {
	b.next = client
}

type binanceResponse struct {
	Symbol string
	Price  string
}
