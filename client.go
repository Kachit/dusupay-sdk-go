package dusupay

import "net/http"

type Client struct {
}

//NewClientFromConfig Create new client from config
func NewClientFromConfig(config *Config, cl *http.Client) *Client {
	if cl == nil {
		cl = &http.Client{}
	}
	return &Client{}
}
