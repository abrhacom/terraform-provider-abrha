package firewall

import (
	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/tag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func firewallSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.NoZeroValues,
		},

		"vm_ids": {
			Type:     schema.TypeSet,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
		},

		"inbound_rule": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     firewallRuleSchema("source"),
		},

		"outbound_rule": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     firewallRuleSchema("destination"),
		},

		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"pending_changes": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"vm_id": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"removing": {
						Type:     schema.TypeBool,
						Optional: true,
					},
					"status": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},

		"tags": tag.TagsSchema(),
	}
}

func firewallRuleSchema(prefix string) *schema.Resource {
	if prefix != "" && prefix[len(prefix)-1:] != "_" {
		prefix += "_"
	}

	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"tcp",
					"udp",
					"icmp",
				}, false),
			},
			"port_range": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			prefix + "addresses": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.NoZeroValues,
				},
				Optional: true,
			},
			prefix + "vm_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			prefix + "load_balancer_uids": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.NoZeroValues,
				},
				Optional: true,
			},
			prefix + "kubernetes_ids": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.NoZeroValues,
				},
				Optional: true,
			},
			prefix + "tags": tag.TagsSchema(),
		},
	}
}

func expandFirewallVmIds(vms []interface{}) []string {
	expandedVms := make([]string, len(vms))
	for i, v := range vms {
		expandedVms[i] = v.(string)
	}

	return expandedVms
}

func expandFirewallRuleStringSet(strings []interface{}) []string {
	expandedStrings := make([]string, len(strings))
	for i, v := range strings {
		expandedStrings[i] = v.(string)
	}

	return expandedStrings
}

func expandFirewallInboundRules(rules []interface{}) []goApiAbrha.InboundRule {
	expandedRules := make([]goApiAbrha.InboundRule, 0, len(rules))
	for _, rawRule := range rules {
		var src goApiAbrha.Sources

		rule := rawRule.(map[string]interface{})

		src.VmIDs = expandFirewallVmIds(rule["source_vm_ids"].(*schema.Set).List())

		src.Addresses = expandFirewallRuleStringSet(rule["source_addresses"].(*schema.Set).List())

		src.LoadBalancerUIDs = expandFirewallRuleStringSet(rule["source_load_balancer_uids"].(*schema.Set).List())

		src.KubernetesIDs = expandFirewallRuleStringSet(rule["source_kubernetes_ids"].(*schema.Set).List())

		src.Tags = tag.ExpandTags(rule["source_tags"].(*schema.Set).List())

		r := goApiAbrha.InboundRule{
			Protocol:  rule["protocol"].(string),
			PortRange: rule["port_range"].(string),
			Sources:   &src,
		}

		expandedRules = append(expandedRules, r)
	}
	return expandedRules
}

func expandFirewallOutboundRules(rules []interface{}) []goApiAbrha.OutboundRule {
	expandedRules := make([]goApiAbrha.OutboundRule, 0, len(rules))
	for _, rawRule := range rules {
		var dest goApiAbrha.Destinations

		rule := rawRule.(map[string]interface{})

		dest.VmIDs = expandFirewallVmIds(rule["destination_vm_ids"].(*schema.Set).List())

		dest.Addresses = expandFirewallRuleStringSet(rule["destination_addresses"].(*schema.Set).List())

		dest.LoadBalancerUIDs = expandFirewallRuleStringSet(rule["destination_load_balancer_uids"].(*schema.Set).List())

		dest.KubernetesIDs = expandFirewallRuleStringSet(rule["destination_kubernetes_ids"].(*schema.Set).List())

		dest.Tags = tag.ExpandTags(rule["destination_tags"].(*schema.Set).List())

		r := goApiAbrha.OutboundRule{
			Protocol:     rule["protocol"].(string),
			PortRange:    rule["port_range"].(string),
			Destinations: &dest,
		}

		expandedRules = append(expandedRules, r)
	}
	return expandedRules
}

