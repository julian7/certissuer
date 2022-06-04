package jsonclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

//ReadJSON reads JSON from HTTP Response and unmarshals it
func ReadJSON(resp *http.Response, data interface{}) error {
	b := new(bytes.Buffer)
	defer resp.Body.Close()

	if _, err := b.ReadFrom(resp.Body); err != nil {
		return fmt.Errorf("reading response: %w", err)
	}

	return json.Unmarshal(b.Bytes(), data)
}

//WriteJSON converts struct to JSON as an IO reader
func WriteJSON(data interface{}) (io.Reader, error) {
	marshaledData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("marshaling json for writing: %w", err)
	}

	return bytes.NewReader(marshaledData), nil
}
