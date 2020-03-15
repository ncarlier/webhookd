package pubkey

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"

	"github.com/ncarlier/webhookd/pkg/logger"
)

type pemTrustStore struct {
	keys map[string]TrustStoreEntry
}

func (ts *pemTrustStore) Get(keyID string) *TrustStoreEntry {
	key, ok := ts.keys[keyID]
	if ok {
		return &key
	}
	return nil
}

func newPEMTrustStore(filename string) (*pemTrustStore, error) {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	result := pemTrustStore{
		keys: make(map[string]TrustStoreEntry),
	}
	for {
		block, rest := pem.Decode(raw)
		if block == nil {
			break
		}
		switch block.Type {
		case "PUBLIC KEY":
			pub, err := x509.ParsePKIXPublicKey(block.Bytes)
			if err != nil {
				return nil, err
			}
			rsaPublicKey, _ := pub.(*rsa.PublicKey)
			keyID, ok := block.Headers["key_id"]
			if !ok {
				keyID = "default"
			}
			result.keys[keyID] = TrustStoreEntry{
				Algorythm: defaultAlgorithm,
				Pubkey:    rsaPublicKey,
			}
			logger.Debug.Printf("public key \"%s\" loaded into the trustore", keyID)
		case "CERTIFICATE":
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return nil, err
			}
			rsaPublicKey, _ := cert.PublicKey.(*rsa.PublicKey)
			keyID := string(cert.Subject.CommonName)
			result.keys[keyID] = TrustStoreEntry{
				Algorythm: defaultAlgorithm,
				Pubkey:    rsaPublicKey,
			}
			logger.Debug.Printf("certificate \"%s\" loaded into the trustore", keyID)
		}
		raw = rest
	}

	if len(result.keys) == 0 {
		return nil, fmt.Errorf("no RSA public key found: %s", filename)
	}
	return &result, nil
}
