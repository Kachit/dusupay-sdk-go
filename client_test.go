package dusupay

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Client_NewClientFromConfigValid(t *testing.T) {
	cfg := BuildStubConfig()
	client, err := NewClientFromConfig(cfg, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, client)
}

func Test_Client_NewClientFromConfigInvalid(t *testing.T) {
	cfg := BuildStubConfig()
	cfg.Uri = ""
	client, err := NewClientFromConfig(cfg, nil)
	assert.Error(t, err)
	assert.Empty(t, client)
}

func Test_Client_GetCollectionsResource(t *testing.T) {
	client, _ := NewClientFromConfig(BuildStubConfig(), nil)
	result := client.Collections()
	assert.NotEmpty(t, result)
}

func Test_Client_GetPayoutsResource(t *testing.T) {
	client, _ := NewClientFromConfig(BuildStubConfig(), nil)
	result := client.Payouts()
	assert.NotEmpty(t, result)
}

func Test_Client_GetProvidersResource(t *testing.T) {
	client, _ := NewClientFromConfig(BuildStubConfig(), nil)
	result := client.Providers()
	assert.NotEmpty(t, result)
}

func Test_Client_GetMerchantsResource(t *testing.T) {
	client, _ := NewClientFromConfig(BuildStubConfig(), nil)
	result := client.Merchants()
	assert.NotEmpty(t, result)
}

func Test_Client_GetRefundsResource(t *testing.T) {
	client, _ := NewClientFromConfig(BuildStubConfig(), nil)
	result := client.Refunds()
	assert.NotEmpty(t, result)
}

func Test_Client_GetBanksResource(t *testing.T) {
	client, _ := NewClientFromConfig(BuildStubConfig(), nil)
	result := client.Banks()
	assert.NotEmpty(t, result)
}
