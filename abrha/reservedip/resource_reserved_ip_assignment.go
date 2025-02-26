package reservedip

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceAbrhaReservedIPAssignment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbrhaReservedIPAssignmentCreate,
		ReadContext:   resourceAbrhaReservedIPAssignmentRead,
		DeleteContext: resourceAbrhaReservedIPAssignmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceAbrhaReservedIPAssignmentImport,
		},

		Schema: map[string]*schema.Schema{
			"ip_address": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsIPv4Address,
			},
			"vm_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
		},
	}
}

func resourceAbrhaReservedIPAssignmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	ipAddress := d.Get("ip_address").(string)
	vmID := d.Get("vm_id").(string)

	log.Printf("[INFO] Assigning the reserved IP (%s) to the Vm %s", ipAddress, vmID)
	action, _, err := client.ReservedIPActions.Assign(context.Background(), ipAddress, vmID)
	if err != nil {
		return diag.Errorf(
			"Error Assigning reserved IP (%s) to the vm: %s", ipAddress, err)
	}

	_, unassignedErr := waitForReservedIPAssignmentReady(ctx, d, "completed", []string{"new", "in-progress"}, "status", meta, action.ID)
	if unassignedErr != nil {
		return diag.Errorf(
			"Error waiting for reserved IP (%s) to be Assigned: %s", ipAddress, unassignedErr)
	}

	d.SetId(id.PrefixedUniqueId(fmt.Sprintf("%s-%s-", vmID, ipAddress)))
	return resourceAbrhaReservedIPAssignmentRead(ctx, d, meta)
}

func resourceAbrhaReservedIPAssignmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	ipAddress := d.Get("ip_address").(string)
	vmID := d.Get("vm_id").(string)

	log.Printf("[INFO] Reading the details of the reserved IP %s", ipAddress)
	reservedIP, _, err := client.ReservedIPs.Get(context.Background(), ipAddress)
	if err != nil {
		return diag.Errorf("Error retrieving reserved IP: %s", err)
	}

	if reservedIP.Vm == nil || reservedIP.Vm.ID != vmID {
		log.Printf("[INFO] A Vm was detected on the reserved IP.")
		d.SetId("")
	}

	return nil
}

func resourceAbrhaReservedIPAssignmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	ipAddress := d.Get("ip_address").(string)
	vmID := d.Get("vm_id").(string)

	log.Printf("[INFO] Reading the details of the reserved IP %s", ipAddress)
	reservedIP, _, err := client.ReservedIPs.Get(context.Background(), ipAddress)
	if err != nil {
		return diag.Errorf("Error retrieving reserved IP: %s", err)
	}

	if reservedIP.Vm.ID == vmID {
		log.Printf("[INFO] Unassigning the reserved IP from the Vm")
		action, _, err := client.ReservedIPActions.Unassign(context.Background(), ipAddress)
		if err != nil {
			return diag.Errorf("Error unassigning reserved IP (%s) from the vm: %s", ipAddress, err)
		}

		_, unassignedErr := waitForReservedIPAssignmentReady(ctx, d, "completed", []string{"new", "in-progress"}, "status", meta, action.ID)
		if unassignedErr != nil {
			return diag.Errorf(
				"Error waiting for reserved IP (%s) to be unassigned: %s", ipAddress, unassignedErr)
		}
	} else {
		log.Printf("[INFO] reserved IP already unassigned, removing from state.")
	}

	d.SetId("")
	return nil
}

func waitForReservedIPAssignmentReady(
	ctx context.Context, d *schema.ResourceData, target string, pending []string, attribute string, meta interface{}, actionID int) (interface{}, error) {
	log.Printf(
		"[INFO] Waiting for reserved IP (%s) to have %s of %s",
		d.Get("ip_address").(string), attribute, target)

	stateConf := &retry.StateChangeConf{
		Pending:    pending,
		Target:     []string{target},
		Refresh:    newReservedIPAssignmentStateRefreshFunc(d, attribute, meta, actionID),
		Timeout:    60 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,

		NotFoundChecks: 60,
	}

	return stateConf.WaitForStateContext(ctx)
}

func newReservedIPAssignmentStateRefreshFunc(
	d *schema.ResourceData, attribute string, meta interface{}, actionID int) retry.StateRefreshFunc {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()
	return func() (interface{}, string, error) {

		log.Printf("[INFO] Refreshing the reserved IP state")
		action, _, err := client.ReservedIPActions.Get(context.Background(), d.Get("ip_address").(string), actionID)
		if err != nil {
			return nil, "", fmt.Errorf("Error retrieving reserved IP (%s) ActionId (%d): %s", d.Get("ip_address").(string), actionID, err)
		}

		log.Printf("[INFO] The reserved IP Action Status is %s", action.Status)
		return &action, action.Status, nil
	}
}

func resourceAbrhaReservedIPAssignmentImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if strings.Contains(d.Id(), ",") {
		s := strings.Split(d.Id(), ",")
		d.SetId(id.PrefixedUniqueId(fmt.Sprintf("%s-%s-", s[1], s[0])))
		d.Set("ip_address", s[0])
		vmID, err := strconv.Atoi(s[1])
		if err != nil {
			return nil, err
		}
		d.Set("vm_id", vmID)
	} else {
		return nil, errors.New("must use the reserved IP and the ID of the Vm joined with a comma (e.g. `ip_address,vm_id`)")
	}

	return []*schema.ResourceData{d}, nil
}
