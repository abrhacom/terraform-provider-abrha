package go_api_abrha

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const vmBasePath = "api/public/v1/vms"

var errNoNetworks = errors.New("no networks have been defined")

// VmsService is an interface for interfacing with the Vm
// endpoints of the Abrha API
// See: https://docs.parspack.com/api/#tag/VMs
type VmsService interface {
	List(context.Context, *ListOptions) ([]Vm, *Response, error)
	ListWithGPUs(context.Context, *ListOptions) ([]Vm, *Response, error)
	ListByName(context.Context, string, *ListOptions) ([]Vm, *Response, error)
	ListByTag(context.Context, string, *ListOptions) ([]Vm, *Response, error)
	Get(context.Context, string) (*Vm, *Response, error)
	Create(context.Context, *VmCreateRequest) (*vmRoot, *Response, error)
	CreateMultiple(context.Context, *VmMultiCreateRequest) ([]Vm, *Response, error)
	Delete(context.Context, string) (*Response, error)
	DeleteByTag(context.Context, string) (*Response, error)
	Kernels(context.Context, string, *ListOptions) ([]Kernel, *Response, error)
	Snapshots(context.Context, string, *ListOptions) ([]Image, *Response, error)
	Backups(context.Context, string, *ListOptions) ([]Image, *Response, error)
	Actions(context.Context, string, *ListOptions) ([]Action, *Response, error)
	Neighbors(context.Context, string) ([]Vm, *Response, error)
	GetBackupPolicy(context.Context, string) (*VmBackupPolicy, *Response, error)
	ListBackupPolicies(context.Context, *ListOptions) (map[int]*VmBackupPolicy, *Response, error)
	ListSupportedBackupPolicies(context.Context) ([]*SupportedBackupPolicy, *Response, error)
}

// VmsServiceOp handles communication with the Vm related methods of the
// Abrha API.
type VmsServiceOp struct {
	client *Client
}

var _ VmsService = &VmsServiceOp{}

// Vm represents a Abrha Vm
type Vm struct {
	ID               string        `json:"id,omitempty"`
	Name             string        `json:"name,omitempty"`
	Memory           int           `json:"memory,omitempty"`
	Vcpus            int           `json:"vcpus,omitempty"`
	Disk             int           `json:"disk,omitempty"`
	Region           *Region       `json:"region,omitempty"`
	Image            *Image        `json:"image,omitempty"`
	Size             *Size         `json:"size,omitempty"`
	SizeSlug         string        `json:"size_slug,omitempty"`
	BackupIDs        []int         `json:"backup_ids,omitempty"`
	NextBackupWindow *BackupWindow `json:"next_backup_window,omitempty"`
	SnapshotIDs      []int         `json:"snapshot_ids,omitempty"`
	Features         []string      `json:"features,omitempty"`
	Locked           bool          `json:"locked,bool,omitempty"`
	Status           string        `json:"status,omitempty"`
	Networks         *Networks     `json:"networks,omitempty"`
	Created          string        `json:"created_at,omitempty"`
	Kernel           *Kernel       `json:"kernel,omitempty"`
	Tags             []string      `json:"tags,omitempty"`
	VolumeIDs        []string      `json:"volume_ids"`
	VPCUUID          string        `json:"vpc_uuid,omitempty"`
}

// PublicIPv4 returns the public IPv4 address for the Vm.
func (d *Vm) PublicIPv4() (string, error) {
	if d.Networks == nil {
		return "", errNoNetworks
	}

	for _, v4 := range d.Networks.V4 {
		if v4.Type == "public" {
			return v4.IPAddress, nil
		}
	}

	return "", nil
}

// PrivateIPv4 returns the private IPv4 address for the Vm.
func (d *Vm) PrivateIPv4() (string, error) {
	if d.Networks == nil {
		return "", errNoNetworks
	}

	for _, v4 := range d.Networks.V4 {
		if v4.Type == "private" {
			return v4.IPAddress, nil
		}
	}

	return "", nil
}

// PublicIPv6 returns the public IPv6 address for the Vm.
func (d *Vm) PublicIPv6() (string, error) {
	if d.Networks == nil {
		return "", errNoNetworks
	}

	for _, v6 := range d.Networks.V6 {
		if v6.Type == "public" {
			return v6.IPAddress, nil
		}
	}

	return "", nil
}

