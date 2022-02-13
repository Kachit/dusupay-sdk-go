package dusupay

import "net/http"

//Base Client
type Client struct {
}

//NewClientFromConfig Create new client from config
func NewClientFromConfig(config *Config, cl *http.Client) *Client {
	if cl == nil {
		cl = &http.Client{}
	}
	return &Client{}
}

//Collections resource
func (c *Client) Collections() *CollectionsResource {
	return &CollectionsResource{}
}

//Payouts resource
func (c *Client) Payouts() *PayoutsResource {
	return &PayoutsResource{}
}

//Providers resource
func (c *Client) Providers() *ProvidersResource {
	return &ProvidersResource{}
}

//Merchants resource
func (c *Client) Merchants() *MerchantsResource {
	return &MerchantsResource{}
}

//Refunds resource
func (c *Client) Refunds() *RefundsResource {
	return &RefundsResource{}
}

//Banks resource
func (c *Client) Banks() *BanksResource {
	return &BanksResource{}
}
