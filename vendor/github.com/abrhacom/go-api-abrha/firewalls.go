package go_api_abrha

import (
	"context"
	"net/http"
	"path"
	"strconv"
)

const firewallsBasePath = "api/public/v1/firewalls"

// FirewallsService is an interface for managing Firewalls with the Abrha API.
// See: https://docs.parspack.com/api/#tag/Firewalls
type FirewallsService interface {
	Get(context.Context, string) (*Firewall, *Response, error)
	Create(context.Context, *FirewallRequest) (*Firewall, *Response, error)
	Update(context.Context, string, *FirewallRequest) (*Firewall, *Response, error)
	Delete(context.Context, string) (*Response, error)
	List(context.Context, *ListOptions) ([]Firewall, *Response, error)
	ListByVm(context.Context, int, *ListOptions) ([]Firewall, *Response, error)
	AddVms(context.Context, string, ...string) (*Response, error)
	RemoveVms(context.Context, string, ...string) (*Response, error)
	AddTags(context.Context, string, ...string) (*Response, error)
	RemoveTags(context.Context, string, ...string) (*Response, error)
	AddRules(context.Context, string, *FirewallRulesRequest) (*Response, error)
	RemoveRules(context.Context, string, *FirewallRulesRequest) (*Response, error)
}

// FirewallsServiceOp handles communication with Firewalls methods of the Abrha API.
type FirewallsServiceOp struct {
	client *Client
}

// Firewall represents a Abrha Firewall configuration.
type Firewall struct {
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	Status         string          `json:"status"`
	InboundRules   []InboundRule   `json:"inbound_rules"`
	OutboundRules  []OutboundRule  `json:"outbound_rules"`
	VmIDs          []string        `json:"vm_ids"`
	Tags           []string        `json:"tags"`
	Created        string          `json:"created_at"`
	PendingChanges []PendingChange `json:"pending_changes"`
}

// String creates a human-readable description of a Firewall.
func (fw Firewall) String() string {
	return Stringify(fw)
}

// URN returns the firewall name in a valid DO API URN form.
func (fw Firewall) URN() string {
	return ToURN("Firewall", fw.ID)
}

// FirewallRequest represents the configuration to be applied to an existing or a new Firewall.
type FirewallRequest struct {
	Name          string         `json:"name"`
	InboundRules  []InboundRule  `json:"inbound_rules"`
	OutboundRules []OutboundRule `json:"outbound_rules"`
	VmIDs         []string       `json:"vm_ids"`
	Tags          []string       `json:"tags"`
}

// FirewallRulesRequest represents rules configuration to be applied to an existing Firewall.
type FirewallRulesRequest struct {
	InboundRules  []InboundRule  `json:"inbound_rules"`
	OutboundRules []OutboundRule `json:"outbound_rules"`
}

// InboundRule represents a Abrha Firewall inbound rule.
type InboundRule struct {
	Protocol  string   `json:"protocol,omitempty"`
	PortRange string   `json:"ports,omitempty"`
	Sources   *Sources `json:"sources"`
}

// OutboundRule represents a Abrha Firewall outbound rule.
type OutboundRule struct {
	Protocol     string        `json:"protocol,omitempty"`
	PortRange    string        `json:"ports,omitempty"`
	Destinations *Destinations `json:"destinations"`
}

// Sources represents a Abrha Firewall InboundRule sources.
type Sources struct {
	Addresses        []string `json:"addresses,omitempty"`
	Tags             []string `json:"tags,omitempty"`
	VmIDs            []string `json:"vm_ids,omitempty"`
	LoadBalancerUIDs []string `json:"load_balancer_uids,omitempty"`
	KubernetesIDs    []string `json:"kubernetes_ids,omitempty"`
}

// PendingChange represents a Abrha Firewall status details.
type PendingChange struct {
	VmID     string `json:"vm_id,omitempty"`
	Removing bool   `json:"removing,omitempty"`
	Status   string `json:"status,omitempty"`
}

// Destinations represents a Abrha Firewall OutboundRule destinations.
type Destinations struct {
	Addresses        []string `json:"addresses,omitempty"`
	Tags             []string `json:"tags,omitempty"`
	VmIDs            []string `json:"vm_ids,omitempty"`
	LoadBalancerUIDs []string `json:"load_balancer_uids,omitempty"`
	KubernetesIDs    []string `json:"kubernetes_ids,omitempty"`
}

var _ FirewallsService = &FirewallsServiceOp{}

// Get an existing Firewall by its identifier.
func (fw *FirewallsServiceOp) Get(ctx context.Context, fID string) (*Firewall, *Response, error) {
	path := path.Join(firewallsBasePath, fID)

	req, err := fw.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(firewallRoot)
	resp, err := fw.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Firewall, resp, err
}

