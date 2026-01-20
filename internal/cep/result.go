package cep

// Address representa os dados de endereço retornados pelas APIs
type Address struct {
	CEP          string `json:"cep"`
	Street       string `json:"street,omitempty"`       // ViaCEP usa "logradouro"
	Neighborhood string `json:"neighborhood,omitempty"` // ViaCEP usa "bairro"
	City         string `json:"city,omitempty"`         // ViaCEP usa "localidade"
	State        string `json:"state,omitempty"`        // ViaCEP usa "uf"
	Source       string `json:"source"`                 // Qual API retornou o resultado
}

// ProviderResponse encapsula a resposta de um provider
type ProviderResponse struct {
	Address *Address
	Error   error
	Source  string
}
