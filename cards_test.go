package dusupay

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CardEncryptorTestSuite struct {
	suite.Suite
}

func buildStubCardData() *CardData {
	return &CardData{
		FirstName:      "John",
		LastName:       "Doe",
		CardNo:         "4622000000005678",
		ExpMonth:       "06",
		ExpYear:        "22",
		CVV:            "123",
		BillingAddress: "Second Street",
		BillingCity:    "San Francisco",
		BillingZip:     "94105",
		BillingState:   "CA",
		BillingCountry: "US",
	}
}

func (suite *CardEncryptorTestSuite) TestEncryptRoundTrip() {
	privateKeyBytes, err := os.ReadFile("main/private.pem")
	if err != nil {
		suite.T().Skip("main/private.pem not available")

		return
	}
	block, _ := pemDecode(privateKeyBytes)
	privateKey, _ := parsePrivateKeyForTest(block.Bytes)

	rawBytes, _ := os.ReadFile("stubs/rsa/public-key.pem")
	encryptor, err := NewCardEncryptor(rawBytes)
	assert.NoError(suite.T(), err)

	cipher, err := encryptor.Encrypt(buildStubCardData())
	assert.NoError(suite.T(), err)

	encrypted, _ := base64.StdEncoding.DecodeString(cipher)
	decrypted, err := rsa.DecryptPKCS1v15(nil, privateKey, encrypted)
	assert.NoError(suite.T(), err)

	var card CardData
	_ = json.Unmarshal(decrypted, &card)
	assert.Equal(suite.T(), *buildStubCardData(), card)
}

func (suite *CardEncryptorTestSuite) TestEncryptWithDusupayKeys() {
	for _, key := range []string{SandboxPublicKey, ProdPublicKey} {
		encryptor, err := NewCardEncryptor([]byte(key))
		assert.NoError(suite.T(), err)
		cipher, err := encryptor.Encrypt(buildStubCardData())
		assert.NoError(suite.T(), err)
		assert.NotEmpty(suite.T(), cipher)
	}
}

func (suite *CardEncryptorTestSuite) TestEncryptInvalidCard() {
	encryptor, _ := NewCardEncryptor([]byte(SandboxPublicKey))
	card := buildStubCardData()
	card.CVV = ""
	cipher, err := encryptor.Encrypt(card)
	assert.EqualError(suite.T(), err, `CardEncryptor.Encrypt error: parameter "cvv" is empty`)
	assert.Empty(suite.T(), cipher)
}

func (suite *CardEncryptorTestSuite) TestCardDataIsValid() {
	cases := map[string]func(c *CardData){
		"first_name": func(c *CardData) { c.FirstName = "" },
		"last_name":  func(c *CardData) { c.LastName = "" },
		"card_no":    func(c *CardData) { c.CardNo = "" },
		"exp_month":  func(c *CardData) { c.ExpMonth = "" },
		"exp_year":   func(c *CardData) { c.ExpYear = "" },
		"cvv":        func(c *CardData) { c.CVV = "" },
	}
	for name, mutate := range cases {
		card := buildStubCardData()
		mutate(card)
		assert.EqualError(suite.T(), card.isValid(), `parameter "`+name+`" is empty`)
	}
	assert.NoError(suite.T(), buildStubCardData().isValid())
}

func (suite *CardEncryptorTestSuite) TestNewCardEncryptorError() {
	encryptor, err := NewCardEncryptor([]byte(`foo`))
	assert.EqualError(suite.T(), err, "wrong public key data")
	assert.Nil(suite.T(), encryptor)
}

func TestCardEncryptorTestSuite(t *testing.T) {
	suite.Run(t, new(CardEncryptorTestSuite))
}
