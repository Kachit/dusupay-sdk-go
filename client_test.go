package dusupay

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ClientTestSuite struct {
	suite.Suite
}

func (suite *ClientTestSuite) TestNewClientFromConfigValid() {
	cfg := BuildStubConfig()
	client, err := NewClientFromConfig(cfg, nil)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), client)
}

func (suite *ClientTestSuite) TestNewClientFromConfigInvalid() {
	cfg := BuildStubConfig()
	cfg.Uri = ""
	client, err := NewClientFromConfig(cfg, nil)
	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), client)
}

func (suite *ClientTestSuite) TestGetCollectionsResource() {
	client, err := NewClientFromConfig(BuildStubConfig(), nil)
	result := client.Collections()
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), result)
}

func (suite *ClientTestSuite) TestGetPayoutsResource() {
	client, err := NewClientFromConfig(BuildStubConfig(), nil)
	result := client.Payouts()
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), result)
}

func (suite *ClientTestSuite) TestGetProvidersResource() {
	client, err := NewClientFromConfig(BuildStubConfig(), nil)
	result := client.Providers()
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), result)
}

func (suite *ClientTestSuite) TestGetMerchantsResource() {
	client, err := NewClientFromConfig(BuildStubConfig(), nil)
	result := client.Merchants()
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), result)
}

func (suite *ClientTestSuite) TestGetRefundsResource() {
	client, err := NewClientFromConfig(BuildStubConfig(), nil)
	result := client.Refunds()
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), result)
}

func (suite *ClientTestSuite) TestGetBanksResource() {
	client, err := NewClientFromConfig(BuildStubConfig(), nil)
	result := client.Banks()
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), result)
}

func (suite *ClientTestSuite) TestGetWebhooksResource() {
	client, err := NewClientFromConfig(BuildStubConfig(), nil)
	result := client.Webhooks()
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), result)
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}
