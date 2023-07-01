package extensions

import (
	"btcRate/infrastructure"
	"fmt"
	"net/http"
	"time"
)

type LoggedHttpClient struct {
	httpClient infrastructure.IHttpClient
}

func NewLoggedHttpClient(httpClient infrastructure.IHttpClient) *LoggedHttpClient {
	return &LoggedHttpClient{httpClient: httpClient}
}

func (c *LoggedHttpClient) SendRequest(req *http.Request) (*infrastructure.HttpResponse, error) {
	url := req.URL.String()
	timestamp := time.Now()

	resp, err := c.httpClient.SendRequest(req)

	if err != nil {
		return nil, err
	}

	fmt.Printf("%s,%s,%d,%s", timestamp.Format("02-01-06 15:04:05.999 Z0700"), url, resp.Code, string(resp.Body))

	return resp, nil
}
