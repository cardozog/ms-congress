package services

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
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
	return ProviderPayment{ExternalID: fmt.Sprintf("sim-%d-%d", orderID, time.Now().UnixNano())}, nil
}

func simulationEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV"))) {
	case "development", "demo", "test":
		return true
	default:
		return false
	}
}
