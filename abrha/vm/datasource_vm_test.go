package vm_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceAbrhaVm_BasicByName(t *testing.T) {
	var vm goApiAbrha.Vm
	name := acceptance.RandomTestName()
	resourceConfig := testAccCheckDataSourceAbrhaVmConfig_basicByName(name)
	dataSourceConfig := `
data "abrha_vm" "foobar" {
  name = abrha_vm.foo.name
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
			},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceAbrhaVmExists("data.abrha_vm.foobar", &vm),
					resource.TestCheckResourceAttr(
						"data.abrha_vm.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"data.abrha_vm.foobar", "image", "ubuntu-22-04-x64"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm.foobar", "ipv6", "true"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm.foobar", "private_networking", "true"),
					resource.TestCheckResourceAttrSet("data.abrha_vm.foobar", "urn"),
					resource.TestCheckResourceAttrSet("data.abrha_vm.foobar", "created_at"),
					resource.TestCheckResourceAttrSet("data.abrha_vm.foobar", "vpc_uuid"),
				),
			},
		},
	})
}

func TestAccDataSourceAbrhaVm_GPUByName(t *testing.T) {
	runGPU := os.Getenv(runGPUEnvVar)
	if runGPU == "" {
		t.Skip("'DO_RUN_GPU_TESTS' env var not set; Skipping tests that requires a GPU VM")
	}

	keyName := acceptance.RandomTestName()
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("abrha@ssh-acceptance-test")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}

	var vm goApiAbrha.Vm
	name := acceptance.RandomTestName()
	resourceConfig := testAccCheckDataSourceAbrhaVmConfig_gpuByName(keyName, publicKeyMaterial, name)
	dataSourceConfig := `
data "abrha_vm" "foobar" {
  name = abrha_vm.foo.name
  gpu  = true
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
			},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceAbrhaVmExists("data.abrha_vm.foobar", &vm),
					resource.TestCheckResourceAttr(
						"data.abrha_vm.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"data.abrha_vm.foobar", "image", gpuImage),
					resource.TestCheckResourceAttr(
						"data.abrha_vm.foobar", "region", "tor1"),
					resource.TestCheckResourceAttrSet("data.abrha_vm.foobar", "urn"),
					resource.TestCheckResourceAttrSet("data.abrha_vm.foobar", "created_at"),
				),
			},
		},
	})
}

func TestAccDataSourceAbrhaVm_BasicById(t *testing.T) {
	var vm goApiAbrha.Vm
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceAbrhaVmConfig_basicById(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceAbrhaVmExists("data.abrha_vm.foobar", &vm),
					resource.TestCheckResourceAttr(
						"data.abrha_vm.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"data.abrha_vm.foobar", "image", "ubuntu-22-04-x64"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm.foobar", "ipv6", "true"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm.foobar", "private_networking", "true"),
					resource.TestCheckResourceAttrSet("data.abrha_vm.foobar", "urn"),
					resource.TestCheckResourceAttrSet("data.abrha_vm.foobar", "created_at"),
				),
			},
		},
	})
}

func TestAccDataSourceAbrhaVm_BasicByTag(t *testing.T) {
	var vm goApiAbrha.Vm
	name := acceptance.RandomTestName()
	tagName := acceptance.RandomTestName("tag")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceAbrhaVmConfig_basicWithTag(tagName, name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foo", &vm),
				),
			},
			{
				Config: testAccCheckDataSourceAbrhaVmConfig_basicByTag(tagName, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceAbrhaVmExists("data.abrha_vm.foobar", &vm),
					resource.TestCheckResourceAttr(
						"data.abrha_vm.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"data.abrha_vm.foobar", "image", "ubuntu-22-04-x64"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm.foobar", "ipv6", "true"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm.foobar", "private_networking", "true"),
					resource.TestCheckResourceAttr(
						"data.abrha_vm.foobar", "tags.#", "1"),
					resource.TestCheckResourceAttrSet("data.abrha_vm.foobar", "urn"),
					resource.TestCheckResourceAttrSet("data.abrha_vm.foobar", "created_at"),
				),
			},
		},
	})
}

func testAccCheckDataSourceAbrhaVmExists(n string, vm *goApiAbrha.Vm) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No vm ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

		id := rs.Primary.ID

		foundVm, _, err := client.Vms.Get(context.Background(), id)

		if err != nil {
			return err
		}

		if foundVm.ID != id {
			return fmt.Errorf("Vm not found")
		}

		*vm = *foundVm

		return nil
	}
}

func testAccCheckDataSourceAbrhaVmConfig_basicByName(name string) string {
	return fmt.Sprintf(`
resource "abrha_vpc" "foobar" {
  name   = "%s"
  region = "nyc3"
}

resource "abrha_vm" "foo" {
  name     = "%s"
  size     = "%s"
  image    = "%s"
  region   = "nyc3"
  ipv6     = true
  vpc_uuid = abrha_vpc.foobar.id
}`, acceptance.RandomTestName(), name, defaultSize, defaultImage)
}

func testAccCheckDataSourceAbrhaVmConfig_gpuByName(keyName, key, name string) string {
	return fmt.Sprintf(`
resource "abrha_ssh_key" "foobar" {
  name       = "%s"
  public_key = "%s"
}

resource "abrha_vm" "foo" {
  name     = "%s"
  size     = "%s"
  image    = "%s"
  region   = "tor1"
  ssh_keys = [abrha_ssh_key.foobar.id]
}`, keyName, key, name, gpuSize, gpuImage)
}

func testAccCheckDataSourceAbrhaVmConfig_basicById(name string) string {
	return fmt.Sprintf(`
resource "abrha_vm" "foo" {
  name   = "%s"
  size   = "%s"
  image  = "%s"
  region = "nyc3"
  ipv6   = true
}

data "abrha_vm" "foobar" {
  id = abrha_vm.foo.id
}
`, name, defaultSize, defaultImage)
}

func testAccCheckDataSourceAbrhaVmConfig_basicWithTag(tagName string, name string) string {
	return fmt.Sprintf(`
resource "abrha_tag" "foo" {
  name = "%s"
}

resource "abrha_vm" "foo" {
  name   = "%s"
  size   = "%s"
  image  = "%s"
  region = "nyc3"
  ipv6   = true
  tags   = [abrha_tag.foo.id]
}
`, tagName, name, defaultSize, defaultImage)
}

func testAccCheckDataSourceAbrhaVmConfig_basicByTag(tagName string, name string) string {
	return fmt.Sprintf(`
resource "abrha_tag" "foo" {
  name = "%s"
}

resource "abrha_vm" "foo" {
  name   = "%s"
  size   = "%s"
  image  = "%s"
  region = "nyc3"
  ipv6   = true
  tags   = [abrha_tag.foo.id]
}

data "abrha_vm" "foobar" {
  tag = abrha_tag.foo.id
}
`, tagName, name, defaultSize, defaultImage)
}
