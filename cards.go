package dusupay

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
)

// CardData card details for direct card payments (S2S).
// See https://developer.dusupay.com/funds-collection/card-payments/direct-card-payment-s2s.md
// Billing fields are required only for cards issued in the USA, Canada and Europe.
type CardData struct {
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	CardNo         string `json:"card_no"`
	ExpMonth       string `json:"exp_month"`
	ExpYear        string `json:"exp_year"`
	CVV            string `json:"cvv"`
	BillingAddress string `json:"billing_address,omitempty"`
	BillingCity    string `json:"billing_city,omitempty"`
	BillingZip     string `json:"billing_zip,omitempty"`
	BillingState   string `json:"billing_state,omitempty"`
	BillingCountry string `json:"billing_country,omitempty"`
}

func (cd *CardData) isValid() error {
	switch {
	case cd.FirstName == "":
		return fmt.Errorf(`parameter "first_name" is empty`)
	case cd.LastName == "":
		return fmt.Errorf(`parameter "last_name" is empty`)
	case cd.CardNo == "":
		return fmt.Errorf(`parameter "card_no" is empty`)
	case cd.ExpMonth == "":
		return fmt.Errorf(`parameter "exp_month" is empty`)
	case cd.ExpYear == "":
		return fmt.Errorf(`parameter "exp_year" is empty`)
	case cd.CVV == "":
		return fmt.Errorf(`parameter "cvv" is empty`)
	}

	return nil
}

// NewCardEncryptor creates a card data encryptor from the DusuPay PEM public key
// (see Config.DusupayPublicKey, SandboxPublicKey, ProdPublicKey)
func NewCardEncryptor(publicKeyBytes []byte) (*CardEncryptor, error) {
	block, _ := pem.Decode(publicKeyBytes)
	if block == nil {
		return nil, errors.New("wrong public key data")
	}

	publicKey, err := parsePublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return &CardEncryptor{publicKey}, nil
}

// CardEncryptor encrypts card data into the card_cipher collection parameter
type CardEncryptor struct {
	publicKey *rsa.PublicKey
}

// Encrypt returns base64 RSA (PKCS#1 v1.5) encrypted card JSON for CollectionRequest.CardCipher
func (ce *CardEncryptor) Encrypt(card *CardData) (string, error) {
	if err := card.isValid(); err != nil {
		return "", fmt.Errorf("CardEncryptor.Encrypt error: %v", err)
	}

	data, err := json.Marshal(card)
	if err != nil {
		return "", fmt.Errorf("CardEncryptor.Encrypt error: %v", err)
	}

	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, ce.publicKey, data)
	if err != nil {
		return "", fmt.Errorf("CardEncryptor.Encrypt error: %v", err)
	}

	return base64.StdEncoding.EncodeToString(encrypted), nil
}
