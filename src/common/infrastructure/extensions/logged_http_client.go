package extensions

import (
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
	repository ILogRepository
}

func NewLoggedHttpClient(httpClient infrastructure.IHttpClient, repository ILogRepository) *LoggedHttpClient {
	return &LoggedHttpClient{httpClient: httpClient, repository: repository}
}

func (c *LoggedHttpClient) SendRequest(req *http.Request) (*infrastructure.HttpResponse, error) {
	url := req.URL.String()
	timestamp := time.Now()

	resp, err := c.httpClient.SendRequest(req)

	var logMessage string

	if err != nil {
		logMessage = fmt.Sprintf("%s,%s,error: '%s'", timestamp.Format("02-01-06 15:04:05.999 Z0700"), url, err.Error())
	} else {
		logMessage = fmt.Sprintf("%s,%s,%d,'%s'", timestamp.Format("02-01-06 15:04:05.999 Z0700"), url, resp.Code, string(resp.Body))
	}

	logErr := c.repository.Log(logMessage)

	if logErr != nil {
		fmt.Printf("error: failed to save the log. %f %s", err, logMessage)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}
