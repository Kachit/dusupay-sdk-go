package dusupay

import "fmt"

//Config structure
type Config struct {
	Uri       string `json:"uri"`
	PublicKey string `json:"public_key"`
	SecretKey string `json:"secret_key"`
}

//Check is sandbox environment
func (c *Config) IsSandbox() bool {
	return c.Uri != ProdAPIUrl
}

//Check is valid Config parameters
func (c *Config) IsValid() error {
	var err error
	if c.Uri == "" {
		err = fmt.Errorf(`parameter "uri" is empty`)
	} else if c.PublicKey == "" {
		err = fmt.Errorf(`parameter "public_key" is empty`)
	} else if c.SecretKey == "" {
		err = fmt.Errorf(`parameter "secret_key" is empty`)
	}
	return err
}

//NewConfig Create new config from credentials
func NewConfig(publicKey string, secretKey string) *Config {
	cfg := &Config{
		Uri:       ProdAPIUrl,
		PublicKey: publicKey,
		SecretKey: secretKey,
	}
	return cfg
}
