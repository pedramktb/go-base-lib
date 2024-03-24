package auth

import (
	"crypto/ed25519"
	"encoding/base64"
	"time"
)

const AuthTimeout = 1 * time.Minute

func Verify(publicKey, message, signature string) bool {
	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return false
	}
	if len(publicKeyBytes) != ed25519.PublicKeySize {
		return false
	}
	signatureBytes, _ := base64.StdEncoding.DecodeString(signature)
	if len(signatureBytes) != ed25519.SignatureSize {
		return false
	}
	return ed25519.Verify(publicKeyBytes, []byte(message), signatureBytes)
}