func firewallPendingChanges(d *schema.ResourceData, firewall *goApiAbrha.Firewall) []interface{} {
	remote := make([]interface{}, 0, len(firewall.PendingChanges))
	for _, change := range firewall.PendingChanges {
		rawChange := map[string]interface{}{
			"vm_id":    change.VmID,
			"removing": change.Removing,
			"status":   change.Status,
		}
		remote = append(remote, rawChange)
	}
	return remote
}

func flattenFirewallVmIds(vms []string) *schema.Set {
	if vms == nil {
		return nil
	}

	flattenedVms := schema.NewSet(schema.HashString, []interface{}{})
	for _, v := range vms {
		flattenedVms.Add(v)
	}

	return flattenedVms
}

func flattenFirewallRuleStringSet(strings []string) *schema.Set {
	flattenedStrings := schema.NewSet(schema.HashString, []interface{}{})
	for _, v := range strings {
		flattenedStrings.Add(v)
	}

	return flattenedStrings
}

func flattenFirewallInboundRules(rules []goApiAbrha.InboundRule) []interface{} {
	if rules == nil {
		return nil
	}

	flattenedRules := make([]interface{}, len(rules))
	for i, rule := range rules {
		sources := rule.Sources
		protocol := rule.Protocol
		portRange := rule.PortRange

		rawRule := map[string]interface{}{
			"protocol": protocol,
		}

		// The API returns 0 when the port range was specified as all.
		// If protocol is `icmp` the API returns 0 for when port was
		// not specified.
		if portRange == "0" {
			if protocol != "icmp" {
				rawRule["port_range"] = "all"
			}
		} else {
			rawRule["port_range"] = portRange
		}

		if sources.Tags != nil {
			rawRule["source_tags"] = tag.FlattenTags(sources.Tags)
		}

		if sources.VmIDs != nil {
			rawRule["source_vm_ids"] = flattenFirewallVmIds(sources.VmIDs)
		}

		if sources.Addresses != nil {
			rawRule["source_addresses"] = flattenFirewallRuleStringSet(sources.Addresses)
		}

		if sources.LoadBalancerUIDs != nil {
			rawRule["source_load_balancer_uids"] = flattenFirewallRuleStringSet(sources.LoadBalancerUIDs)
		}

		if sources.KubernetesIDs != nil {
			rawRule["source_kubernetes_ids"] = flattenFirewallRuleStringSet(sources.KubernetesIDs)
		}

		flattenedRules[i] = rawRule
	}

	return flattenedRules
}

func flattenFirewallOutboundRules(rules []goApiAbrha.OutboundRule) []interface{} {
	if rules == nil {
		return nil
	}

	flattenedRules := make([]interface{}, len(rules))
	for i, rule := range rules {
		destinations := rule.Destinations
		protocol := rule.Protocol
		portRange := rule.PortRange

		rawRule := map[string]interface{}{
			"protocol": protocol,
		}

		// The API returns 0 when the port range was specified as all.
		// If protocol is `icmp` the API returns 0 for when port was
		// not specified.
		if portRange == "0" {
			if protocol != "icmp" {
				rawRule["port_range"] = "all"
			}
		} else {
			rawRule["port_range"] = portRange
		}

		if destinations.Tags != nil {
			rawRule["destination_tags"] = tag.FlattenTags(destinations.Tags)
		}

		if destinations.VmIDs != nil {
			rawRule["destination_vm_ids"] = flattenFirewallVmIds(destinations.VmIDs)
		}

		if destinations.Addresses != nil {
			rawRule["destination_addresses"] = flattenFirewallRuleStringSet(destinations.Addresses)
		}

		if destinations.LoadBalancerUIDs != nil {
			rawRule["destination_load_balancer_uids"] = flattenFirewallRuleStringSet(destinations.LoadBalancerUIDs)
		}

		if destinations.KubernetesIDs != nil {
			rawRule["destination_kubernetes_ids"] = flattenFirewallRuleStringSet(destinations.KubernetesIDs)
		}

		flattenedRules[i] = rawRule
	}

	return flattenedRules
}
