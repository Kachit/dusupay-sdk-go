package dusupay

import (
	"context"
	"fmt"
)

//ResourceAbstract base resource
type ResourceAbstract struct {
	tr  *Transport
	cfg *Config
}

//Get HTTP method wrapper
func (r *ResourceAbstract) Get(ctx context.Context, path string, query map[string]interface{}) (*Response, error) {
	rsp, err := r.tr.Get(ctx, path, query)
	if err != nil {
		return nil, fmt.Errorf("ResourceAbstract@get request: %v", err)
	}
	return NewResponse(rsp), nil
}

//Post HTTP method wrapper
func (r *ResourceAbstract) Post(ctx context.Context, path string, body map[string]interface{}, query map[string]interface{}) (*Response, error) {
	rsp, err := r.tr.Post(ctx, path, body, query)
	if err != nil {
		return nil, fmt.Errorf("ResourceAbstract@post request: %v", err)
	}
	return NewResponse(rsp), nil
}

//NewResourceAbstract Create new resource abstract
func NewResourceAbstract(transport *Transport) *ResourceAbstract {
	return &ResourceAbstract{tr: transport, cfg: transport.rb.cfg}
}
