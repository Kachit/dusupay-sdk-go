package dusupay

//Config structure
type Config struct {
	Uri       string
	PublicKey string
	SecretKey string
}

func (c *Config) IsSandbox() bool {
	return c.Uri != ProdAPIUrl
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
