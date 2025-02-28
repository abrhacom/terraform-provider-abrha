package monitoring

import (
	"context"
	"log"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/abrhacom/terraform-provider-abrha/abrha/tag"
	"github.com/abrhacom/terraform-provider-abrha/abrha/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceAbrhaMonitorAlert() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbrhaMonitorAlertCreate,
		ReadContext:   resourceAbrhaMonitorAlertRead,
		UpdateContext: resourceAbrhaMonitorAlertUpdate,
		DeleteContext: resourceAbrhaMonitorAlertDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					goApiAbrha.VmCPUUtilizationPercent,
					goApiAbrha.VmMemoryUtilizationPercent,
					goApiAbrha.VmDiskUtilizationPercent,
					goApiAbrha.VmPublicOutboundBandwidthRate,
					goApiAbrha.VmPublicInboundBandwidthRate,
					goApiAbrha.VmPrivateOutboundBandwidthRate,
					goApiAbrha.VmPrivateInboundBandwidthRate,
					goApiAbrha.VmDiskReadRate,
					goApiAbrha.VmDiskWriteRate,
					goApiAbrha.VmOneMinuteLoadAverage,
					goApiAbrha.VmFiveMinuteLoadAverage,
					goApiAbrha.VmFifteenMinuteLoadAverage,
					goApiAbrha.LoadBalancerCPUUtilizationPercent,
					goApiAbrha.LoadBalancerConnectionUtilizationPercent,
					goApiAbrha.LoadBalancerVmHealth,
					goApiAbrha.LoadBalancerTLSUtilizationPercent,
					goApiAbrha.LoadBalancerIncreaseInHTTPErrorRatePercentage4xx,
					goApiAbrha.LoadBalancerIncreaseInHTTPErrorRatePercentage5xx,
					goApiAbrha.LoadBalancerIncreaseInHTTPErrorRateCount4xx,
					goApiAbrha.LoadBalancerIncreaseInHTTPErrorRateCount5xx,
					goApiAbrha.LoadBalancerHighHttpResponseTime,
					goApiAbrha.LoadBalancerHighHttpResponseTime50P,
					goApiAbrha.LoadBalancerHighHttpResponseTime95P,
					goApiAbrha.LoadBalancerHighHttpResponseTime99P,
					goApiAbrha.DbaasFifteenMinuteLoadAverage,
					goApiAbrha.DbaasMemoryUtilizationPercent,
					goApiAbrha.DbaasDiskUtilizationPercent,
					goApiAbrha.DbaasCPUUtilizationPercent,
				}, false),
			},

			"compare": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(goApiAbrha.GreaterThan),
					string(goApiAbrha.LessThan),
				}, false),
				Description: "The comparison operator to use for value",
			},

			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Description of the alert policy",
			},

			"enabled": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},

			"value": {
				Type:         schema.TypeFloat,
				Required:     true,
				ValidateFunc: validation.FloatAtLeast(0),
			},

			"tags": tag.TagsSchema(),

			"alerts": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "List with details how to notify about the alert. Support for Slack or email.",
				MaxItems:    1,
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

			"entities": {
				Type:        schema.TypeSet,
				Optional:    true,
				MinItems:    1,
				Description: "The vms to apply the alert policy to",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"window": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"5m", "10m", "30m", "1h",
				}, false),
			},
		},
	}
}

func resourceAbrhaMonitorAlertCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	alertCreateRequest := &goApiAbrha.AlertPolicyCreateRequest{
		Type:        d.Get("type").(string),
		Enabled:     goApiAbrha.PtrTo(d.Get("enabled").(bool)),
		Description: d.Get("description").(string),
		Tags:        tag.ExpandTags(d.Get("tags").(*schema.Set).List()),
		Compare:     goApiAbrha.AlertPolicyComp(d.Get("compare").(string)),
		Window:      d.Get("window").(string),
		Value:       float32(d.Get("value").(float64)),
		Entities:    expandEntities(d.Get("entities").(*schema.Set).List()),
		Alerts:      expandAlerts(d.Get("alerts").([]interface{})),
	}

	log.Printf("[DEBUG] Alert Policy create configuration: %#v", alertCreateRequest)
	alertPolicy, _, err := client.Monitoring.CreateAlertPolicy(context.Background(), alertCreateRequest)
	if err != nil {
		return diag.Errorf("Error creating Alert Policy: %s", err)
	}

	d.SetId(alertPolicy.UUID)
	log.Printf("[INFO] Alert Policy created, ID: %s", d.Id())

	return resourceAbrhaMonitorAlertRead(ctx, d, meta)
}

