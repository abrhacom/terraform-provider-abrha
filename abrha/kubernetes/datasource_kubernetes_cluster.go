package kubernetes

import (
	"context"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/abrhacom/terraform-provider-abrha/abrha/tag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceAbrhaKubernetesCluster() *schema.Resource {
	dsNodePoolSchema := nodePoolSchema(false)

	for _, k := range dsNodePoolSchema {
		k.Computed = true
		k.Required = false
		k.Default = nil
		k.ValidateFunc = nil
	}

	return &schema.Resource{
		ReadContext: dataSourceAbrhaKubernetesClusterRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"surge_upgrade": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"ha": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"vpc_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"cluster_subnet": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"service_subnet": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ipv4_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tag.TagsSchema(),

			"maintenance_policy": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"day": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"duration": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"node_pool": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: dsNodePoolSchema,
				},
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"kube_config": kubernetesConfigSchema(),

			"auto_upgrade": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAbrhaKubernetesClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	clusters, resp, err := client.Kubernetes.List(context.Background(), &goApiAbrha.ListOptions{})
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			return diag.Errorf("No clusters found")
		}

		return diag.Errorf("Error listing Kubernetes clusters: %s", err)
	}

	// select the correct cluster
	for _, c := range clusters {
		if c.Name == d.Get("name").(string) {
			d.SetId(c.ID)

			return parspackKubernetesClusterRead(client, c, d)
		}
	}

	return diag.Errorf("Unable to find cluster with name: %s", d.Get("name").(string))
}
