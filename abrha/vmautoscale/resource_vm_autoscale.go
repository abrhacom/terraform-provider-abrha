package vmautoscale

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceAbrhaVmAutoscale() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbrhaVmAutoscaleCreate,
		ReadContext:   resourceAbrhaVmAutoscaleRead,
		UpdateContext: resourceAbrhaVmAutoscaleUpdate,
		DeleteContext: resourceAbrhaVmAutoscaleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the Vm autoscale pool",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Vm autoscale pool",
			},
			"config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"min_instances": {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "Min number of members",
							ValidateFunc: validation.All(validation.NoZeroValues),
						},
						"max_instances": {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "Max number of members",
							ValidateFunc: validation.All(validation.NoZeroValues),
						},
						"target_cpu_utilization": {
							Type:         schema.TypeFloat,
							Optional:     true,
							Description:  "CPU target threshold",
							ValidateFunc: validation.All(validation.FloatBetween(0, 1)),
						},
						"target_memory_utilization": {
							Type:         schema.TypeFloat,
							Optional:     true,
							Description:  "Memory target threshold",
							ValidateFunc: validation.All(validation.FloatBetween(0, 1)),
						},
						"cooldown_minutes": {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "Cooldown duration",
							ValidateFunc: validation.All(validation.NoZeroValues),
						},
						"target_number_instances": {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "Target number of members",
							ValidateFunc: validation.All(validation.NoZeroValues),
						},
					},
				},
			},
			"vm_template": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Vm size",
						},
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Vm region",
						},
						"image": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Vm image",
						},
						"tags": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Description: "Vm tags",
						},
						"ssh_keys": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Required:    true,
							Description: "Vm SSH keys",
						},
						"vpc_uuid": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Vm VPC UUID",
						},
						"with_vm_agent": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Enable vm agent",
						},
						"project_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Vm project ID",
						},
						"ipv6": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Enable vm IPv6",
						},
						"user_data": {
							Type:        schema.TypeString,
							Optional:    true,
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

func expandConfig(config []interface{}) *goApiAbrha.VmAutoscaleConfiguration {
	if len(config) > 0 {
		poolConfig := config[0].(map[string]interface{})
		return &goApiAbrha.VmAutoscaleConfiguration{
			MinInstances:            uint64(poolConfig["min_instances"].(int)),
			MaxInstances:            uint64(poolConfig["max_instances"].(int)),
			TargetCPUUtilization:    poolConfig["target_cpu_utilization"].(float64),
			TargetMemoryUtilization: poolConfig["target_memory_utilization"].(float64),
			CooldownMinutes:         uint32(poolConfig["cooldown_minutes"].(int)),
			TargetNumberInstances:   uint64(poolConfig["target_number_instances"].(int)),
		}
	}
	return nil
}

func expandTemplate(template []interface{}) *goApiAbrha.VmAutoscaleResourceTemplate {
	if len(template) > 0 {
		poolTemplate := template[0].(map[string]interface{})

		var tags []string
		if v, ok := poolTemplate["tags"]; ok {
			for _, tag := range v.(*schema.Set).List() {
				tags = append(tags, tag.(string))
			}
		}

		var sshKeys []string
		if v, ok := poolTemplate["ssh_keys"]; ok {
			for _, key := range v.(*schema.Set).List() {
				sshKeys = append(sshKeys, key.(string))
			}
		}

		return &goApiAbrha.VmAutoscaleResourceTemplate{
			Size:        poolTemplate["size"].(string),
			Region:      poolTemplate["region"].(string),
			Image:       poolTemplate["image"].(string),
			Tags:        tags,
			SSHKeys:     sshKeys,
			VpcUUID:     poolTemplate["vpc_uuid"].(string),
			WithVmAgent: poolTemplate["with_vm_agent"].(bool),
			ProjectID:   poolTemplate["project_id"].(string),
			IPV6:        poolTemplate["ipv6"].(bool),
			UserData:    poolTemplate["user_data"].(string),
		}
	}
	return nil
}

func resourceAbrhaVmAutoscaleCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	pool, _, err := client.VmAutoscale.Create(context.Background(), &goApiAbrha.VmAutoscalePoolRequest{
		Name:       d.Get("name").(string),
		Config:     expandConfig(d.Get("config").([]interface{})),
		VmTemplate: expandTemplate(d.Get("vm_template").([]interface{})),
	})
	if err != nil {
		return diag.Errorf("Error creating Vm autoscale pool: %v", err)
	}
	d.SetId(pool.ID)

	// Setup to poll for autoscale pool scaling up to the desired count
	stateConf := &retry.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{"provisioning"},
		Target:     []string{"active"},
		Refresh:    vmAutoscaleRefreshFunc(client, d.Id()),
		MinTimeout: 15 * time.Second,
		Timeout:    15 * time.Minute,
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("Error waiting for Vm autoscale pool (%s) to become active: %v", pool.Name, err)
	}

	return resourceAbrhaVmAutoscaleRead(ctx, d, meta)
}

func resourceAbrhaVmAutoscaleRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	pool, _, err := client.VmAutoscale.Get(context.Background(), d.Id())
	if err != nil {
		if strings.Contains(err.Error(), fmt.Sprintf("autoscale group with id %s not found", d.Id())) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error retrieving Vm autoscale pool: %v", err)
	}

	d.Set("name", pool.Name)
	d.Set("config", flattenConfig(pool.Config))
	d.Set("current_utilization", flattenUtilization(pool.CurrentUtilization))
	d.Set("status", pool.Status)
	d.Set("created_at", pool.CreatedAt.UTC().String())
	d.Set("updated_at", pool.UpdatedAt.UTC().String())

	// Persist existing image specification (id/slug) if it exists
	if t, ok := d.GetOk("vm_template"); ok {
		tList := t.([]interface{})
		if len(tList) > 0 {
			tMap := tList[0].(map[string]interface{})
			if tMap["image"] != "" {
				pool.VmTemplate.Image = tMap["image"].(string)
			}
		}
	}
	d.Set("vm_template", flattenTemplate(pool.VmTemplate))

	return nil
}

func resourceAbrhaVmAutoscaleUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	_, _, err := client.VmAutoscale.Update(context.Background(), d.Id(), &goApiAbrha.VmAutoscalePoolRequest{
		Name:       d.Get("name").(string),
		Config:     expandConfig(d.Get("config").([]interface{})),
		VmTemplate: expandTemplate(d.Get("vm_template").([]interface{})),
	})
	if err != nil {
		return diag.Errorf("Error updating Vm autoscale pool: %v", err)
	}

	return resourceAbrhaVmAutoscaleRead(ctx, d, meta)
}

func resourceAbrhaVmAutoscaleDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	_, err := client.VmAutoscale.DeleteDangerous(context.Background(), d.Id())
	if err != nil {
		return diag.Errorf("Error updating Vm autoscale pool: %v", err)
	}

	// Setup to poll for autoscale pool deletion
	stateConf := &retry.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{http.StatusText(http.StatusOK)},
		Target:     []string{http.StatusText(http.StatusNotFound)},
		Refresh:    vmAutoscaleRefreshFunc(client, d.Id()),
		MinTimeout: 5 * time.Second,
		Timeout:    1 * time.Minute,
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("Error waiting for Vm autoscale pool (%s) to become be deleted: %v", d.Get("name"), err)
	}

	d.SetId("")
	return nil
}

func vmAutoscaleRefreshFunc(client *goApiAbrha.Client, poolID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		// Check autoscale pool status
		pool, _, err := client.VmAutoscale.Get(context.Background(), poolID)
		if err != nil {
			if strings.Contains(err.Error(), fmt.Sprintf("autoscale group with id %s not found", poolID)) {
				return pool, http.StatusText(http.StatusNotFound), nil
			}
			return nil, "", fmt.Errorf("Error retrieving Vm autoscale pool: %v", err)
		}
		if pool.Status != "active" {
			return pool, pool.Status, nil
		}
		members := make([]*goApiAbrha.VmAutoscaleResource, 0)
		opts := &goApiAbrha.ListOptions{
			Page:    1,
			PerPage: 100,
		}
		// Paginate through autoscale pool members and validate status
		for {
			m, resp, err := client.VmAutoscale.ListMembers(context.Background(), poolID, opts)
			if err != nil {
				return nil, "", fmt.Errorf("Error listing Vm autoscale pool members: %v", err)
			}
			members = append(members, m...)
			if resp.Links.IsLastPage() {
				break
			}
			page, err := resp.Links.CurrentPage()
			if err != nil {
				break
			}
			opts.Page = page + 1
		}
		// Scan through the list to find a non-active provision state
		for i := range members {
			if members[i].Status != "active" {
				return members, members[i].Status, nil
			}
		}
		return members, "active", nil
	}
}
