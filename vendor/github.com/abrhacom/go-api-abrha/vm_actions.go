package go_api_abrha

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// ActionRequest represents Abrha Action Request
type ActionRequest map[string]interface{}

// VmActionsService is an interface for interfacing with the Vm actions
// endpoints of the Abrha API
// See: https://docs.parspack.com/api/#tag/VM-Actions
type VmActionsService interface {
	Shutdown(context.Context, string) (*Action, *Response, error)
	ShutdownByTag(context.Context, string) ([]Action, *Response, error)
	PowerOff(context.Context, string) (*Action, *Response, error)
	PowerOffByTag(context.Context, string) ([]Action, *Response, error)
	PowerOn(context.Context, string) (*Action, *Response, error)
	PowerOnByTag(context.Context, string) ([]Action, *Response, error)
	PowerCycle(context.Context, string) (*Action, *Response, error)
	PowerCycleByTag(context.Context, string) ([]Action, *Response, error)
	Reboot(context.Context, string) (*Action, *Response, error)
	Restore(context.Context, string, int) (*Action, *Response, error)
	Resize(context.Context, string, string, bool) (*Action, *Response, error)
	Rename(context.Context, string, string) (*Action, *Response, error)
	Snapshot(context.Context, string, string) (*Action, *Response, error)
	SnapshotByTag(context.Context, string, string) ([]Action, *Response, error)
	EnableBackups(context.Context, string) (*Action, *Response, error)
	EnableBackupsByTag(context.Context, string) ([]Action, *Response, error)
	EnableBackupsWithPolicy(context.Context, string, *VmBackupPolicyRequest) (*Action, *Response, error)
	ChangeBackupPolicy(context.Context, string, *VmBackupPolicyRequest) (*Action, *Response, error)
	DisableBackups(context.Context, string) (*Action, *Response, error)
	DisableBackupsByTag(context.Context, string) ([]Action, *Response, error)
	PasswordReset(context.Context, string) (*Action, *Response, error)
	RebuildByImageID(context.Context, string, int) (*Action, *Response, error)
	RebuildByImageSlug(context.Context, string, string) (*Action, *Response, error)
	ChangeKernel(context.Context, string, int) (*Action, *Response, error)
	EnableIPv6(context.Context, string) (*Action, *Response, error)
	EnableIPv6ByTag(context.Context, string) ([]Action, *Response, error)
	EnablePrivateNetworking(context.Context, string) (*Action, *Response, error)
	EnablePrivateNetworkingByTag(context.Context, string) ([]Action, *Response, error)
	Get(context.Context, string, int) (*Action, *Response, error)
	GetByURI(context.Context, string) (*Action, *Response, error)
}

// VmActionsServiceOp handles communication with the Vm action related
// methods of the Abrha API.
type VmActionsServiceOp struct {
	client *Client
}

var _ VmActionsService = &VmActionsServiceOp{}

// Shutdown a Vm
func (s *VmActionsServiceOp) Shutdown(ctx context.Context, id string) (*Action, *Response, error) {
	request := &ActionRequest{"type": "shutdown"}
	return s.doAction(ctx, id, request)
}

// ShutdownByTag shuts down Vms matched by a Tag.
func (s *VmActionsServiceOp) ShutdownByTag(ctx context.Context, tag string) ([]Action, *Response, error) {
	request := &ActionRequest{"type": "shutdown"}
	return s.doActionByTag(ctx, tag, request)
}

// PowerOff a Vm
func (s *VmActionsServiceOp) PowerOff(ctx context.Context, id string) (*Action, *Response, error) {
	request := &ActionRequest{"type": "power_off"}
	return s.doAction(ctx, id, request)
}

// PowerOffByTag powers off Vms matched by a Tag.
func (s *VmActionsServiceOp) PowerOffByTag(ctx context.Context, tag string) ([]Action, *Response, error) {
	request := &ActionRequest{"type": "power_off"}
	return s.doActionByTag(ctx, tag, request)
}

// PowerOn a Vm
func (s *VmActionsServiceOp) PowerOn(ctx context.Context, id string) (*Action, *Response, error) {
	request := &ActionRequest{"type": "power_on"}
	return s.doAction(ctx, id, request)
}

// PowerOnByTag powers on Vms matched by a Tag.
func (s *VmActionsServiceOp) PowerOnByTag(ctx context.Context, tag string) ([]Action, *Response, error) {
	request := &ActionRequest{"type": "power_on"}
	return s.doActionByTag(ctx, tag, request)
}

