package vmautoscale_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccParspackVmAutoscale_Static(t *testing.T) {
	var autoscalePool goApiAbrha.VmAutoscalePool
	name := acceptance.RandomTestName()

	createConfig := testAccCheckParspackVmAutoscaleConfig_static(name, 1)
	updateConfig := strings.ReplaceAll(createConfig, "target_number_instances = 1", "target_number_instances = 2")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckParspackVmAutoscaleDestroy,
		Steps: []resource.TestStep{
			{
				// Test create
				Config: createConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckParspackVmAutoscaleExists("abrha_vm_autoscale.foobar", &autoscalePool),
					resource.TestCheckResourceAttrSet("abrha_vm_autoscale.foobar", "id"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "config.#", "1"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "config.0.min_instances", "0"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "config.0.max_instances", "0"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "config.0.target_cpu_utilization", "0"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "config.0.target_memory_utilization", "0"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "config.0.cooldown_minutes", "0"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "config.0.target_number_instances", "1"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "vm_template.#", "1"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "vm_template.0.size", "c-2"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "vm_template.0.region", "nyc3"),
					resource.TestCheckResourceAttrSet(
						"abrha_vm_autoscale.foobar", "vm_template.0.image"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "vm_template.0.with_vm_agent", "true"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "vm_template.0.ipv6", "true"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "vm_template.0.user_data", "\n#cloud-config\nruncmd:\n- apt-get update\n- apt-get install -y stress-ng\n"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "vm_template.0.tags.#", "2"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "vm_template.0.ssh_keys.#", "2"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "status", "active"),
					resource.TestCheckResourceAttrSet(
						"abrha_vm_autoscale.foobar", "created_at"),
					resource.TestCheckResourceAttrSet(
						"abrha_vm_autoscale.foobar", "updated_at"),
				),
			},
			{
				// Test update (static scale up)
				Config: updateConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckParspackVmAutoscaleExists("abrha_vm_autoscale.foobar", &autoscalePool),
					resource.TestCheckResourceAttrSet("abrha_vm_autoscale.foobar", "id"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "config.#", "1"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "config.0.min_instances", "0"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "config.0.max_instances", "0"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "config.0.target_cpu_utilization", "0"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "config.0.target_memory_utilization", "0"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "config.0.cooldown_minutes", "0"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "config.0.target_number_instances", "2"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "vm_template.#", "1"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "vm_template.0.size", "c-2"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "vm_template.0.region", "nyc3"),
					resource.TestCheckResourceAttrSet(
						"abrha_vm_autoscale.foobar", "vm_template.0.image"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "vm_template.0.with_vm_agent", "true"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "vm_template.0.ipv6", "true"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "vm_template.0.user_data", "\n#cloud-config\nruncmd:\n- apt-get update\n- apt-get install -y stress-ng\n"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "vm_template.0.tags.#", "2"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "vm_template.0.ssh_keys.#", "2"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "status", "active"),
					resource.TestCheckResourceAttrSet(
						"abrha_vm_autoscale.foobar", "created_at"),
					resource.TestCheckResourceAttrSet(
						"abrha_vm_autoscale.foobar", "updated_at"),
				),
			},
		},
	})
}

func TestAccParspackVmAutoscale_Dynamic(t *testing.T) {
	var autoscalePool goApiAbrha.VmAutoscalePool
	name := acceptance.RandomTestName()

	createConfig := testAccCheckParspackVmAutoscaleConfig_dynamic(name, 1)
	updateConfig := strings.ReplaceAll(createConfig, "min_instances             = 1", "min_instances             = 2")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckParspackVmAutoscaleDestroy,
		Steps: []resource.TestStep{
			{
				// Test create
				Config: createConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckParspackVmAutoscaleExists("abrha_vm_autoscale.foobar", &autoscalePool),
					resource.TestCheckResourceAttrSet("abrha_vm_autoscale.foobar", "id"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "config.#", "1"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "config.0.min_instances", "1"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "config.0.max_instances", "3"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "config.0.target_cpu_utilization", "0.5"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "config.0.target_memory_utilization", "0.5"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "config.0.cooldown_minutes", "5"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "config.0.target_number_instances", "0"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "vm_template.#", "1"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "vm_template.0.size", "c-2"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "vm_template.0.region", "nyc3"),
					resource.TestCheckResourceAttrSet(
						"abrha_vm_autoscale.foobar", "vm_template.0.image"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "vm_template.0.with_vm_agent", "true"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "vm_template.0.ipv6", "true"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "vm_template.0.user_data", "\n#cloud-config\nruncmd:\n- apt-get update\n- apt-get install -y stress-ng\n"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "vm_template.0.tags.#", "2"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "vm_template.0.ssh_keys.#", "2"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "status", "active"),
					resource.TestCheckResourceAttrSet(
						"abrha_vm_autoscale.foobar", "created_at"),
					resource.TestCheckResourceAttrSet(
						"abrha_vm_autoscale.foobar", "updated_at"),
				),
			},
			{
				// Test update (dynamic scale up)
				Config: updateConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckParspackVmAutoscaleExists("abrha_vm_autoscale.foobar", &autoscalePool),
					resource.TestCheckResourceAttrSet("abrha_vm_autoscale.foobar", "id"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "config.#", "1"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "config.0.min_instances", "2"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "config.0.max_instances", "3"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "config.0.target_cpu_utilization", "0.5"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "config.0.target_memory_utilization", "0.5"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "config.0.cooldown_minutes", "5"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "config.0.target_number_instances", "0"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "vm_template.#", "1"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "vm_template.0.size", "c-2"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "vm_template.0.region", "nyc3"),
					resource.TestCheckResourceAttrSet(
						"abrha_vm_autoscale.foobar", "vm_template.0.image"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "vm_template.0.with_vm_agent", "true"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "vm_template.0.ipv6", "true"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "vm_template.0.user_data", "\n#cloud-config\nruncmd:\n- apt-get update\n- apt-get install -y stress-ng\n"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "vm_template.0.tags.#", "2"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "vm_template.0.ssh_keys.#", "2"),
					resource.TestCheckResourceAttr(
						"abrha_vm_autoscale.foobar", "status", "active"),
					resource.TestCheckResourceAttrSet(
						"abrha_vm_autoscale.foobar", "created_at"),
					resource.TestCheckResourceAttrSet(
						"abrha_vm_autoscale.foobar", "updated_at"),
				),
			},
		},
	})
}

