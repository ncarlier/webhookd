package signature

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/ncarlier/webhookd/pkg/truststore"
)

const (
	defaultKeyID        = "default"
	xSignatureEd25519   = "X-Signature-Ed25519"
	xSignatureTimestamp = "X-Signature-Timestamp"
)

// IsEd25519SignatureRequest test if HTTP headers contains Ed25519 Signature
func IsEd25519SignatureRequest(headers http.Header) bool {
	return headers.Get(xSignatureEd25519) != ""
}

// Ed25519SignatureHandler validate request HTTP signature
func Ed25519SignatureHandler(r *http.Request, ts truststore.TrustStore) error {
	pubkey := ts.GetPublicKey(defaultKeyID)
	if pubkey == nil {
		return fmt.Errorf("public key not found: %s", defaultKeyID)
	}

	key, ok := pubkey.(ed25519.PublicKey)
	if !ok {
		return errors.New("invalid public key: verify the algorithm")
	}

	value := r.Header.Get(xSignatureEd25519)
	timestamp := r.Header.Get(xSignatureTimestamp)
	if value == "" || timestamp == "" {
		return errors.New("missing signature header")
	}

	sig, err := hex.DecodeString(value)
	if err != nil || len(sig) != ed25519.SignatureSize || sig[63]&224 != 0 {
		return fmt.Errorf("invalid signature format: %s", sig)
	}

	var msg bytes.Buffer
	msg.WriteString(timestamp)

	defer r.Body.Close()
	var body bytes.Buffer

	// Copy the original body back into the request after finishing.
	defer func() {
		r.Body = io.NopCloser(&body)
	}()

	// Copy body into buffers
	_, err = io.Copy(&msg, io.TeeReader(r.Body, &body))
	if err != nil {
		return err
	}

	if !ed25519.Verify(key, msg.Bytes(), sig) {
		return errors.New("invalid signature")
	}
	return nil
}
