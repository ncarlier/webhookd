package test

import (
	"testing"

	"github.com/go-fed/httpsig"
	"github.com/ncarlier/webhookd/pkg/assert"
	"github.com/ncarlier/webhookd/pkg/logger"
	"github.com/ncarlier/webhookd/pkg/pubkey"
)

func TestTrustStoreWithNoKeyID(t *testing.T) {
	logger.Init("warn")

	ts, err := pubkey.NewTrustStore("test-key-01.pem")
	assert.Nil(t, err, "")
	assert.NotNil(t, ts, "")
	entry := ts.Get("test")
	assert.True(t, entry == nil, "")
	entry = ts.Get("default")
	assert.NotNil(t, entry, "")
	assert.Equal(t, httpsig.RSA_SHA256, entry.Algorithm, "")
}

func TestTrustStoreWithKeyID(t *testing.T) {
	logger.Init("warn")

	ts, err := pubkey.NewTrustStore("test-key-02.pem")
	assert.Nil(t, err, "")
	assert.NotNil(t, ts, "")
	entry := ts.Get("test")
	assert.NotNil(t, entry, "")
	assert.Equal(t, httpsig.RSA_SHA256, entry.Algorithm, "")
}

func TestTrustStoreWithCertificate(t *testing.T) {
	logger.Init("warn")

	ts, err := pubkey.NewTrustStore("test-cert.pem")
	assert.Nil(t, err, "")
	assert.NotNil(t, ts, "")
	entry := ts.Get("test.localnet")
	assert.NotNil(t, entry, "")
	assert.Equal(t, httpsig.RSA_SHA256, entry.Algorithm, "")
}

func TestTrustStoreWithMultipleEntries(t *testing.T) {
	logger.Init("warn")

	ts, err := pubkey.NewTrustStore("test-multi.pem")
	assert.Nil(t, err, "")
	assert.NotNil(t, ts, "")
	entry := ts.Get("test.localnet")
	assert.NotNil(t, entry, "")
	assert.Equal(t, httpsig.RSA_SHA256, entry.Algorithm, "")
	entry = ts.Get("foo")
	assert.NotNil(t, entry, "")
	assert.Equal(t, httpsig.RSA_SHA256, entry.Algorithm, "")
}
