package vm_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAbrhaVms_Basic(t *testing.T) {
	name1 := acceptance.RandomTestName("01")
	name2 := acceptance.RandomTestName("02")

	resourcesConfig := fmt.Sprintf(`
resource "abrha_vm" "foo" {
  name   = "%s"
  size   = "%s"
  image  = "%s"
  region = "nyc3"
}

resource "abrha_vm" "bar" {
  name   = "%s"
  size   = "%s"
  image  = "%s"
  region = "nyc3"
}
`, name1, defaultSize, defaultImage, name2, defaultSize, defaultImage)

	datasourceConfig := fmt.Sprintf(`
data "abrha_vms" "result" {
  filter {
    key    = "name"
    values = ["%s"]
  }
}
`, name1)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourcesConfig,
			},
			{
				Config: resourcesConfig + datasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.abrha_vms.result", "vms.#", "1"),
					resource.TestCheckResourceAttr("data.abrha_vms.result", "vms.0.name", name1),
					resource.TestCheckResourceAttrPair("data.abrha_vms.result", "vms.0.id", "abrha_vm.foo", "id"),
				),
			},
			{
				Config: resourcesConfig,
			},
		},
	})
}

func TestAccDataSourceAbrhaVms_GPUVm(t *testing.T) {
	runGPU := os.Getenv(runGPUEnvVar)
	if runGPU == "" {
		t.Skip("'DO_RUN_GPU_TESTS' env var not set; Skipping tests that requires a GPU Vm")
	}

	keyName := acceptance.RandomTestName()
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("abrha@ssh-acceptance-test")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}

	name1 := acceptance.RandomTestName("gpu")
	name2 := acceptance.RandomTestName("regular")

	resourcesConfig := fmt.Sprintf(`
resource "abrha_ssh_key" "foobar" {
  name       = "%s"
  public_key = "%s"
}

resource "abrha_vm" "gpu" {
  name     = "%s"
  size     = "%s"
  image    = "%s"
  region   = "nyc2"
  ssh_keys = [abrha_ssh_key.foobar.id]
}

resource "abrha_vm" "regular" {
  name   = "%s"
  size   = "%s"
  image  = "%s"
  region = "nyc2"
}
`, keyName, publicKeyMaterial, name1, gpuSize, gpuImage, name2, defaultSize, defaultImage)

	datasourceConfig := `
data "abrha_vms" "result" {
  gpus = true
}
`
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourcesConfig,
			},
			{
				Config: resourcesConfig + datasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.abrha_vms.result", "vms.#", "1"),
					resource.TestCheckResourceAttr("data.abrha_vms.result", "vms.0.name", name1),
					resource.TestCheckResourceAttrPair("data.abrha_vms.result", "vms.0.id", "abrha_vm.gpu", "id"),
				),
			},
			{
				Config: resourcesConfig,
			},
		},
	})
}
