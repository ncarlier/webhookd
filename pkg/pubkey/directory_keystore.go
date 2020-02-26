package pubkey

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-fed/httpsig"
	"github.com/ncarlier/webhookd/pkg/logger"
)

const defaultAlgorithm = httpsig.RSA_SHA256

type directoryKeyStore struct {
	algorithm string
	keys      map[string]crypto.PublicKey
}

func (ks *directoryKeyStore) Get(keyID string) (crypto.PublicKey, httpsig.Algorithm, error) {
	key, ok := ks.keys[keyID]
	if !ok {
		return nil, defaultAlgorithm, fmt.Errorf("public key not found: %s", keyID)
	}
	return key, defaultAlgorithm, nil
}

func newDirectoryKeyStore(root string) (*directoryKeyStore, error) {
	store := &directoryKeyStore{
		algorithm: "",
		keys:      make(map[string]crypto.PublicKey),
	}

	walkErr := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".pem" {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			block, _ := pem.Decode(data)
			if block == nil {
				return fmt.Errorf("invalid PEM file: %s", path)
			}

			pub, err := x509.ParsePKIXPublicKey(block.Bytes)
			if err != nil {
				return nil
			}
			rsaPublicKey, ok := pub.(*rsa.PublicKey)
			if !ok {
				return fmt.Errorf("unable to cast public key to RSA public key")
			}

			keyID, ok := block.Headers["key_id"]
			if ok {
				store.keys[keyID] = rsaPublicKey
				logger.Debug.Println("HTTP signature public key loaded: ", path)
			}
		}
		return nil
	})
	return store, walkErr
}
