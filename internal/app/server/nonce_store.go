package server

import (
	"sync"
	"time"
)

type NonceStore struct {
	store map[string]time.Time
	ttl   time.Duration
	mu    sync.Mutex
}

func NewNonceStore(ttl time.Duration) *NonceStore {
	return &NonceStore{
		store: make(map[string]time.Time),
		ttl:   ttl,
		mu:    sync.Mutex{},
	}
}

func (ns *NonceStore) Save(nonce string) error {
	ns.mu.Lock()
	defer ns.mu.Unlock()
	ns.store[nonce] = time.Now()

	return nil
}

func (ns *NonceStore) IsValid(nonce string) bool {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	timestamp, exists := ns.store[nonce]
	if !exists {
		return false
	}

	if time.Since(timestamp) > ns.ttl {
		delete(ns.store, nonce)

		return false
	}

	return true
}

func (ns *NonceStore) MarkAsUsed(nonce string) {
	ns.mu.Lock()
	defer ns.mu.Unlock()
	delete(ns.store, nonce)
}
