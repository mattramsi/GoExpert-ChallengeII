//go:build integration
// +build integration

package cep

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_SearchCEP_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	client := NewClient()
	validCEP := "01153000"

	start := time.Now()
	address, err := client.SearchCEP(validCEP)
	duration := time.Since(start)

	require.NoError(t, err, "Should not return error")
	require.NotNil(t, address, "Address should not be nil")
	assert.NotEmpty(t, address.CEP, "CEP should be populated")
	assert.NotEmpty(t, address.Source, "Source should be populated")
	assert.Less(t, duration, 1*time.Second+100*time.Millisecond, "Response time should be less than 1.1s")

	t.Logf("CEP encontrado: %+v", address)
	t.Logf("Fonte: %s", address.Source)
	t.Logf("Tempo de resposta: %v", duration)
}

func TestClient_SearchCEP_InvalidCEP_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	client := NewClient()
	invalidCEP := "00000000"

	address, err := client.SearchCEP(invalidCEP)

	require.Error(t, err, "Should return error for invalid CEP")
	assert.Nil(t, address, "Address should be nil for invalid CEP")
}

func TestClient_SearchCEP_Timeout_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	slowProvider := &mockProvider{
		name:  "SlowProvider",
		delay: 2 * time.Second,
		err:   nil,
	}

	client := NewClientWithProviders(slowProvider)

	start := time.Now()
	address, err := client.SearchCEP("01153000")
	duration := time.Since(start)

	require.Error(t, err, "Should return timeout error")
	assert.Nil(t, address, "Address should be nil on timeout")
	assert.LessOrEqual(t, duration, 1*time.Second+200*time.Millisecond, "Timeout should be ~1s")
	assert.GreaterOrEqual(t, duration, 1*time.Second-100*time.Millisecond, "Timeout should be ~1s (not too fast)")
}
