package test

import (
	"crypto/rsa"
	"testing"

	"github.com/ncarlier/webhookd/pkg/assert"
	"github.com/ncarlier/webhookd/pkg/truststore"
)

func TestTrustStoreWithP12(t *testing.T) {
	t.Skip()

	ts, err := truststore.New("test.p12")
	assert.Nil(t, err, "")
	assert.NotNil(t, ts, "")
	pubkey := ts.GetPublicKey("test.localnet")
	assert.NotNil(t, pubkey, "")
	_, ok := pubkey.(*rsa.PublicKey)
	assert.True(t, ok, "")
}
