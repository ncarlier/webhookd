package signature

import (
	"fmt"
	"net/http"

	"github.com/go-fed/httpsig"
	"github.com/ncarlier/webhookd/pkg/truststore"
)

// HTTPSignatureHandler validate request HTTP signature
func HTTPSignatureHandler(r *http.Request, ts truststore.TrustStore) error {
	verifier, err := httpsig.NewVerifier(r)
	if err != nil {
		return err
	}
	pubkeyID := verifier.KeyId()
	pubkey := ts.GetPublicKey(pubkeyID)
	if pubkey == nil {
		return fmt.Errorf("public key not found: %s", pubkeyID)
	}
	// TODO dynamic algo
	err = verifier.Verify(pubkey, httpsig.RSA_SHA256)
	if err != nil {
		return err
	}
	return nil
}
