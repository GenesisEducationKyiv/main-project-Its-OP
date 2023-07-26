package extensions

import (
	"btcRate/common/application"
	"btcRate/common/infrastructure"
	"net/http"
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

	resp, err := c.httpClient.SendRequest(req)

	if err != nil {
		c.logger.Debug("http request executed", "status", "Failure", "url", url, "error", err)
	} else {
		c.logger.Debug("http request executed", "status", "Success", "code", resp.Code, "url", url, "responseBody", string(resp.Body))
	}

	return resp, err
}
