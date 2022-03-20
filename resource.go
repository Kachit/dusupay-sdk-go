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

//get HTTP method wrapper
func (r *ResourceAbstract) get(ctx context.Context, path string, query map[string]interface{}) (*Response, error) {
	rsp, err := r.tr.Get(ctx, path, query)
	if err != nil {
		return nil, fmt.Errorf("ResourceAbstract.get request: %v", err)
	}
	return NewResponse(rsp), nil
}

//post HTTP method wrapper
func (r *ResourceAbstract) post(ctx context.Context, path string, body map[string]interface{}, query map[string]interface{}) (*Response, error) {
	rsp, err := r.tr.Post(ctx, path, body, query)
	if err != nil {
		return nil, fmt.Errorf("ResourceAbstract.post request: %v", err)
	}
	return NewResponse(rsp), nil
}

//NewResourceAbstract Create new resource abstract
func NewResourceAbstract(transport *Transport, config *Config) *ResourceAbstract {
	return &ResourceAbstract{tr: transport, cfg: config}
}
