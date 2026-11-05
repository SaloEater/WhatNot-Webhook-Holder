package discord

import (
	"crypto/ed25519"
	"encoding/hex"
)

// VerifyInteraction checks the Ed25519 signature Discord puts on every interaction (over timestamp + raw body).
func VerifyInteraction(publicKeyHex, signatureHex, timestamp string, body []byte) bool {
	publicKey, err := hex.DecodeString(publicKeyHex)
	if err != nil || len(publicKey) != ed25519.PublicKeySize {
		return false
	}
	signature, err := hex.DecodeString(signatureHex)
	if err != nil || len(signature) != ed25519.SignatureSize {
		return false
	}
	msg := make([]byte, 0, len(timestamp)+len(body))
	msg = append(msg, timestamp...)
	msg = append(msg, body...)
	return ed25519.Verify(publicKey, msg, signature)
}
