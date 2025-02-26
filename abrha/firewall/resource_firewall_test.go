package firewall_test

import (
	"context"
	"fmt"
	"testing"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAbrhaFirewall_AllowOnlyInbound(t *testing.T) {
	rName := acceptance.RandomTestName()
	var firewall goApiAbrha.Firewall

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAbrhaFirewallConfig_OnlyInbound(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaFirewallExists("abrha_firewall.foobar", &firewall),
					resource.TestCheckResourceAttr("abrha_firewall.foobar", "inbound_rule.#", "1"),
				),
			},
		},
	})
}

func TestAccAbrhaFirewall_AllowMultipleInbound(t *testing.T) {
	rName := acceptance.RandomTestName()
	var firewall goApiAbrha.Firewall

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAbrhaFirewallConfig_OnlyMultipleInbound(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaFirewallExists("abrha_firewall.foobar", &firewall),
					resource.TestCheckResourceAttr("abrha_firewall.foobar", "inbound_rule.#", "2"),
				),
			},
		},
	})
}

func TestAccAbrhaFirewall_AllowOnlyOutbound(t *testing.T) {
	rName := acceptance.RandomTestName()
	var firewall goApiAbrha.Firewall

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAbrhaFirewallConfig_OnlyOutbound(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaFirewallExists("abrha_firewall.foobar", &firewall),
					resource.TestCheckResourceAttr("abrha_firewall.foobar", "outbound_rule.#", "1"),
				),
			},
		},
	})
}

func TestAccAbrhaFirewall_AllowMultipleOutbound(t *testing.T) {
	rName := acceptance.RandomTestName()
	var firewall goApiAbrha.Firewall

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAbrhaFirewallConfig_OnlyMultipleOutbound(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaFirewallExists("abrha_firewall.foobar", &firewall),
					resource.TestCheckResourceAttr("abrha_firewall.foobar", "outbound_rule.#", "2"),
				),
			},
		},
	})
}

func TestAccAbrhaFirewall_MultipleInboundAndOutbound(t *testing.T) {
	rName := acceptance.RandomTestName()
	tagName := acceptance.RandomTestName("tag")
	var firewall goApiAbrha.Firewall

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAbrhaFirewallConfig_MultipleInboundAndOutbound(tagName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaFirewallExists("abrha_firewall.foobar", &firewall),
					resource.TestCheckResourceAttr("abrha_firewall.foobar", "inbound_rule.#", "2"),
					resource.TestCheckResourceAttr("abrha_firewall.foobar", "outbound_rule.#", "2"),
				),
			},
		},
	})
}

func TestAccAbrhaFirewall_fullPortRange(t *testing.T) {
	rName := acceptance.RandomTestName()
	var firewall goApiAbrha.Firewall

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAbrhaFirewallConfig_fullPortRange(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaFirewallExists("abrha_firewall.foobar", &firewall),
					resource.TestCheckResourceAttr("abrha_firewall.foobar", "inbound_rule.#", "1"),
					resource.TestCheckResourceAttr("abrha_firewall.foobar", "outbound_rule.#", "1"),
				),
			},
		},
	})
}

func TestAccAbrhaFirewall_icmp(t *testing.T) {
	rName := acceptance.RandomTestName()
	var firewall goApiAbrha.Firewall

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAbrhaFirewallConfig_icmp(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaFirewallExists("abrha_firewall.foobar", &firewall),
					resource.TestCheckResourceAttr("abrha_firewall.foobar", "inbound_rule.#", "1"),
					resource.TestCheckResourceAttr("abrha_firewall.foobar", "outbound_rule.#", "1"),
				),
			},
		},
	})
}

