package vaultapi

import (
	"net/http"
	"path"
)

func (api *VaultAPI) endpoint2url(endpoint string) string {
	u := *api.URL
	u.Path = path.Join(u.Path, "v1", endpoint)

	return u.String()
}

//GetJSON makes a standard GET request to Vault API, and returns JSON response
func (api *VaultAPI) GetJSON(endpoint string, data interface{}) error {
	return api.Client.GetJSON(api.endpoint2url(endpoint), data)
}

// PostJSON sends a POST request with JSON payload, returns raw HTTP response
func (api *VaultAPI) PostJSON(endpoint string, data interface{}) (*http.Response, error) {
	return api.Client.PostJSON(api.endpoint2url(endpoint), data)
}

// PostAndReceiveJSON sends a POST request with JSON payload, and returns received JSON object marshaled into a struct
func (api *VaultAPI) PostAndReceiveJSON(endpoint string, data interface{}, ret interface{}) error {
	return api.Client.PostAndReceiveJSON(api.endpoint2url(endpoint), data, ret)
}
