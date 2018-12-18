package auth

import (
	"testing"

	"github.com/ncarlier/webhookd/pkg/assert"
)

func TestValidateCredentials(t *testing.T) {
	htpasswdFile, err := NewHtpasswdFromFile("test.htpasswd")
	assert.Nil(t, err, ".htpasswd file should be loaded")
	assert.NotNil(t, htpasswdFile, ".htpasswd file should be loaded")
	assert.Equal(t, true, htpasswdFile.validateCredentials("foo", "bar"), "credentials should be valid")
	assert.Equal(t, false, htpasswdFile.validateCredentials("foo", "bir"), "credentials should not be valid")
}
