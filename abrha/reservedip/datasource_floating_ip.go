package reservedip

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceAbrhaFloatingIP() *schema.Resource {
	return &schema.Resource{
		// TODO: Uncomment when dates for deprecation timeline are set.
		// DeprecationMessage: "This data source is deprecated and will be removed in a future release. Please use abrha_reserved_ip instead.",
		ReadContext: dataSourceAbrhaFloatingIPRead,
		Schema: map[string]*schema.Schema{
			"ip_address": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "floating ip address",
				ValidateFunc: validation.NoZeroValues,
			},
			// computed attributes
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the uniform resource name for the floating ip",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the region that the floating ip is reserved to",
			},
			"vm_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the vm id that the floating ip has been assigned to.",
			},
		},
	}
}

func dataSourceAbrhaFloatingIPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := dataSourceAbrhaReservedIPRead(ctx, d, meta)
	if err != nil {
		return err
	}
	reservedIPURNtoFloatingIPURN(d)

	return nil
}
