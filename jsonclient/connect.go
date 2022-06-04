package jsonclient

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

//Request builds a new request
func (client *JSONClient) Request(method, url string, data io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, data)
	if err != nil {
		return nil, fmt.Errorf("cannot make %s request: %w", method, err)
	}
	err = client.callHooks("ALL", req)
	if err != nil {
		return nil, err
	}
	err = client.callHooks(method, req)
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}

//Get makes a standard GET request to vault API
func (client *JSONClient) Get(url string) (*http.Response, error) {
	return client.Request("GET", url, nil)
}

//GetJSON makes a standard GET request to Vault API, and returns JSON response
func (client *JSONClient) GetJSON(url string, data interface{}) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("GET request %s returned with %s", url, resp.Status)
	}
	return ReadJSON(resp, data)
}

//Post makes a standard POST request to vault API
func (client *JSONClient) Post(url string, data io.Reader) (*http.Response, error) {
	return client.Request("POST", url, data)
}

// PostJSON returns a standard HTTP response coming from the API server's POST endpoint
func (client *JSONClient) PostJSON(url string, data interface{}) (*http.Response, error) {
	dataReader, err := WriteJSON(data)
	if err != nil {
		return nil, fmt.Errorf("cannot make POST request body: %w", err)
	}
	resp, err := client.Post(url, dataReader)
	if err != nil {
		return nil, fmt.Errorf("cannot create new POST request: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b := new(bytes.Buffer)
		defer resp.Body.Close()

		_, err := b.ReadFrom(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", resp.Status, err)
		}

		return nil, fmt.Errorf("%s: %s", resp.Status, b.String())
	}
	return resp, nil
}

//PostAndReceiveJSON posts a JSON payload to the API server, and marshals
// the received JSON object into a struct
func (client *JSONClient) PostAndReceiveJSON(url string, data interface{}, ret interface{}) error {
	resp, err := client.PostJSON(url, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return ReadJSON(resp, ret)
}
