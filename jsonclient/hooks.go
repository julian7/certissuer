package jsonclient

import (
	"fmt"
	"net/http"
)

//AddHook adds a new hook into the hooks slice. Possible methods are
//the standard HTTP methods (eg. GET, POST, PUT, etc.), or ALL for
//hooking all methods. ALL hooks run first, then method-specific hooks.
func (client *JSONClient) AddHook(method string, hook RequestHook) {
	if _, ok := client.hooks[method]; !ok {
		client.hooks[method] = make([]RequestHook, 0, 5)
	}
	client.hooks[method] = append(client.hooks[method], hook)
}

func (client *JSONClient) callHooks(method string, req *http.Request) error {
	for _, item := range client.hooks[method] {
		err := item(req)
		if err != nil {
			return fmt.Errorf("halted by request hook: %w", err)
		}
	}
	return nil
}
