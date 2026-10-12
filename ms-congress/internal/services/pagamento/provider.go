package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"
)

type ProviderPayment struct {
	ExternalID string
}

type Provider interface {
	CreatePayment(orderID uint64, amount float64, idempotencyKey string) (ProviderPayment, error)
}

type SimulatedProvider struct{}

func (SimulatedProvider) CreatePayment(orderID uint64, _ float64, _ string) (ProviderPayment, error) {
	if !simulationEnabled() {
		return ProviderPayment{}, errors.New("provedor simulado desabilitado fora de desenvolvimento")
	}
	random := make([]byte, 24)
	if _, err := rand.Read(random); err != nil {
		return ProviderPayment{}, fmt.Errorf("gerar referência do pagamento: %w", err)
	}
	return ProviderPayment{ExternalID: fmt.Sprintf("sim-%d-%s", orderID, base64.RawURLEncoding.EncodeToString(random))}, nil
}

func simulationEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV"))) {
	case "development", "demo", "test":
		return true
	default:
		return false
	}
}
