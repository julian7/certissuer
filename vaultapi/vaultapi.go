package vaultapi

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/julian7/certissuer/jsonclient"
)

//HTTPClient is a HTTP Client which can do HTTP requests
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

//VaultAPI is the interface of a vault service
type VaultAPI struct {
	jsonclient.Client
	URL       *url.URL
	TLSVerify bool
	token     string
}

//New initializes a new VaultAPI object
func New(vaultaddr, token string, tlsverify bool) (*VaultAPI, error) {

	vaulturl, err := url.Parse(vaultaddr)
	if err != nil {
		return nil, fmt.Errorf("initializing VaultAPI: %w", err)
	}

	client := jsonclient.New(tlsverify)
	ret := &VaultAPI{
		URL:       vaulturl,
		Client:    client,
		TLSVerify: tlsverify,
		token:     token,
	}
	client.AddHook("ALL", ret.clientHookAddToken())
	client.AddHook("POST", func(req *http.Request) error {
		req.Header.Set("Content-type", "application/json")
		return nil
	})

	return ret, nil
}

func (api *VaultAPI) clientHookAddToken() func(*http.Request) error {
	return func(req *http.Request) error {
		if len(api.token) > 0 {
			req.Header.Add("X-Vault-Token", api.token)
		}
		return nil
	}
}
