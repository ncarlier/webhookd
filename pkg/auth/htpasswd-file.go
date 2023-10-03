package auth

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/csv"
	"net/http"
	"os"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

var (
	shaRe = regexp.MustCompile(`^{SHA}`)
	bcrRe = regexp.MustCompile(`^\$2b\$|^\$2a\$|^\$2y\$`)
)

// HtpasswdFile is a map for usernames to passwords.
type HtpasswdFile struct {
	path  string
	users map[string]string
}

// NewHtpasswdFromFile reads the users and passwords from a htpasswd file and returns them.
func NewHtpasswdFromFile(path string) (*HtpasswdFile, error) {
	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	cr := csv.NewReader(r)
	cr.Comma = ':'
	cr.Comment = '#'
	cr.TrimLeadingSpace = true

	records, err := cr.ReadAll()
	if err != nil {
		return nil, err
	}

	users := make(map[string]string)
	for _, record := range records {
		users[record[0]] = record[1]
	}

	return &HtpasswdFile{
		path:  path,
		users: users,
	}, nil
}

// Validate HTTP request credentials
func (h *HtpasswdFile) Validate(r *http.Request) (ok bool, username string) {
	username, password, ok := r.BasicAuth()
	ok = ok && h.validateCredentials(username, password)
	return
}

func (h *HtpasswdFile) validateCredentials(user, password string) bool {
	pwd, exists := h.users[user]
	if !exists {
		return false
	}

	switch {
	case shaRe.MatchString(pwd):
		d := sha1.New()
		_, _ = d.Write([]byte(password))
		if pwd[5:] == base64.StdEncoding.EncodeToString(d.Sum(nil)) {
			return true
		}
	case bcrRe.MatchString(pwd):
		err := bcrypt.CompareHashAndPassword([]byte(pwd), []byte(password))
		if err == nil {
			return true
		}
	}
	return false
}
