package database

import (
	"context"

	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceAbrhaDatabaseCA() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAbrhaDatabaseCARead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"certificate": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAbrhaDatabaseCARead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()
	clusterID := d.Get("cluster_id").(string)
	d.SetId(clusterID)

	ca, _, err := client.Databases.GetCA(context.Background(), clusterID)
	if err != nil {
		return diag.Errorf("Error retrieving database CA certificate: %s", err)
	}

	d.Set("certificate", string(ca.Certificate))

	return nil
}
