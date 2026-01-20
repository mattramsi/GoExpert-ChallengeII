package cep

type Address struct {
	CEP          string `json:"cep"`
	Street       string `json:"street,omitempty"`
	Neighborhood string `json:"neighborhood,omitempty"`
	City         string `json:"city,omitempty"`
	State        string `json:"state,omitempty"`
	Source       string `json:"source"`
}

type ProviderResponse struct {
	Address *Address
	Error   error
	Source  string
}
