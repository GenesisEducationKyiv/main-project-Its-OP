package integrations

import "net/http"

type Code int
type Body []byte

type IExtendedHttpClient interface {
	SendRequest(req *http.Request) (Body, Code, error)
}
