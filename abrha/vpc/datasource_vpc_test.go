package vpc_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAbrhaVPC_ByName(t *testing.T) {
	vpcName := acceptance.RandomTestName()
	vpcDesc := "A description for the VPC"
	resourceConfig := fmt.Sprintf(testAccCheckDataSourceAbrhaVPCConfig_Basic, vpcName, vpcDesc)
	dataSourceConfig := `
data "abrha_vpc" "foobar" {
  name = abrha_vpc.foobar.name
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
					testAccCheckAbrhaVPCExists("data.abrha_vpc.foobar"),
					resource.TestCheckResourceAttr(
						"data.abrha_vpc.foobar", "name", vpcName),
					resource.TestCheckResourceAttr(
						"data.abrha_vpc.foobar", "description", vpcDesc),
					resource.TestCheckResourceAttrSet(
						"data.abrha_vpc.foobar", "default"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_vpc.foobar", "ip_range"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_vpc.foobar", "created_at"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_vpc.foobar", "urn"),
				),
			},
		},
	})
}

func TestAccDataSourceAbrhaVPC_RegionDefault(t *testing.T) {
	vpcVmName := acceptance.RandomTestName()
	vpcConfigRegionDefault := fmt.Sprintf(testAccCheckDataSourceAbrhaVPCConfig_RegionDefault, vpcVmName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: vpcConfigRegionDefault,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaVPCExists("data.abrha_vpc.foobar"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_vpc.foobar", "name"),
					resource.TestCheckResourceAttr(
						"data.abrha_vpc.foobar", "default", "true"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_vpc.foobar", "created_at"),
				),
			},
		},
	})
}

func TestAccDataSourceAbrhaVPC_ExpectErrors(t *testing.T) {
	vpcName := acceptance.RandomTestName()
	vpcNotExist := fmt.Sprintf(testAccCheckDataSourceAbrhaVPCConfig_DoesNotExist, vpcName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckDataSourceAbrhaVPCConfig_MissingRegionDefault,
				ExpectError: regexp.MustCompile(`unable to find default VPC in foo region`),
			},
			{
				Config:      vpcNotExist,
				ExpectError: regexp.MustCompile(`no VPCs found with name`),
			},
		},
	})
}

const testAccCheckDataSourceAbrhaVPCConfig_Basic = `
resource "abrha_vpc" "foobar" {
  name        = "%s"
  description = "%s"
  region      = "nyc3"
}`

const testAccCheckDataSourceAbrhaVPCConfig_RegionDefault = `
// Create Vm to ensure default VPC exists
resource "abrha_vm" "foo" {
  image              = "ubuntu-22-04-x64"
  name               = "%s"
  region             = "nyc3"
  size               = "s-1vcpu-1gb"
  private_networking = "true"
}

data "abrha_vpc" "foobar" {
  region = "nyc3"
}
`

const testAccCheckDataSourceAbrhaVPCConfig_MissingRegionDefault = `
data "abrha_vpc" "foobar" {
  region = "foo"
}
`

const testAccCheckDataSourceAbrhaVPCConfig_DoesNotExist = `
data "abrha_vpc" "foobar" {
  name = "%s"
}
`
