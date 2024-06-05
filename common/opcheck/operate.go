package opcheck

import (
	"crypto/sha256"
	"encoding/hex"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
)

// SignMessage signs a message using the provided private key
func SignMessage(message string, privKey cryptotypes.PrivKey) (string, error) {
	hash := sha256.Sum256([]byte(message))
	signature, err := privKey.Sign(hash[:])
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(signature), nil
}

// VerifySignature verifies a signed message using the provided public key
func VerifySignature(message, signatureHex string, pubKey cryptotypes.PubKey) (bool, error) {
	hash := sha256.Sum256([]byte(message))
	signature, err := hex.DecodeString(signatureHex)
	if err != nil {
		return false, err
	}
	return pubKey.VerifySignature(hash[:], signature), nil
}