// Kernel object
type Kernel struct {
	ID      int    `json:"id,float64,omitempty"`
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

// BackupWindow object
type BackupWindow struct {
	Start *Timestamp `json:"start,omitempty"`
	End   *Timestamp `json:"end,omitempty"`
}

// Convert Vm to a string
func (d Vm) String() string {
	return Stringify(d)
}

// URN returns the vm ID in a valid DO API URN form.
func (d Vm) URN() string {
	return ToURN("vm", d.ID)
}

// VmRoot represents a Vm root
type vmRoot struct {
	Vm    *Vm    `json:"vm"`
	Links *Links `json:"links,omitempty"`
}

type vmsRoot struct {
	Vms   []Vm   `json:"vms"`
	Links *Links `json:"links"`
	Meta  *Meta  `json:"meta"`
}

type kernelsRoot struct {
	Kernels []Kernel `json:"kernels,omitempty"`
	Links   *Links   `json:"links"`
	Meta    *Meta    `json:"meta"`
}

type vmSnapshotsRoot struct {
	Snapshots []Image `json:"snapshots,omitempty"`
	Links     *Links  `json:"links"`
	Meta      *Meta   `json:"meta"`
}

type backupsRoot struct {
	Backups []Image `json:"backups,omitempty"`
	Links   *Links  `json:"links"`
	Meta    *Meta   `json:"meta"`
}

// VmCreateImage identifies an image for the create request. It prefers slug over ID.
type VmCreateImage struct {
	ID   int
	Slug string
}

// MarshalJSON returns either the slug or id of the image. It returns the id
// if the slug is empty.
func (d VmCreateImage) MarshalJSON() ([]byte, error) {
	if d.Slug != "" {
		return json.Marshal(d.Slug)
	}

	return json.Marshal(d.ID)
}

// VmCreateVolume identifies a volume to attach for the create request.
type VmCreateVolume struct {
	ID string
	// Deprecated: You must pass the volume's ID when creating a Vm.
	Name string
}

// MarshalJSON returns an object with either the ID or name of the volume. It
// prefers the ID over the name.
func (d VmCreateVolume) MarshalJSON() ([]byte, error) {
	if d.ID != "" {
		return json.Marshal(struct {
			ID string `json:"id"`
		}{ID: d.ID})
	}

	return json.Marshal(struct {
		Name string `json:"name"`
	}{Name: d.Name})
}

// VmCreateSSHKey identifies a SSH Key for the create request. It prefers fingerprint over ID.
type VmCreateSSHKey struct {
	ID          int
	Fingerprint string
}

// MarshalJSON returns either the fingerprint or id of the ssh key. It returns
// the id if the fingerprint is empty.
func (d VmCreateSSHKey) MarshalJSON() ([]byte, error) {
	if d.Fingerprint != "" {
		return json.Marshal(d.Fingerprint)
	}

	return json.Marshal(d.ID)
}

// VmCreateRequest represents a request to create a Vm.
type VmCreateRequest struct {
	Name              string                 `json:"name"`
	Region            string                 `json:"region"`
	Size              string                 `json:"size"`
	Image             VmCreateImage          `json:"image"`
	SSHKeys           []VmCreateSSHKey       `json:"ssh_keys"`
	Backups           bool                   `json:"backups"`
	IPv6              bool                   `json:"ipv6"`
	PrivateNetworking bool                   `json:"private_networking"`
	Monitoring        bool                   `json:"monitoring"`
	UserData          string                 `json:"user_data,omitempty"`
	Volumes           []VmCreateVolume       `json:"volumes,omitempty"`
	Tags              []string               `json:"tags"`
	VPCUUID           string                 `json:"vpc_uuid,omitempty"`
	WithVmAgent       *bool                  `json:"with_vm_agent,omitempty"`
	BackupPolicy      *VmBackupPolicyRequest `json:"backup_policy,omitempty"`
}

// VmMultiCreateRequest is a request to create multiple Vms.
type VmMultiCreateRequest struct {
	Names             []string               `json:"names"`
	Region            string                 `json:"region"`
	Size              string                 `json:"size"`
	Image             VmCreateImage          `json:"image"`
	SSHKeys           []VmCreateSSHKey       `json:"ssh_keys"`
	Backups           bool                   `json:"backups"`
	IPv6              bool                   `json:"ipv6"`
	PrivateNetworking bool                   `json:"private_networking"`
	Monitoring        bool                   `json:"monitoring"`
	UserData          string                 `json:"user_data,omitempty"`
	Tags              []string               `json:"tags"`
	VPCUUID           string                 `json:"vpc_uuid,omitempty"`
	WithVmAgent       *bool                  `json:"with_vm_agent,omitempty"`
	BackupPolicy      *VmBackupPolicyRequest `json:"backup_policy,omitempty"`
}

// VmBackupPolicyRequest defines the backup policy when creating a Vm.
type VmBackupPolicyRequest struct {
	Plan     string `json:"plan,omitempty"`
	Weekday  string `json:"weekday,omitempty"`
	Monthday int    `json:"monthday,omitempty"`
	Hour     *int   `json:"hour,omitempty"`
}

func (d VmCreateRequest) String() string {
	return Stringify(d)
}

func (d VmMultiCreateRequest) String() string {
	return Stringify(d)
}

// Networks represents the Vm's Networks.
type Networks struct {
	V4 []NetworkV4 `json:"v4,omitempty"`
	V6 []NetworkV6 `json:"v6,omitempty"`
}

// NetworkV4 represents a Abrha IPv4 Network.
type NetworkV4 struct {
	IPAddress string `json:"ip_address,omitempty"`
	Netmask   string `json:"netmask,omitempty"`
	Gateway   string `json:"gateway,omitempty"`
	Type      string `json:"type,omitempty"`
}

func (n NetworkV4) String() string {
	return Stringify(n)
}

// NetworkV6 represents a Abrha IPv6 network.
type NetworkV6 struct {
	IPAddress string `json:"ip_address,omitempty"`
	Netmask   int    `json:"netmask,omitempty"`
	Gateway   string `json:"gateway,omitempty"`
	Type      string `json:"type,omitempty"`
}

func (n NetworkV6) String() string {
	return Stringify(n)
}

// Performs a list request given a path.
func (s *VmsServiceOp) list(ctx context.Context, path string) ([]Vm, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(vmsRoot)
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

	return root.Vms, resp, err
}

// List all Vms.
func (s *VmsServiceOp) List(ctx context.Context, opt *ListOptions) ([]Vm, *Response, error) {
	path := vmBasePath
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	return s.list(ctx, path)
}

// ListWithGPUs lists all Vms with GPUs.
func (s *VmsServiceOp) ListWithGPUs(ctx context.Context, opt *ListOptions) ([]Vm, *Response, error) {
	path := fmt.Sprintf("%s?type=gpus", vmBasePath)
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	return s.list(ctx, path)
}

// ListByName lists all Vms filtered by name returning only exact matches.
// It is case-insensitive
func (s *VmsServiceOp) ListByName(ctx context.Context, name string, opt *ListOptions) ([]Vm, *Response, error) {
	path := fmt.Sprintf("%s?name=%s", vmBasePath, name)
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	return s.list(ctx, path)
}

// ListByTag lists all Vms matched by a Tag.
func (s *VmsServiceOp) ListByTag(ctx context.Context, tag string, opt *ListOptions) ([]Vm, *Response, error) {
	path := fmt.Sprintf("%s?tag_name=%s", vmBasePath, tag)
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	return s.list(ctx, path)
}

// Get individual Vm.
func (s *VmsServiceOp) Get(ctx context.Context, vmID string) (*Vm, *Response, error) {
	path := fmt.Sprintf("%s/%s", vmBasePath, vmID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(vmRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Vm, resp, err
}

// Create Vm
func (s *VmsServiceOp) Create(ctx context.Context, createRequest *VmCreateRequest) (*vmRoot, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	path := vmBasePath

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(vmRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root, resp, err
}

// CreateMultiple creates multiple Vms.
func (s *VmsServiceOp) CreateMultiple(ctx context.Context, createRequest *VmMultiCreateRequest) ([]Vm, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	path := vmBasePath

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(vmsRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root.Vms, resp, err
}

// Performs a delete request given a path
func (s *VmsServiceOp) delete(ctx context.Context, path string) (*Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)

	return resp, err
}

// Delete Vm.
func (s *VmsServiceOp) Delete(ctx context.Context, vmID string) (*Response, error) {
	path := fmt.Sprintf("%s/%s", vmBasePath, vmID)

	return s.delete(ctx, path)
}

// DeleteByTag deletes Vms matched by a Tag.
func (s *VmsServiceOp) DeleteByTag(ctx context.Context, tag string) (*Response, error) {
	if tag == "" {
		return nil, NewArgError("tag", "cannot be empty")
	}

	path := fmt.Sprintf("%s?tag_name=%s", vmBasePath, tag)

	return s.delete(ctx, path)
}

// Kernels lists kernels available for a Vm.
func (s *VmsServiceOp) Kernels(ctx context.Context, vmID string, opt *ListOptions) ([]Kernel, *Response, error) {

	path := fmt.Sprintf("%s/%s/kernels", vmBasePath, vmID)
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(kernelsRoot)
	resp, err := s.client.Do(ctx, req, root)
	if l := root.Links; l != nil {
		resp.Links = l
	}
	if m := root.Meta; m != nil {
		resp.Meta = m
	}

	return root.Kernels, resp, err
}

// Actions lists the actions for a Vm.
func (s *VmsServiceOp) Actions(ctx context.Context, vmID string, opt *ListOptions) ([]Action, *Response, error) {
	path := fmt.Sprintf("%s/%s/actions", vmBasePath, vmID)
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

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

// Backups lists the backups for a Vm.
func (s *VmsServiceOp) Backups(ctx context.Context, vmID string, opt *ListOptions) ([]Image, *Response, error) {
	path := fmt.Sprintf("%s/%s/backups", vmBasePath, vmID)
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(backupsRoot)
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

	return root.Backups, resp, err
}

// Snapshots lists the snapshots available for a Vm.
func (s *VmsServiceOp) Snapshots(ctx context.Context, vmID string, opt *ListOptions) ([]Image, *Response, error) {

	path := fmt.Sprintf("%s/%s/snapshots", vmBasePath, vmID)
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(vmSnapshotsRoot)
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

	return root.Snapshots, resp, err
}

// Neighbors lists the neighbors for a Vm.
func (s *VmsServiceOp) Neighbors(ctx context.Context, vmID string) ([]Vm, *Response, error) {

	path := fmt.Sprintf("%s/%s/neighbors", vmBasePath, vmID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(vmsRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Vms, resp, err
}

func (s *VmsServiceOp) vmActionStatus(ctx context.Context, uri string) (string, error) {
	action, _, err := s.client.VmActions.GetByURI(ctx, uri)

	if err != nil {
		return "", err
	}

	return action.Status, nil
}

// VmBackupPolicy defines the information about a vm's backup policy.
type VmBackupPolicy struct {
	VmID             string                `json:"vm_id,omitempty"`
	BackupEnabled    bool                  `json:"backup_enabled,omitempty"`
	BackupPolicy     *VmBackupPolicyConfig `json:"backup_policy,omitempty"`
	NextBackupWindow *BackupWindow         `json:"next_backup_window,omitempty"`
}

// VmBackupPolicyConfig defines the backup policy for a Vm.
type VmBackupPolicyConfig struct {
	Plan                string `json:"plan,omitempty"`
	Weekday             string `json:"weekday,omitempty"`
	Monthday            int    `json:"monthday,omitempty"`
	Hour                int    `json:"hour,omitempty"`
	WindowLengthHours   int    `json:"window_length_hours,omitempty"`
	RetentionPeriodDays int    `json:"retention_period_days,omitempty"`
}

// vmBackupPolicyRoot represents a VmBackupPolicy root
type vmBackupPolicyRoot struct {
	VmBackupPolicy *VmBackupPolicy `json:"policy,omitempty"`
}

type vmBackupPoliciesRoot struct {
	VmBackupPolicies map[int]*VmBackupPolicy `json:"policies,omitempty"`
	Links            *Links                  `json:"links,omitempty"`
	Meta             *Meta                   `json:"meta"`
}

// Get individual vm backup policy.
func (s *VmsServiceOp) GetBackupPolicy(ctx context.Context, vmID string) (*VmBackupPolicy, *Response, error) {

	path := fmt.Sprintf("%s/%s/backups/policy", vmBasePath, vmID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(vmBackupPolicyRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.VmBackupPolicy, resp, err
}

// List all vm backup policies.
func (s *VmsServiceOp) ListBackupPolicies(ctx context.Context, opt *ListOptions) (map[int]*VmBackupPolicy, *Response, error) {
	path := fmt.Sprintf("%s/backups/policies", vmBasePath)
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(vmBackupPoliciesRoot)
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

	return root.VmBackupPolicies, resp, nil
}

type SupportedBackupPolicy struct {
	Name                 string   `json:"name,omitempty"`
	PossibleWindowStarts []int    `json:"possible_window_starts,omitempty"`
	WindowLengthHours    int      `json:"window_length_hours,omitempty"`
	RetentionPeriodDays  int      `json:"retention_period_days,omitempty"`
	PossibleDays         []string `json:"possible_days,omitempty"`
}

type vmSupportedBackupPoliciesRoot struct {
	SupportedBackupPolicies []*SupportedBackupPolicy `json:"supported_policies,omitempty"`
}

// List supported vm backup policies.
func (s *VmsServiceOp) ListSupportedBackupPolicies(ctx context.Context) ([]*SupportedBackupPolicy, *Response, error) {
	path := fmt.Sprintf("%s/backups/supported_policies", vmBasePath)
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(vmSupportedBackupPoliciesRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.SupportedBackupPolicies, resp, nil
}
