package pubkey

import (
	"crypto"
	"fmt"
	"net/url"
	"strings"

	"github.com/go-fed/httpsig"
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
		return nil, fmt.Errorf("invalid KeyStore URI: %s", uri)
	}
	switch u.Scheme {
	case "file":
		store, err = newDirectoryKeyStore(strings.TrimPrefix(uri, "file://"))
	default:
		err = fmt.Errorf("non supported KeyStore URI: %s", uri)
	}

	return
}
