package reservedipv6

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

func ResourceAbrhaReservedIPV6Assignment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbrhaReservedIPV6AssignmentCreate,
		ReadContext:   resourceAbrhaReservedIPV6AssignmentRead,
		DeleteContext: resourceAbrhaReservedIPV6AssignmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceAbrhaReservedIPV6AssignmentImport,
		},

		Schema: map[string]*schema.Schema{
			"ip": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsIPv6Address,
			},
			"vm_id": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
		},
	}
}

func resourceAbrhaReservedIPV6AssignmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	ipAddress := d.Get("ip").(string)
	vmID := d.Get("vm_id").(string)

	log.Printf("[INFO] Assigning the reserved IPv6 (%s) to the Vm %s", ipAddress, vmID)
	action, _, err := client.ReservedIPV6Actions.Assign(context.Background(), ipAddress, vmID)
	if err != nil {
		return diag.Errorf(
			"Error Assigning reserved IPv6 (%s) to the vm: %s", ipAddress, err)
	}

	_, assignedErr := waitForReservedIPV6AssignmentReady(ctx, d, "completed", []string{"new", "in-progress"}, "status", meta, action.ID)
	if assignedErr != nil {
		return diag.Errorf(
			"Error waiting for reserved IPv6 (%s) to be Assigned: %s", ipAddress, assignedErr)
	}

	d.SetId(id.PrefixedUniqueId(fmt.Sprintf("%s-%s-", vmID, ipAddress)))

	return resourceAbrhaReservedIPV6AssignmentRead(ctx, d, meta)

}

func resourceAbrhaReservedIPV6AssignmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	ipAddress := d.Get("ip").(string)
	vmID := d.Get("vm_id")

	log.Printf("[INFO] Reading the details of the reserved IPv6 %s", ipAddress)
	reservedIPv6, _, err := client.ReservedIPV6s.Get(context.Background(), ipAddress)
	if err != nil {
		return diag.Errorf("Error retrieving reserved IPv6: %s", err)
	}

	if reservedIPv6.Vm == nil || reservedIPv6.Vm.ID != vmID {
		// log.Printf("[INFO] Vm assignment was unsuccessful on the reserved IPv6.")
		return diag.Errorf("Error assigning reserved IPv6 %s to vmID %s", ipAddress, vmID)
	}

	return nil
}

func resourceAbrhaReservedIPV6AssignmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	ipAddress := d.Get("ip").(string)
	vmID := d.Get("vm_id")

	log.Printf("[INFO] Reading the details of the reserved IPv6 %s", ipAddress)
	reservedIPv6, _, err := client.ReservedIPV6s.Get(context.Background(), ipAddress)
	if err != nil {
		return diag.Errorf("Error retrieving reserved IPv6: %s", err)
	}

	if reservedIPv6.Vm.ID == vmID {
		log.Printf("[INFO] Unassigning the reserved IPv6 from the Vm")
		action, _, err := client.ReservedIPV6Actions.Unassign(context.Background(), ipAddress)
		if err != nil {
			return diag.Errorf("Error unassigning reserved IPv6 (%s) from the vm: %s", ipAddress, err)
		}

		_, unassignedErr := waitForReservedIPV6AssignmentReady(ctx, d, "completed", []string{"new", "in-progress"}, "status", meta, action.ID)
		if unassignedErr != nil {
			return diag.Errorf(
				"Error waiting for reserved IPv6 (%s) to be unassigned: %s", ipAddress, unassignedErr)
		}
	} else {
		log.Printf("[INFO] reserved IPv6 already unassigned, removing from state.")
	}

	d.SetId("")
	return nil
}

func waitForReservedIPV6AssignmentReady(
	ctx context.Context, d *schema.ResourceData, target string, pending []string, attribute string, meta interface{}, actionID int) (interface{}, error) {
	log.Printf(
		"[INFO] Waiting for reserved IPv6 (%s) to have %s of %s",
		d.Get("ip").(string), attribute, target)

	stateConf := &retry.StateChangeConf{
		Pending:    pending,
		Target:     []string{target},
		Refresh:    newReservedIPV6AssignmentStateRefreshFunc(d, meta, actionID),
		Timeout:    60 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,

		NotFoundChecks: 60,
	}

	return stateConf.WaitForStateContext(ctx)
}

func newReservedIPV6AssignmentStateRefreshFunc(
	d *schema.ResourceData, meta interface{}, actionID int) retry.StateRefreshFunc {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()
	return func() (interface{}, string, error) {

		log.Printf("[INFO] Refreshing the reserved IPv6 state")
		action, _, err := client.Actions.Get(context.Background(), actionID)
		if err != nil {
			return nil, "", fmt.Errorf("error retrieving reserved IPv6 (%s) ActionId (%d): %s", d.Get("ip_address"), actionID, err)
		}

		log.Printf("[INFO] The reserved IPv6 Action Status is %s", action.Status)
		return &action, action.Status, nil
	}
}

func resourceAbrhaReservedIPV6AssignmentImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if strings.Contains(d.Id(), ",") {
		s := strings.Split(d.Id(), ",")
		d.SetId(id.PrefixedUniqueId(fmt.Sprintf("%s-%s-", s[1], s[0])))
		d.Set("ip", s[0])
		vmID, err := strconv.Atoi(s[1])
		if err != nil {
			return nil, err
		}
		d.Set("vm_id", vmID)
	} else {
		return nil, errors.New("must use the reserved IPv6 and the ID of the Vm joined with a comma (e.g. `ip,vm_id`)")
	}

	return []*schema.ResourceData{d}, nil
}
