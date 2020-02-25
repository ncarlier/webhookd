package pubkey

import (
	"crypto"
	"fmt"
	"github.com/go-fed/httpsig"
	"net/url"
)

// KeyStore is a generic interface to retrieve a public key
type KeyStore interface {
	Get(keyID string) (crypto.PublicKey, httpsig.Algorithm, error)
}

// NewKeyStore creates new Key Store from URI
func NewKeyStore(uri string) (store KeyStore, err error) {
	if uri == "" {
		return nil, nil
	}
	u, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("invalid KeyStore URL: %s", uri)
	}
	switch u.Scheme {
	case "file":
		store, err = newDirectoryKeyStore(u.RawPath)
	default:
		err = fmt.Errorf("non supported KeyStore URL: %s", uri)
	}

	return store, nil
}
