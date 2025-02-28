package go_api_abrha

import (
	"context"
	"fmt"
	"net/http"
)

// FloatingIPActionsService is an interface for interfacing with the
// floating IPs actions endpoints of the Pars Pack API.
// See: https://docs.parspack.com/api/#tag/Floating-IP-Actions
type FloatingIPActionsService interface {
	Assign(ctx context.Context, ip string, vmID string) (*Action, *Response, error)
	Unassign(ctx context.Context, ip string) (*Action, *Response, error)
	Get(ctx context.Context, ip string, actionID int) (*Action, *Response, error)
	List(ctx context.Context, ip string, opt *ListOptions) ([]Action, *Response, error)
}

// FloatingIPActionsServiceOp handles communication with the floating IPs
// action related methods of the Abrha API.
type FloatingIPActionsServiceOp struct {
	client *Client
}

// Assign a floating IP to a vm.
func (s *FloatingIPActionsServiceOp) Assign(ctx context.Context, ip string, vmID string) (*Action, *Response, error) {
	request := &ActionRequest{
		"type":  "assign",
		"vm_id": vmID,
	}
	return s.doAction(ctx, ip, request)
}

// Unassign a floating IP from the vm it is currently assigned to.
func (s *FloatingIPActionsServiceOp) Unassign(ctx context.Context, ip string) (*Action, *Response, error) {
	request := &ActionRequest{"type": "unassign"}
	return s.doAction(ctx, ip, request)
}

// Get an action for a particular floating IP by id.
func (s *FloatingIPActionsServiceOp) Get(ctx context.Context, ip string, actionID int) (*Action, *Response, error) {
	path := fmt.Sprintf("%s/%d", floatingIPActionPath(ip), actionID)
	return s.get(ctx, path)
}

// List the actions for a particular floating IP.
func (s *FloatingIPActionsServiceOp) List(ctx context.Context, ip string, opt *ListOptions) ([]Action, *Response, error) {
	path := floatingIPActionPath(ip)
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	return s.list(ctx, path)
}

func (s *FloatingIPActionsServiceOp) doAction(ctx context.Context, ip string, request *ActionRequest) (*Action, *Response, error) {
	path := floatingIPActionPath(ip)

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

func (s *FloatingIPActionsServiceOp) get(ctx context.Context, path string) (*Action, *Response, error) {
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

func (s *FloatingIPActionsServiceOp) list(ctx context.Context, path string) ([]Action, *Response, error) {
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

	return root.Actions, resp, err
}

func floatingIPActionPath(ip string) string {
	return fmt.Sprintf("%s/%s/actions", floatingBasePath, ip)
}
