package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

const (
	RSAPrivateKey = "RSA PRIVATE KEY"
	RSAPublicKey  = "RSA PUBLIC KEY"
)

func GenerateKeyPair() ([]byte, []byte, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	publicKey := privateKey.Public()

	privateKeyBytes := pem.EncodeToMemory(
		&pem.Block{Type: RSAPrivateKey, Bytes: x509.MarshalPKCS1PrivateKey(privateKey)},
	)

	pkixPublicKey, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, nil, err
	}
	publicKeyBytes := pem.EncodeToMemory(
		&pem.Block{Type: RSAPublicKey, Bytes: pkixPublicKey},
	)
	return privateKeyBytes, publicKeyBytes, nil
}
