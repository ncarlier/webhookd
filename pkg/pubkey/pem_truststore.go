package pubkey

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"

	"github.com/go-fed/httpsig"
)

type pemTrustStore struct {
	key crypto.PublicKey
}

func (ts *pemTrustStore) Get(keyID string) (crypto.PublicKey, httpsig.Algorithm, error) {
	return ts.key, defaultAlgorithm, nil
}

func newPEMTrustStore(filename string) (*pemTrustStore, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("invalid PEM file: %s", filename)
	}

	var rsaPublicKey *rsa.PublicKey
	switch block.Type {
	case "PUBLIC KEY":
		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rsaPublicKey, _ = pub.(*rsa.PublicKey)
	case "CERTIFICATE":
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, err
		}
		rsaPublicKey, _ = cert.PublicKey.(*rsa.PublicKey)
	}

	if rsaPublicKey == nil {
		return nil, fmt.Errorf("no RSA public key found: %s", filename)
	}
	return &pemTrustStore{
		key: rsaPublicKey,
	}, nil
}
