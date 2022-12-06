package dusupay

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"testing"
)

type SignatureTestSuite struct {
	suite.Suite
}

func (suite *SignatureTestSuite) TestValidateSignatureCollectionWebhookSuccess() {
	rawBytes, _ := ioutil.ReadFile("stubs/rsa/public-key.pem")
	webhook := &CollectionWebhook{ID: 226, InternalReference: "DUSUPAY405GZM1G5JXGA71IK", TransactionStatus: "COMPLETED"}
	validator, _ := NewSignatureValidator(rawBytes)
	err := validator.ValidateSignature(webhook, "https://www.sample-url.com/callback", stubSignature)
	assert.NoError(suite.T(), err)
}

func (suite *SignatureTestSuite) TestValidateSignatureCollectionWebhookWrongPayload() {
	rawBytes, _ := ioutil.ReadFile("stubs/rsa/public-key.pem")
	webhook := &CollectionWebhook{ID: 225, InternalReference: "DUSUPAY405GZM1G5JXGA71IK", TransactionStatus: "COMPLETED"}
	validator, _ := NewSignatureValidator(rawBytes)
	err := validator.ValidateSignature(webhook, "https://www.sample-url.com/callback", stubSignature)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "crypto/rsa: verification error", err.Error())
}

func (suite *SignatureTestSuite) TestValidateSignatureCollectionWebhookNonBase64Signature() {
	rawBytes, _ := ioutil.ReadFile("stubs/rsa/public-key.pem")
	webhook := &CollectionWebhook{ID: 225, InternalReference: "DUSUPAY405GZM1G5JXGA71IK", TransactionStatus: "COMPLETED"}
	validator, _ := NewSignatureValidator(rawBytes)
	err := validator.ValidateSignature(webhook, "https://www.sample-url.com/callback", "qwerty")
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "illegal base64 data at input byte 4", err.Error())
}

func (suite *SignatureTestSuite) TestValidateSignaturePayoutWebhookSuccess() {
	rawBytes, _ := ioutil.ReadFile("stubs/rsa/public-key.pem")
	webhook := &PayoutWebhook{ID: 226, InternalReference: "DUSUPAY405GZM1G5JXGA71IK", TransactionStatus: "COMPLETED"}
	validator, _ := NewSignatureValidator(rawBytes)
	err := validator.ValidateSignature(webhook, "https://www.sample-url.com/callback", stubSignature)
	assert.NoError(suite.T(), err)
}

func (suite *SignatureTestSuite) TestValidateSignatureRefundWebhookSuccess() {
	rawBytes, _ := ioutil.ReadFile("stubs/rsa/public-key.pem")
	webhook := &RefundWebhook{ID: 226, InternalReference: "DUSUPAY405GZM1G5JXGA71IK", TransactionStatus: "COMPLETED"}
	validator, _ := NewSignatureValidator(rawBytes)
	err := validator.ValidateSignature(webhook, "https://www.sample-url.com/callback", stubSignature)
	assert.NoError(suite.T(), err)
}

func (suite *SignatureTestSuite) TestParsePublicKeyError() {
	pk, err := parsePublicKey([]byte(`foo`))
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "parsePublicKey wrong parse PKIX: asn1: structure error: tags don't match")
	assert.Nil(suite.T(), pk)
}

func (suite *SignatureTestSuite) TestNewSignatureValidatorError() {
	validator, err := NewSignatureValidator([]byte(`foo`))
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "wrong public key data")
	assert.Nil(suite.T(), validator)
}

func TestSignatureTestSuite(t *testing.T) {
	suite.Run(t, new(SignatureTestSuite))
}
