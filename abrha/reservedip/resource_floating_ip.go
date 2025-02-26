package reservedip

import (
	"context"
	"strings"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceAbrhaFloatingIP() *schema.Resource {
	return &schema.Resource{
		// TODO: Uncomment when dates for deprecation timeline are set.
		// DeprecationMessage: "This resource is deprecated and will be removed in a future release. Please use abrha_reserved_ip instead.",
		CreateContext: resourceAbrhaFloatingIPCreate,
		UpdateContext: resourceAbrhaFloatingIPUpdate,
		ReadContext:   resourceAbrhaFloatingIPRead,
		DeleteContext: resourceAbrhaReservedIPDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceAbrhaFloatingIPImport,
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
				Description: "the uniform resource name for the floating ip",
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

func resourceAbrhaFloatingIPCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := resourceAbrhaReservedIPCreate(ctx, d, meta)
	if err != nil {
		return err
	}
	reservedIPURNtoFloatingIPURN(d)

	return nil
}

func resourceAbrhaFloatingIPUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := resourceAbrhaReservedIPUpdate(ctx, d, meta)
	if err != nil {
		return err
	}
	reservedIPURNtoFloatingIPURN(d)

	return nil
}

func resourceAbrhaFloatingIPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := resourceAbrhaReservedIPRead(ctx, d, meta)
	if err != nil {
		return err
	}
	reservedIPURNtoFloatingIPURN(d)

	return nil
}

func resourceAbrhaFloatingIPImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	_, err := resourceAbrhaReservedIPImport(ctx, d, meta)
	if err != nil {
		return nil, err
	}
	reservedIPURNtoFloatingIPURN(d)

	return []*schema.ResourceData{d}, nil
}

// reservedIPURNtoFloatingIPURN re-formats a reserved IP URN as floating IP URN.
// TODO: Remove when the projects' API changes return values.
func reservedIPURNtoFloatingIPURN(d *schema.ResourceData) {
	ip := d.Get("ip_address")
	d.Set("urn", goApiAbrha.FloatingIP{IP: ip.(string)}.URN())
}
