package go_api_abrha

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	vmAutoscaleBasePath = "/api/public/v1/vms/autoscale"
)

// VmAutoscaleService defines an interface for managing vm autoscale pools through Abrha API
type VmAutoscaleService interface {
	Create(context.Context, *VmAutoscalePoolRequest) (*VmAutoscalePool, *Response, error)
	Get(context.Context, string) (*VmAutoscalePool, *Response, error)
	List(context.Context, *ListOptions) ([]*VmAutoscalePool, *Response, error)
	ListMembers(context.Context, string, *ListOptions) ([]*VmAutoscaleResource, *Response, error)
	ListHistory(context.Context, string, *ListOptions) ([]*VmAutoscaleHistoryEvent, *Response, error)
	Update(context.Context, string, *VmAutoscalePoolRequest) (*VmAutoscalePool, *Response, error)
	Delete(context.Context, string) (*Response, error)
	DeleteDangerous(context.Context, string) (*Response, error)
}

// VmAutoscalePool represents a Abrha vm autoscale pool
type VmAutoscalePool struct {
	ID                 string                          `json:"id"`
	Name               string                          `json:"name"`
	Config             *VmAutoscaleConfiguration       `json:"config"`
	VmTemplate         *VmAutoscaleResourceTemplate    `json:"vm_template"`
	CreatedAt          time.Time                       `json:"created_at"`
	UpdatedAt          time.Time                       `json:"updated_at"`
	CurrentUtilization *VmAutoscaleResourceUtilization `json:"current_utilization,omitempty"`
	Status             string                          `json:"status"`
}

// VmAutoscaleConfiguration represents a Abrha vm autoscale pool configuration
type VmAutoscaleConfiguration struct {
	MinInstances            uint64  `json:"min_instances,omitempty"`
	MaxInstances            uint64  `json:"max_instances,omitempty"`
	TargetCPUUtilization    float64 `json:"target_cpu_utilization,omitempty"`
	TargetMemoryUtilization float64 `json:"target_memory_utilization,omitempty"`
	CooldownMinutes         uint32  `json:"cooldown_minutes,omitempty"`
	TargetNumberInstances   uint64  `json:"target_number_instances,omitempty"`
}

// VmAutoscaleResourceTemplate represents a Abrha vm autoscale pool resource template
type VmAutoscaleResourceTemplate struct {
	Size        string   `json:"size"`
	Region      string   `json:"region"`
	Image       string   `json:"image"`
	Tags        []string `json:"tags"`
	SSHKeys     []string `json:"ssh_keys"`
	VpcUUID     string   `json:"vpc_uuid"`
	WithVmAgent bool     `json:"with_vm_agent"`
	ProjectID   string   `json:"project_id"`
	IPV6        bool     `json:"ipv6"`
	UserData    string   `json:"user_data"`
}

// VmAutoscaleResourceUtilization represents a Abrha vm autoscale pool resource utilization
type VmAutoscaleResourceUtilization struct {
	Memory float64 `json:"memory,omitempty"`
	CPU    float64 `json:"cpu,omitempty"`
}

// VmAutoscaleResource represents a Abrha vm autoscale pool resource
type VmAutoscaleResource struct {
	VmID               string                          `json:"vm_id"`
	CreatedAt          time.Time                       `json:"created_at"`
	UpdatedAt          time.Time                       `json:"updated_at"`
	HealthStatus       string                          `json:"health_status"`
	UnhealthyReason    string                          `json:"unhealthy_reason,omitempty"`
	Status             string                          `json:"status"`
	CurrentUtilization *VmAutoscaleResourceUtilization `json:"current_utilization,omitempty"`
}

