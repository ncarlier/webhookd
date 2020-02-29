package middleware

import (
	"net/http"

	"github.com/go-fed/httpsig"
	"github.com/ncarlier/webhookd/pkg/pubkey"
)

// HTTPSignature is a middleware to checks HTTP request signature
func HTTPSignature(keyStore pubkey.KeyStore) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			verifier, err := httpsig.NewVerifier(r)
			if err != nil {
				w.WriteHeader(400)
				w.Write([]byte("invalid HTTP signature: " + err.Error()))
				return
			}
			pubKeyID := verifier.KeyId()
			pubKey, algo, err := keyStore.Get(pubKeyID)
			if err != nil {
				w.WriteHeader(400)
				w.Write([]byte("invalid HTTP signature: " + err.Error()))
				return
			}
			err = verifier.Verify(pubKey, algo)
			if err != nil {
				w.WriteHeader(400)
				w.Write([]byte("invalid HTTP signature: " + err.Error()))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
