package truststore

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log/slog"
	"os"
)

func newPEMTrustStore(filename string) (TrustStore, error) {
	raw, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	result := &InMemoryTrustStore{
		Keys: make(map[string]crypto.PublicKey),
	}
	for {
		block, rest := pem.Decode(raw)
		if block == nil {
			break
		}
		switch block.Type {
		case "PUBLIC KEY":
			keyID, ok := block.Headers["key_id"]
			if !ok {
				keyID = "default"
			}

			key, err := x509.ParsePKIXPublicKey(block.Bytes)
			if err != nil {
				return nil, err
			}

			result.Keys[keyID] = key
			slog.Debug("public key loaded into the trustore", "id", keyID)
		case "CERTIFICATE":
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return nil, err
			}
			keyID := string(cert.Subject.CommonName)
			result.Keys[keyID] = cert.PublicKey
			slog.Debug("certificate loaded into the trustore", "id", keyID)
		}
		raw = rest
	}

	if len(result.Keys) == 0 {
		return nil, fmt.Errorf("no RSA public key found: %s", filename)
	}
	return result, nil
}
