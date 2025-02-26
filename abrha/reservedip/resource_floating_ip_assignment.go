package reservedip

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceAbrhaFloatingIPAssignment() *schema.Resource {
	return &schema.Resource{
		// TODO: Uncomment when dates for deprecation timeline are set.
		// DeprecationMessage: "This resource is deprecated and will be removed in a future release. Please use abrha_reserved_ip_assignment instead.",
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