// Create a new Firewall with a given configuration.
func (fw *FirewallsServiceOp) Create(ctx context.Context, fr *FirewallRequest) (*Firewall, *Response, error) {
	req, err := fw.client.NewRequest(ctx, http.MethodPost, firewallsBasePath, fr)
	if err != nil {
		return nil, nil, err
	}

	root := new(firewallRoot)
	resp, err := fw.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Firewall, resp, err
}

// Update an existing Firewall with new configuration.
func (fw *FirewallsServiceOp) Update(ctx context.Context, fID string, fr *FirewallRequest) (*Firewall, *Response, error) {
	path := path.Join(firewallsBasePath, fID)

	req, err := fw.client.NewRequest(ctx, "PUT", path, fr)
	if err != nil {
		return nil, nil, err
	}

	root := new(firewallRoot)
	resp, err := fw.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Firewall, resp, err
}

// Delete a Firewall by its identifier.
func (fw *FirewallsServiceOp) Delete(ctx context.Context, fID string) (*Response, error) {
	path := path.Join(firewallsBasePath, fID)
	return fw.createAndDoReq(ctx, http.MethodDelete, path, nil)
}

// List Firewalls.
func (fw *FirewallsServiceOp) List(ctx context.Context, opt *ListOptions) ([]Firewall, *Response, error) {
	path, err := addOptions(firewallsBasePath, opt)
	if err != nil {
		return nil, nil, err
	}

	return fw.listHelper(ctx, path)
}

// ListByVm Firewalls.
func (fw *FirewallsServiceOp) ListByVm(ctx context.Context, dID int, opt *ListOptions) ([]Firewall, *Response, error) {
	basePath := path.Join(vmBasePath, strconv.Itoa(dID), "firewalls")
	path, err := addOptions(basePath, opt)
	if err != nil {
		return nil, nil, err
	}

	return fw.listHelper(ctx, path)
}

// AddVms to a Firewall.
func (fw *FirewallsServiceOp) AddVms(ctx context.Context, fID string, vmIDs ...string) (*Response, error) {
	path := path.Join(firewallsBasePath, fID, "vms")
	return fw.createAndDoReq(ctx, http.MethodPost, path, &vmsRequest{IDs: vmIDs})
}

// RemoveVms from a Firewall.
func (fw *FirewallsServiceOp) RemoveVms(ctx context.Context, fID string, vmIDs ...string) (*Response, error) {
	path := path.Join(firewallsBasePath, fID, "vms")
	return fw.createAndDoReq(ctx, http.MethodDelete, path, &vmsRequest{IDs: vmIDs})
}

// AddTags to a Firewall.
func (fw *FirewallsServiceOp) AddTags(ctx context.Context, fID string, tags ...string) (*Response, error) {
	path := path.Join(firewallsBasePath, fID, "tags")
	return fw.createAndDoReq(ctx, http.MethodPost, path, &tagsRequest{Tags: tags})
}

// RemoveTags from a Firewall.
func (fw *FirewallsServiceOp) RemoveTags(ctx context.Context, fID string, tags ...string) (*Response, error) {
	path := path.Join(firewallsBasePath, fID, "tags")
	return fw.createAndDoReq(ctx, http.MethodDelete, path, &tagsRequest{Tags: tags})
}

// AddRules to a Firewall.
func (fw *FirewallsServiceOp) AddRules(ctx context.Context, fID string, rr *FirewallRulesRequest) (*Response, error) {
	path := path.Join(firewallsBasePath, fID, "rules")
	return fw.createAndDoReq(ctx, http.MethodPost, path, rr)
}

// RemoveRules from a Firewall.
func (fw *FirewallsServiceOp) RemoveRules(ctx context.Context, fID string, rr *FirewallRulesRequest) (*Response, error) {
	path := path.Join(firewallsBasePath, fID, "rules")
	return fw.createAndDoReq(ctx, http.MethodDelete, path, rr)
}

type vmsRequest struct {
	IDs []string `json:"vm_ids"`
}

type tagsRequest struct {
	Tags []string `json:"tags"`
}

type firewallRoot struct {
	Firewall *Firewall `json:"firewall"`
}

type firewallsRoot struct {
	Firewalls []Firewall `json:"firewalls"`
	Links     *Links     `json:"links"`
	Meta      *Meta      `json:"meta"`
}

func (fw *FirewallsServiceOp) createAndDoReq(ctx context.Context, method, path string, v interface{}) (*Response, error) {
	req, err := fw.client.NewRequest(ctx, method, path, v)
	if err != nil {
		return nil, err
	}

	return fw.client.Do(ctx, req, nil)
}

func (fw *FirewallsServiceOp) listHelper(ctx context.Context, path string) ([]Firewall, *Response, error) {
	req, err := fw.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(firewallsRoot)
	resp, err := fw.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if l := root.Links; l != nil {
		resp.Links = l
	}
	if m := root.Meta; m != nil {
		resp.Meta = m
	}

	return root.Firewalls, resp, err
}
