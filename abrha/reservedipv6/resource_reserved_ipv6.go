package reservedipv6

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

func ResourceAbrhaReservedIPV6() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbrhaReservedIPV6Create,
		ReadContext:   resourceAbrhaReservedIPV6Read,
		DeleteContext: resourceAbrhaReservedIPV6Delete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceAbrhaReservedIPV6Import,
		},

		Schema: map[string]*schema.Schema{
			"region_slug": {
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
				Description: "the uniform resource name for the reserved ipv6",
			},
			"ip": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IsIPv6Address,
			},
			"vm_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAbrhaReservedIPV6Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	log.Printf("[INFO] Creating a reserved IPv6 in a region")
	regionOpts := &goApiAbrha.ReservedIPV6CreateRequest{
		Region: d.Get("region_slug").(string),
	}

	log.Printf("[DEBUG] Reserved IPv6 create: %#v", regionOpts)
	reservedIP, _, err := client.ReservedIPV6s.Create(context.Background(), regionOpts)
	if err != nil {
		return diag.Errorf("Error creating reserved IPv6: %s", err)
	}

	d.SetId(reservedIP.IP)

	return resourceAbrhaReservedIPV6Read(ctx, d, meta)
}

func resourceAbrhaReservedIPV6Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	log.Printf("[INFO] Reading the details of the reserved IPv6 %s", d.Id())
	reservedIP, resp, err := client.ReservedIPV6s.Get(context.Background(), d.Id())
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[WARN] Reserved IPv6 (%s) not found", d.Id())
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving reserved IPv6: %s", err)
	}

	if _, ok := d.GetOk("vm_id"); ok && reservedIP.Vm != nil {
		d.Set("region_slug", reservedIP.Vm.Region.Slug)
		d.Set("vm_id", reservedIP.Vm.ID)
	} else {
		d.Set("region_slug", reservedIP.RegionSlug)
	}

	d.Set("ip", reservedIP.IP)
	d.Set("urn", reservedIP.URN())

	return nil
}

func resourceAbrhaReservedIPV6Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	if _, ok := d.GetOk("vm_id"); ok {
		log.Printf("[INFO] Unassigning the reserved IPv6 from the Vm")
		action, resp, err := client.ReservedIPV6Actions.Unassign(context.Background(), d.Id())
		if resp.StatusCode != 422 {
			if err != nil {
				return diag.Errorf(
					"Error unassigning reserved IPv6 (%s) from the vm: %s", d.Id(), err)
			}

			_, unassignedErr := waitForReservedIPV6Ready(ctx, d, "completed", []string{"new", "in-progress"}, "status", meta, action.ID)
			if unassignedErr != nil {
				return diag.Errorf(
					"Error waiting for reserved IPv6 (%s) to be unassigned: %s", d.Id(), unassignedErr)
			}
		} else {
			log.Printf("[DEBUG] Couldn't unassign reserved IPv6 (%s) from vm, possibly out of sync: %s", d.Id(), err)
		}
	}

	log.Printf("[INFO] Deleting reserved IPv6: %s", d.Id())
	_, err := client.ReservedIPV6s.Delete(context.Background(), d.Id())
	if err != nil && strings.Contains(err.Error(), "404") {
		return diag.Errorf("Error deleting reserved IPv6: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceAbrhaReservedIPV6Import(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()
	reservedIP, resp, err := client.ReservedIPV6s.Get(context.Background(), d.Id())
	if resp.StatusCode != 404 {
		if err != nil {
			return nil, err
		}

		d.Set("ip", reservedIP.IP)
		d.Set("urn", reservedIP.URN())
		d.Set("region_slug", reservedIP.RegionSlug)

		if reservedIP.Vm != nil {
			d.Set("vm_id", reservedIP.Vm.ID)
		}
	}

	return []*schema.ResourceData{d}, nil
}

func waitForReservedIPV6Ready(
	ctx context.Context, d *schema.ResourceData, target string, pending []string, attribute string, meta interface{}, actionID int) (interface{}, error) {
	log.Printf(
		"[INFO] Waiting for reserved IPv6 (%s) to have %s of %s",
		d.Id(), attribute, target)

	stateConf := &retry.StateChangeConf{
		Pending:    pending,
		Target:     []string{target},
		Refresh:    newReservedIPV6StateRefreshFunc(d, meta, actionID),
		Timeout:    60 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,

		NotFoundChecks: 60,
	}

	return stateConf.WaitForStateContext(ctx)
}

func newReservedIPV6StateRefreshFunc(
	d *schema.ResourceData, meta interface{}, actionID int) retry.StateRefreshFunc {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()
	return func() (interface{}, string, error) {

		log.Printf("[INFO] Assigning the reserved IPv6 to the Vm")
		action, _, err := client.Actions.Get(context.Background(), actionID)
		if err != nil {
			return nil, "", fmt.Errorf("error retrieving reserved IPv6 (%s) ActionId (%d): %s", d.Id(), actionID, err)
		}

		log.Printf("[INFO] The reserved IPv6 Action Status is %s", action.Status)
		return &action, action.Status, nil
	}
}
