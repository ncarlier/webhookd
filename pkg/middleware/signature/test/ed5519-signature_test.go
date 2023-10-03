package test

import (
	"bytes"
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/ncarlier/webhookd/pkg/assert"
	"github.com/ncarlier/webhookd/pkg/middleware/signature"
	"github.com/ncarlier/webhookd/pkg/truststore"
)

func TestEd5519Signature(t *testing.T) {
	pubkey, privkey, err := ed25519.GenerateKey(rand.Reader)
	assert.Nil(t, err, "")

	ts := &truststore.InMemoryTrustStore{
		Keys: map[string]crypto.PublicKey{
			"default": pubkey,
		},
	}

	body := "this is a test"
	req, err := http.NewRequest("POST", "/", bytes.NewBufferString(body))
	assert.Nil(t, err, "")

	now := time.Now()
	timestamp := strconv.FormatInt(now.Unix(), 10)

	var msg bytes.Buffer
	msg.WriteString(timestamp)
	msg.WriteString(body)
	s := ed25519.Sign(privkey, msg.Bytes())
	req.Header.Set("X-Signature-Ed25519", hex.EncodeToString(s[:ed25519.SignatureSize]))
	req.Header.Set("X-Signature-Timestamp", timestamp)
	req.Header.Add("date", now.UTC().Format(http.TimeFormat))
	req.Header.Set("Content-Type", "text/plain")

	err = signature.Ed25519SignatureHandler(req, ts)
	assert.Nil(t, err, "")
}
