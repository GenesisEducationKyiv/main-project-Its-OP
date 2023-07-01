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

type CoinbaseClient struct {
	client  infrastructure.IHttpClient
	baseURL *url.URL
	next    application.ICoinClient
}

func NewCoinbaseClient(client infrastructure.IHttpClient) *CoinbaseClient {
	baseUrl := &url.URL{Scheme: "https", Host: "api.coinbase.com", Path: "/v2"}
	return &CoinbaseClient{client: client, baseURL: baseUrl}
}

func (c *CoinbaseClient) GetRate(currency string, coin string) (float64, time.Time, error) {
	path := fmt.Sprintf("/prices/%s-%s/spot", coin, currency)
	url := c.baseURL.String() + path

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, time.Time{}, err
	}

	resp, err := c.client.SendRequest(req)
	if err != nil || resp.Code != http.StatusOK {
		if c.next != nil {
			return c.next.GetRate(currency, coin)
		}

		return 0.0, time.Time{}, &domain.EndpointInaccessibleError{Message: endpointInaccessibleErrorMessage}
	}

	timestamp := time.Now()

	var result coinbaseResponse
	err = json.Unmarshal(resp.Body, &result)
	if err != nil {
		if c.next != nil {
			return c.next.GetRate(currency, coin)
		}

		return 0, time.Time{}, err
	}

	price, err := strconv.ParseFloat(result.Data.Amount, 64)
	if err != nil {
		return 0.0, time.Time{}, err
	}

	return price, timestamp, nil
}

func (b *CoinbaseClient) SetNext(client application.ICoinClient) {
	b.next = client
}

type coinbaseResponse struct {
	Data struct {
		Amount string `json:"amount"`
	} `json:"data"`
}
