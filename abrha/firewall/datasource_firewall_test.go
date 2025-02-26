package firewall_test

import (
	"testing"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAbrhaFirewall_Basic(t *testing.T) {
	fwDataConfig := `
data "abrha_firewall" "foobar" {
  firewall_id = abrha_firewall.foobar.id
}`

	var firewall goApiAbrha.Firewall
	fwName := acceptance.RandomTestName()

	fwCreateConfig := testAccAbrhaFirewallConfig_OnlyInbound(fwName)
	updatedFWCreateConfig := testAccAbrhaFirewallConfig_OnlyMultipleInbound(fwName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fwCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaFirewallExists("abrha_firewall.foobar", &firewall),
				),
			},
			{
				Config: fwCreateConfig + fwDataConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.abrha_firewall.foobar", "name", fwName),
					resource.TestCheckResourceAttrPair("abrha_firewall.foobar", "id",
						"data.abrha_firewall.foobar", "firewall_id"),
					resource.TestCheckResourceAttrPair("abrha_firewall.foobar", "vm_ids",
						"data.abrha_firewall.foobar", "vm_ids"),
					resource.TestCheckResourceAttrPair("abrha_firewall.foobar", "inbound_rule",
						"data.abrha_firewall.foobar", "inbound_rule"),
					resource.TestCheckResourceAttrPair("abrha_firewall.foobar", "outbound_rule",
						"data.abrha_firewall.foobar", "outbound_rule"),
					resource.TestCheckResourceAttrPair("abrha_firewall.foobar", "status",
						"data.abrha_firewall.foobar", "status"),
					resource.TestCheckResourceAttrPair("abrha_firewall.foobar", "created_at",
						"data.abrha_firewall.foobar", "created_at"),
					resource.TestCheckResourceAttrPair("abrha_firewall.foobar", "pending_changes",
						"data.abrha_firewall.foobar", "pending_changes"),
					resource.TestCheckResourceAttrPair("abrha_firewall.foobar", "tags",
						"data.abrha_firewall.foobar", "tags"),
					resource.TestCheckResourceAttr("data.abrha_firewall.foobar", "inbound_rule.0.protocol", "tcp"),
					resource.TestCheckResourceAttr("data.abrha_firewall.foobar", "inbound_rule.0.port_range", "22"),
					resource.TestCheckResourceAttr("data.abrha_firewall.foobar", "inbound_rule.0.source_addresses.0", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("data.abrha_firewall.foobar", "inbound_rule.0.source_addresses.1", "::/0"),
				),
			},
			{
				Config: updatedFWCreateConfig,
			},
			{
				Config: updatedFWCreateConfig + fwDataConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.abrha_firewall.foobar", "inbound_rule.0.protocol", "tcp"),
					resource.TestCheckResourceAttr("data.abrha_firewall.foobar", "inbound_rule.0.port_range", "22"),
					resource.TestCheckResourceAttr("data.abrha_firewall.foobar", "inbound_rule.0.source_addresses.0", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("data.abrha_firewall.foobar", "inbound_rule.0.source_addresses.1", "::/0"),
					resource.TestCheckResourceAttr("data.abrha_firewall.foobar", "inbound_rule.1.protocol", "tcp"),
					resource.TestCheckResourceAttr("data.abrha_firewall.foobar", "inbound_rule.1.port_range", "80"),
					resource.TestCheckResourceAttr("data.abrha_firewall.foobar", "inbound_rule.1.source_addresses.0", "1.2.3.0/24"),
					resource.TestCheckResourceAttr("data.abrha_firewall.foobar", "inbound_rule.1.source_addresses.1", "2002::/16"),
				),
			},
		},
	})
}
