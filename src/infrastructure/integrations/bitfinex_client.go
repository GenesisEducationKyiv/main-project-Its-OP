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

type BitfinexClient struct {
	client  infrastructure.IHttpClient
	baseURL *url.URL
	next    application.ICoinClient
}

func NewBitfinexClient(client infrastructure.IHttpClient) *BitfinexClient {
	baseUrl := &url.URL{Scheme: "https", Host: "api.bitfinex.com", Path: "/v1"}
	return &BitfinexClient{client: client, baseURL: baseUrl}
}

func (b *BitfinexClient) GetRate(currency string, coin string) (float64, time.Time, error) {
	path := fmt.Sprintf("/pubticker/%s%s", coin, currency)
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

	var result bitfinexResponse
	err = json.Unmarshal(resp.Body, &result)
	if err != nil {
		if b.next != nil {
			return b.next.GetRate(currency, coin)
		}

		return 0, time.Time{}, err
	}

	price, err := strconv.ParseFloat(result.Amount, 64)
	if err != nil {
		return 0.0, time.Time{}, err
	}

	return price, timestamp, nil
}

func (b *BitfinexClient) SetNext(client application.ICoinClient) {
	b.next = client
}

type bitfinexResponse struct {
	Amount string `json:"ask"`
}
