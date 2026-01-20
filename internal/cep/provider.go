package cep

// Provider define a interface para providers de CEP
// Seguindo o princípio de inversão de dependência (SOLID)
type Provider interface {
	FetchCEP(cep string) (*Address, error)
	GetName() string
}
