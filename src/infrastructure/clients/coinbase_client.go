package clients

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

type CoinbaseClient struct {
	client  *http.Client
	baseURL *url.URL
}

func NewCoinbaseClient() *BinanceClient {
	baseUrl := &url.URL{Scheme: "https", Host: "api.coinbase.com", Path: "/v2"}
	return &BinanceClient{client: http.DefaultClient, baseURL: baseUrl}
}

func (c *CoinbaseClient) GetRate(currency string, coin string) (float64, time.Time, error) {
	path := fmt.Sprintf("/prices/%s-%s/spot", coin, currency)
	url := c.baseURL.String() + path

	resp, err := c.client.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return 0.0, time.Time{}, &domain.EndpointInaccessibleError{Message: "Couldn't access the Coinbase endpoint"}
	}

	timestamp := time.Now()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0.0, time.Time{}, err
	}

	var result coinbaseResponse
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return 0.0, time.Time{}, err
	}

	price, err := strconv.ParseFloat(result.Data.Amount, 64)
	if err != nil {
		return 0.0, time.Time{}, err
	}

	return price, timestamp, nil
}

type coinbaseResponse struct {
	Data struct {
		Base     string `json:"base"`
		Currency string `json:"currency"`
		Amount   string `json:"amount"`
	} `json:"data"`
}
