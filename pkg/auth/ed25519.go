package auth

import (
	"crypto/ed25519"
	"encoding/base64"
	"time"

	"github.com/ez-as/ironlink-base-lib/pkg/logging"
)

const SignatureTimeout = 1 * time.Minute

func ED25519Verify(publicKey, message, signature string) bool {
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

func ED25519Sign(privateKey, message string) string {
	privateKeyBytes, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return ""
	}
	if len(privateKeyBytes) != ed25519.PrivateKeySize {
		return ""
	}
	return base64.StdEncoding.EncodeToString(ed25519.Sign(privateKeyBytes, []byte(message)))
}

type ED25519Verifier struct {
	masterPublicKey ed25519.PublicKey
}

func NewED25519Verifier(masterPublicKey string) *ED25519Verifier {
	masterPubKey, _ := base64.StdEncoding.DecodeString(masterPublicKey)
	if len(masterPubKey) != ed25519.PublicKeySize {
		logging.Logger().Fatal("invalid master public key")
	}

	return &ED25519Verifier{
		masterPublicKey: masterPubKey,
	}
}

func (v *ED25519Verifier) Verify(message, signature string) bool {
	signatureBytes, _ := base64.StdEncoding.DecodeString(signature)
	if len(signatureBytes) != ed25519.SignatureSize {
		return false
	}
	return ed25519.Verify(v.masterPublicKey, []byte(message), signatureBytes)
}

type ED25519Signer struct {
	masterPrivateKey ed25519.PrivateKey
}

func NewED25519Signer(masterPrivateKey string) *ED25519Signer {
	masterSeed, _ := base64.StdEncoding.DecodeString(masterPrivateKey)
	if len(masterSeed) != ed25519.SeedSize {
		logging.Logger().Fatal("invalid master private key")
	}

	return &ED25519Signer{
		masterPrivateKey: ed25519.NewKeyFromSeed(masterSeed),
	}
}

func (s *ED25519Signer) Sign(message string) string {
	return base64.StdEncoding.EncodeToString(ed25519.Sign(s.masterPrivateKey, []byte(message)))
}

func (s *ED25519Signer) Verify(message, signature string) bool {
	signatureBytes, _ := base64.StdEncoding.DecodeString(signature)
	if len(signatureBytes) != ed25519.SignatureSize {
		return false
	}
	return ed25519.Verify(s.masterPrivateKey.Public().(ed25519.PublicKey), []byte(message), signatureBytes)
}
