package vmautoscale_test

import (
	"testing"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAbrhaVmAutoscale_Static(t *testing.T) {
	var autoscalePool goApiAbrha.VmAutoscalePool
	name := acceptance.RandomTestName()

	createConfig := testAccCheckParspackVmAutoscaleConfig_static(name, 1)
	dataSourceIDConfig := `
data "abrha_vm_autoscale" "foo" {
  id = abrha_vm_autoscale.foobar.id
}`
	dataSourceNameConfig := `
data "abrha_vm_autoscale" "foo" {
  name = abrha_vm_autoscale.foobar.name
}`

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
				),
			},
			{
				// Import by id
				Config: createConfig + dataSourceIDConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckParspackVmAutoscaleExists("data.abrha_vm_autoscale.foo", &autoscalePool),
					resource.TestCheckResourceAttrSet("data.abrha_vm_autoscale.foo", "id"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "name", name),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "config.#", "1"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "config.0.min_instances", "0"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "config.0.max_instances", "0"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "config.0.target_cpu_utilization", "0"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "config.0.target_memory_utilization", "0"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "config.0.cooldown_minutes", "0"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "config.0.target_number_instances", "1"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "vm_template.#", "1"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "vm_template.0.size", "c-2"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "vm_template.0.region", "nyc3"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_vm_autoscale.foo", "vm_template.0.image"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "vm_template.0.with_vm_agent", "true"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "vm_template.0.ipv6", "true"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "vm_template.0.user_data", "\n#cloud-config\nruncmd:\n- apt-get update\n- apt-get install -y stress-ng\n"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "vm_template.0.tags.#", "2"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "vm_template.0.ssh_keys.#", "2"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "status", "active"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_vm_autoscale.foo", "created_at"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_vm_autoscale.foo", "updated_at"),
				),
			},
			{
				// Import by name
				Config: createConfig + dataSourceNameConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckParspackVmAutoscaleExists("data.abrha_vm_autoscale.foo", &autoscalePool),
					resource.TestCheckResourceAttrSet("data.abrha_vm_autoscale.foo", "id"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "name", name),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "config.#", "1"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "config.0.min_instances", "0"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "config.0.max_instances", "0"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "config.0.target_cpu_utilization", "0"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "config.0.target_memory_utilization", "0"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "config.0.cooldown_minutes", "0"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "config.0.target_number_instances", "1"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "vm_template.#", "1"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "vm_template.0.size", "c-2"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "vm_template.0.region", "nyc3"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_vm_autoscale.foo", "vm_template.0.image"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "vm_template.0.with_vm_agent", "true"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "vm_template.0.ipv6", "true"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "vm_template.0.user_data", "\n#cloud-config\nruncmd:\n- apt-get update\n- apt-get install -y stress-ng\n"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "vm_template.0.tags.#", "2"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "vm_template.0.ssh_keys.#", "2"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "status", "active"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_vm_autoscale.foo", "created_at"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_vm_autoscale.foo", "updated_at"),
				),
			},
		},
	})
}

func TestAccDataSourceAbrhaVmAutoscale_Dynamic(t *testing.T) {
	var autoscalePool goApiAbrha.VmAutoscalePool
	name := acceptance.RandomTestName()

	createConfig := testAccCheckParspackVmAutoscaleConfig_dynamic(name, 1)
	dataSourceIDConfig := `
data "abrha_vm_autoscale" "foo" {
  id = abrha_vm_autoscale.foobar.id
}`
	dataSourceNameConfig := `
data "abrha_vm_autoscale" "foo" {
  name = abrha_vm_autoscale.foobar.name
}`

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
				),
			},
			{
				// Import by id
				Config: createConfig + dataSourceIDConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckParspackVmAutoscaleExists("data.abrha_vm_autoscale.foo", &autoscalePool),
					resource.TestCheckResourceAttrSet("data.abrha_vm_autoscale.foo", "id"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "name", name),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "config.#", "1"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "config.0.min_instances", "1"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "config.0.max_instances", "3"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "config.0.target_cpu_utilization", "0.5"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "config.0.target_memory_utilization", "0.5"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "config.0.cooldown_minutes", "5"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "config.0.target_number_instances", "0"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "vm_template.#", "1"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "vm_template.0.size", "c-2"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "vm_template.0.region", "nyc3"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_vm_autoscale.foo", "vm_template.0.image"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "vm_template.0.with_vm_agent", "true"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "vm_template.0.ipv6", "true"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "vm_template.0.user_data", "\n#cloud-config\nruncmd:\n- apt-get update\n- apt-get install -y stress-ng\n"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "vm_template.0.tags.#", "2"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "vm_template.0.ssh_keys.#", "2"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "status", "active"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_vm_autoscale.foo", "created_at"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_vm_autoscale.foo", "updated_at"),
				),
			},
			{
				// Import by name
				Config: createConfig + dataSourceNameConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckParspackVmAutoscaleExists("data.abrha_vm_autoscale.foo", &autoscalePool),
					resource.TestCheckResourceAttrSet("data.abrha_vm_autoscale.foo", "id"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "name", name),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "config.#", "1"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "config.0.min_instances", "1"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "config.0.max_instances", "3"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "config.0.target_cpu_utilization", "0.5"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "config.0.target_memory_utilization", "0.5"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "config.0.cooldown_minutes", "5"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "config.0.target_number_instances", "0"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "vm_template.#", "1"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "vm_template.0.size", "c-2"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "vm_template.0.region", "nyc3"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_vm_autoscale.foo", "vm_template.0.image"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "vm_template.0.with_vm_agent", "true"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "vm_template.0.ipv6", "true"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "vm_template.0.user_data", "\n#cloud-config\nruncmd:\n- apt-get update\n- apt-get install -y stress-ng\n"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "vm_template.0.tags.#", "2"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "vm_template.0.ssh_keys.#", "2"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm_autoscale.foo", "status", "active"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_vm_autoscale.foo", "created_at"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_vm_autoscale.foo", "updated_at"),
				),
			},
		},
	})
}
