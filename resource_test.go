package dusupay

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Resource_NewResourceAbstract(t *testing.T) {
	config := BuildStubConfig()
	transport := NewHttpTransport(config, nil)
	resource := NewResourceAbstract(transport, config)
	assert.NotEmpty(t, resource)
}
