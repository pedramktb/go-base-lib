package auth

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/ssh"
)

// ED25519Verify verifies a message with a public key and a signature and returns true if the signature is valid
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

// ED25519Sign signs a message with a private key and returns the base64 encoded signature
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

// ED25519Verifier is a struct that verifies messages with a master public key
type ED25519Verifier struct {
	masterPublicKey ed25519.PublicKey
}

// NewED25519Verifier creates a new ED25519Verifier with a master public key
func NewED25519Verifier(masterPublicKey string) (*ED25519Verifier, error) {
	masterPubKey, _ := base64.StdEncoding.DecodeString(masterPublicKey)
	if len(masterPubKey) != ed25519.PublicKeySize {
		return nil, fmt.Errorf("invalid master public key size")
	}

	return &ED25519Verifier{
		masterPublicKey: masterPubKey,
	}, nil
}

// Verify verifies a message with a signature using the master public key and returns true if the signature is valid
func (v *ED25519Verifier) Verify(message, signature string) bool {
	signatureBytes, _ := base64.StdEncoding.DecodeString(signature)
	if len(signatureBytes) != ed25519.SignatureSize {
		return false
	}
	return ed25519.Verify(v.masterPublicKey, []byte(message), signatureBytes)
}

// ED25519Signer is a struct that signs messages with a master private key
type ED25519Signer struct {
	masterPrivateKey ed25519.PrivateKey
	sshSigner        ssh.Signer
}

// NewED25519Signer creates a new ED25519Signer (and ssh.Signer) with a master private key
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

// PublicKey returns the master private key's public key as a base64 encoded string
func (s *ED25519Signer) PublicKey() string {
	return base64.StdEncoding.EncodeToString(s.masterPrivateKey.Public().(ed25519.PublicKey))
}

// Sign signs a message with the master private key and returns the base64 encoded signature
func (s *ED25519Signer) Sign(message string) string {
	return base64.StdEncoding.EncodeToString(ed25519.Sign(s.masterPrivateKey, []byte(message)))
}

// Verify verifies a message with a signature using the master private key's public key
func (s *ED25519Signer) Verify(message, signature string) bool {
	signatureBytes, _ := base64.StdEncoding.DecodeString(signature)
	if len(signatureBytes) != ed25519.SignatureSize {
		return false
	}
	return ed25519.Verify(s.masterPrivateKey.Public().(ed25519.PublicKey), []byte(message), signatureBytes)
}

// SSHSigner returns the ssh.Signer based on the master private key
func (s *ED25519Signer) SSHSigner() ssh.Signer {
	return s.sshSigner
}
