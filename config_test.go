package dusupay

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
	testable *Config
}

func (suite *ConfigTestSuite) SetupTest() {
	suite.testable = NewConfig("foo", "bar")
}

func (suite *ConfigTestSuite) TestNewConfigByDefault() {
	assert.Equal(suite.T(), ProdAPIUrl, suite.testable.URI)
	assert.Equal(suite.T(), "foo", suite.testable.PublicKey)
	assert.Equal(suite.T(), "bar", suite.testable.SecretKey)
}

func (suite *ConfigTestSuite) TestNewConfigSandbox() {
	result := NewConfigSandbox("foo", "bar")
	assert.Equal(suite.T(), SandboxAPIUrl, result.URI)
	assert.Equal(suite.T(), "foo", result.PublicKey)
	assert.Equal(suite.T(), "bar", result.SecretKey)
}

func (suite *ConfigTestSuite) TestIsSandbox() {
	assert.False(suite.T(), suite.testable.IsSandbox())
	suite.testable.URI = SandboxAPIUrl
	assert.True(suite.T(), suite.testable.IsSandbox())
}

func (suite *ConfigTestSuite) TestIsValidSuccess() {
	assert.Nil(suite.T(), suite.testable.IsValid())
	assert.NoError(suite.T(), suite.testable.IsValid())
}

func (suite *ConfigTestSuite) TestIsValidEmptyUri() {
	suite.testable.URI = ""
	result := suite.testable.IsValid()
	assert.Error(suite.T(), result)
	assert.Equal(suite.T(), `parameter "uri" is empty`, result.Error())
}

func (suite *ConfigTestSuite) TestIsValidEmptyPublicKey() {
	suite.testable.PublicKey = ""
	result := suite.testable.IsValid()
	assert.Error(suite.T(), result)
	assert.Equal(suite.T(), `parameter "public_key" is empty`, result.Error())
}

func (suite *ConfigTestSuite) TestIsValidEmptySecretKey() {
	suite.testable.SecretKey = ""
	result := suite.testable.IsValid()
	assert.Error(suite.T(), result)
	assert.Equal(suite.T(), `parameter "secret_key" is empty`, result.Error())
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}
