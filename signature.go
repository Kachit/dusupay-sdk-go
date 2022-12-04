package dusupay

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

//IncomingWebhookInterface interface
type IncomingWebhookInterface interface {
	BuildPayloadString(url string) string
}

func NewSignatureValidator(publicKeyBytes []byte) (*SignatureValidator, error) {
	block, _ := pem.Decode(publicKeyBytes)
	publicKey, err := parsePublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return &SignatureValidator{publicKey}, nil
}

type SignatureValidator struct {
	publicKey *rsa.PublicKey
}

func (sv *SignatureValidator) ValidateSignature(webhook IncomingWebhookInterface, webhookUrl string, signature string) error {
	messageBytes := bytes.NewBufferString(webhook.BuildPayloadString(webhookUrl))
	hash := sha512.New()
	hash.Write(messageBytes.Bytes())
	digest := hash.Sum(nil)

	data, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return err
	}
	return rsa.VerifyPKCS1v15(sv.publicKey, crypto.SHA512, digest, data)
}

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
