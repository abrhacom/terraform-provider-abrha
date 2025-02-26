package vm

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/abrhacom/terraform-provider-abrha/abrha/tag"
	"github.com/abrhacom/terraform-provider-abrha/abrha/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var (
	errVmBackupPolicy = errors.New("backup_policy can only be set when backups are enabled")
)

func ResourceAbrhaVm() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbrhaVmCreate,
		ReadContext:   resourceAbrhaVmRead,
		UpdateContext: resourceAbrhaVmUpdate,
		DeleteContext: resourceAbrhaVmDelete,
		Importer: &schema.ResourceImporter{
			State: resourceAbrhaVmImport,
		},
		MigrateState:  ResourceAbrhaVmMigrateState,
		SchemaVersion: 1,

		// We are using these timeouts to be the minimum timeout for an operation.
		// This is how long an operation will wait for a state update, however
		// implementation of updates and deletes contain multiple instances of waiting for a state update
		// so the true timeout of an operation could be a multiple of the set value.
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Second),
		},

		Schema: map[string]*schema.Schema{
			"image": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				//StateFunc: func(val interface{}) string {
				//	// DO API V2 region slug is always lowercase
				//	return strings.ToLower(val.(string))
				//},
				ValidateFunc: validation.NoZeroValues,
			},

			"size": {
				Type:     schema.TypeString,
				Required: true,
				//StateFunc: func(val interface{}) string {
				//	// DO API V2 size slug is always lowercase
				//	return strings.ToLower(val.(string))
				//},
				ValidateFunc: validation.NoZeroValues,
			},

			"graceful_shutdown": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"urn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"disk": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"vcpus": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"memory": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"price_hourly": {
				Type:     schema.TypeFloat,
				Computed: true,
			},

			"price_monthly": {
				Type:     schema.TypeFloat,
				Computed: true,
			},

			"resize_disk": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"locked": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"backups": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"backup_policy": {
				Type:         schema.TypeList,
				Optional:     true,
				MaxItems:     1,
				RequiredWith: []string{"backups"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"plan": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"daily",
								"weekly",
								"monthly",
							}, false),
						},
						"monthday": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 28),
						},
						"weekday": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"SUN", "MON", "TUE", "WED", "THU", "FRI", "SAT",
							}, false),
						},
						"hour": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 20),
						},
					},
				},
			},

			"ipv6": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"ipv6_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"private_networking": {
				Type:       schema.TypeBool,
				Optional:   true,
				Computed:   true,
				Deprecated: "This parameter has been deprecated. Use `vpc_uuid` instead to specify a VPC network for the Vm. If no `vpc_uuid` is provided, the Vm will be placed in your account's default VPC for the region.",
			},

			"ipv4_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ipv4_address_private": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ssh_keys": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.NoZeroValues,
				},
			},

			"user_data": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
				StateFunc:    util.HashStringStateFunc(),
				// In order to support older statefiles with fully saved user data
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new != "" && old == d.Get("user_data")
				},
			},

			"volume_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
			},

			"monitoring": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},

			"vm_agent": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"tags": tag.TagsSchema(),

			"vpc_uuid": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.NoZeroValues,
			},
		},

		CustomizeDiff: customdiff.All(
			// If the `ipv6` attribute is changed to `true`, we need to mark the
			// `ipv6_address` attribute as changing in the plan. If not, the plan
			// will become inconsistent once the address is known when referenced
			// in another resource such as a domain record, e.g.:
			// https://github.com/parspack/terraform-provider-parspack/issues/981
			customdiff.IfValueChange("ipv6",
				func(ctx context.Context, old, new, meta interface{}) bool {
					return !old.(bool) && new.(bool)
				},
				customdiff.ComputedIf("ipv6_address", func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) bool {
					return d.Get("ipv6").(bool)
				}),
			),
			// Forces replacement when IPv6 has attribute changes to `false`
			// https://github.com/parspack/terraform-provider-parspack/issues/1104
			customdiff.ForceNewIfChange("ipv6",
				func(ctx context.Context, old, new, meta interface{}) bool {
					return old.(bool) && !new.(bool)
				},
			),
		),
	}
}

func resourceAbrhaVmCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	image := d.Get("image").(string)

	// Build up our creation options
	opts := &goApiAbrha.VmCreateRequest{
		Image:  goApiAbrha.VmCreateImage{},
		Name:   d.Get("name").(string),
		Region: d.Get("region").(string),
		Size:   d.Get("size").(string),
		Tags:   tag.ExpandTags(d.Get("tags").(*schema.Set).List()),
	}

	imageId, err := strconv.Atoi(image)
	if err == nil {
		// The image field is provided as an ID (number).
		opts.Image.ID = imageId
	} else {
		opts.Image.Slug = image
	}

	if attr, ok := d.GetOk("backups"); ok {
		_, exist := d.GetOk("backup_policy")
		if exist && !attr.(bool) { // Check there is no backup_policy specified when backups are disabled.
			return diag.FromErr(errVmBackupPolicy)
		}
		opts.Backups = attr.(bool)
	}

	if attr, ok := d.GetOk("ipv6"); ok {
		opts.IPv6 = attr.(bool)
	}

	if attr, ok := d.GetOk("private_networking"); ok {
		opts.PrivateNetworking = attr.(bool)
	}

	if attr, ok := d.GetOk("user_data"); ok {
		opts.UserData = attr.(string)
	}

	if attr, ok := d.GetOk("volume_ids"); ok {
		for _, id := range attr.(*schema.Set).List() {
			if id == nil {
				continue
			}
			volumeId := id.(string)
			if volumeId == "" {
				continue
			}

			opts.Volumes = append(opts.Volumes, goApiAbrha.VmCreateVolume{
				ID: volumeId,
			})
		}
	}

	if attr, ok := d.GetOk("monitoring"); ok {
		opts.Monitoring = attr.(bool)
	}

	if attr, ok := d.GetOkExists("vm_agent"); ok {
		opts.WithVmAgent = goApiAbrha.PtrTo(attr.(bool))
	}

	if attr, ok := d.GetOk("vpc_uuid"); ok {
		opts.VPCUUID = attr.(string)
	}

	// Get configured ssh_keys
	if v, ok := d.GetOk("ssh_keys"); ok {
		expandedSshKeys, err := expandSshKeys(v.(*schema.Set).List())
		if err != nil {
			return diag.FromErr(err)
		}
		opts.SSHKeys = expandedSshKeys
	}

	log.Printf("[DEBUG] Vm create configuration: %#v", opts)
	// Get configured backup_policy
	if policy, ok := d.GetOk("backup_policy"); ok {
		if !d.Get("backups").(bool) {
			return diag.FromErr(errVmBackupPolicy)
		}

		backupPolicy, err := expandBackupPolicy(policy)
		if err != nil {
			return diag.FromErr(err)
		}
		opts.BackupPolicy = backupPolicy
	}

	log.Printf("[DEBUG] Vm create configuration: %#v", opts)

	vmRoot, _, err := client.Vms.Create(context.Background(), opts)
	if err != nil {
		return diag.Errorf("Error creating Vm: %s", err)
	}

	vm := vmRoot.Vm

	// Assign the vms id
	d.SetId(vm.ID)
	log.Printf("[INFO] Vm ID: %s", d.Id())

	log.Printf("[INFO] Created vm action_id %d", vmRoot.Links.Actions[0].ID)
	action, _, err := client.VmActions.Get(context.Background(), d.Id(), vmRoot.Links.Actions[0].ID)
	if err != nil {
		return diag.Errorf("error getting VM action: %s", err)
	}

	// wait for job to complete
	actionId := vmRoot.Links.Actions[0].ID
	log.Printf("[DEBUG] Wating for create vm action (%d) to success...", actionId)
	if err := util.WaitForAction(client, action); err != nil {
		return diag.Errorf("Error waiting for create vm action for vm (%s) and action (%d) to complate: %s", d.Id(), actionId, err)
	}

	//Ensure Vm status has moved to "active."
	_, err = waitForVmAttribute(ctx, d, "active", []string{"new"}, "status", schema.TimeoutCreate, meta)

	if err != nil {
		return diag.Errorf("Error waiting for vm (%s) to become ready: %s", d.Id(), err)
	}

	// waitForVmAttribute updates the Vm's state and calls setVmAttributes.
	// So there is no need to call resourceAbrhaVmRead and add additional API calls.
	return nil
}

func resourceAbrhaVmRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	id := d.Id()

	// Retrieve the vm properties for updating the state
	vm, resp, err := client.Vms.Get(context.Background(), id)
	if err != nil {
		// check if the vm no longer exists.
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[WARN] Abrha Vm (%s) not found", d.Id())
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving Vm: %s", err)
	}

	err = setVmAttributes(d, vm)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func setVmAttributes(d *schema.ResourceData, vm *goApiAbrha.Vm) error {
	// Note that the image attribute is not set here. It is intentionally allowed
	// to drift once the Vm has been created. This is to workaround the fact that
	// image slugs can move to point to images with a different ID. Image slugs are helpers
	// that always point to the most recent version of an image.
	// See: https://github.com/parspack/terraform-provider-parspack/issues/152
	d.Set("name", vm.Name)
	d.Set("urn", vm.URN())
	d.Set("region", vm.Region.Slug)
	d.Set("size", vm.Size.Slug)
	d.Set("price_hourly", vm.Size.PriceHourly)
	d.Set("price_monthly", vm.Size.PriceMonthly)
	d.Set("disk", vm.Disk)
	d.Set("vcpus", vm.Vcpus)
	d.Set("memory", vm.Memory)
	d.Set("status", vm.Status)
	d.Set("locked", vm.Locked)
	d.Set("created_at", vm.Created)
	d.Set("vpc_uuid", vm.VPCUUID)

	d.Set("ipv4_address", FindIPv4AddrByType(vm, "public"))
	d.Set("ipv4_address_private", FindIPv4AddrByType(vm, "private"))
	d.Set("ipv6_address", strings.ToLower(FindIPv6AddrByType(vm, "public")))

	if features := vm.Features; features != nil {
		d.Set("backups", containsAbrhaVmFeature(features, "backups"))
		d.Set("ipv6", containsAbrhaVmFeature(features, "ipv6"))
		d.Set("private_networking", containsAbrhaVmFeature(features, "private_networking"))
		d.Set("monitoring", containsAbrhaVmFeature(features, "monitoring"))
	}

	if err := d.Set("volume_ids", flattenAbrhaVmVolumeIds(vm.VolumeIDs)); err != nil {
		return fmt.Errorf("Error setting `volume_ids`: %+v", err)
	}

	if err := d.Set("tags", tag.FlattenTags(vm.Tags)); err != nil {
		return fmt.Errorf("Error setting `tags`: %+v", err)
	}

	// Initialize the connection info
	d.SetConnInfo(map[string]string{
		"type": "ssh",
		"host": FindIPv4AddrByType(vm, "public"),
	})

	return nil
}

func resourceAbrhaVmImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	// Retrieve the image from API during import
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()
	id := d.Id()

	vm, _, err := client.Vms.Get(context.Background(), id)
	if err != nil {
		return nil, fmt.Errorf("Error importing vm: %s", err)
	}

	if vm.Image.Slug != "" {
		d.Set("image", vm.Image.Slug)
	} else {
		d.Set("image", goApiAbrha.Stringify(vm.Image.ID))
	}

	// This is a non API attribute. So set to the default setting in the schema.
	d.Set("resize_disk", true)

	return []*schema.ResourceData{d}, nil
}

func FindIPv6AddrByType(d *goApiAbrha.Vm, addrType string) string {
	for _, addr := range d.Networks.V6 {
		if addr.Type == addrType {
			if ip := net.ParseIP(addr.IPAddress); ip != nil {
				return strings.ToLower(addr.IPAddress)
			}
		}
	}
	return ""
}

func FindIPv4AddrByType(d *goApiAbrha.Vm, addrType string) string {
	for _, addr := range d.Networks.V4 {
		if addr.Type == addrType {
			if ip := net.ParseIP(addr.IPAddress); ip != nil {
				return addr.IPAddress
			}
		}
	}
	return ""
}

func resourceAbrhaVmUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()
	var warnings []diag.Diagnostic

	id := d.Id()

	if d.HasChange("size") {
		newSize := d.Get("size")
		resizeDisk := d.Get("resize_disk").(bool)

		_, _, err := client.VmActions.PowerOff(context.Background(), id)
		if err != nil && !strings.Contains(err.Error(), "Vm is already powered off") {
			return diag.Errorf(
				"Error powering off Vm (%s): %s", d.Id(), err)
		}

		// Wait for power off
		_, err = waitForVmAttribute(ctx, d, "off", []string{"active"}, "status", schema.TimeoutUpdate, meta)
		if err != nil {
			return diag.Errorf(
				"Error waiting for vm (%s) to become powered off: %s", d.Id(), err)
		}

		// Resize the vm
		var action *goApiAbrha.Action
		action, _, err = client.VmActions.Resize(context.Background(), id, newSize.(string), resizeDisk)
		if err != nil {
			newErr := powerOnAndWait(ctx, d, meta)
			if newErr != nil {
				return diag.Errorf(
					"Error powering on vm (%s) after failed resize: %s", d.Id(), err)
			}
			return diag.Errorf(
				"Error resizing vm (%s): %s", d.Id(), err)
		}

		// Wait for the resize action to complete.
		if err = util.WaitForAction(client, action); err != nil {
			newErr := powerOnAndWait(ctx, d, meta)
			if newErr != nil {
				return diag.Errorf(
					"Error powering on vm (%s) after waiting for resize to finish: %s", d.Id(), err)
			}
			return diag.Errorf(
				"Error waiting for resize vm (%s) to finish: %s", d.Id(), err)
		}

		_, _, err = client.VmActions.PowerOn(context.Background(), id)

		if err != nil {
			return diag.Errorf(
				"Error powering on vm (%s) after resize: %s", d.Id(), err)
		}

		// Wait for power on
		_, err = waitForVmAttribute(ctx, d, "active", []string{"off"}, "status", schema.TimeoutUpdate, meta)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("name") {
		oldName, newName := d.GetChange("name")

		// Rename the vm
		_, _, err := client.VmActions.Rename(context.Background(), id, newName.(string))

		if err != nil {
			return diag.Errorf(
				"Error renaming vm (%s): %s", d.Id(), err)
		}

		// Wait for the name to change
		_, err = waitForVmAttribute(
			ctx, d, newName.(string), []string{"", oldName.(string)}, "name", schema.TimeoutUpdate, meta)

		if err != nil {
			return diag.Errorf(
				"Error waiting for rename vm (%s) to finish: %s", d.Id(), err)
		}
	}

	if d.HasChange("backups") {
		if d.Get("backups").(bool) {
			// Enable backups on vm
			var action *goApiAbrha.Action
			// Apply backup_policy if specified, otherwise use the default policy
			policy, ok := d.GetOk("backup_policy")
			if ok {
				backupPolicy, err := expandBackupPolicy(policy)
				if err != nil {
					return diag.FromErr(err)
				}
				action, _, err = client.VmActions.EnableBackupsWithPolicy(context.Background(), id, backupPolicy)
				if err != nil {
					return diag.Errorf(
						"Error enabling backups on vm (%s): %s", d.Id(), err)
				}
			} else {
				var err error
				action, _, err = client.VmActions.EnableBackups(context.Background(), id)
				if err != nil {
					return diag.Errorf(
						"Error enabling backups on vm (%s): %s", d.Id(), err)
				}
			}
			if err := util.WaitForAction(client, action); err != nil {
				return diag.Errorf("Error waiting for backups to be enabled for vm (%s): %s", d.Id(), err)
			}
		} else {
			// Disable backups on vm
			// Check there is no backup_policy specified
			_, ok := d.GetOk("backup_policy")
			if ok {
				return diag.FromErr(errVmBackupPolicy)
			}
			action, _, err := client.VmActions.DisableBackups(context.Background(), id)
			if err != nil {
				return diag.Errorf(
					"Error disabling backups on vm (%s): %s", d.Id(), err)
			}

			if err := util.WaitForAction(client, action); err != nil {
				return diag.Errorf("Error waiting for backups to be disabled for vm (%s): %s", d.Id(), err)
			}
		}
	}

	if d.HasChange("backup_policy") {
		_, ok := d.GetOk("backup_policy")
		if ok {
			if !d.Get("backups").(bool) {
				return diag.FromErr(errVmBackupPolicy)
			}

			_, new := d.GetChange("backup_policy")
			newPolicy, err := expandBackupPolicy(new)
			if err != nil {
				return diag.FromErr(err)
			}

			action, _, err := client.VmActions.ChangeBackupPolicy(context.Background(), id, newPolicy)
			if err != nil {
				return diag.Errorf(
					"error changing backup policy on vm (%s): %s", d.Id(), err)
			}

			if err := util.WaitForAction(client, action); err != nil {
				return diag.Errorf("error waiting for backup policy to be changed for vm (%s): %s", d.Id(), err)
			}
		}
	}

	// As there is no way to disable private networking,
	// we only check if it needs to be enabled
	if d.HasChange("private_networking") && d.Get("private_networking").(bool) {
		_, _, err := client.VmActions.EnablePrivateNetworking(context.Background(), id)

		if err != nil {
			return diag.Errorf(
				"Error enabling private networking for vm (%s): %s", d.Id(), err)
		}

		// Wait for the private_networking to turn on
		_, err = waitForVmAttribute(
			ctx, d, "true", []string{"", "false"}, "private_networking", schema.TimeoutUpdate, meta)

		if err != nil {
			return diag.Errorf(
				"Error waiting for private networking to be enabled on for vm (%s): %s", d.Id(), err)
		}
	}

	// As there is no way to disable IPv6, we only check if it needs to be enabled
	if d.HasChange("ipv6") && d.Get("ipv6").(bool) {
		_, _, err := client.VmActions.EnableIPv6(context.Background(), id)
		if err != nil {
			return diag.Errorf(
				"Error turning on ipv6 for vm (%s): %s", d.Id(), err)
		}

		// Wait for ipv6 to turn on
		_, err = waitForVmAttribute(
			ctx, d, "true", []string{"", "false"}, "ipv6", schema.TimeoutUpdate, meta)

		if err != nil {
			return diag.Errorf(
				"Error waiting for ipv6 to be turned on for vm (%s): %s", d.Id(), err)
		}

		warnings = append(warnings, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Enabling IPv6 requires additional OS-level configuration",
			Detail:   "When enabling IPv6 on an existing vm, additional OS-level configuration is required.",
		})
	}

	if d.HasChange("tags") {
		err := tag.SetTags(client, d, goApiAbrha.VmResourceType)
		if err != nil {
			return diag.Errorf("Error updating tags: %s", err)
		}
	}

	if d.HasChange("volume_ids") {
		oldIDs, newIDs := d.GetChange("volume_ids")
		newSet := func(ids []interface{}) map[string]struct{} {
			out := make(map[string]struct{}, len(ids))
			for _, id := range ids {
				out[id.(string)] = struct{}{}
			}
			return out
		}
		// leftDiff returns all elements in Left that are not in Right
		leftDiff := func(left, right map[string]struct{}) map[string]struct{} {
			out := make(map[string]struct{})
			for l := range left {
				if _, ok := right[l]; !ok {
					out[l] = struct{}{}
				}
			}
			return out
		}
		oldIDSet := newSet(oldIDs.(*schema.Set).List())
		newIDSet := newSet(newIDs.(*schema.Set).List())
		for volumeID := range leftDiff(newIDSet, oldIDSet) {
			action, _, err := client.StorageActions.Attach(context.Background(), volumeID, id)
			if err != nil {
				return diag.Errorf("Error attaching volume %q to vm (%s): %s", volumeID, d.Id(), err)
			}
			// can't fire >1 action at a time, so waiting for each is OK
			if err := util.WaitForAction(client, action); err != nil {
				return diag.Errorf("Error waiting for volume %q to attach to vm (%s): %s", volumeID, d.Id(), err)
			}
		}
		for volumeID := range leftDiff(oldIDSet, newIDSet) {
			err := detachVolumeIDOnVm(d, volumeID, meta)
			if err != nil {
				return diag.Errorf("Error detaching volume %q on vm %s: %s", volumeID, d.Id(), err)

			}
		}
	}

	readErr := resourceAbrhaVmRead(ctx, d, meta)
	if readErr != nil {
		warnings = append(warnings, readErr...)
	}

	return warnings
}

func resourceAbrhaVmDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	id := d.Id()

	_, err := waitForVmAttribute(
		ctx, d, "false", []string{"", "true"}, "locked", schema.TimeoutDelete, meta)

	if err != nil {
		return diag.Errorf(
			"Error waiting for vm to be unlocked for destroy (%s): %s", d.Id(), err)
	}

	shutdown := d.Get("graceful_shutdown").(bool)
	if shutdown {
		log.Printf("[INFO] Shutting down vm: %s", d.Id())

		// Shutdown the vm
		// DO API doesn't return an error if we try to shutdown an already shutdown vm
		_, _, err = client.VmActions.Shutdown(context.Background(), id)
		if err != nil {
			return diag.Errorf(
				"Error shutting down the the vm (%s): %s", d.Id(), err)
		}

		// Wait for shutdown
		_, err = waitForVmAttribute(ctx, d, "off", []string{"active"}, "status", schema.TimeoutDelete, meta)
		if err != nil {
			return diag.Errorf("Error waiting for vm (%s) to become off: %s", d.Id(), err)
		}
	}

	log.Printf("[INFO] Trying to Detach Storage Volumes (if any) from vm: %s", d.Id())
	err = detachVolumesFromVm(d, meta)
	if err != nil {
		return diag.Errorf(
			"Error detaching the volumes from the vm (%s): %s", d.Id(), err)
	}

	log.Printf("[INFO] Deleting vm: %s", d.Id())

	// Destroy the vm
	resp, err := client.Vms.Delete(context.Background(), id)

	// Handle already destroyed vms
	if err != nil && resp.StatusCode == 404 {
		return nil
	}

	_, err = waitForVmDestroy(ctx, d, meta)
	if err != nil && strings.Contains(err.Error(), "404") {
		return nil
	} else if err != nil {
		return diag.Errorf("Error deleting vm: %s", err)
	}

	return nil
}

