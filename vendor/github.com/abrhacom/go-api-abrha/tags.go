package go_api_abrha

import (
	"context"
	"fmt"
	"net/http"
)

const tagsBasePath = "api/public/v1/tags"

// TagsService is an interface for interfacing with the tags
// endpoints of the Abrha API
// See: https://docs.parspack.com/api/#tag/Tags
type TagsService interface {
	List(context.Context, *ListOptions) ([]Tag, *Response, error)
	Get(context.Context, string) (*Tag, *Response, error)
	Create(context.Context, *TagCreateRequest) (*Tag, *Response, error)
	Delete(context.Context, string) (*Response, error)

	TagResources(context.Context, string, *TagResourcesRequest) (*Response, error)
	UntagResources(context.Context, string, *UntagResourcesRequest) (*Response, error)
}

// TagsServiceOp handles communication with tag related method of the
// Abrha API.
type TagsServiceOp struct {
	client *Client
}

var _ TagsService = &TagsServiceOp{}

// ResourceType represents a class of resource, currently only vm are supported
type ResourceType string

const (
	// VmResourceType holds the string representing our ResourceType of Vm.
	VmResourceType ResourceType = "vm"
	// ImageResourceType holds the string representing our ResourceType of Image.
	ImageResourceType ResourceType = "image"
	// VolumeResourceType holds the string representing our ResourceType of Volume.
	VolumeResourceType ResourceType = "volume"
	// LoadBalancerResourceType holds the string representing our ResourceType of LoadBalancer.
	LoadBalancerResourceType ResourceType = "load_balancer"
	// VolumeSnapshotResourceType holds the string representing our ResourceType for storage Snapshots.
	VolumeSnapshotResourceType ResourceType = "volume_snapshot"
	// DatabaseResourceType holds the string representing our ResourceType of Database.
	DatabaseResourceType ResourceType = "database"
)

// Resource represent a single resource for associating/disassociating with tags
type Resource struct {
	ID   string       `json:"resource_id,omitempty"`
	Type ResourceType `json:"resource_type,omitempty"`
}

// TaggedResources represent the set of resources a tag is attached to
type TaggedResources struct {
	Count           int                             `json:"count"`
	LastTaggedURI   string                          `json:"last_tagged_uri,omitempty"`
	Vms             *TaggedvmsResources             `json:"vms,omitempty"`
	Images          *TaggedImagesResources          `json:"images"`
	Volumes         *TaggedVolumesResources         `json:"volumes"`
	VolumeSnapshots *TaggedVolumeSnapshotsResources `json:"volume_snapshots"`
	Databases       *TaggedDatabasesResources       `json:"databases"`
}

// TaggedvmsResources represent the vm resources a tag is attached to
type TaggedvmsResources struct {
	Count         int    `json:"count,float64,omitempty"`
	LastTagged    *Vm    `json:"last_tagged,omitempty"`
	LastTaggedURI string `json:"last_tagged_uri,omitempty"`
}

// TaggedResourcesData represent the generic resources a tag is attached to
type TaggedResourcesData struct {
	Count         int    `json:"count,float64,omitempty"`
	LastTaggedURI string `json:"last_tagged_uri,omitempty"`
}

// TaggedImagesResources represent the image resources a tag is attached to
type TaggedImagesResources TaggedResourcesData

// TaggedVolumesResources represent the volume resources a tag is attached to
type TaggedVolumesResources TaggedResourcesData

// TaggedVolumeSnapshotsResources represent the volume snapshot resources a tag is attached to
type TaggedVolumeSnapshotsResources TaggedResourcesData

// TaggedDatabasesResources represent the database resources a tag is attached to
type TaggedDatabasesResources TaggedResourcesData

// Tag represent Abrha tag
type Tag struct {
	Name      string           `json:"name,omitempty"`
	Resources *TaggedResources `json:"resources,omitempty"`
}

// TagCreateRequest represents the JSON structure of a request of that type.
type TagCreateRequest struct {
	Name string `json:"name"`
}

// TagResourcesRequest represents the JSON structure of a request of that type.
type TagResourcesRequest struct {
	Resources []Resource `json:"resources"`
}

// UntagResourcesRequest represents the JSON structure of a request of that type.
type UntagResourcesRequest struct {
	Resources []Resource `json:"resources"`
}

type tagsRoot struct {
	Tags  []Tag  `json:"tags"`
	Links *Links `json:"links"`
	Meta  *Meta  `json:"meta"`
}

type tagRoot struct {
	Tag *Tag `json:"tag"`
}

// List all tags
func (s *TagsServiceOp) List(ctx context.Context, opt *ListOptions) ([]Tag, *Response, error) {
	path := tagsBasePath
	path, err := addOptions(path, opt)

	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(tagsRoot)
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

	return root.Tags, resp, err
}

// Get a single tag
func (s *TagsServiceOp) Get(ctx context.Context, name string) (*Tag, *Response, error) {
	path := fmt.Sprintf("%s/%s", tagsBasePath, name)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(tagRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Tag, resp, err
}

// Create a new tag
func (s *TagsServiceOp) Create(ctx context.Context, createRequest *TagCreateRequest) (*Tag, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, tagsBasePath, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(tagRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Tag, resp, err
}

// Delete an existing tag
func (s *TagsServiceOp) Delete(ctx context.Context, name string) (*Response, error) {
	if name == "" {
		return nil, NewArgError("name", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s", tagsBasePath, name)
	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)

	return resp, err
}

// TagResources associates resources with a given Tag.
func (s *TagsServiceOp) TagResources(ctx context.Context, name string, tagRequest *TagResourcesRequest) (*Response, error) {
	if name == "" {
		return nil, NewArgError("name", "cannot be empty")
	}

	if tagRequest == nil {
		return nil, NewArgError("tagRequest", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%s/resources", tagsBasePath, name)
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, tagRequest)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)

	return resp, err
}

// UntagResources dissociates resources with a given Tag.
func (s *TagsServiceOp) UntagResources(ctx context.Context, name string, untagRequest *UntagResourcesRequest) (*Response, error) {
	if name == "" {
		return nil, NewArgError("name", "cannot be empty")
	}

	if untagRequest == nil {
		return nil, NewArgError("tagRequest", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%s/resources", tagsBasePath, name)
	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, untagRequest)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)

	return resp, err
}
