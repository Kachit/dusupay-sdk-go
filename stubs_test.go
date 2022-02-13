package dusupay

func BuildStubConfig() *Config {
	return &Config{
		Uri:       SandboxAPIUrl,
		PublicKey: "PublicKey",
		SecretKey: "SecretKey",
	}
}
