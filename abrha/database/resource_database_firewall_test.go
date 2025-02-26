package database_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAbrhaDatabaseFirewall_Basic(t *testing.T) {
	databaseClusterName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaDatabaseFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseFirewallConfigBasic, databaseClusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"abrha_database_firewall.example", "rule.#", "1"),
				),
			},
			// Add a new rule
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseFirewallConfigAddRule, databaseClusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"abrha_database_firewall.example", "rule.#", "2"),
				),
			},
			// Remove an existing rule
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseFirewallConfigBasic, databaseClusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"abrha_database_firewall.example", "rule.#", "1"),
				),
			},
		},
	})
}

func TestAccAbrhaDatabaseFirewall_MultipleResourceTypes(t *testing.T) {
	dbName := acceptance.RandomTestName()
	vmName := acceptance.RandomTestName()
	tagName := acceptance.RandomTestName()
	appName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaDatabaseFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseFirewallConfigMultipleResourceTypes,
					dbName, vmName, tagName, appName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"abrha_database_firewall.example", "rule.#", "4"),
				),
			},
		},
	})
}

func testAccCheckAbrhaDatabaseFirewallDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "abrha_database_firewall" {
			continue
		}

		clusterId := rs.Primary.Attributes["cluster_id"]

		_, _, err := client.Databases.GetFirewallRules(context.Background(), clusterId)
		if err == nil {
			return fmt.Errorf("DatabaseFirewall still exists")
		}
	}

	return nil
}

const testAccCheckAbrhaDatabaseFirewallConfigBasic = `
resource "abrha_database_cluster" "foobar" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

resource "abrha_database_firewall" "example" {
  cluster_id = abrha_database_cluster.foobar.id

  rule {
    type  = "ip_addr"
    value = "192.168.1.1"
  }
}
`

const testAccCheckAbrhaDatabaseFirewallConfigAddRule = `
resource "abrha_database_cluster" "foobar" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

resource "abrha_database_firewall" "example" {
  cluster_id = abrha_database_cluster.foobar.id

  rule {
    type  = "ip_addr"
    value = "192.168.1.1"
  }

  rule {
    type  = "ip_addr"
    value = "192.0.2.0"
  }
}
`

const testAccCheckAbrhaDatabaseFirewallConfigMultipleResourceTypes = `
resource "abrha_database_cluster" "foobar" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

resource "abrha_vm" "foobar" {
  name   = "%s"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = "nyc3"
}

resource "abrha_tag" "foobar" {
  name = "%s"
}

resource "abrha_app" "foobar" {
  spec {
    name   = "%s"
    region = "nyc"

    service {
      name               = "go-service"
      environment_slug   = "go"
      instance_count     = 1
      instance_size_slug = "basic-xxs"

      git {
        repo_clone_url = "https://github.com/abrha/sample-golang.git"
        branch         = "main"
      }
    }
  }
}

resource "abrha_database_firewall" "example" {
  cluster_id = abrha_database_cluster.foobar.id

  rule {
    type  = "ip_addr"
    value = "192.168.1.1"
  }

  rule {
    type  = "vm"
    value = abrha_vm.foobar.id
  }

  rule {
    type  = "tag"
    value = abrha_tag.foobar.name
  }

  rule {
    type  = "app"
    value = abrha_app.foobar.id
  }
}
`
