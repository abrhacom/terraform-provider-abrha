package region

import (
	"context"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceAbrhaRegion() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAbrhaRegionRead,
		Schema: map[string]*schema.Schema{
			"slug": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sizes": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"features": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"available": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceAbrhaRegionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	regions, err := getAbrhaRegions(meta, nil)
	if err != nil {
		return diag.Errorf("Unable to load regions: %s", err)
	}

	slug := d.Get("slug").(string)

	var regionForSlug *interface{}
	for _, region := range regions {
		if region.(goApiAbrha.Region).Slug == slug {
			regionForSlug = &region
			break
		}
	}

	if regionForSlug == nil {
		return diag.Errorf("Region does not exist: %s", slug)
	}

	flattenedRegion, err := flattenRegion(*regionForSlug, meta, nil)
	if err != nil {
		return nil
	}

	if err := util.SetResourceDataFromMap(d, flattenedRegion); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id.UniqueId())
	return nil
}
