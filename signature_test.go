package dusupay

import (
	"crypto"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SignatureTestSuite struct {
	suite.Suite
}

func (suite *SignatureTestSuite) TestValidateSignatureCollectionCallbackSuccess() {
	rawBytes, _ := os.ReadFile("stubs/rsa/public-key.pem")
	payload := &CollectionCallbackPayload{
		Event:             "transaction.completed",
		MerchantReference: "76859aae-f148-48c5-9901-2e474cf19b71",
		InternalReference: "DUSUPAY405GZM1G5JXGA71IK",
		TransactionType:   "COLLECTION",
		TransactionStatus: "COMPLETED",
	}
	validator, err := NewSignatureValidator(rawBytes)
	assert.NoError(suite.T(), err)
	assert.NoError(suite.T(), validator.ValidateSignature(payload, stubSignature))
}

func (suite *SignatureTestSuite) TestValidateSignatureWrongPayload() {
	rawBytes, _ := os.ReadFile("stubs/rsa/public-key.pem")
	payload := &CollectionCallbackPayload{
		Event:             "transaction.completed",
		MerchantReference: "wrong-reference",
		InternalReference: "DUSUPAY405GZM1G5JXGA71IK",
		TransactionType:   "COLLECTION",
		TransactionStatus: "COMPLETED",
	}
	validator, _ := NewSignatureValidator(rawBytes)
	err := validator.ValidateSignature(payload, stubSignature)
	assert.EqualError(suite.T(), err, "crypto/rsa: verification error")
}

func (suite *SignatureTestSuite) TestValidateSignatureNonBase64() {
	rawBytes, _ := os.ReadFile("stubs/rsa/public-key.pem")
	payload := &CollectionCallbackPayload{Event: "transaction.completed", MerchantReference: "ref", InternalReference: "int", TransactionType: "COLLECTION", TransactionStatus: "COMPLETED"}
	validator, _ := NewSignatureValidator(rawBytes)
	err := validator.ValidateSignature(payload, "qwerty")
	assert.EqualError(suite.T(), err, "illegal base64 data at input byte 4")
}

func (suite *SignatureTestSuite) TestValidateSignaturePayoutCallbackSuccess() {
	rawBytes, _ := os.ReadFile("stubs/rsa/public-key.pem")
	privateKeyBytes, err := os.ReadFile("main/private.pem")
	if err != nil {
		suite.T().Skip("main/private.pem not available")

		return
	}
	block, _ := pemDecode(privateKeyBytes)
	privateKey, _ := parsePrivateKeyForTest(block.Bytes)

	payload := &PayoutCallbackPayload{
		Event:             "transaction.completed",
		MerchantReference: "76859aae-f148-48c5-9901-2e474cf19b72",
		InternalReference: "DUSUPAY405GZM1G5JXGA71IZ",
		TransactionType:   "PAYOUT",
		TransactionStatus: "COMPLETED",
	}
	hash := sha256.New()
	hash.Write([]byte(payload.BuildPayloadString()))
	digest := hash.Sum(nil)
	sig, _ := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, digest)
	signature := base64.StdEncoding.EncodeToString(sig)

	validator, _ := NewSignatureValidator(rawBytes)
	assert.NoError(suite.T(), validator.ValidateSignature(payload, signature))
}

func (suite *SignatureTestSuite) TestParsePublicKeyError() {
	pk, err := parsePublicKey([]byte(`foo`))
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "parsePublicKey wrong parse PKIX")
	assert.Nil(suite.T(), pk)
}

func (suite *SignatureTestSuite) TestNewSignatureValidatorError() {
	validator, err := NewSignatureValidator([]byte(`foo`))
	assert.EqualError(suite.T(), err, "wrong public key data")
	assert.Nil(suite.T(), validator)
}

func TestSignatureTestSuite(t *testing.T) {
	suite.Run(t, new(SignatureTestSuite))
}

type HMACSignatureTestSuite struct {
	suite.Suite
}

func (suite *HMACSignatureTestSuite) TestValidateSignatureSuccess() {
	payload := &CollectionCallbackPayload{
		Event:             "transaction.completed",
		MerchantReference: "76859aae-f148-48c5-9901-2e474cf19b71",
		InternalReference: "DUSUPAY405GZM1G5JXGA71IK",
		TransactionType:   "COLLECTION",
		TransactionStatus: "COMPLETED",
	}
	key := "test-signing-key"
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(payload.BuildPayloadString()))
	hashVal := hex.EncodeToString(mac.Sum(nil))
	signature := "t=1234567890,s=" + hashVal

	validator := NewHMACSignatureValidator(key)
	assert.NoError(suite.T(), validator.ValidateSignature(payload, signature))
}

func (suite *HMACSignatureTestSuite) TestValidateSignatureMismatch() {
	payload := &CollectionCallbackPayload{Event: "transaction.completed", MerchantReference: "ref", InternalReference: "int", TransactionType: "COLLECTION", TransactionStatus: "COMPLETED"}
	validator := NewHMACSignatureValidator("test-signing-key")
	err := validator.ValidateSignature(payload, "t=123,s=wronghash")
	assert.EqualError(suite.T(), err, "HMAC signature mismatch")
}

func (suite *HMACSignatureTestSuite) TestValidateSignatureMalformed() {
	payload := &CollectionCallbackPayload{Event: "transaction.completed", MerchantReference: "ref", InternalReference: "int", TransactionType: "COLLECTION", TransactionStatus: "COMPLETED"}
	validator := NewHMACSignatureValidator("test-signing-key")
	err := validator.ValidateSignature(payload, "t=123")
	assert.EqualError(suite.T(), err, "malformed HMAC signature: missing 's' field")
}

func TestHMACSignatureTestSuite(t *testing.T) {
	suite.Run(t, new(HMACSignatureTestSuite))
}