func waitForVmDestroy(ctx context.Context, d *schema.ResourceData, meta interface{}) (interface{}, error) {
	log.Printf("[INFO] Waiting for vm (%s) to be destroyed", d.Id())

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"active", "off"},
		Target:     []string{"archived"},
		Refresh:    vmStateRefreshFunc(ctx, d, "status", meta),
		Timeout:    60 * time.Second,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	return stateConf.WaitForStateContext(ctx)
}

func waitForVmAttribute(
	ctx context.Context, d *schema.ResourceData, target string, pending []string, attribute string, timeoutKey string, meta interface{}) (interface{}, error) {
	// Wait for the vm so we can get the networking attributes
	// that show up after a while
	log.Printf(
		"[INFO] Waiting for vm (%s) to have %s of %s",
		d.Id(), attribute, target)

	stateConf := &retry.StateChangeConf{
		Pending:    pending,
		Target:     []string{target},
		Refresh:    vmStateRefreshFunc(ctx, d, attribute, meta),
		Timeout:    d.Timeout(timeoutKey),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,

		// This is a hack around DO API strangeness.
		// https://github.com/hashicorp/terraform/issues/481
		//
		NotFoundChecks: 60,
	}

	return stateConf.WaitForStateContext(ctx)
}

// TODO This function still needs a little more refactoring to make it
// cleaner and more efficient
func vmStateRefreshFunc(
	ctx context.Context, d *schema.ResourceData, attribute string, meta interface{}) retry.StateRefreshFunc {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()
	return func() (interface{}, string, error) {
		id := d.Id()

		// Retrieve the vm properties
		vm, _, err := client.Vms.Get(context.Background(), id)
		if err != nil {
			return nil, "", fmt.Errorf("Error retrieving vm: %s", err)
		}

		err = setVmAttributes(d, vm)
		if err != nil {
			return nil, "", err
		}

		// If the vm is locked, continue waiting. We can
		// only perform actions on unlocked vms, so it's
		// pointless to look at that status
		if d.Get("locked").(bool) {
			log.Println("[DEBUG] Vm is locked, skipping status check and retrying")
			return nil, "", nil
		}

		// See if we can access our attribute
		if attr, ok := d.GetOkExists(attribute); ok {
			switch attr := attr.(type) {
			case bool:
				return &vm, strconv.FormatBool(attr), nil
			default:
				return &vm, attr.(string), nil
			}
		}

		return nil, "", nil
	}
}

