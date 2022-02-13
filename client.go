package dusupay

import "net/http"

//Base Client
type Client struct {
	transport *Transport
}

//NewClientFromConfig Create new client from config
func NewClientFromConfig(config *Config, cl *http.Client) (*Client, error) {
	err := config.IsValid()
	if err != nil {
		return nil, err
	}
	if cl == nil {
		cl = &http.Client{}
	}
	transport := NewHttpTransport(config, cl)
	return &Client{transport}, nil
}

//Collections resource
func (c *Client) Collections() *CollectionsResource {
	return &CollectionsResource{ResourceAbstract: NewResourceAbstract(c.transport)}
}

//Payouts resource
func (c *Client) Payouts() *PayoutsResource {
	return &PayoutsResource{ResourceAbstract: NewResourceAbstract(c.transport)}
}

//Providers resource
func (c *Client) Providers() *ProvidersResource {
	return &ProvidersResource{ResourceAbstract: NewResourceAbstract(c.transport)}
}

//Merchants resource
func (c *Client) Merchants() *MerchantsResource {
	return &MerchantsResource{ResourceAbstract: NewResourceAbstract(c.transport)}
}

//Refunds resource
func (c *Client) Refunds() *RefundsResource {
	return &RefundsResource{ResourceAbstract: NewResourceAbstract(c.transport)}
}

//Banks resource
func (c *Client) Banks() *BanksResource {
	return &BanksResource{ResourceAbstract: NewResourceAbstract(c.transport)}
}
