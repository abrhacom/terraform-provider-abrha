package go_api_abrha

import (
	"context"
	"fmt"
	"net/http"
)

// StorageActionsService is an interface for interfacing with the
// storage actions endpoints of the Pars Pack API.
// See: https://docs.parspack.com/reference/api/cloud-server/#tag/Block-Storage-Actions
type StorageActionsService interface {
	Attach(ctx context.Context, volumeID string, vmID string) (*Action, *Response, error)
	DetachByVmID(ctx context.Context, volumeID string, vmID string) (*Action, *Response, error)
	Get(ctx context.Context, volumeID string, actionID int) (*Action, *Response, error)
	List(ctx context.Context, volumeID string, opt *ListOptions) ([]Action, *Response, error)
	Resize(ctx context.Context, volumeID string, sizeGigabytes int, regionSlug string) (*Action, *Response, error)
}

// StorageActionsServiceOp handles communication with the storage volumes
// action related methods of the Abrha API.
type StorageActionsServiceOp struct {
	client *Client
}

// StorageAttachment represents the attachment of a block storage
// volume to a specific Vm under the device name.
type StorageAttachment struct {
	VmID string `json:"vm_id"`
}

// Attach a storage volume to a Vm.
func (s *StorageActionsServiceOp) Attach(ctx context.Context, volumeID string, vmID string) (*Action, *Response, error) {
	request := &ActionRequest{
		"type":  "attach",
		"vm_id": vmID,
	}
	return s.doAction(ctx, volumeID, request)
}

// DetachByVmID a storage volume from a Vm by Vm ID.
func (s *StorageActionsServiceOp) DetachByVmID(ctx context.Context, volumeID string, vmID string) (*Action, *Response, error) {
	request := &ActionRequest{
		"type":  "detach",
		"vm_id": vmID,
	}
	return s.doAction(ctx, volumeID, request)
}

// Get an action for a particular storage volume by id.
func (s *StorageActionsServiceOp) Get(ctx context.Context, volumeID string, actionID int) (*Action, *Response, error) {
	path := fmt.Sprintf("%s/%d", storageAllocationActionPath(volumeID), actionID)
	return s.get(ctx, path)
}

// List the actions for a particular storage volume.
func (s *StorageActionsServiceOp) List(ctx context.Context, volumeID string, opt *ListOptions) ([]Action, *Response, error) {
	path := storageAllocationActionPath(volumeID)
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	return s.list(ctx, path)
}

// Resize a storage volume.
func (s *StorageActionsServiceOp) Resize(ctx context.Context, volumeID string, sizeGigabytes int, regionSlug string) (*Action, *Response, error) {
	request := &ActionRequest{
		"type":           "resize",
		"size_gigabytes": sizeGigabytes,
		"region":         regionSlug,
	}
	return s.doAction(ctx, volumeID, request)
}

func (s *StorageActionsServiceOp) doAction(ctx context.Context, volumeID string, request *ActionRequest) (*Action, *Response, error) {
	path := storageAllocationActionPath(volumeID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, request)
	if err != nil {
		return nil, nil, err
	}

	root := new(actionRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Event, resp, err
}

func (s *StorageActionsServiceOp) get(ctx context.Context, path string) (*Action, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(actionRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Event, resp, err
}

func (s *StorageActionsServiceOp) list(ctx context.Context, path string) ([]Action, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(actionsRoot)
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

	return root.Actions, resp, err
}

func storageAllocationActionPath(volumeID string) string {
	return fmt.Sprintf("%s/%s/actions", storageAllocPath, volumeID)
}
