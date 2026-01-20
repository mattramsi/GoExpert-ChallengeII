//go:build integration
// +build integration

package cep

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestClient_SearchCEP_Integration testa a busca real de CEP nas APIs
// Para executar: go test -tags=integration
func TestClient_SearchCEP_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	client := NewClient()
	validCEP := "01153000" // CEP válido conhecido

	// Act
	start := time.Now()
	address, err := client.SearchCEP(validCEP)
	duration := time.Since(start)

	// Assert
	require.NoError(t, err, "Should not return error")
	require.NotNil(t, address, "Address should not be nil")
	assert.NotEmpty(t, address.CEP, "CEP should be populated")
	assert.NotEmpty(t, address.Source, "Source should be populated")
	assert.Less(t, duration, 1*time.Second+100*time.Millisecond, "Response time should be less than 1.1s")

	t.Logf("CEP encontrado: %+v", address)
	t.Logf("Fonte: %s", address.Source)
	t.Logf("Tempo de resposta: %v", duration)
}

// TestClient_SearchCEP_InvalidCEP_Integration testa com CEP inválido
func TestClient_SearchCEP_InvalidCEP_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	client := NewClient()
	invalidCEP := "00000000" // CEP inválido

	// Act
	address, err := client.SearchCEP(invalidCEP)

	// Assert
	require.Error(t, err, "Should return error for invalid CEP")
	assert.Nil(t, address, "Address should be nil for invalid CEP")
}

// TestClient_SearchCEP_Timeout_Integration testa o timeout
func TestClient_SearchCEP_Timeout_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	// Usa um provider mock muito lento para testar timeout
	slowProvider := &mockProvider{
		name:  "SlowProvider",
		delay: 2 * time.Second,
		err:   nil,
	}

	client := NewClientWithProviders(slowProvider)

	// Act
	start := time.Now()
	address, err := client.SearchCEP("01153000")
	duration := time.Since(start)

	// Assert
	require.Error(t, err, "Should return timeout error")
	assert.Nil(t, address, "Address should be nil on timeout")
	assert.LessOrEqual(t, duration, 1*time.Second+200*time.Millisecond, "Timeout should be ~1s")
	assert.GreaterOrEqual(t, duration, 1*time.Second-100*time.Millisecond, "Timeout should be ~1s (not too fast)")
}