// VmAutoscaleHistoryEvent represents a Abrha vm autoscale pool history event
type VmAutoscaleHistoryEvent struct {
	HistoryEventID       string    `json:"history_event_id"`
	CurrentInstanceCount uint64    `json:"current_instance_count"`
	DesiredInstanceCount uint64    `json:"desired_instance_count"`
	Reason               string    `json:"reason"`
	Status               string    `json:"status"`
	ErrorReason          string    `json:"error_reason,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// VmAutoscalePoolRequest represents a Abrha vm autoscale pool create/update request
type VmAutoscalePoolRequest struct {
	Name       string                       `json:"name"`
	Config     *VmAutoscaleConfiguration    `json:"config"`
	VmTemplate *VmAutoscaleResourceTemplate `json:"vm_template"`
}

type vmAutoscalePoolRoot struct {
	AutoscalePool *VmAutoscalePool `json:"autoscale_pool"`
}

type vmAutoscalePoolsRoot struct {
	AutoscalePools []*VmAutoscalePool `json:"autoscale_pools"`
	Links          *Links             `json:"links"`
	Meta           *Meta              `json:"meta"`
}

type vmAutoscaleMembersRoot struct {
	Vms   []*VmAutoscaleResource `json:"vms"`
	Links *Links                 `json:"links"`
	Meta  *Meta                  `json:"meta"`
}

type vmAutoscaleHistoryEventsRoot struct {
	History []*VmAutoscaleHistoryEvent `json:"history"`
	Links   *Links                     `json:"links"`
	Meta    *Meta                      `json:"meta"`
}

// VmAutoscaleServiceOp handles communication with vm autoscale-related methods of the Abrha API
type VmAutoscaleServiceOp struct {
	client *Client
}

var _ VmAutoscaleService = &VmAutoscaleServiceOp{}

// Create a new vm autoscale pool
func (d *VmAutoscaleServiceOp) Create(ctx context.Context, createReq *VmAutoscalePoolRequest) (*VmAutoscalePool, *Response, error) {
	req, err := d.client.NewRequest(ctx, http.MethodPost, vmAutoscaleBasePath, createReq)
	if err != nil {
		return nil, nil, err
	}
	root := new(vmAutoscalePoolRoot)
	resp, err := d.client.Do(ctx, req, root)
	if err != nil {
		return nil, nil, err
	}
	return root.AutoscalePool, resp, nil
}

// Get an existing vm autoscale pool
func (d *VmAutoscaleServiceOp) Get(ctx context.Context, id string) (*VmAutoscalePool, *Response, error) {
	req, err := d.client.NewRequest(ctx, http.MethodGet, fmt.Sprintf("%s/%s", vmAutoscaleBasePath, id), nil)
	if err != nil {
		return nil, nil, err
	}
	root := new(vmAutoscalePoolRoot)
	resp, err := d.client.Do(ctx, req, root)
	if err != nil {
		return nil, nil, err
	}
	return root.AutoscalePool, resp, err
}

// List all existing vm autoscale pools
func (d *VmAutoscaleServiceOp) List(ctx context.Context, opts *ListOptions) ([]*VmAutoscalePool, *Response, error) {
	path, err := addOptions(vmAutoscaleBasePath, opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := d.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	root := new(vmAutoscalePoolsRoot)
	resp, err := d.client.Do(ctx, req, root)
	if err != nil {
		return nil, nil, err
	}
	if root.Links != nil {
		resp.Links = root.Links
	}
	if root.Meta != nil {
		resp.Meta = root.Meta
	}
	return root.AutoscalePools, resp, err
}

// ListMembers all members for an existing vm autoscale pool
func (d *VmAutoscaleServiceOp) ListMembers(ctx context.Context, id string, opts *ListOptions) ([]*VmAutoscaleResource, *Response, error) {
	path, err := addOptions(fmt.Sprintf("%s/%s/members", vmAutoscaleBasePath, id), opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := d.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	root := new(vmAutoscaleMembersRoot)
	resp, err := d.client.Do(ctx, req, root)
	if err != nil {
		return nil, nil, err
	}
	if root.Links != nil {
		resp.Links = root.Links
	}
	if root.Meta != nil {
		resp.Meta = root.Meta
	}
	return root.Vms, resp, err
}

// ListHistory all history events for an existing vm autoscale pool
func (d *VmAutoscaleServiceOp) ListHistory(ctx context.Context, id string, opts *ListOptions) ([]*VmAutoscaleHistoryEvent, *Response, error) {
	path, err := addOptions(fmt.Sprintf("%s/%s/history", vmAutoscaleBasePath, id), opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := d.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	root := new(vmAutoscaleHistoryEventsRoot)
	resp, err := d.client.Do(ctx, req, root)
	if err != nil {
		return nil, nil, err
	}
	if root.Links != nil {
		resp.Links = root.Links
	}
	if root.Meta != nil {
		resp.Meta = root.Meta
	}
	return root.History, resp, err
}

// Update an existing autoscale pool
func (d *VmAutoscaleServiceOp) Update(ctx context.Context, id string, updateReq *VmAutoscalePoolRequest) (*VmAutoscalePool, *Response, error) {
	req, err := d.client.NewRequest(ctx, http.MethodPut, fmt.Sprintf("%s/%s", vmAutoscaleBasePath, id), updateReq)
	if err != nil {
		return nil, nil, err
	}
	root := new(vmAutoscalePoolRoot)
	resp, err := d.client.Do(ctx, req, root)
	if err != nil {
		return nil, nil, err
	}
	return root.AutoscalePool, resp, nil
}

// Delete an existing autoscale pool
func (d *VmAutoscaleServiceOp) Delete(ctx context.Context, id string) (*Response, error) {
	req, err := d.client.NewRequest(ctx, http.MethodDelete, fmt.Sprintf("%s/%s", vmAutoscaleBasePath, id), nil)
	if err != nil {
		return nil, err
	}
	return d.client.Do(ctx, req, nil)
}

// DeleteDangerous deletes an existing autoscale pool with all underlying resources
func (d *VmAutoscaleServiceOp) DeleteDangerous(ctx context.Context, id string) (*Response, error) {
	req, err := d.client.NewRequest(ctx, http.MethodDelete, fmt.Sprintf("%s/%s/dangerous", vmAutoscaleBasePath, id), nil)
	req.Header.Set("X-Dangerous", "true")
	if err != nil {
		return nil, err
	}
	return d.client.Do(ctx, req, nil)
}
