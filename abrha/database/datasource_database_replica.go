package database

import (
	"context"

	"github.com/abrhacom/terraform-provider-abrha/abrha/tag"
	"github.com/abrhacom/terraform-provider-abrha/abrha/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceAbrhaDatabaseReplica() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAbrhaDatabaseReplicaRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"private_network_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"host": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"private_host": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"uri": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"private_uri": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"database": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"user": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"password": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: tag.ValidateTag,
				},
				Set: util.HashStringIgnoreCase,
			},

			"storage_size_mib": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAbrhaDatabaseReplicaRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clusterId := d.Get("cluster_id").(string)
	name := d.Get("name").(string)
	d.SetId(makeReplicaId(clusterId, name))
	return resourceAbrhaDatabaseReplicaRead(ctx, d, meta)
}
