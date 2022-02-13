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
