package integrations

import "net/http"

type IExtendedHttpClient interface {
	SendRequest(req *http.Request) ([]byte, int, error)
}
