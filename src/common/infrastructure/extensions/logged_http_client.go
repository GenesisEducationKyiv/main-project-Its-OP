package extensions

import (
	"btcRate/coin/infrastructure"
	"fmt"
	"net/http"
	"time"
)

type ILogRepository interface {
	Log(data string) error
}

type LoggedHttpClient struct {
	httpClient infrastructure.IHttpClient
	repository ILogRepository
}

func NewLoggedHttpClient(httpClient infrastructure.IHttpClient, repository ILogRepository) *LoggedHttpClient {
	return &LoggedHttpClient{httpClient: httpClient, repository: repository}
}

func (c *LoggedHttpClient) SendRequest(req *http.Request) (*infrastructure.HttpResponse, error) {
	url := req.URL.String()
	timestamp := time.Now()

	resp, err := c.httpClient.SendRequest(req)

	if err != nil {
		return nil, err
	}

	logMessage := fmt.Sprintf("%s,%s,%d,'%s'", timestamp.Format("02-01-06 15:04:05.999 Z0700"), url, resp.Code, string(resp.Body))

	err = c.repository.Log(logMessage)

	if err != nil {
		fmt.Printf("error: failed to save the log. %f %s", err, logMessage)
	}

	return resp, nil
}
