package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var (
	ErrInvalidNonce     = errors.New("nonce is invalid or already used")
	ErrNonceExpired     = errors.New("nonce expired")
	ErrInvalidSignature = errors.New("invalid nonce signature")
)

type nonceStore struct {
	store     map[string]time.Time
	ttl       time.Duration
	secretKey []byte
	mu        sync.Mutex
}

func newNonceStore(ttl time.Duration, secretKey []byte) *nonceStore {
	return &nonceStore{
		store:     make(map[string]time.Time),
		ttl:       ttl,
		secretKey: secretKey,
	}
}

func (ns *nonceStore) generateSignedNonce() (string, string, error) {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	nonce := fmt.Sprintf("%x", rand.Int63())
	ns.store[nonce] = time.Now()

	mac := hmac.New(sha256.New, ns.secretKey)
	if _, err := mac.Write([]byte(nonce)); err != nil {
		return "", "", fmt.Errorf("%w: failed to compute HMAC", ErrInvalidSignature)
	}
	signature := hex.EncodeToString(mac.Sum(nil))

	return nonce, signature, nil
}

func (ns *nonceStore) validateNonce(nonce, signature string) error {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	timestamp, exists := ns.store[nonce]
	if !exists {
		return ErrInvalidNonce
	}

	if time.Since(timestamp) > ns.ttl {
		delete(ns.store, nonce)
		return ErrNonceExpired
	}

	mac := hmac.New(sha256.New, ns.secretKey)
	mac.Write([]byte(nonce))
	expectedSignature := hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(expectedSignature), []byte(signature)) {
		return ErrInvalidSignature
	}

	delete(ns.store, nonce)
	return nil
}
