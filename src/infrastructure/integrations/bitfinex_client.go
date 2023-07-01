package integrations

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

type BitfinexClient struct {
	client  *http.Client
	baseURL *url.URL
}

func NewByBitClient() *BitfinexClient {
	baseUrl := &url.URL{Scheme: "https", Host: "api.bitfinex.com", Path: "/v1"}
	return &BitfinexClient{client: http.DefaultClient, baseURL: baseUrl}
}

func (b *BitfinexClient) GetRate(currency string, coin string) (float64, time.Time, error) {
	path := fmt.Sprintf("/pubticker/%s%s", coin, currency)
	url := b.baseURL.String() + path

	resp, err := b.client.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return 0.0, time.Time{}, &domain.EndpointInaccessibleError{Message: "Couldn't access the Bitfinex endpoint"}
	}

	timestamp := time.Now()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0.0, time.Time{}, err
	}

	var result bitfinexResponse
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return 0.0, time.Time{}, err
	}

	price, err := strconv.ParseFloat(result.Amount, 64)
	if err != nil {
		return 0.0, time.Time{}, err
	}

	return price, timestamp, nil
}

type bitfinexResponse struct {
	Amount string `json:"ask"`
}
