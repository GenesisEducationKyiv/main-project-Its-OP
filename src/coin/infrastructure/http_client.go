package infrastructure

import (
	"io"
	"net/http"
)

type HttpClient struct {
	client *http.Client
}

func NewHttpClient(client *http.Client) *HttpClient {
	if client == nil {
		return &HttpClient{client: &http.Client{}}
	}

	return &HttpClient{client: client}
}

func (c *HttpClient) SendRequest(req *http.Request) (*HttpResponse, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err = resp.Body.Close(); err != nil {
		return nil, err
	}

	return &HttpResponse{Code: resp.StatusCode, Body: body}, nil
}
