package cep

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockProvider é um provider mock para testes
type mockProvider struct {
	name      string
	address   *Address
	err       error
	delay     time.Duration
	fetchFunc func(string) (*Address, error)
}

func (m *mockProvider) FetchCEP(cep string) (*Address, error) {
	if m.delay > 0 {
		time.Sleep(m.delay)
	}
	if m.fetchFunc != nil {
		return m.fetchFunc(cep)
	}
	if m.err != nil {
		return nil, m.err
	}
	return m.address, nil
}

func (m *mockProvider) GetName() string {
	return m.name
}

func TestClient_SearchCEP_FirstProviderSuccess(t *testing.T) {
	// Arrange
	fastProvider := &mockProvider{
		name: "FastProvider",
		address: &Address{
			CEP:    "01153000",
			Street: "Rua Teste",
			Source: "FastProvider",
		},
		delay: 10 * time.Millisecond,
	}

	slowProvider := &mockProvider{
		name: "SlowProvider",
		address: &Address{
			CEP:    "01153000",
			Street: "Rua Teste",
			Source: "SlowProvider",
		},
		delay: 100 * time.Millisecond,
	}

	client := NewClientWithProviders(fastProvider, slowProvider)

	// Act
	address, err := client.SearchCEP("01153000")

	// Assert
	require.NoError(t, err, "Should not return error")
	require.NotNil(t, address, "Address should not be nil")
	assert.Equal(t, "FastProvider", address.Source, "Should return result from fast provider")
	assert.Equal(t, "01153000", address.CEP, "CEP should match")
}

func TestClient_SearchCEP_SecondProviderSuccess(t *testing.T) {
	// Arrange
	fastProvider := &mockProvider{
		name:  "FastProvider",
		err:   errors.New("API error"),
		delay: 10 * time.Millisecond,
	}

	slowProvider := &mockProvider{
		name: "SlowProvider",
		address: &Address{
			CEP:    "01153000",
			Street: "Rua Teste",
			Source: "SlowProvider",
		},
		delay: 50 * time.Millisecond,
	}

	client := NewClientWithProviders(fastProvider, slowProvider)

	// Act
	address, err := client.SearchCEP("01153000")

	// Assert
	require.NoError(t, err, "Should not return error")
	require.NotNil(t, address, "Address should not be nil")
	assert.Equal(t, "SlowProvider", address.Source, "Should return result from slow provider when fast provider fails")
}

func TestClient_SearchCEP_Timeout(t *testing.T) {
	// Arrange
	slowProvider1 := &mockProvider{
		name:  "SlowProvider1",
		delay: 2 * time.Second, // Mais que o timeout de 1 segundo
		err:   errors.New("timeout"),
	}

	slowProvider2 := &mockProvider{
		name:  "SlowProvider2",
		delay: 2 * time.Second, // Mais que o timeout de 1 segundo
		err:   errors.New("timeout"),
	}

	client := NewClientWithProviders(slowProvider1, slowProvider2)

	// Act
	address, err := client.SearchCEP("01153000")

	// Assert
	require.Error(t, err, "Should return timeout error")
	assert.Nil(t, address, "Address should be nil on timeout")
	assert.Contains(t, strings.ToLower(err.Error()), "timeout", "Error message should contain 'timeout'")
}

func TestClient_SearchCEP_EmptyCEP(t *testing.T) {
	// Arrange
	client := NewClient()

	// Act
	address, err := client.SearchCEP("")

	// Assert
	require.Error(t, err, "Should return error for empty CEP")
	assert.Nil(t, address, "Address should be nil for empty CEP")
}

func TestClient_SearchCEP_AllProvidersFail(t *testing.T) {
	// Arrange
	provider1 := &mockProvider{
		name:  "Provider1",
		err:   errors.New("API error 1"),
		delay: 10 * time.Millisecond,
	}

	provider2 := &mockProvider{
		name:  "Provider2",
		err:   errors.New("API error 2"),
		delay: 20 * time.Millisecond,
	}

	client := NewClientWithProviders(provider1, provider2)

	// Act
	address, err := client.SearchCEP("00000000")

	// Assert
	require.Error(t, err, "Should return error when all providers fail")
	assert.Nil(t, address, "Address should be nil when all providers fail")
}

func TestBrasilAPIProvider_GetName(t *testing.T) {
	// Arrange
	provider := NewBrasilAPIProvider()

	// Act
	name := provider.GetName()

	// Assert
	assert.Equal(t, "BrasilAPI", name, "Provider name should be 'BrasilAPI'")
}

func TestViaCEPProvider_GetName(t *testing.T) {
	// Arrange
	provider := NewViaCEPProvider()

	// Act
	name := provider.GetName()

	// Assert
	assert.Equal(t, "ViaCEP", name, "Provider name should be 'ViaCEP'")
}
