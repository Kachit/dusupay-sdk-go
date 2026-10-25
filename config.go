package dusupay

import "fmt"

// Config structure
type Config struct {
	URI        string `json:"uri"`
	PublicKey  string `json:"public_key"`
	SecretKey  string `json:"secret_key"`
	SigningKey string `json:"signing_key"`
	// Language optional Accept-Language header value for API messages (en, fr, sw)
	Language LanguageCode `json:"language,omitempty"`
}

// IsSandbox check is sandbox environment
func (c *Config) IsSandbox() bool {
	return c.URI != ProdAPIUrl
}

// DusupayPublicKey returns the DusuPay PEM public key for the configured environment
func (c *Config) DusupayPublicKey() []byte {
	if c.IsSandbox() {
		return []byte(SandboxPublicKey)
	}

	return []byte(ProdPublicKey)
}

// IsValid check is valid config parameters
func (c *Config) IsValid() error {
	switch {
	case c.URI == "":
		return fmt.Errorf(`parameter "uri" is empty`)
	case c.PublicKey == "":
		return fmt.Errorf(`parameter "public_key" is empty`)
	case c.SecretKey == "":
		return fmt.Errorf(`parameter "secret_key" is empty`)
	}

	return nil
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