// Powers on the vm and waits for it to be active
func powerOnAndWait(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	id := d.Id()

	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()
	_, _, err := client.VmActions.PowerOn(context.Background(), id)
	if err != nil {
		return err
	}
	// this method is only used for vm updates so use that as the timeout parameter
	// Wait for power on
	_, err = waitForVmAttribute(ctx, d, "active", []string{"off"}, "status", schema.TimeoutUpdate, meta)
	if err != nil {
		return err
	}

	return nil
}

// Detach volumes from vm
func detachVolumesFromVm(d *schema.ResourceData, meta interface{}) error {
	var errors []error
	if attr, ok := d.GetOk("volume_ids"); ok {
		errors = make([]error, 0, attr.(*schema.Set).Len())
		for _, volumeID := range attr.(*schema.Set).List() {
			err := detachVolumeIDOnVm(d, volumeID.(string), meta)
			if err != nil {
				return err
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("Error detaching one or more volumes: %v", errors)
	}

	return nil
}

func detachVolumeIDOnVm(d *schema.ResourceData, volumeID string, meta interface{}) error {
	id := d.Id()

	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()
	action, _, err := client.StorageActions.DetachByVmID(context.Background(), volumeID, id)
	if err != nil {
		return fmt.Errorf("Error detaching volume %q from vm (%s): %s", volumeID, d.Id(), err)
	}
	// can't fire >1 action at a time, so waiting for each is OK
	if err := util.WaitForAction(client, action); err != nil {
		return fmt.Errorf("Error waiting for volume %q to detach from vm (%s): %s", volumeID, d.Id(), err)
	}

	return nil
}

func containsAbrhaVmFeature(features []string, name string) bool {
	for _, v := range features {
		if v == name {
			return true
		}
	}
	return false
}

func expandSshKeys(sshKeys []interface{}) ([]goApiAbrha.VmCreateSSHKey, error) {
	expandedSshKeys := make([]goApiAbrha.VmCreateSSHKey, len(sshKeys))
	for i, s := range sshKeys {
		sshKey := s.(string)

		var expandedSshKey goApiAbrha.VmCreateSSHKey
		if id, err := strconv.Atoi(sshKey); err == nil {
			expandedSshKey.ID = id
		} else {
			expandedSshKey.Fingerprint = sshKey
		}

		expandedSshKeys[i] = expandedSshKey
	}

	return expandedSshKeys, nil
}

func flattenAbrhaVmVolumeIds(volumeids []string) *schema.Set {
	flattenedVolumes := schema.NewSet(schema.HashString, []interface{}{})
	for _, v := range volumeids {
		flattenedVolumes.Add(v)
	}

	return flattenedVolumes
}

func expandBackupPolicy(v interface{}) (*goApiAbrha.VmBackupPolicyRequest, error) {
	var policy goApiAbrha.VmBackupPolicyRequest
	policyList := v.([]interface{})

	for _, rawPolicy := range policyList {
		policyMap, ok := rawPolicy.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("vm backup policy type assertion failed: expected map[string]interface{}, got %T", rawPolicy)
		}

		planVal, exists := policyMap["plan"]
		if !exists {
			return nil, errors.New("backup_policy plan key does not exist")
		}
		plan, ok := planVal.(string)
		if !ok {
			return nil, errors.New("backup_policy plan is not a string")
		}
		policy.Plan = plan

		weekdayVal, exists := policyMap["weekday"]
		if !exists {
			return nil, errors.New("backup_policy weekday key does not exist")
		}
		weekday, ok := weekdayVal.(string)
		if !ok {
			return nil, errors.New("backup_policy weekday is not a string")
		}
		policy.Weekday = weekday

		monthdayVal, exists := policyMap["monthday"]
		if !exists {
			return nil, errors.New("backup_policy monthday key does not exist")
		}
		monthday, ok := monthdayVal.(int)
		if !ok {
			return nil, errors.New("backup_policy monthday is not a string")
		}
		policy.Monthday = monthday

		hourVal, exists := policyMap["hour"]
		if !exists {
			return nil, errors.New("backup_policy hour key does not exist")
		}
		hour, ok := hourVal.(int)
		if !ok {
			return nil, errors.New("backup_policy hour is not an int")
		}
		policy.Hour = &hour
	}

	return &policy, nil
}
