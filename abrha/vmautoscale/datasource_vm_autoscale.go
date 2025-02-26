package vmautoscale

import (
	"context"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceAbrhaVmAutoscale() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAbrhaVmAutoscaleRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "ID of the Vm autoscale pool",
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Name of the Vm autoscale pool",
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
			},
			"config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"min_instances": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Min number of members",
						},
						"max_instances": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Max number of members",
						},
						"target_cpu_utilization": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "CPU target threshold",
						},
						"target_memory_utilization": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Memory target threshold",
						},
						"cooldown_minutes": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Cooldown duration",
						},
						"target_number_instances": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Target number of members",
						},
					},
				},
			},
			"vm_template": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Vm size",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Vm region",
						},
						"image": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Vm image",
						},
						"tags": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "Vm tags",
						},
						"ssh_keys": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "Vm SSH keys",
						},
						"vpc_uuid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Vm VPC UUID",
						},
						"with_vm_agent": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Enable vm agent",
						},
						"project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Vm project ID",
						},
						"ipv6": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Enable vm IPv6",
						},
						"user_data": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Vm user data",
						},
					},
				},
			},
			"current_utilization": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"memory": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Average Memory utilization",
						},
						"cpu": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Average CPU utilization",
						},
					},
				},
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Vm autoscale pool status",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Vm autoscale pool create timestamp",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Vm autoscale pool update timestamp",
			},
		},
	}
}

func dataSourceAbrhaVmAutoscaleRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	var foundVmAutoscalePool *goApiAbrha.VmAutoscalePool
	if id, ok := d.GetOk("id"); ok {
		pool, _, err := client.VmAutoscale.Get(context.Background(), id.(string))
		if err != nil {
			return diag.Errorf("Error retrieving Vm autoscale pool: %v", err)
		}
		foundVmAutoscalePool = pool
	} else if name, ok := d.GetOk("name"); ok {
		vmAutoscalePoolList := make([]*goApiAbrha.VmAutoscalePool, 0)
		opts := &goApiAbrha.ListOptions{
			Page:    1,
			PerPage: 100,
		}
		// Paginate through all active resources
		for {
			pools, resp, err := client.VmAutoscale.List(context.Background(), opts)
			if err != nil {
				return diag.Errorf("Error listing Vm autoscale pools: %v", err)
			}
			vmAutoscalePoolList = append(vmAutoscalePoolList, pools...)
			if resp.Links.IsLastPage() {
				break
			}
			page, err := resp.Links.CurrentPage()
			if err != nil {
				break
			}
			opts.Page = page + 1
		}
		// Scan through the list to find a resource name match
		for i := range vmAutoscalePoolList {
			if vmAutoscalePoolList[i].Name == name {
				foundVmAutoscalePool = vmAutoscalePoolList[i]
				break
			}
		}
	} else {
		return diag.Errorf("Need to specify either a name or an id to look up the Vm autoscale pool")
	}
	if foundVmAutoscalePool == nil {
		return diag.Errorf("Vm autoscale pool not found")
	}

	d.SetId(foundVmAutoscalePool.ID)
	d.Set("name", foundVmAutoscalePool.Name)
	d.Set("config", flattenConfig(foundVmAutoscalePool.Config))
	d.Set("vm_template", flattenTemplate(foundVmAutoscalePool.VmTemplate))
	d.Set("current_utilization", flattenUtilization(foundVmAutoscalePool.CurrentUtilization))
	d.Set("status", foundVmAutoscalePool.Status)
	d.Set("created_at", foundVmAutoscalePool.CreatedAt.UTC().String())
	d.Set("updated_at", foundVmAutoscalePool.UpdatedAt.UTC().String())

	return nil
}

func flattenConfig(config *goApiAbrha.VmAutoscaleConfiguration) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, 1)
	if config != nil {
		r := make(map[string]interface{})
		r["min_instances"] = config.MinInstances
		r["max_instances"] = config.MaxInstances
		r["target_cpu_utilization"] = config.TargetCPUUtilization
		r["target_memory_utilization"] = config.TargetMemoryUtilization
		r["cooldown_minutes"] = config.CooldownMinutes
		r["target_number_instances"] = config.TargetNumberInstances
		result = append(result, r)
	}
	return result
}

func flattenTemplate(template *goApiAbrha.VmAutoscaleResourceTemplate) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, 1)
	if template != nil {
		r := make(map[string]interface{})
		r["size"] = template.Size
		r["region"] = template.Region
		r["image"] = template.Image
		r["vpc_uuid"] = template.VpcUUID
		r["with_vm_agent"] = template.WithVmAgent
		r["project_id"] = template.ProjectID
		r["ipv6"] = template.IPV6
		r["user_data"] = template.UserData

		tagSet := schema.NewSet(schema.HashString, []interface{}{})
		for _, tag := range template.Tags {
			tagSet.Add(tag)
		}
		r["tags"] = tagSet

		keySet := schema.NewSet(schema.HashString, []interface{}{})
		for _, key := range template.SSHKeys {
			keySet.Add(key)
		}
		r["ssh_keys"] = keySet
		result = append(result, r)
	}
	return result
}

func flattenUtilization(util *goApiAbrha.VmAutoscaleResourceUtilization) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, 1)
	if util != nil {
		r := make(map[string]interface{})
		r["memory"] = util.Memory
		r["cpu"] = util.CPU
		result = append(result, r)
	}
	return result
}
