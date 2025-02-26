package reservedip

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceAbrhaReservedIP() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAbrhaReservedIPRead,
		Schema: map[string]*schema.Schema{
			"ip_address": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "reserved ip address",
				ValidateFunc: validation.NoZeroValues,
			},
			// computed attributes
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the uniform resource name for the reserved ip",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the region that the reserved ip is reserved to",
			},
			"vm_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the vm id that the reserved ip has been assigned to.",
			},
		},
	}
}

func dataSourceAbrhaReservedIPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ipAddress := d.Get("ip_address").(string)
	d.SetId(ipAddress)

	return resourceAbrhaReservedIPRead(ctx, d, meta)
}
