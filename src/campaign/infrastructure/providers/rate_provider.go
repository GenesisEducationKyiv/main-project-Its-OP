package providers

import (
	"btcRate/campaign/domain"
	"btcRate/common/infrastructure"
	"net/url"
	"time"
)

type RateProvider struct {
	client  infrastructure.IHttpClient
	baseURL *url.URL
}

func NewRateProvider(client infrastructure.IHttpClient) *RateProvider {
	// TODO: set the url of the feature-coin container
	baseUrl := &url.URL{Host: "coin", Path: "path"}
	return &RateProvider{client: client, baseURL: baseUrl}
}

func (*RateProvider) GetRate(currency string, coin string) (domain.Rate, error) {
	// TODO: implement proper rate fetching
	return domain.Rate{Amount: 1000, Timestamp: time.Now()}, nil
}
