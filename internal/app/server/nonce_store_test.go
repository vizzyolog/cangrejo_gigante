package server_test

import (
	"errors"
	"testing"
	"time"

	"cangrejo_gigante/internal/app/server"
)

func TestNonceStore(t *testing.T) {
	t.Parallel()

	store := server.NewNonceStore(1*time.Second, []byte("test-key"))

	nonce, signature, err := store.GenerateSignedNonce()

	if err != nil {
		t.Fatalf("Failed to generate nonce: %v", err)
	}

	if err := store.ValidateNonce(nonce, signature); err != nil {
		t.Errorf("Nonce validation failed: %v", err)
	}

	if err := store.ValidateNonce(nonce, signature); err == nil {
		t.Errorf("Expected error for reused nonce, got nil")
	}

	time.Sleep(2 * time.Second)

	if err := store.ValidateNonce(nonce, signature); err == nil {
		t.Errorf("Expected error for expired nonce, got nil")
	}
}

func TestGenerateSignedNonce(t *testing.T) {
	t.Parallel()

	nonceStore := server.NewNonceStore(5*time.Second, []byte("test-key"))

	nonce, signature, err := nonceStore.GenerateSignedNonce()
	if err != nil {
		t.Fatalf("Failed to generate nonce: %v", err)
	}

	if nonce == "" || signature == "" {
		t.Errorf("Nonce or signature should not be empty")
	}
}

func TestReuseNonce(t *testing.T) {
	t.Parallel()

	nonceStore := server.NewNonceStore(5*time.Second, []byte("test-key"))

	nonce, signature, err := nonceStore.GenerateSignedNonce()
	if err != nil {
		t.Fatalf("Failed to generate nonce: %v", err)
	}

	err = nonceStore.ValidateNonce(nonce, signature)
	if err != nil {
		t.Errorf("Nonce validation failed: %v", err)
	}

	err = nonceStore.ValidateNonce(nonce, signature)
	if !errors.Is(err, server.ErrInvalidNonce) {
		t.Errorf("Expected nonce reuse error, got: %v", err)
	}
}

func TestInvalidSignature(t *testing.T) {
	t.Parallel()

	nonceStore := server.NewNonceStore(5*time.Second, []byte("test-key"))

	nonce, _, err := nonceStore.GenerateSignedNonce()
	if err != nil {
		t.Fatalf("Failed to generate nonce: %v", err)
	}

	err = nonceStore.ValidateNonce(nonce, "invalid-signature")
	if !errors.Is(err, server.ErrInvalidSignature) {
		t.Errorf("Expected invalid signature error, got: %v", err)
	}
}
