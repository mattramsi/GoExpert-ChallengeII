package cep

import (
	"context"
	"fmt"
	"time"
)

const (
	defaultTimeout = 1 * time.Second
)

type Client struct {
	providers []Provider
	timeout   time.Duration
}

func NewClient() *Client {
	return &Client{
		providers: []Provider{
			NewBrasilAPIProvider(),
			NewViaCEPProvider(),
		},
		timeout: defaultTimeout,
	}
}

func NewClientWithProviders(providers ...Provider) *Client {
	return &Client{
		providers: providers,
		timeout:   defaultTimeout,
	}
}

func (c *Client) SearchCEP(cep string) (*Address, error) {
	if cep == "" {
		return nil, fmt.Errorf("CEP cannot be empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	return c.SearchCEPWithContext(ctx, cep)
}

func (c *Client) SearchCEPWithContext(ctx context.Context, cep string) (*Address, error) {
	if cep == "" {
		return nil, fmt.Errorf("CEP cannot be empty")
	}

	responseChan := make(chan ProviderResponse, len(c.providers))

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
				return
			}
		}(provider)
	}

	var firstError error
	responsesReceived := 0

	for responsesReceived < len(c.providers) {
		select {
		case response := <-responseChan:
			responsesReceived++

			if response.Error == nil && response.Address != nil {
				return response.Address, nil
			}

			if firstError == nil {
				firstError = response.Error
			}

		case <-ctx.Done():
			if firstError != nil {
				return nil, fmt.Errorf("timeout: nenhuma API respondeu em %v (%w)", c.timeout, firstError)
			}
			return nil, fmt.Errorf("timeout: nenhuma API respondeu em %v", c.timeout)
		}
	}

	if firstError != nil {
		return nil, fmt.Errorf("todas as APIs falharam: %w", firstError)
	}
	return nil, fmt.Errorf("nenhuma resposta recebida")
}
