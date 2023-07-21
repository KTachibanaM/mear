package ssh

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"golang.org/x/crypto/ssh"
)

// Generate a new RSA key pair
// and return the private key in PEM format and the public key in OpenSSH format
func Keygen() ([]byte, []byte, error) {
	// Generate a new RSA key pair
	private_key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate private key: %v", err)
	}

	// Encode the private key in PEM format
	private_key_pem := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(private_key),
	}

	// Encode the public key in OpenSSH format
	public_key, err := ssh.NewPublicKey(&private_key.PublicKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate public key: %v", err)
	}

	return pem.EncodeToMemory(private_key_pem), ssh.MarshalAuthorizedKey(public_key), nil
}
