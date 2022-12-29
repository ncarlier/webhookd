package middleware

import (
	"net/http"

	"github.com/ncarlier/webhookd/pkg/middleware/signature"
	"github.com/ncarlier/webhookd/pkg/truststore"
)

// Signature is a middleware to checks HTTP request signature
func Signature(ts truststore.TrustStore) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler := signature.HTTPSignatureHandler
			if signature.IsEd25519SignatureRequest(r.Header) {
				handler = signature.Ed25519SignatureHandler
			}
			if err := handler(r, ts); err != nil {
				w.WriteHeader(401)
				w.Write([]byte("401 Unauthorized: " + err.Error()))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
