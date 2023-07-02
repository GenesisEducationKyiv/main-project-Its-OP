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

func (b *BitfinexClient) GetRate(currency string, coin string) (*application.SpotPrice, error) {
	path := fmt.Sprintf("/pubticker/%s%s", coin, currency)
	url := b.baseURL.String() + path

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := b.client.SendRequest(req)
	if err != nil || resp.Code != http.StatusOK {
		if b.next != nil {
			return b.next.GetRate(currency, coin)
		}

		return nil, &domain.EndpointInaccessibleError{Message: endpointInaccessibleErrorMessage}
	}

	timestamp := time.Now()

	var result bitfinexResponse
	err = json.Unmarshal(resp.Body, &result)
	if err != nil {
		if b.next != nil {
			return b.next.GetRate(currency, coin)
		}

		return nil, err
	}

	price, err := strconv.ParseFloat(result.Amount, 64)
	if err != nil {
		return nil, err
	}

	return &application.SpotPrice{Amount: price, Timestamp: timestamp}, nil
}

func (b *BitfinexClient) SetNext(client application.ICoinClient) {
	b.next = client
}

type bitfinexResponse struct {
	Amount string `json:"ask"`
}
