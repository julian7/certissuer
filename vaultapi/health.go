package vaultapi

import (
	"errors"
	"fmt"
)

//HealthResponse is the response object from /sys/health endpoint
type HealthResponse struct {
	Initialized bool `json:"initialized"`
	Sealed      bool `json:"sealed"`
}

//GetHealth returns whether the system is healthy
func (api *VaultAPI) GetHealth() error {
	var resp HealthResponse
	err := api.GetJSON("sys/health", &resp)
	if err != nil {
		return fmt.Errorf("cannot get health API endpoint: %w", err)
	}
	if !resp.Initialized {
		return errors.New("Vault is not initialized")
	}
	if resp.Sealed {
		return errors.New("Vault is sealed")
	}
	return nil
}
