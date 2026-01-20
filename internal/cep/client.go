package cep

import (
	"context"
	"fmt"
	"time"
)

const (
	// defaultTimeout define o timeout padrão de 1 segundo para as requisições
	defaultTimeout = 1 * time.Second
)

// Client representa o cliente para buscar CEP de múltiplos providers
// Seguindo o princípio de responsabilidade única (SOLID)
type Client struct {
	providers []Provider
	timeout   time.Duration
}

// NewClient cria uma nova instância do Client com os providers padrão
func NewClient() *Client {
	return &Client{
		providers: []Provider{
			NewBrasilAPIProvider(),
			NewViaCEPProvider(),
		},
		timeout: defaultTimeout,
	}
}

// NewClientWithProviders cria uma nova instância do Client com providers customizados
// Útil para testes e extensibilidade (Open/Closed Principle - SOLID)
func NewClientWithProviders(providers ...Provider) *Client {
	return &Client{
		providers: providers,
		timeout:   defaultTimeout,
	}
}

// SearchCEP busca o CEP nos providers disponíveis e retorna o resultado mais rápido
// Usa goroutines e channels para fazer requisições simultâneas
func (c *Client) SearchCEP(cep string) (*Address, error) {
	if cep == "" {
		return nil, fmt.Errorf("CEP cannot be empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	return c.SearchCEPWithContext(ctx, cep)
}

// SearchCEPWithContext busca o CEP usando um context customizado
// Permite maior controle sobre cancelamento e timeout
func (c *Client) SearchCEPWithContext(ctx context.Context, cep string) (*Address, error) {
	if cep == "" {
		return nil, fmt.Errorf("CEP cannot be empty")
	}

	// Channel para receber as respostas dos providers
	responseChan := make(chan ProviderResponse, len(c.providers))

	// Lança goroutines para cada provider simultaneamente
	for _, provider := range c.providers {
		go func(p Provider) {
			address, err := p.FetchCEP(cep)

			select {
			case responseChan <- ProviderResponse{
				Address: address,
				Error:   err,
				Source:  p.GetName(),
			}:
			case <-ctx.Done():
				// Context foi cancelado, descarta a resposta
				return
			}
		}(provider)
	}

	// Aguarda a primeira resposta bem-sucedida ou timeout
	var firstError error
	responsesReceived := 0

	for responsesReceived < len(c.providers) {
		select {
		case response := <-responseChan:
			responsesReceived++

			// Se a resposta for bem-sucedida, retorna imediatamente
			// (descartando qualquer outra resposta que ainda esteja pendente)
			if response.Error == nil && response.Address != nil {
				return response.Address, nil
			}

			// Armazena o primeiro erro encontrado (caso todas falhem)
			if firstError == nil {
				firstError = response.Error
			}

		case <-ctx.Done():
			// Timeout atingido
			if firstError != nil {
				return nil, fmt.Errorf("timeout: nenhuma API respondeu em %v (%w)", c.timeout, firstError)
			}
			return nil, fmt.Errorf("timeout: nenhuma API respondeu em %v", c.timeout)
		}
	}

	// Se chegou aqui, nenhuma resposta foi bem-sucedida
	if firstError != nil {
		return nil, fmt.Errorf("todas as APIs falharam: %w", firstError)
	}
	return nil, fmt.Errorf("nenhuma resposta recebida")
}
