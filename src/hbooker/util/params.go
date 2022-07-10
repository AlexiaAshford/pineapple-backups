package util

import (
	"net/url"
)

// GetRequest Manage the HTTP GET request parameters
type GetRequest struct {
	urls url.Values
}

// Init Initializer
func (p *GetRequest) Init() *GetRequest {
	p.urls = url.Values{}
	return p
}

// InitFrom Initialized from another instance
func (p *GetRequest) InitFrom(reqParams *GetRequest) *GetRequest {
	if reqParams != nil {
		p.urls = reqParams.urls
	} else {
		p.urls = url.Values{}
	}
	return p
}

// AddParam Add URL escape property and value pair
func (p *GetRequest) AddParam(property string, value string) *GetRequest {
	if property != "" && value != "" {
		p.urls.Add(property, value)
	}
	return p
}

// BuildParams Concat the property and value pair
func (p *GetRequest) BuildParams() string {
	return "?" + p.urls.Encode()
}
