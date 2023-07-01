package infrastructure

import (
	"io"
	"net/http"
)

type ExtendedHttpClient struct {
	client *http.Client
}

func NewExtendedHttpClient(client *http.Client) *ExtendedHttpClient {
	if client == nil {
		return &ExtendedHttpClient{client: &http.Client{}}
	}

	return &ExtendedHttpClient{client: client}
}

func (c *ExtendedHttpClient) SendRequest(req *http.Request) ([]byte, int, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, 0, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	if err = resp.Body.Close(); err != nil {
		return nil, 0, err
	}

	return body, resp.StatusCode, nil
}