func testAccCheckParspackVmAutoscaleExists(n string, autoscalePool *goApiAbrha.VmAutoscalePool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %v", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Resource ID not set")
		}
		// Check for valid ID response to validate that the resource has been created
		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()
		pool, _, err := client.VmAutoscale.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}
		if pool.ID != rs.Primary.ID {
			return fmt.Errorf("Vm autoscale pool not found")
		}
		*autoscalePool = *pool
		return nil
	}
}

func testAccCheckParspackVmAutoscaleDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "abrha_vm_autoscale" {
			continue
		}
		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()
		_, _, err := client.VmAutoscale.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			if strings.Contains(err.Error(), fmt.Sprintf("autoscale group with id %s not found", rs.Primary.ID)) {
				return nil
			}
			return fmt.Errorf("Vm autoscale pool still exists")
		}
	}
	return nil
}

func testAccCheckParspackVmAutoscaleConfig_static(name string, size int) string {
	pubKey1, _, err := acctest.RandSSHKeyPair("abrha@acceptance-test")
	if err != nil {
		fmt.Println("Unable to generate public key", err)
		return ""
	}

	pubKey2, _, err := acctest.RandSSHKeyPair("abrha@acceptance-test")
	if err != nil {
		fmt.Println("Unable to generate public key", err)
		return ""
	}

	return fmt.Sprintf(`
resource "abrha_ssh_key" "foo" {
  name       = "%s"
  public_key = "%s"
}

resource "abrha_ssh_key" "bar" {
  name       = "%s"
  public_key = "%s"
}

resource "abrha_tag" "foo" {
  name = "%s"
}

resource "abrha_tag" "bar" {
  name = "%s"
}

resource "abrha_vm_autoscale" "foobar" {
  name = "%s"

  config {
    target_number_instances = %d
  }

  vm_template {
    size          = "c-2"
    region        = "nyc3"
    image         = "ubuntu-24-04-x64"
    tags          = [abrha_tag.foo.id, abrha_tag.bar.id]
    ssh_keys      = [abrha_ssh_key.foo.id, abrha_ssh_key.bar.id]
    with_vm_agent = true
    ipv6          = true
    user_data     = "\n#cloud-config\nruncmd:\n- apt-get update\n- apt-get install -y stress-ng\n"
  }
}`, acceptance.RandomTestName("sshKey1"), pubKey1,
		acceptance.RandomTestName("sshKey2"), pubKey2,
		acceptance.RandomTestName("tag1"),
		acceptance.RandomTestName("tag2"),
		name, size)
}

func testAccCheckParspackVmAutoscaleConfig_dynamic(name string, size int) string {
	pubKey1, _, err := acctest.RandSSHKeyPair("abrha@acceptance-test")
	if err != nil {
		fmt.Println("Unable to generate public key", err)
		return ""
	}

	pubKey2, _, err := acctest.RandSSHKeyPair("abrha@acceptance-test")
	if err != nil {
		fmt.Println("Unable to generate public key", err)
		return ""
	}

	return fmt.Sprintf(`
resource "abrha_ssh_key" "foo" {
  name       = "%s"
  public_key = "%s"
}

resource "abrha_ssh_key" "bar" {
  name       = "%s"
  public_key = "%s"
}

resource "abrha_tag" "foo" {
  name = "%s"
}

resource "abrha_tag" "bar" {
  name = "%s"
}

resource "abrha_vm_autoscale" "foobar" {
  name = "%s"

  config {
    min_instances             = %d
    max_instances             = 3
    target_cpu_utilization    = 0.5
    target_memory_utilization = 0.5
    cooldown_minutes          = 5
  }

  vm_template {
    size          = "c-2"
    region        = "nyc3"
    image         = "ubuntu-24-04-x64"
    tags          = [abrha_tag.foo.id, abrha_tag.bar.id]
    ssh_keys      = [abrha_ssh_key.foo.id, abrha_ssh_key.bar.id]
    with_vm_agent = true
    ipv6          = true
    user_data     = "\n#cloud-config\nruncmd:\n- apt-get update\n- apt-get install -y stress-ng\n"
  }
}`, acceptance.RandomTestName("sshKey1"), pubKey1,
		acceptance.RandomTestName("sshKey2"), pubKey2,
		acceptance.RandomTestName("tag1"),
		acceptance.RandomTestName("tag2"),
		name, size)
}
