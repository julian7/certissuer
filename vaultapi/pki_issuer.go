package vaultapi

import (
	"fmt"
)

//CertRequest contains the payload of a PKI's issue request
type CertRequest struct {
	CommonName string `json:"common_name"`
	AltNames   string `json:"alt_names"`
	IPSans     string `json:"ip_sans"`
	TTL        string `json:"ttl"`
}

//CertResponse contains all data around issues certificate
type CertResponse struct {
	LeaseID       string   `json:"lease_id"`
	Renewable     bool     `json:"renewable"`
	LeaseDuration uint     `json:"lease_duration"`
	Data          CertData `json:"data"`
	Warnings      string   `json:"warnings"`
}

//CertData contains issued certificate
type CertData struct {
	Certificate    string   `json:"certificate"`
	IssuingCA      string   `json:"issuing_ca"`
	CAChain        []string `json:"ca_chain"`
	PrivateKey     string   `json:"private_key"`
	PrivateKeyType string   `json:"private_key_type"`
	SerialNumber   string   `json:"serial_number"`
}

//IssueCert requests Vault to issue a TLS certificate
func (api *VaultAPI) IssueCert(pki, role string, request *CertRequest) (*CertResponse, error) {
	endpoint := fmt.Sprintf("%s/issue/%s", pki, role)
	resp := new(CertResponse)

	err := api.PostAndReceiveJSON(endpoint, request, resp)
	if err != nil {
		return nil, fmt.Errorf("cannot issue %s cert in %s PKI: %w", role, pki, err)
	}
	return resp, nil
}
