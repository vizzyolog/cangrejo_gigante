package server

import (
	"errors"
	"testing"
	"time"
)

func TestNonceStore(t *testing.T) {
	store := newNonceStore(1*time.Second, []byte("test-key"))

	nonce, signature, err := store.generateSignedNonce()
	if err != nil {
		t.Fatalf("Failed to generate nonce: %v", err)
	}

	if err := store.validateNonce(nonce, signature); err != nil {
		t.Errorf("Nonce validation failed: %v", err)
	}

	if err := store.validateNonce(nonce, signature); err == nil {
		t.Errorf("Expected error for reused nonce, got nil")
	}

	time.Sleep(2 * time.Second)
	if err := store.validateNonce(nonce, signature); err == nil {
		t.Errorf("Expected error for expired nonce, got nil")
	}
}

func TestGenerateSignedNonce(t *testing.T) {
	ns := newNonceStore(5*time.Second, []byte("test-key"))

	nonce, signature, err := ns.generateSignedNonce()
	if err != nil {
		t.Fatalf("Failed to generate nonce: %v", err)
	}

	if nonce == "" || signature == "" {
		t.Errorf("Nonce or signature should not be empty")
	}
}

func TestReuseNonce(t *testing.T) {
	ns := newNonceStore(5*time.Second, []byte("test-key"))

	nonce, signature, err := ns.generateSignedNonce()
	if err != nil {
		t.Fatalf("Failed to generate nonce: %v", err)
	}

	err = ns.validateNonce(nonce, signature)
	if err != nil {
		t.Errorf("Nonce validation failed: %v", err)
	}

	err = ns.validateNonce(nonce, signature)
	if !errors.Is(err, ErrInvalidNonce) {
		t.Errorf("Expected nonce reuse error, got: %v", err)
	}
}

func TestInvalidSignature(t *testing.T) {
	ns := newNonceStore(5*time.Second, []byte("test-key"))

	nonce, _, err := ns.generateSignedNonce()
	if err != nil {
		t.Fatalf("Failed to generate nonce: %v", err)
	}

	err = ns.validateNonce(nonce, "invalid-signature")
	if !errors.Is(err, ErrInvalidSignature) {
		t.Errorf("Expected invalid signature error, got: %v", err)
	}
}
