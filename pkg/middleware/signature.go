package middleware

import (
	"net/http"

	"github.com/go-fed/httpsig"
	"github.com/ncarlier/webhookd/pkg/pubkey"
)

// HTTPSignature is a middleware to checks HTTP request signature
func HTTPSignature(trustStore pubkey.TrustStore) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			verifier, err := httpsig.NewVerifier(r)
			if err != nil {
				w.WriteHeader(400)
				w.Write([]byte("invalid HTTP signature: " + err.Error()))
				return
			}
			pubKeyID := verifier.KeyId()
			entry := trustStore.Get(pubKeyID)
			if entry == nil {
				w.WriteHeader(400)
				w.Write([]byte("invalid HTTP signature: public key not found: " + pubKeyID))
				return
			}
			err = verifier.Verify(entry.Pubkey, entry.Algorithm)
			if err != nil {
				w.WriteHeader(400)
				w.Write([]byte("invalid HTTP signature: " + err.Error()))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
