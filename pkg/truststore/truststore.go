package truststore

import (
	"crypto"
	"fmt"
	"log/slog"
	"path/filepath"
)

// TrustStore is a generic interface to retrieve a public key
type TrustStore interface {
	GetPublicKey(keyID string) crypto.PublicKey
}

// InMemoryTrustStore is a in memory storage for public keys
type InMemoryTrustStore struct {
	Keys map[string]crypto.PublicKey
}

// GetPublicKey returns the public key with this key ID
func (ts *InMemoryTrustStore) GetPublicKey(keyID string) crypto.PublicKey {
	if key, ok := ts.Keys[keyID]; ok {
		return key
	}
	return nil
}

// New creates new Trust Store from URI
func New(filename string) (store TrustStore, err error) {
	if filename == "" {
		return nil, nil
	}

	slog.Debug("loading trust store...", "filname", filename)
	switch filepath.Ext(filename) {
	case ".pem":
		store, err = newPEMTrustStore(filename)
	case ".p12":
		store, err = newP12TrustStore(filename)
	default:
		err = fmt.Errorf("unsupported trust store file format: %s", filename)
	}

	return
}
