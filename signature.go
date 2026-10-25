package dusupay

import (
	"bytes"
	"crypto"
	"crypto/hmac"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
)

// NewSignatureValidator creates an RSA signature validator from a PEM public key
func NewSignatureValidator(publicKeyBytes []byte) (*SignatureValidator, error) {
	block, _ := pem.Decode(publicKeyBytes)
	if block == nil {
		return nil, errors.New("wrong public key data")
	}

	publicKey, err := parsePublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return &SignatureValidator{publicKey}, nil
}

// SignatureValidator validates RSA-SHA256 webhook signatures
type SignatureValidator struct {
	publicKey *rsa.PublicKey
}

// ValidateSignature verifies the rsa-signature header.
// See https://developer.dusupay.com/callbacks/rsa-signature-verification.md
// Payload format: event:merchant_reference:internal_reference:transaction_type:transaction_status
func (sv *SignatureValidator) ValidateSignature(payload CallbackPayloadInterface, signature string) error {
	messageBytes := bytes.NewBufferString(payload.BuildPayloadString())
	hash := sha256.New()
	hash.Write(messageBytes.Bytes())
	digest := hash.Sum(nil)

	data, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return err
	}

	return rsa.VerifyPKCS1v15(sv.publicKey, crypto.SHA256, digest, data)
}

// parsePublicKey parses a DER-encoded RSA public key
func parsePublicKey(rawBytes []byte) (*rsa.PublicKey, error) {
	key, err := x509.ParsePKIXPublicKey(rawBytes)
	if err != nil {
		return nil, fmt.Errorf("parsePublicKey wrong parse PKIX: %v", err)
	}

	switch pk := key.(type) {
	case *rsa.PublicKey:
		return pk, nil
	default:
		return nil, errors.New("parsePublicKey: PublicKey must be of type rsa.PublicKey")
	}
}

// NewHMACSignatureValidator creates an HMAC-SHA256 signature validator
func NewHMACSignatureValidator(signingKey string) *HMACSignatureValidator {
	return &HMACSignatureValidator{signingKey: signingKey}
}

// HMACSignatureValidator validates HMAC-SHA256 webhook signatures
type HMACSignatureValidator struct {
	signingKey string
}

// ValidateSignature verifies the hmac-signature header.
// See https://developer.dusupay.com/callbacks/hmac-signature-verification.md
// Signature format: t=timestamp,s=hmac_hash
// Payload format: event:merchant_reference:internal_reference:transaction_type:transaction_status
func (hv *HMACSignatureValidator) ValidateSignature(payload CallbackPayloadInterface, signature string) error {
	s, err := parseHMACHash(signature)
	if err != nil {
		return err
	}

	mac := hmac.New(sha256.New, []byte(hv.signingKey))
	mac.Write([]byte(payload.BuildPayloadString()))

	calculated := hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(calculated), []byte(s)) {
		return errors.New("HMAC signature mismatch")
	}

	return nil
}

func parseHMACHash(signature string) (string, error) {
	for _, part := range strings.Split(signature, ",") {
		if strings.HasPrefix(part, "s=") {
			return strings.TrimPrefix(part, "s="), nil
		}
	}

	return "", errors.New("malformed HMAC signature: missing 's' field")
}
