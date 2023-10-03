package truststore

import (
	"crypto"
	"log/slog"
	"os"

	"golang.org/x/crypto/pkcs12"
)

func newP12TrustStore(filename string) (TrustStore, error) {
	data, err := os.ReadFile(filename)
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
	slog.Debug("certificate loaded into the trustore", "id", keyID)

	return result, nil
}
