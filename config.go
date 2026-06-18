package dusupay

import "fmt"

// Config structure
type Config struct {
	URI        string `json:"uri"`
	PublicKey  string `json:"public_key"`
	SecretKey  string `json:"secret_key"`
	SigningKey string `json:"signing_key"`
}

// IsSandbox check is sandbox environment
func (c *Config) IsSandbox() bool {
	return c.URI != ProdAPIUrl
}

// IsValid check is valid config parameters
func (c *Config) IsValid() error {
	var err error
	if c.URI == "" {
		err = fmt.Errorf(`parameter "uri" is empty`)
	} else if c.PublicKey == "" {
		err = fmt.Errorf(`parameter "public_key" is empty`)
	} else if c.SecretKey == "" {
		err = fmt.Errorf(`parameter "secret_key" is empty`)
	}
	return err
}

// NewConfig Create new config from credentials (Prod version)
func NewConfig(publicKey string, secretKey string) *Config {
	cfg := &Config{
		URI:       ProdAPIUrl,
		PublicKey: publicKey,
		SecretKey: secretKey,
	}
	return cfg
}

// NewConfigSandbox Create new config from credentials (Sandbox version)
func NewConfigSandbox(publicKey string, secretKey string) *Config {
	cfg := &Config{
		URI:       SandboxAPIUrl,
		PublicKey: publicKey,
		SecretKey: secretKey,
	}
	return cfg
}
