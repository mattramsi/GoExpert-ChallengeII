package cep

type Provider interface {
	FetchCEP(cep string) (*Address, error)
	GetName() string
}
