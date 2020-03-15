package test

import (
	"testing"

	"github.com/go-fed/httpsig"
	"github.com/ncarlier/webhookd/pkg/assert"
	"github.com/ncarlier/webhookd/pkg/logger"
	"github.com/ncarlier/webhookd/pkg/pubkey"
)

func TestKeyStore(t *testing.T) {
	logger.Init("warn")

	ks, err := pubkey.NewTrustStore("test-key.pem")
	assert.Nil(t, err, "")
	assert.NotNil(t, ks, "")

	pk, algo, err := ks.Get("test")
	assert.Nil(t, err, "")
	assert.NotNil(t, pk, "")
	assert.Equal(t, httpsig.RSA_SHA256, algo, "")
}
