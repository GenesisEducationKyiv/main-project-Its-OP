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

type BitfinexClient struct {
	client  IExtendedHttpClient
	baseURL *url.URL
}

func NewBitfinexClient(client IExtendedHttpClient) *BitfinexClient {
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

	respBody, code, err := b.client.SendRequest(req)
	if err != nil || code != http.StatusOK {
		return 0.0, time.Time{}, &domain.EndpointInaccessibleError{Message: "Couldn't access the Bitfinex endpoint"}
	}

	timestamp := time.Now()

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
