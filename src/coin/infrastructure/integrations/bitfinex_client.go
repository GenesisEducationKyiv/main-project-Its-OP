package integrations

import (
	"coin/application/services"
	"coin/domain"
	"coin/infrastructure"
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
}

func NewBitfinexClient(client infrastructure.IHttpClient) *BitfinexClient {
	baseUrl := &url.URL{Scheme: "https", Host: "api.bitfinex.com", Path: "/v1"}
	return &BitfinexClient{client: client, baseURL: baseUrl}
}

func (b *BitfinexClient) GetRate(currency string, coin string) (*services.SpotPrice, error) {
	path := fmt.Sprintf("/pubticker/%s%s", coin, currency)
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

	var result bitfinexResponse
	err = json.Unmarshal(resp.Body, &result)
	if err != nil {
		return nil, err
	}

	price, err := strconv.ParseFloat(result.Amount, 64)
	if err != nil {
		return nil, err
	}

	return &services.SpotPrice{Amount: price, Timestamp: timestamp}, nil
}

type bitfinexResponse struct {
	Amount string `json:"ask"`
}
