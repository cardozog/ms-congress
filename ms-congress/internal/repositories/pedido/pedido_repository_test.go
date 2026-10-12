package repositories

import (
	"testing"
)

func TestNewCredentialCode(t *testing.T) {
	first, err := newCredentialCode()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	second, err := newCredentialCode()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if first == second {
		t.Fatal("credential codes must be unique")
	}
	if len(first) < 40 {
		t.Fatalf("credential code is too short: %d", len(first))
	}
}
