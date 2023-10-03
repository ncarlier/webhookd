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

	req, err := http.NewRequest("POST", "http://localhost:8080", http.NoBody)
	assert.Nil(t, err, "")
	req.SetBasicAuth("foo", "bar")
	ok, username := htpasswdFile.Validate(req)
	assert.Equal(t, true, ok, "credentials should be valid")
	assert.Equal(t, "foo", username, "invalid username")

	req.SetBasicAuth("foo", "bad")
	ok, username = htpasswdFile.Validate(req)
	assert.Equal(t, false, ok, "credentials should be invalid")
	assert.Equal(t, "foo", username, "invalid username")
}
