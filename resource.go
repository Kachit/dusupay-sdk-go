package dusupay

//ResourceAbstract base resource
type ResourceAbstract struct {
	tr  *Transport
	cfg *Config
}

//NewResourceAbstract Create new resource abstract
func NewResourceAbstract(transport *Transport, config *Config) *ResourceAbstract {
	return &ResourceAbstract{tr: transport, cfg: config}
}
