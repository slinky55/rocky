package logic

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

func GenerateChallenge() ([]byte, error) {
	challenge := []byte("challenge")
	return challenge, nil
}

func GenerateSessionToken() (string, error) {
	challenge := make([]byte, 32)
	_, err := rand.Read(challenge)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(challenge), nil
}

func VerifyChallenge(challenge []byte, signature []byte, publicKey string) error {
	pub, err := loadRsaPublicKey(publicKey)
	if err != nil {
		return err
	}

	err = rsa.VerifyPSS(pub, crypto.SHA256, challenge, signature, nil)
	if err != nil {
		return err
	}

	return nil
}

func loadRsaPublicKey(publicKey string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		return nil, errors.New("failed to decode PEM block containing public key")
	}

	key, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	return key, nil
}
