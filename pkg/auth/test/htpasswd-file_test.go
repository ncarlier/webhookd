package test

import (
	"net/http"
	"testing"

	"github.com/ncarlier/webhookd/pkg/assert"
	"github.com/ncarlier/webhookd/pkg/auth"
)

func TestValidateCredentials(t *testing.T) {
	htpasswdFile, err := auth.NewHtpasswdFromFile("test.htpasswd")
	assert.Nil(t, err, ".htpasswd file should be loaded")
	assert.NotNil(t, htpasswdFile, ".htpasswd file should be loaded")

	req, err := http.NewRequest("POST", "http://localhost:8080", nil)
	assert.Nil(t, err, "")
	req.SetBasicAuth("foo", "bar")
	assert.Equal(t, true, htpasswdFile.Validate(req), "credentials should be valid")

	req.SetBasicAuth("foo", "bad")
	assert.Equal(t, false, htpasswdFile.Validate(req), "credentials should not be valid")
}
