package dusupay

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Config_NewConfig(t *testing.T) {
	result := NewConfig("foo", "bar")
	assert.Equal(t, ProdAPIUrl, result.Uri)
	assert.Equal(t, "foo", result.PublicKey)
	assert.Equal(t, "bar", result.SecretKey)
}

func Test_Config_IsSandbox(t *testing.T) {
	result := NewConfig("foo", "bar")
	assert.False(t, result.IsSandbox())
	result.Uri = SandboxAPIUrl
	assert.True(t, result.IsSandbox())
}

func Test_Config_IsValidSuccess(t *testing.T) {
	config := Config{Uri: ProdAPIUrl, PublicKey: "foo", SecretKey: "bar"}
	assert.Nil(t, config.IsValid())
	assert.NoError(t, config.IsValid())
}

func Test_Config_IsValidEmptyUri(t *testing.T) {
	filter := Config{PublicKey: "foo", SecretKey: "bar"}
	result := filter.IsValid()
	assert.Error(t, result)
	assert.Equal(t, `parameter "uri" is empty`, result.Error())
}

func Test_Config_IsValidEmptyPublicKey(t *testing.T) {
	filter := Config{Uri: ProdAPIUrl, PublicKey: "", SecretKey: "bar"}
	result := filter.IsValid()
	assert.Error(t, result)
	assert.Equal(t, `parameter "public_key" is empty`, result.Error())
}

func Test_Config_IsValidEmptySecretKey(t *testing.T) {
	filter := Config{Uri: ProdAPIUrl, PublicKey: "foo", SecretKey: ""}
	result := filter.IsValid()
	assert.Error(t, result)
	assert.Equal(t, `parameter "secret_key" is empty`, result.Error())
}
