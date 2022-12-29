package truststore

import (
	"crypto"
	"fmt"
	"path/filepath"

	"github.com/ncarlier/webhookd/pkg/logger"
)

// TrustStore is a generic interface to retrieve a public key
type TrustStore interface {
	GetPublicKey(keyID string) crypto.PublicKey
}

// InMemoryTrustStore is a in memory storage for public keys
type InMemoryTrustStore struct {
	Keys map[string]crypto.PublicKey
}

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

	logger.Debug.Printf("loading trust store: %s", filename)
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
