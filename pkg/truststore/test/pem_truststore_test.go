package test

import (
	"crypto/rsa"
	"testing"

	"github.com/ncarlier/webhookd/pkg/assert"
	"github.com/ncarlier/webhookd/pkg/truststore"
)

func TestTrustStoreWithNoKeyID(t *testing.T) {
	ts, err := truststore.New("test-key-01.pem")
	assert.Nil(t, err, "")
	assert.NotNil(t, ts, "")
	pubkey := ts.GetPublicKey("test")
	assert.True(t, pubkey == nil, "")
	pubkey = ts.GetPublicKey("default")
	assert.NotNil(t, pubkey, "")
	_, ok := pubkey.(*rsa.PublicKey)
	assert.True(t, ok, "")
}

func TestTrustStoreWithKeyID(t *testing.T) {
	ts, err := truststore.New("test-key-02.pem")
	assert.Nil(t, err, "")
	assert.NotNil(t, ts, "")
	pubkey := ts.GetPublicKey("test")
	assert.NotNil(t, pubkey, "")
	_, ok := pubkey.(*rsa.PublicKey)
	assert.True(t, ok, "")
}

func TestTrustStoreWithCertificate(t *testing.T) {
	ts, err := truststore.New("test-cert.pem")
	assert.Nil(t, err, "")
	assert.NotNil(t, ts, "")
	pubkey := ts.GetPublicKey("test.localnet")
	assert.NotNil(t, pubkey, "")
	_, ok := pubkey.(*rsa.PublicKey)
	assert.True(t, ok, "")
}

func TestTrustStoreWithMultipleEntries(t *testing.T) {
	ts, err := truststore.New("test-multi.pem")
	assert.Nil(t, err, "")
	assert.NotNil(t, ts, "")
	pubkey := ts.GetPublicKey("test.localnet")
	assert.NotNil(t, pubkey, "")
	_, ok := pubkey.(*rsa.PublicKey)
	assert.True(t, ok, "")
	pubkey = ts.GetPublicKey("foo")
	assert.NotNil(t, pubkey, "")
	_, ok = pubkey.(*rsa.PublicKey)
	assert.True(t, ok, "")
}
