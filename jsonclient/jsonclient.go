package jsonclient

import (
	"crypto/tls"
	"io"
	"net/http"
)

//RequestHook is a single item in the hook slice to call whenever a call is made
type RequestHook func(request *http.Request) error

//Client is a HTTP Client which can do HTTP requests
type Client interface {
	AddHook(method string, callback RequestHook)
	Do(*http.Request) (*http.Response, error)
	Request(method string, url string, payload io.Reader) (*http.Response, error)
	Get(url string) (*http.Response, error)
	GetJSON(url string, data interface{}) error
	Post(url string, payload io.Reader) (*http.Response, error)
	PostJSON(url string, data interface{}) (*http.Response, error)
	PostAndReceiveJSON(url string, data interface{}, result interface{}) error
}

//JSONClient is a HTTP Client implementation
type JSONClient struct {
	*http.Client
	hooks map[string][]RequestHook
}

//New creates a new JSONClient
func New(tlsverify bool) *JSONClient {
	client := new(http.Client)
	if !tlsverify {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	return &JSONClient{Client: client, hooks: make(map[string][]RequestHook)}
}
