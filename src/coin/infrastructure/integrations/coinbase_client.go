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

type CoinbaseClient struct {
	client  infrastructure.IHttpClient
	baseURL *url.URL
}

func NewCoinbaseClient(client infrastructure.IHttpClient) *CoinbaseClient {
	baseUrl := &url.URL{Scheme: "https", Host: "api.coinbase.com", Path: "/v2"}
	return &CoinbaseClient{client: client, baseURL: baseUrl}
}

func (c *CoinbaseClient) GetRate(currency string, coin string) (*application.SpotPrice, error) {
	path := fmt.Sprintf("/prices/%s-%s/spot", coin, currency)
	url := c.baseURL.String() + path

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.SendRequest(req)
	if err != nil || resp.Code != http.StatusOK {
		return nil, &domain.EndpointInaccessibleError{Message: endpointInaccessibleErrorMessage}
	}

	timestamp := time.Now()

	var result coinbaseResponse
	err = json.Unmarshal(resp.Body, &result)
	if err != nil {
		return nil, err
	}

	price, err := strconv.ParseFloat(result.Data.Amount, 64)
	if err != nil {
		return nil, err
	}

	return &application.SpotPrice{Amount: price, Timestamp: timestamp}, err
}

type coinbaseResponse struct {
	Data struct {
		Amount string `json:"amount"`
	} `json:"data"`
}
