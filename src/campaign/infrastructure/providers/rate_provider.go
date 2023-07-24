package providers

import (
	"btcRate/campaign/domain"
	"btcRate/common/infrastructure"
	"encoding/json"
	"net/http"
	"net/url"
)

const endpointInaccessibleErrorMessage = "Couldn't access the rate provider"

type RateProvider struct {
	client  infrastructure.IHttpClient
	baseURL *url.URL
}

func NewRateProvider(client infrastructure.IHttpClient, baseURL *url.URL) *RateProvider {
	return &RateProvider{client: client, baseURL: baseURL}
}

func (r *RateProvider) GetRate() (*domain.Rate, error) {
	url := r.baseURL.String() + domain.GetRate

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := r.client.SendRequest(req)
	if err != nil || resp.Code != http.StatusOK {
		return nil, &domain.EndpointInaccessibleError{Message: endpointInaccessibleErrorMessage}
	}

	var result domain.Rate
	err = json.Unmarshal(resp.Body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
