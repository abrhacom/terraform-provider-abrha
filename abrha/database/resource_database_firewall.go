package database

import (
	"context"
	"log"
	"time"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceAbrhaDatabaseFirewall() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbrhaDatabaseFirewallCreate,
		ReadContext:   resourceAbrhaDatabaseFirewallRead,
		UpdateContext: resourceAbrhaDatabaseFirewallUpdate,
		DeleteContext: resourceAbrhaDatabaseFirewallDelete,
		Importer: &schema.ResourceImporter{
			State: resourceAbrhaDatabaseFirewallImport,
		},

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"rule": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"ip_addr",
								"vm",
								"k8s",
								"tag",
								"app",
							}, false),
						},

						"value": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.NoZeroValues,
						},

						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAbrhaDatabaseFirewallCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()
	clusterID := d.Get("cluster_id").(string)

	rules := buildDatabaseFirewallRequest(d.Get("rule").(*schema.Set).List())

	_, err := client.Databases.UpdateFirewallRules(context.TODO(), clusterID, &rules)
	if err != nil {
		return diag.Errorf("Error creating DatabaseFirewall: %s", err)
	}

	d.SetId(id.PrefixedUniqueId(clusterID + "-"))

	return resourceAbrhaDatabaseFirewallRead(ctx, d, meta)
}

func resourceAbrhaDatabaseFirewallRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()
	clusterID := d.Get("cluster_id").(string)

	rules, resp, err := client.Databases.GetFirewallRules(context.TODO(), clusterID)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error retrieving DatabaseFirewall: %s", err)
	}

	err = d.Set("rule", flattenDatabaseFirewallRules(rules))
	if err != nil {
		return diag.FromErr(err)
	} else {
		return nil
	}
}

func resourceAbrhaDatabaseFirewallUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()
	clusterID := d.Get("cluster_id").(string)

	rules := buildDatabaseFirewallRequest(d.Get("rule").(*schema.Set).List())

	_, err := client.Databases.UpdateFirewallRules(context.TODO(), clusterID, &rules)
	if err != nil {
		return diag.Errorf("Error updating DatabaseFirewall: %s", err)
	}

	return resourceAbrhaDatabaseFirewallRead(ctx, d, meta)
}

func resourceAbrhaDatabaseFirewallDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()
	clusterID := d.Get("cluster_id").(string)

	log.Printf("[INFO] Deleting DatabaseFirewall: %s", d.Id())
	req := goApiAbrha.DatabaseUpdateFirewallRulesRequest{
		Rules: []*goApiAbrha.DatabaseFirewallRule{},
	}

	_, err := client.Databases.UpdateFirewallRules(context.TODO(), clusterID, &req)
	if err != nil {
		return diag.Errorf("Error deleting DatabaseFirewall: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceAbrhaDatabaseFirewallImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	clusterID := d.Id()
	d.Set("cluster_id", clusterID)
	d.SetId(id.PrefixedUniqueId(clusterID + "-"))

	return []*schema.ResourceData{d}, nil
}

func buildDatabaseFirewallRequest(rules []interface{}) goApiAbrha.DatabaseUpdateFirewallRulesRequest {
	expandedRules := make([]*goApiAbrha.DatabaseFirewallRule, 0, len(rules))
	for _, rawRule := range rules {
		rule := rawRule.(map[string]interface{})

		r := goApiAbrha.DatabaseFirewallRule{
			Type:  rule["type"].(string),
			Value: rule["value"].(string),
		}

		if rule["uuid"].(string) != "" {
			r.UUID = rule["uuid"].(string)
		}

		expandedRules = append(expandedRules, &r)
	}

	return goApiAbrha.DatabaseUpdateFirewallRulesRequest{
		Rules: expandedRules,
	}
}

func flattenDatabaseFirewallRules(rules []goApiAbrha.DatabaseFirewallRule) []interface{} {
	if rules == nil {
		return nil
	}

	flattenedRules := make([]interface{}, len(rules))
	for i, rule := range rules {
		rawRule := map[string]interface{}{
			"uuid":       rule.UUID,
			"type":       rule.Type,
			"value":      rule.Value,
			"created_at": rule.CreatedAt.Format(time.RFC3339),
		}

		flattenedRules[i] = rawRule
	}

	return flattenedRules
}
