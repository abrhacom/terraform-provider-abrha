package uptime

import (
	"context"
	"errors"
	"log"
	"strings"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/abrhacom/terraform-provider-abrha/abrha/monitoring"
	"github.com/abrhacom/terraform-provider-abrha/abrha/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceAbrhaUptimeAlert() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbrhaUptimeAlertCreate,
		ReadContext:   resourceAbrhaUptimeAlertRead,
		UpdateContext: resourceAbrhaUptimeAlertUpdate,
		DeleteContext: resourceAbrhaUptimeAlertDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceAbrhaUptimeAlertImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "A human-friendly display name for the alert.",
				Required:    true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"check_id": {
				Type:        schema.TypeString,
				Description: "A unique identifier for a check.",
				Required:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "The type of health check to perform. Enum: 'latency' 'down' 'down_global' 'ssl_expiry'",
				ValidateFunc: validation.StringInSlice([]string{
					"latency",
					"down",
					"down_global",
					"ssl_expiry",
				}, false),
				Required: true,
			},
			"threshold": {
				Type:        schema.TypeInt,
				Description: "The threshold at which the alert will enter a trigger state. The specific threshold is dependent on the alert type.",
				Optional:    true,
			},
			"comparison": {
				Type:        schema.TypeString,
				Description: "The comparison operator used against the alert's threshold. Enum: 'greater_than' 'less_than",
				ValidateFunc: validation.StringInSlice([]string{
					"greater_than",
					"less_than",
				}, false),
				Optional: true,
			},
			"period": {
				Type:        schema.TypeString,
				Description: "Period of time the threshold must be exceeded to trigger the alert. Enum '2m' '3m' '5m' '10m' '15m' '30m' '1h'",
				ValidateFunc: validation.StringInSlice([]string{
					"2m",
					"3m",
					"5m",
					"10m",
					"15m",
					"30m",
					"1h",
				}, false),
				Optional: true,
			},
			"notifications": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "The notification settings for a trigger alert.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"slack": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"channel": {
										Type:             schema.TypeString,
										Required:         true,
										DiffSuppressFunc: util.CaseSensitive,
										Description:      "The Slack channel to send alerts to",
										ValidateFunc:     validation.StringIsNotEmpty,
									},
									"url": {
										Type:             schema.TypeString,
										Required:         true,
										DiffSuppressFunc: util.CaseSensitive,
										Description:      "The webhook URL for Slack",
										ValidateFunc:     validation.StringIsNotEmpty,
									},
								},
							},
						},
						"email": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of email addresses to sent notifications to",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func resourceAbrhaUptimeAlertCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	checkID := d.Get("check_id").(string)

	opts := &goApiAbrha.CreateUptimeAlertRequest{
		Name:          d.Get("name").(string),
		Type:          d.Get("type").(string),
		Notifications: expandNotifications(d.Get("notifications").([]interface{})),
		Comparison:    goApiAbrha.UptimeAlertComp(d.Get("comparison").(string)),
		Threshold:     d.Get("threshold").(int),
		Period:        d.Get("period").(string),
	}

	log.Printf("[DEBUG] Uptime alert create configuration: %#v", opts)
	alert, _, err := client.UptimeChecks.CreateAlert(ctx, checkID, opts)
	if err != nil {
		return diag.Errorf("Error creating Uptime Alert: %s", err)
	}

	d.SetId(alert.ID)
	log.Printf("[INFO] Uptime Alert name: %s", alert.Name)

	return resourceAbrhaUptimeAlertRead(ctx, d, meta)
}

func expandNotifications(config []interface{}) *goApiAbrha.Notifications {
	alertConfig := config[0].(map[string]interface{})
	alerts := &goApiAbrha.Notifications{
		Slack: monitoring.ExpandSlack(alertConfig["slack"].([]interface{})),
		Email: monitoring.ExpandEmail(alertConfig["email"].([]interface{})),
	}
	return alerts
}

func resourceAbrhaUptimeAlertUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	checkID := d.Get("check_id").(string)

	opts := &goApiAbrha.UpdateUptimeAlertRequest{
		Name:          d.Get("name").(string),
		Type:          d.Get("type").(string),
		Notifications: expandNotifications(d.Get("notifications").([]interface{})),
	}

	if v, ok := d.GetOk("comparison"); ok {
		opts.Comparison = goApiAbrha.UptimeAlertComp(v.(string))
	}
	if v, ok := d.GetOk("threshold"); ok {
		opts.Threshold = v.(int)
	}
	if v, ok := d.GetOk("period"); ok {
		opts.Period = v.(string)
	}

	log.Printf("[DEBUG] Uptime alert update configuration: %#v", opts)

	alert, _, err := client.UptimeChecks.UpdateAlert(ctx, checkID, d.Id(), opts)
	if err != nil {
		return diag.Errorf("Error updating Alert: %s", err)
	}

	log.Printf("[INFO] Uptime Alert name: %s", alert.Name)

	return resourceAbrhaUptimeAlertRead(ctx, d, meta)
}

func resourceAbrhaUptimeAlertDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	checkID := d.Get("check_id").(string)

	log.Printf("[INFO] Deleting uptime alert: %s", d.Id())

	// Delete the uptime alert
	_, err := client.UptimeChecks.DeleteAlert(ctx, checkID, d.Id())

	if err != nil {
		return diag.Errorf("Error deleting uptime alerts: %s", err)
	}

	return nil
}

func resourceAbrhaUptimeAlertRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	checkID := d.Get("check_id").(string)

	alert, resp, err := client.UptimeChecks.GetAlert(ctx, checkID, d.Id())
	if err != nil {
		// If the check is somehow already destroyed, mark as
		// successfully gone
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving check: %s", err)
	}

	d.SetId(alert.ID)
	d.Set("name", alert.Name)
	d.Set("type", alert.Type)
	d.Set("threshold", alert.Threshold)
	d.Set("notifications", flattenNotifications(alert.Notifications))
	d.Set("comparison", alert.Comparison)
	d.Set("period", alert.Period)

	return nil
}

func flattenNotifications(alerts *goApiAbrha.Notifications) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"email": monitoring.FlattenEmail(alerts.Email),
			"slack": monitoring.FlattenSlack(alerts.Slack),
		},
	}
}

func resourceAbrhaUptimeAlertImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if strings.Contains(d.Id(), ",") {
		s := strings.Split(d.Id(), ",")

		d.SetId(s[1])
		d.Set("check_id", s[0])
	} else {
		return nil, errors.New("must use the IDs of the check and alert joined with a comma (e.g. `check_id,alert_id`)")
	}

	return []*schema.ResourceData{d}, nil
}