func TestAccAbrhaFirewall_ImportMultipleRules(t *testing.T) {
	resourceName := "abrha_firewall.foobar"
	rName := acceptance.RandomTestName()
	tagName := acceptance.RandomTestName("tag")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAbrhaFirewallConfig_MultipleInboundAndOutbound(tagName, rName),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAbrhaFirewallConfig_OnlyInbound(rName string) string {
	return fmt.Sprintf(`
resource "abrha_firewall" "foobar" {
  name = "%s"
  inbound_rule {
    protocol         = "tcp"
    port_range       = "22"
    source_addresses = ["0.0.0.0/0", "::/0"]
  }

}
	`, rName)
}

func testAccAbrhaFirewallConfig_OnlyOutbound(rName string) string {
	return fmt.Sprintf(`
resource "abrha_firewall" "foobar" {
  name = "%s"
  outbound_rule {
    protocol              = "tcp"
    port_range            = "22"
    destination_addresses = ["0.0.0.0/0", "::/0"]
  }

}
	`, rName)
}

func testAccAbrhaFirewallConfig_OnlyMultipleInbound(rName string) string {
	return fmt.Sprintf(`
resource "abrha_firewall" "foobar" {
  name = "%s"
  inbound_rule {
    protocol         = "tcp"
    port_range       = "22"
    source_addresses = ["0.0.0.0/0", "::/0"]
  }
  inbound_rule {
    protocol         = "tcp"
    port_range       = "80"
    source_addresses = ["1.2.3.0/24", "2002::/16"]
  }

}
	`, rName)
}

func testAccAbrhaFirewallConfig_OnlyMultipleOutbound(rName string) string {
	return fmt.Sprintf(`
resource "abrha_firewall" "foobar" {
  name = "%s"
  outbound_rule {
    protocol              = "tcp"
    port_range            = "22"
    destination_addresses = ["192.168.1.0/24", "2002:1001::/48"]
  }
  outbound_rule {
    protocol              = "udp"
    port_range            = "53"
    destination_addresses = ["1.2.3.0/24", "2002::/16"]
  }

}
	`, rName)
}

func testAccAbrhaFirewallConfig_MultipleInboundAndOutbound(tagName string, rName string) string {
	return fmt.Sprintf(`
resource "abrha_tag" "foobar" {
  name = "%s"
}

resource "abrha_firewall" "foobar" {
  name = "%s"
  inbound_rule {
    protocol         = "tcp"
    port_range       = "22"
    source_addresses = ["0.0.0.0/0", "::/0"]
  }
  inbound_rule {
    protocol         = "tcp"
    port_range       = "443"
    source_addresses = ["192.168.1.0/24", "2002:1001:1:2::/64"]
    source_tags      = ["%s"]
  }
  outbound_rule {
    protocol              = "tcp"
    port_range            = "443"
    destination_addresses = ["192.168.1.0/24", "2002:1001:1:2::/64"]
    destination_tags      = ["%s"]
  }
  outbound_rule {
    protocol              = "udp"
    port_range            = "53"
    destination_addresses = ["0.0.0.0/0", "::/0"]
  }

}
	`, tagName, rName, tagName, tagName)
}

func testAccAbrhaFirewallConfig_fullPortRange(rName string) string {
	return fmt.Sprintf(`
resource "abrha_firewall" "foobar" {
  name = "%s"
  inbound_rule {
    protocol         = "tcp"
    port_range       = "all"
    source_addresses = ["192.168.1.1/32"]
  }
  outbound_rule {
    protocol              = "tcp"
    port_range            = "all"
    destination_addresses = ["192.168.1.2/32"]
  }
}
`, rName)
}

func testAccAbrhaFirewallConfig_icmp(rName string) string {
	return fmt.Sprintf(`
resource "abrha_firewall" "foobar" {
  name = "%s"
  inbound_rule {
    protocol         = "icmp"
    source_addresses = ["192.168.1.1/32"]
  }
  outbound_rule {
    protocol              = "icmp"
    port_range            = "1-65535"
    destination_addresses = ["192.168.1.2/32"]
  }
}
`, rName)
}

func testAccCheckAbrhaFirewallDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "abrha_firewall" {
			continue
		}

		// Try to find the firewall
		_, _, err := client.Firewalls.Get(context.Background(), rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Firewall still exists")
		}
	}

	return nil
}

func testAccCheckAbrhaFirewallExists(n string, firewall *goApiAbrha.Firewall) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

		foundFirewall, _, err := client.Firewalls.Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundFirewall.ID != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		*firewall = *foundFirewall

		return nil
	}
}
