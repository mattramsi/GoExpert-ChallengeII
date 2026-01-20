package cep

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	viaCEPBaseURL = "http://viacep.com.br/ws"
	viaCEPName    = "ViaCEP"
)

type ViaCEPProvider struct {
	client  *http.Client
	baseURL string
}

func NewViaCEPProvider() *ViaCEPProvider {
	return &ViaCEPProvider{
		client: &http.Client{
			Timeout: httpTimeout,
		},
		baseURL: viaCEPBaseURL,
	}
}

func (v *ViaCEPProvider) GetName() string {
	return viaCEPName
}

type viaCEPResponse struct {
	CEP         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	UF          string `json:"uf"`
	IBGE        string `json:"ibge"`
	GIA         string `json:"gia"`
	DDD         string `json:"ddd"`
	SIAFI       string `json:"siafi"`
	Erro        bool   `json:"erro"`
}

func (v *ViaCEPProvider) FetchCEP(cep string) (*Address, error) {
	url := fmt.Sprintf("%s/%s/json/", v.baseURL, cep)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := v.client.Do(req)
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

	var apiResp viaCEPResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if apiResp.Erro {
		return nil, fmt.Errorf("CEP not found")
	}

	return &Address{
		CEP:          apiResp.CEP,
		Street:       apiResp.Logradouro,
		Neighborhood: apiResp.Bairro,
		City:         apiResp.Localidade,
		State:        apiResp.UF,
		Source:       viaCEPName,
	}, nil
}
