package key

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/pkg/errors"
	"golang.org/x/crypto/curve25519"
)

const keyLength = 32

// PrivateKeyToPublicKey generates wireguard public key from private key
func PrivateKeyToPublicKey(key string) (string, error) {
	k, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return "", err
	}

	var pub [keyLength]byte
	var priv [keyLength]byte
	copy(priv[:], k[:keyLength])
	curve25519.ScalarBaseMult(&pub, &priv)

	return base64.StdEncoding.EncodeToString(pub[:]), nil
}

// GeneratePrivateKey generates a private key
func GeneratePrivateKey() (string, error) {
	randomBytes := make([]byte, keyLength)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", errors.Wrapf(err, "failed to generate random bytes for private key")
	}

	// https://cr.yp.to/ecdh.html
	randomBytes[0] &= 248
	randomBytes[31] &= 127
	randomBytes[31] |= 64

	return base64.StdEncoding.EncodeToString(randomBytes), nil
}
