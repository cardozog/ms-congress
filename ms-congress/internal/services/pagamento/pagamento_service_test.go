package services

import "testing"

func TestValidateCardAcceptsFormattedTestCard(t *testing.T) {
	details, message := validateCard(CardData{
		Number: "4111 1111 1111 1111", CPF: "529.982.247-25",
		Expiry: "1230", SecurityCode: "123",
	})
	if message != "" {
		t.Fatalf("expected valid card, got %q", message)
	}
	if details.CardLast4 != "1111" || details.ExpiryMonth != 12 || details.ExpiryYear != 2030 {
		t.Fatalf("unexpected persisted card details: %+v", details)
	}
}

func TestValidateCardRejectsInvalidCPF(t *testing.T) {
	_, message := validateCard(CardData{
		Number: "4111111111111111", CPF: "22222222222",
		Expiry: "12/30", SecurityCode: "123",
	})
	if message != "CPF inválido" {
		t.Fatalf("expected CPF validation error, got %q", message)
	}
}

func TestValidateCardRejectsExpiredCard(t *testing.T) {
	_, message := validateCard(CardData{
		Number: "4111111111111111", CPF: "52998224725",
		Expiry: "01/20", SecurityCode: "123",
	})
	if message != "cartão vencido" {
		t.Fatalf("expected expiry validation error, got %q", message)
	}
}
