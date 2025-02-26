package reservedip

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceAbrhaReservedIP() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbrhaReservedIPCreate,
		UpdateContext: resourceAbrhaReservedIPUpdate,
		ReadContext:   resourceAbrhaReservedIPRead,
		DeleteContext: resourceAbrhaReservedIPDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceAbrhaReservedIPImport,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				StateFunc: func(val interface{}) string {
					// DO API V2 region slug is always lowercase
					return strings.ToLower(val.(string))
				},
			},
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the uniform resource name for the reserved ip",
			},
			"ip_address": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IsIPv4Address,
			},
			"vm_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAbrhaReservedIPCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	log.Printf("[INFO] Creating a reserved IP in a region")
	regionOpts := &goApiAbrha.ReservedIPCreateRequest{
		Region: d.Get("region").(string),
	}

	log.Printf("[DEBUG] Reserved IP create: %#v", regionOpts)
	reservedIP, _, err := client.ReservedIPs.Create(context.Background(), regionOpts)
	if err != nil {
		return diag.Errorf("Error creating reserved IP: %s", err)
	}

	d.SetId(reservedIP.IP)

	if v, ok := d.GetOk("vm_id"); ok {
		log.Printf("[INFO] Assigning the reserved IP to the Vm %s", v.(string))
		action, _, err := client.ReservedIPActions.Assign(context.Background(), d.Id(), v.(string))
		if err != nil {
			return diag.Errorf(
				"Error Assigning reserved IP (%s) to the Vm: %s", d.Id(), err)
		}

		_, unassignedErr := waitForReservedIPReady(ctx, d, "completed", []string{"new", "in-progress"}, "status", meta, action.ID)
		if unassignedErr != nil {
			return diag.Errorf(
				"Error waiting for reserved IP (%s) to be assigned: %s", d.Id(), unassignedErr)
		}
	}

	return resourceAbrhaReservedIPRead(ctx, d, meta)
}

func resourceAbrhaReservedIPUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	if d.HasChange("vm_id") {
		if v, ok := d.GetOk("vm_id"); ok {
			log.Printf("[INFO] Assigning the reserved IP %s to the Vm %d", d.Id(), v.(int))
			action, _, err := client.ReservedIPActions.Assign(context.Background(), d.Id(), v.(string))
			if err != nil {
				return diag.Errorf(
					"Error assigning reserved IP (%s) to the Vm: %s", d.Id(), err)
			}

			_, unassignedErr := waitForReservedIPReady(ctx, d, "completed", []string{"new", "in-progress"}, "status", meta, action.ID)
			if unassignedErr != nil {
				return diag.Errorf(
					"Error waiting for reserved IP (%s) to be Assigned: %s", d.Id(), unassignedErr)
			}
		} else {
			log.Printf("[INFO] Unassigning the reserved IP %s", d.Id())
			action, _, err := client.ReservedIPActions.Unassign(context.Background(), d.Id())
			if err != nil {
				return diag.Errorf(
					"Error unassigning reserved IP (%s): %s", d.Id(), err)
			}

			_, unassignedErr := waitForReservedIPReady(ctx, d, "completed", []string{"new", "in-progress"}, "status", meta, action.ID)
			if unassignedErr != nil {
				return diag.Errorf(
					"Error waiting for reserved IP (%s) to be Unassigned: %s", d.Id(), unassignedErr)
			}
		}
	}

	return resourceAbrhaReservedIPRead(ctx, d, meta)
}

func resourceAbrhaReservedIPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	log.Printf("[INFO] Reading the details of the reserved IP %s", d.Id())
	reservedIP, resp, err := client.ReservedIPs.Get(context.Background(), d.Id())
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[WARN] Reserved IP (%s) not found", d.Id())
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving reserved IP: %s", err)
	}

	if _, ok := d.GetOk("vm_id"); ok && reservedIP.Vm != nil {
		d.Set("region", reservedIP.Vm.Region.Slug)
		d.Set("vm_id", reservedIP.Vm.ID)
	} else {
		d.Set("region", reservedIP.Region.Slug)
	}

	d.Set("ip_address", reservedIP.IP)
	d.Set("urn", reservedIP.URN())

	return nil
}

func resourceAbrhaReservedIPDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	if _, ok := d.GetOk("vm_id"); ok {
		log.Printf("[INFO] Unassigning the reserved IP from the Vm")
		action, resp, err := client.ReservedIPActions.Unassign(context.Background(), d.Id())
		if resp.StatusCode != 422 {
			if err != nil {
				return diag.Errorf(
					"Error unassigning reserved IP (%s) from the vm: %s", d.Id(), err)
			}

			_, unassignedErr := waitForReservedIPReady(ctx, d, "completed", []string{"new", "in-progress"}, "status", meta, action.ID)
			if unassignedErr != nil {
				return diag.Errorf(
					"Error waiting for reserved IP (%s) to be unassigned: %s", d.Id(), unassignedErr)
			}
		} else {
			log.Printf("[DEBUG] Couldn't unassign reserved IP (%s) from vm, possibly out of sync: %s", d.Id(), err)
		}
	}

	log.Printf("[INFO] Deleting reserved IP: %s", d.Id())
	_, err := client.ReservedIPs.Delete(context.Background(), d.Id())
	if err != nil {
		return diag.Errorf("Error deleting reserved IP: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceAbrhaReservedIPImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()
	reservedIP, resp, err := client.ReservedIPs.Get(context.Background(), d.Id())
	if resp.StatusCode != 404 {
		if err != nil {
			return nil, err
		}

		d.Set("ip_address", reservedIP.IP)
		d.Set("urn", reservedIP.URN())
		d.Set("region", reservedIP.Region.Slug)

		if reservedIP.Vm != nil {
			d.Set("vm_id", reservedIP.Vm.ID)
		}
	}

	return []*schema.ResourceData{d}, nil
}

func waitForReservedIPReady(
	ctx context.Context, d *schema.ResourceData, target string, pending []string, attribute string, meta interface{}, actionID int) (interface{}, error) {
	log.Printf(
		"[INFO] Waiting for reserved IP (%s) to have %s of %s",
		d.Id(), attribute, target)

	stateConf := &retry.StateChangeConf{
		Pending:    pending,
		Target:     []string{target},
		Refresh:    newReservedIPStateRefreshFunc(d, attribute, meta, actionID),
		Timeout:    60 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,

		NotFoundChecks: 60,
	}

	return stateConf.WaitForStateContext(ctx)
}

func newReservedIPStateRefreshFunc(
	d *schema.ResourceData, attribute string, meta interface{}, actionID int) retry.StateRefreshFunc {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()
	return func() (interface{}, string, error) {

		log.Printf("[INFO] Assigning the reserved IP to the Vm")
		action, _, err := client.ReservedIPActions.Get(context.Background(), d.Id(), actionID)
		if err != nil {
			return nil, "", fmt.Errorf("Error retrieving reserved IP (%s) ActionId (%d): %s", d.Id(), actionID, err)
		}

		log.Printf("[INFO] The reserved IP Action Status is %s", action.Status)
		return &action, action.Status, nil
	}
}