// PowerCycle a Vm
func (s *VmActionsServiceOp) PowerCycle(ctx context.Context, id string) (*Action, *Response, error) {
	request := &ActionRequest{"type": "power_cycle"}
	return s.doAction(ctx, id, request)
}

// PowerCycleByTag power cycles Vms matched by a Tag.
func (s *VmActionsServiceOp) PowerCycleByTag(ctx context.Context, tag string) ([]Action, *Response, error) {
	request := &ActionRequest{"type": "power_cycle"}
	return s.doActionByTag(ctx, tag, request)
}

// Reboot a Vm
func (s *VmActionsServiceOp) Reboot(ctx context.Context, id string) (*Action, *Response, error) {
	request := &ActionRequest{"type": "reboot"}
	return s.doAction(ctx, id, request)
}

// Restore an image to a Vm
func (s *VmActionsServiceOp) Restore(ctx context.Context, id string, imageID int) (*Action, *Response, error) {
	requestType := "restore"
	request := &ActionRequest{
		"type":  requestType,
		"image": float64(imageID),
	}
	return s.doAction(ctx, id, request)
}

// Resize a Vm
func (s *VmActionsServiceOp) Resize(ctx context.Context, id string, sizeSlug string, resizeDisk bool) (*Action, *Response, error) {
	requestType := "resize"
	request := &ActionRequest{
		"type": requestType,
		"size": sizeSlug,
		"disk": resizeDisk,
	}
	return s.doAction(ctx, id, request)
}

// Rename a Vm
func (s *VmActionsServiceOp) Rename(ctx context.Context, id string, name string) (*Action, *Response, error) {
	requestType := "rename"
	request := &ActionRequest{
		"type": requestType,
		"name": name,
	}
	return s.doAction(ctx, id, request)
}

// Snapshot a Vm.
func (s *VmActionsServiceOp) Snapshot(ctx context.Context, id string, name string) (*Action, *Response, error) {
	requestType := "snapshot"
	request := &ActionRequest{
		"type": requestType,
		"name": name,
	}
	return s.doAction(ctx, id, request)
}

// SnapshotByTag snapshots Vms matched by a Tag.
func (s *VmActionsServiceOp) SnapshotByTag(ctx context.Context, tag string, name string) ([]Action, *Response, error) {
	requestType := "snapshot"
	request := &ActionRequest{
		"type": requestType,
		"name": name,
	}
	return s.doActionByTag(ctx, tag, request)
}

// EnableBackups enables backups for a Vm.
func (s *VmActionsServiceOp) EnableBackups(ctx context.Context, id string) (*Action, *Response, error) {
	request := &ActionRequest{"type": "enable_backups"}
	return s.doAction(ctx, id, request)
}

// EnableBackupsByTag enables backups for Vms matched by a Tag.
func (s *VmActionsServiceOp) EnableBackupsByTag(ctx context.Context, tag string) ([]Action, *Response, error) {
	request := &ActionRequest{"type": "enable_backups"}
	return s.doActionByTag(ctx, tag, request)
}

// EnableBackupsWithPolicy enables vm's backup with a backup policy applied.
func (s *VmActionsServiceOp) EnableBackupsWithPolicy(ctx context.Context, id string, policy *VmBackupPolicyRequest) (*Action, *Response, error) {
	if policy == nil {
		return nil, nil, NewArgError("policy", "policy can't be nil")
	}

	policyMap := map[string]interface{}{
		"plan":     policy.Plan,
		"weekday":  policy.Weekday,
		"monthday": policy.Monthday,
	}
	if policy.Hour != nil {
		policyMap["hour"] = policy.Hour
	}

	request := &ActionRequest{"type": "enable_backups", "backup_policy": policyMap}
	return s.doAction(ctx, id, request)
}

// ChangeBackupPolicy updates a backup policy when backups are enabled.
func (s *VmActionsServiceOp) ChangeBackupPolicy(ctx context.Context, id string, policy *VmBackupPolicyRequest) (*Action, *Response, error) {
	if policy == nil {
		return nil, nil, NewArgError("policy", "policy can't be nil")
	}

	policyMap := map[string]interface{}{
		"plan":     policy.Plan,
		"weekday":  policy.Weekday,
		"monthday": policy.Monthday,
	}
	if policy.Hour != nil {
		policyMap["hour"] = policy.Hour
	}

	request := &ActionRequest{"type": "change_backup_policy", "backup_policy": policyMap}
	return s.doAction(ctx, id, request)
}

// DisableBackups disables backups for a Vm.
func (s *VmActionsServiceOp) DisableBackups(ctx context.Context, id string) (*Action, *Response, error) {
	request := &ActionRequest{"type": "disable_backups"}
	return s.doAction(ctx, id, request)
}

