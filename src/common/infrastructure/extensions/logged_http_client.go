package extensions

import (
	"btcRate/common/application"
	"btcRate/common/infrastructure"
	"fmt"
	"net/http"
	"time"
)

type ILogRepository interface {
	Log(data string) error
}

type LoggedHttpClient struct {
	httpClient infrastructure.IHttpClient
	logger     application.ILogger
}

func NewLoggedHttpClient(httpClient infrastructure.IHttpClient, logger application.ILogger) *LoggedHttpClient {
	return &LoggedHttpClient{httpClient: httpClient, logger: logger}
}

func (c *LoggedHttpClient) SendRequest(req *http.Request) (*infrastructure.HttpResponse, error) {
	url := req.URL.String()
	timestamp := time.Now()

	resp, err := c.httpClient.SendRequest(req)

	var logMessage string
	if err != nil {
		logMessage = fmt.Sprintf("%s,%s", timestamp.Format("02-01-06 15:04:05.999 Z0700"), url)
		c.logger.Error(logMessage, err)
	} else {
		c.logger.Debug("http request executed", "status", "Success", "code", resp.Code, "url", url, "responseBody", string(resp.Body))
	}

	return resp, err
}
