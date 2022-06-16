package dusupay

import "net/http"

//Client struct
type Client struct {
	transport *Transport
	config    *Config
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
	return &Client{transport, config}, nil
}

//Collections resource
func (c *Client) Collections() *CollectionsResource {
	return &CollectionsResource{NewResourceAbstract(c.transport, c.config)}
}

//Payouts resource
func (c *Client) Payouts() *PayoutsResource {
	return &PayoutsResource{NewResourceAbstract(c.transport, c.config)}
}

//Providers resource
func (c *Client) Providers() *ProvidersResource {
	return &ProvidersResource{NewResourceAbstract(c.transport, c.config)}
}

//Merchants resource
func (c *Client) Merchants() *MerchantsResource {
	return &MerchantsResource{NewResourceAbstract(c.transport, c.config)}
}

//Refunds resource
func (c *Client) Refunds() *RefundsResource {
	return &RefundsResource{NewResourceAbstract(c.transport, c.config)}
}

//Banks resource
func (c *Client) Banks() *BanksResource {
	return &BanksResource{NewResourceAbstract(c.transport, c.config)}
}

//Webhooks resource
func (c *Client) Webhooks() *WebhooksResource {
	return &WebhooksResource{NewResourceAbstract(c.transport, c.config)}
}