// DisableBackupsByTag disables backups for Vm matched by a Tag.
func (s *VmActionsServiceOp) DisableBackupsByTag(ctx context.Context, tag string) ([]Action, *Response, error) {
	request := &ActionRequest{"type": "disable_backups"}
	return s.doActionByTag(ctx, tag, request)
}

// PasswordReset resets the password for a Vm.
func (s *VmActionsServiceOp) PasswordReset(ctx context.Context, id string) (*Action, *Response, error) {
	request := &ActionRequest{"type": "password_reset"}
	return s.doAction(ctx, id, request)
}

// RebuildByImageID rebuilds a Vm from an image with a given id.
func (s *VmActionsServiceOp) RebuildByImageID(ctx context.Context, id string, imageID int) (*Action, *Response, error) {
	request := &ActionRequest{"type": "rebuild", "image": imageID}
	return s.doAction(ctx, id, request)
}

// RebuildByImageSlug rebuilds a Vm from an Image matched by a given Slug.
func (s *VmActionsServiceOp) RebuildByImageSlug(ctx context.Context, id string, slug string) (*Action, *Response, error) {
	request := &ActionRequest{"type": "rebuild", "image": slug}
	return s.doAction(ctx, id, request)
}

// ChangeKernel changes the kernel for a Vm.
func (s *VmActionsServiceOp) ChangeKernel(ctx context.Context, id string, kernelID int) (*Action, *Response, error) {
	request := &ActionRequest{"type": "change_kernel", "kernel": kernelID}
	return s.doAction(ctx, id, request)
}

// EnableIPv6 enables IPv6 for a Vm.
func (s *VmActionsServiceOp) EnableIPv6(ctx context.Context, id string) (*Action, *Response, error) {
	request := &ActionRequest{"type": "enable_ipv6"}
	return s.doAction(ctx, id, request)
}

// EnableIPv6ByTag enables IPv6 for Vms matched by a Tag.
func (s *VmActionsServiceOp) EnableIPv6ByTag(ctx context.Context, tag string) ([]Action, *Response, error) {
	request := &ActionRequest{"type": "enable_ipv6"}
	return s.doActionByTag(ctx, tag, request)
}

// EnablePrivateNetworking enables private networking for a Vm.
func (s *VmActionsServiceOp) EnablePrivateNetworking(ctx context.Context, id string) (*Action, *Response, error) {
	request := &ActionRequest{"type": "enable_private_networking"}
	return s.doAction(ctx, id, request)
}

// EnablePrivateNetworkingByTag enables private networking for Vms matched by a Tag.
func (s *VmActionsServiceOp) EnablePrivateNetworkingByTag(ctx context.Context, tag string) ([]Action, *Response, error) {
	request := &ActionRequest{"type": "enable_private_networking"}
	return s.doActionByTag(ctx, tag, request)
}

func (s *VmActionsServiceOp) doAction(ctx context.Context, id string, request *ActionRequest) (*Action, *Response, error) {
	if request == nil {
		return nil, nil, NewArgError("request", "request can't be nil")
	}

	path := vmActionPath(id)

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

func (s *VmActionsServiceOp) doActionByTag(ctx context.Context, tag string, request *ActionRequest) ([]Action, *Response, error) {
	if tag == "" {
		return nil, nil, NewArgError("tag", "cannot be empty")
	}

	if request == nil {
		return nil, nil, NewArgError("request", "request can't be nil")
	}

	path := vmActionPathByTag(tag)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, request)
	if err != nil {
		return nil, nil, err
	}

	root := new(actionsRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Actions, resp, err
}

// Get an action for a particular Vm by id.
func (s *VmActionsServiceOp) Get(ctx context.Context, vmID string, actionID int) (*Action, *Response, error) {
	if actionID < 1 {
		return nil, nil, NewArgError("actionID", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d", vmActionPath(vmID), actionID)
	return s.get(ctx, path)
}

// GetByURI gets an action for a particular Vm by URI.
func (s *VmActionsServiceOp) GetByURI(ctx context.Context, rawurl string) (*Action, *Response, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, nil, err
	}

	return s.get(ctx, u.Path)

}

func (s *VmActionsServiceOp) get(ctx context.Context, path string) (*Action, *Response, error) {
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

func vmActionPath(vmID string) string {
	return fmt.Sprintf("api/public/v1/vms/%s/actions", vmID)
}

func vmActionPathByTag(tag string) string {
	return fmt.Sprintf("api/public/v1/vms/actions?tag_name=%s", tag)
}
