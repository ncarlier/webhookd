package truststore

import (
	"crypto"
	"io/ioutil"

	"github.com/ncarlier/webhookd/pkg/logger"
	"golang.org/x/crypto/pkcs12"
)

func newP12TrustStore(filename string) (TrustStore, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	_, cert, err := pkcs12.Decode(data, "test")
	if err != nil {
		return nil, err
	}

	result := &InMemoryTrustStore{
		Keys: make(map[string]crypto.PublicKey),
	}

	keyID := string(cert.Subject.CommonName)
	result.Keys[keyID] = cert.PublicKey
	logger.Debug.Printf("certificate \"%s\" loaded into the trustore", keyID)

	return result, nil
}
