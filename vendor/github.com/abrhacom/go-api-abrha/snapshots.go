package go_api_abrha

import (
	"context"
	"fmt"
	"net/http"
)

const snapshotBasePath = "api/public/v1/snapshots"

// SnapshotsService is an interface for interfacing with the snapshots
// endpoints of the Abrha API
// See: https://docs.parspack.com/reference/api/cloud-server/#tag/Snapshots
type SnapshotsService interface {
	List(context.Context, *ListOptions) ([]Snapshot, *Response, error)
	ListVolume(context.Context, *ListOptions) ([]Snapshot, *Response, error)
	ListVm(context.Context, *ListOptions) ([]Snapshot, *Response, error)
	Get(context.Context, string) (*Snapshot, *Response, error)
	Delete(context.Context, string) (*Response, error)
}

// SnapshotsServiceOp handles communication with the snapshot related methods of the
// Abrha API.
type SnapshotsServiceOp struct {
	client *Client
}

var _ SnapshotsService = &SnapshotsServiceOp{}

// Snapshot represents a Abrha Snapshot
type Snapshot struct {
	ID            string   `json:"id,omitempty"`
	Name          string   `json:"name,omitempty"`
	ResourceID    string   `json:"resource_id,omitempty"`
	ResourceType  string   `json:"resource_type,omitempty"`
	Regions       []string `json:"regions,omitempty"`
	MinDiskSize   int      `json:"min_disk_size,omitempty"`
	SizeGigaBytes float64  `json:"size_gigabytes,omitempty"`
	Created       string   `json:"created_at,omitempty"`
	Tags          []string `json:"tags,omitempty"`
}

type snapshotRoot struct {
	Snapshot *Snapshot `json:"snapshot"`
}

type snapshotsRoot struct {
	Snapshots []Snapshot `json:"snapshots"`
	Links     *Links     `json:"links,omitempty"`
	Meta      *Meta      `json:"meta,omitempty"`
}

type listSnapshotOptions struct {
	ResourceType string `url:"resource_type,omitempty"`
}

func (s Snapshot) String() string {
	return Stringify(s)
}

// List lists all the snapshots available.
func (s *SnapshotsServiceOp) List(ctx context.Context, opt *ListOptions) ([]Snapshot, *Response, error) {
	return s.list(ctx, opt, nil)
}

// ListVm lists all the Vm snapshots.
func (s *SnapshotsServiceOp) ListVm(ctx context.Context, opt *ListOptions) ([]Snapshot, *Response, error) {
	listOpt := listSnapshotOptions{ResourceType: "vm"}
	return s.list(ctx, opt, &listOpt)
}

// ListVolume lists all the volume snapshots.
func (s *SnapshotsServiceOp) ListVolume(ctx context.Context, opt *ListOptions) ([]Snapshot, *Response, error) {
	listOpt := listSnapshotOptions{ResourceType: "volume"}
	return s.list(ctx, opt, &listOpt)
}

// Get retrieves a snapshot by id.
func (s *SnapshotsServiceOp) Get(ctx context.Context, snapshotID string) (*Snapshot, *Response, error) {
	return s.get(ctx, snapshotID)
}

// Delete an snapshot.
func (s *SnapshotsServiceOp) Delete(ctx context.Context, snapshotID string) (*Response, error) {
	path := fmt.Sprintf("%s/%s", snapshotBasePath, snapshotID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)

	return resp, err
}

// Helper method for getting an individual snapshot
func (s *SnapshotsServiceOp) get(ctx context.Context, ID string) (*Snapshot, *Response, error) {
	path := fmt.Sprintf("%s/%s", snapshotBasePath, ID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(snapshotRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Snapshot, resp, err
}

// Helper method for listing snapshots
func (s *SnapshotsServiceOp) list(ctx context.Context, opt *ListOptions, listOpt *listSnapshotOptions) ([]Snapshot, *Response, error) {
	path := snapshotBasePath
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}
	path, err = addOptions(path, listOpt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(snapshotsRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root.Snapshots, resp, err
}
