package dusupay

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Common_NewConfig(t *testing.T) {
	result := NewConfig("foo", "bar")
	assert.Equal(t, ProdAPIUrl, result.Uri)
	assert.Equal(t, "foo", result.PublicKey)
	assert.Equal(t, "bar", result.SecretKey)
}
