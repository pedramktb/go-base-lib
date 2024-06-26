package auth

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"
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

func NewED25519Verifier(masterPublicKey string) (*ED25519Verifier, error) {
	masterPubKey, _ := base64.StdEncoding.DecodeString(masterPublicKey)
	if len(masterPubKey) != ed25519.PublicKeySize {
		return nil, fmt.Errorf("invalid master public key size")
	}

	return &ED25519Verifier{
		masterPublicKey: masterPubKey,
	}, nil
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
	sshSigner        ssh.Signer
}

func NewED25519Signer(masterPrivateKey string) (*ED25519Signer, error) {
	masterSeed, _ := base64.StdEncoding.DecodeString(masterPrivateKey)
	if len(masterSeed) != ed25519.SeedSize {
		return nil, fmt.Errorf("invalid master private key size")
	}

	key := ed25519.NewKeyFromSeed(masterSeed)

	sshSigner, err := ssh.NewSignerFromSigner(key)
	if err != nil {
		return nil, fmt.Errorf("error creating ssh signer: %v", err)
	}

	return &ED25519Signer{
		masterPrivateKey: key,
		sshSigner:        sshSigner,
	}, nil
}

func (s *ED25519Signer) PublicKey() string {
	return base64.StdEncoding.EncodeToString(s.masterPrivateKey.Public().(ed25519.PublicKey))
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

func (s *ED25519Signer) SSHSigner() ssh.Signer {
	return s.sshSigner
}
