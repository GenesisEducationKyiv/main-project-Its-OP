package infrastructure

import "net/http"

type IHttpClient interface {
	SendRequest(req *http.Request) (*HttpResponse, error)
}

type HttpResponse struct {
	Body []byte
	Code int
}
