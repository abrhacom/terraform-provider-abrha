package go_api_abrha

import (
	"context"
	"net/http"
)

// RegionsService is an interface for interfacing with the regions
// endpoints of the Abrha API
// See: https://docs.parspack.com/reference/api/cloud-server/#tag/Regions
type RegionsService interface {
	List(context.Context, *ListOptions) ([]Region, *Response, error)
}

// RegionsServiceOp handles communication with the region related methods of the
// Abrha API.
type RegionsServiceOp struct {
	client *Client
}

var _ RegionsService = &RegionsServiceOp{}

// Region represents a Abrha Region
type Region struct {
	Slug      string   `json:"slug,omitempty"`
	Name      string   `json:"name,omitempty"`
	Sizes     []string `json:"sizes,omitempty"`
	Available bool     `json:"available,omitempty"`
	Features  []string `json:"features,omitempty"`
}

type regionsRoot struct {
	Regions []Region
	Links   *Links `json:"links"`
	Meta    *Meta  `json:"meta"`
}

func (r Region) String() string {
	return Stringify(r)
}

// List all regions
func (s *RegionsServiceOp) List(ctx context.Context, opt *ListOptions) ([]Region, *Response, error) {
	path := "api/public/v1/regions"
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(regionsRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if l := root.Links; l != nil {
		resp.Links = l
	}
	if m := root.Meta; m != nil {
		resp.Meta = m
	}

	return root.Regions, resp, err
}
