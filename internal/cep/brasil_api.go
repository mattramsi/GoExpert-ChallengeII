package cep

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	brasilAPIBaseURL = "https://brasilapi.com.br/api/cep/v1"
	brasilAPIName    = "BrasilAPI"
)

// BrasilAPIProvider implementa o Provider para a API BrasilAPI
type BrasilAPIProvider struct {
	client  *http.Client
	baseURL string
}

// NewBrasilAPIProvider cria uma nova instância do provider BrasilAPI
func NewBrasilAPIProvider() *BrasilAPIProvider {
	return &BrasilAPIProvider{
		client: &http.Client{
			Timeout: httpTimeout,
		},
		baseURL: brasilAPIBaseURL,
	}
}

// GetName retorna o nome do provider
func (b *BrasilAPIProvider) GetName() string {
	return brasilAPIName
}

// brasilAPIResponse representa a estrutura de resposta da BrasilAPI
type brasilAPIResponse struct {
	CEP          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

// FetchCEP busca o CEP na API BrasilAPI
func (b *BrasilAPIProvider) FetchCEP(cep string) (*Address, error) {
	url := fmt.Sprintf("%s/%s", b.baseURL, cep)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := b.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResp brasilAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &Address{
		CEP:          apiResp.CEP,
		Street:       apiResp.Street,
		Neighborhood: apiResp.Neighborhood,
		City:         apiResp.City,
		State:        apiResp.State,
		Source:       brasilAPIName,
	}, nil
}
