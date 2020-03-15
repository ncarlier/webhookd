package pubkey

import (
	"crypto"
	"fmt"
	"path/filepath"

	"github.com/go-fed/httpsig"
)

const defaultAlgorithm = httpsig.RSA_SHA256

// TrustStore is a generic interface to retrieve a public key
type TrustStore interface {
	Get(keyID string) (crypto.PublicKey, httpsig.Algorithm, error)
}

// NewTrustStore creates new Key Store from URI
func NewTrustStore(filename string) (store TrustStore, err error) {
	if filename == "" {
		return nil, nil
	}

	switch filepath.Ext(filename) {
	case ".pem":
		store, err = newPEMTrustStore(filename)
	default:
		err = fmt.Errorf("unsupported TrustStore file format: %s", filename)
	}

	return
}
