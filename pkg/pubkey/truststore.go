package pubkey

import (
	"crypto"
	"fmt"
	"path/filepath"

	"github.com/go-fed/httpsig"
	"github.com/ncarlier/webhookd/pkg/logger"
)

const defaultAlgorithm = httpsig.RSA_SHA256

// TrustStoreEntry is a trust store entry
type TrustStoreEntry struct {
	Pubkey    crypto.PublicKey
	Algorithm httpsig.Algorithm
}

// TrustStore is a generic interface to retrieve a public key
type TrustStore interface {
	Get(keyID string) *TrustStoreEntry
}

// NewTrustStore creates new Key Store from URI
func NewTrustStore(filename string) (store TrustStore, err error) {
	if filename == "" {
		return nil, nil
	}

	logger.Debug.Printf("loading trust store: %s", filename)
	switch filepath.Ext(filename) {
	case ".pem":
		store, err = newPEMTrustStore(filename)
	default:
		err = fmt.Errorf("unsupported trust store file format: %s", filename)
	}

	return
}