func expandAlerts(config []interface{}) goApiAbrha.Alerts {
	alertConfig := config[0].(map[string]interface{})
	alerts := goApiAbrha.Alerts{
		Slack: ExpandSlack(alertConfig["slack"].([]interface{})),
		Email: ExpandEmail(alertConfig["email"].([]interface{})),
	}
	return alerts
}

func flattenAlerts(alerts goApiAbrha.Alerts) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"email": FlattenEmail(alerts.Email),
			"slack": FlattenSlack(alerts.Slack),
		},
	}
}

func ExpandSlack(slackChannels []interface{}) []goApiAbrha.SlackDetails {
	if len(slackChannels) == 0 {
		return nil
	}

	expandedSlackChannels := make([]goApiAbrha.SlackDetails, 0, len(slackChannels))

	for _, slackChannel := range slackChannels {
		slack := slackChannel.(map[string]interface{})
		n := goApiAbrha.SlackDetails{
			Channel: slack["channel"].(string),
			URL:     slack["url"].(string),
		}

		expandedSlackChannels = append(expandedSlackChannels, n)
	}

	return expandedSlackChannels
}

func FlattenSlack(slackChannels []goApiAbrha.SlackDetails) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(slackChannels))

	for _, slackChannel := range slackChannels {
		item := make(map[string]interface{})
		item["url"] = slackChannel.URL
		item["channel"] = slackChannel.Channel
		result = append(result, item)
	}

	return result
}

func ExpandEmail(config []interface{}) []string {
	if len(config) == 0 {
		return nil
	}
	emailList := make([]string, len(config))

	for i, v := range config {
		emailList[i] = v.(string)
	}

	return emailList
}

func FlattenEmail(emails []string) []string {
	if len(emails) == 0 {
		return nil
	}

	flattenedEmails := make([]string, 0)
	for _, v := range emails {
		if v != "" {
			flattenedEmails = append(flattenedEmails, v)
		}
	}

	return flattenedEmails
}

func expandEntities(config []interface{}) []string {
	alertEntities := make([]string, len(config))

	for i, v := range config {
		alertEntities[i] = v.(string)
	}

	return alertEntities
}

func resourceAbrhaMonitorAlertUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	updateRequest := &goApiAbrha.AlertPolicyUpdateRequest{
		Type:        d.Get("type").(string),
		Enabled:     goApiAbrha.PtrTo(d.Get("enabled").(bool)),
		Description: d.Get("description").(string),
		Tags:        tag.ExpandTags(d.Get("tags").(*schema.Set).List()),
		Compare:     goApiAbrha.AlertPolicyComp(d.Get("compare").(string)),
		Window:      d.Get("window").(string),
		Value:       float32(d.Get("value").(float64)),
		Entities:    expandEntities(d.Get("entities").(*schema.Set).List()),
		Alerts:      expandAlerts(d.Get("alerts").([]interface{})),
	}

	_, _, err := client.Monitoring.UpdateAlertPolicy(ctx, d.Id(), updateRequest)
	if err != nil {
		return diag.Errorf("Error updating monitoring alert: %s", err)
	}

	return resourceAbrhaMonitorAlertRead(ctx, d, meta)
}

func resourceAbrhaMonitorAlertRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	alert, resp, err := client.Monitoring.GetAlertPolicy(ctx, d.Id())

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[DEBUG] Alert (%s) was not found - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error reading Alert: %s", err)
	}

	d.SetId(alert.UUID)
	d.Set("description", alert.Description)
	d.Set("enabled", alert.Enabled)
	d.Set("compare", alert.Compare)
	d.Set("alerts", flattenAlerts(alert.Alerts))
	d.Set("value", alert.Value)
	d.Set("window", alert.Window)
	d.Set("entities", alert.Entities)
	d.Set("tags", tag.FlattenTags(alert.Tags))
	d.Set("type", alert.Type)

	return nil
}

func resourceAbrhaMonitorAlertDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	log.Printf("[INFO] Deleting the monitor alert")
	_, err := client.Monitoring.DeleteAlertPolicy(context.Background(), d.Id())
	if err != nil {
		return diag.Errorf("Error deleting monitor alert: %s", err)
	}
	d.SetId("")
	return nil
}
