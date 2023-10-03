package test

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"testing"
	"time"

	"github.com/go-fed/httpsig"
	"github.com/ncarlier/webhookd/pkg/assert"
	"github.com/ncarlier/webhookd/pkg/middleware/signature"
	"github.com/ncarlier/webhookd/pkg/truststore"
)

func assertSigner(t *testing.T) httpsig.Signer {
	prefs := []httpsig.Algorithm{httpsig.RSA_SHA256}
	digestAlgorithm := httpsig.DigestSha256
	headers := []string{httpsig.RequestTarget, "date"}
	signer, _, err := httpsig.NewSigner(prefs, digestAlgorithm, headers, httpsig.Signature, 0)
	assert.Nil(t, err, "")
	return signer
}

func TestHTTPSignature(t *testing.T) {
	privkey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.Nil(t, err, "")
	pubkey := &privkey.PublicKey

	ts := &truststore.InMemoryTrustStore{
		Keys: map[string]crypto.PublicKey{
			"default": pubkey,
		},
	}

	signer := assertSigner(t)
	var body []byte
	req, err := http.NewRequest("GET", "/", http.NoBody)
	assert.Nil(t, err, "")
	req.Header.Add("date", time.Now().UTC().Format(http.TimeFormat))
	err = signer.SignRequest(privkey, "default", req, body)
	assert.Nil(t, err, "")

	err = signature.HTTPSignatureHandler(req, ts)
	assert.Nil(t, err, "")
}
