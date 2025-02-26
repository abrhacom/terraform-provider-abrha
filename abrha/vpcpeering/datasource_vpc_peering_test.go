package vpcpeering_test

import (
	"fmt"
	"regexp"
	"testing"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAbrhaVPCPeering_ByID(t *testing.T) {
	var vpcPeering goApiAbrha.VPCPeering
	vpcPeeringName := acceptance.RandomTestName()
	vpcName1 := acceptance.RandomTestName()
	vpcName2 := acceptance.RandomTestName()
	resourceConfig := fmt.Sprintf(testAccCheckDataSourceAbrhaVPCPeeringConfig_Basic, vpcName1, vpcName2, vpcPeeringName)
	dataSourceConfig := `
data "abrha_vpc_peering" "foobar" {
  id = abrha_vpc_peering.foobar.id
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
					testAccCheckAbrhaVPCPeeringExists("data.abrha_vpc_peering.foobar", &vpcPeering),
					resource.TestCheckResourceAttr(
						"data.abrha_vpc_peering.foobar", "name", vpcPeeringName),
					resource.TestCheckResourceAttr(
						"data.abrha_vpc_peering.foobar", "vpc_ids.#", "2"),
					resource.TestCheckResourceAttrPair(
						"data.abrha_vpc_peering.foobar", "vpc_ids.0", "abrha_vpc.vpc1", "id"),
					resource.TestCheckResourceAttrPair(
						"data.abrha_vpc_peering.foobar", "vpc_ids.1", "abrha_vpc.vpc2", "id"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_vpc_peering.foobar", "created_at"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_vpc_peering.foobar", "status"),
				),
			},
		},
	})
}

func TestAccDataSourceAbrhaVPCPeering_ByName(t *testing.T) {
	var vpcPeering goApiAbrha.VPCPeering
	vpcPeeringName := acceptance.RandomTestName()
	vpcName1 := acceptance.RandomTestName()
	vpcName2 := acceptance.RandomTestName()
	resourceConfig := fmt.Sprintf(testAccCheckDataSourceAbrhaVPCPeeringConfig_Basic, vpcName1, vpcName2, vpcPeeringName)
	dataSourceConfig := `
data "abrha_vpc_peering" "foobar" {
  name = abrha_vpc_peering.foobar.name
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
					testAccCheckAbrhaVPCPeeringExists("data.abrha_vpc_peering.foobar", &vpcPeering),
					resource.TestCheckResourceAttr(
						"abrha_vpc_peering.foobar", "name", vpcPeeringName),
					resource.TestCheckResourceAttr(
						"abrha_vpc_peering.foobar", "vpc_ids.#", "2"),
					resource.TestCheckResourceAttrPair(
						"abrha_vpc_peering.foobar", "vpc_ids.0", "abrha_vpc.vpc1", "id"),
					resource.TestCheckResourceAttrPair(
						"abrha_vpc_peering.foobar", "vpc_ids.1", "abrha_vpc.vpc2", "id"),
					resource.TestCheckResourceAttrSet(
						"abrha_vpc_peering.foobar", "created_at"),
					resource.TestCheckResourceAttrSet(
						"abrha_vpc_peering.foobar", "status"),
				),
			},
		},
	})
}

func TestAccDataSourceAbrhaVPCPeering_ExpectErrors(t *testing.T) {
	vpcPeeringName := acceptance.RandomTestName()
	vpcPeeringNotExist := fmt.Sprintf(testAccCheckDataSourceAbrhaVPCPeeringConfig_DoesNotExist, vpcPeeringName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      vpcPeeringNotExist,
				ExpectError: regexp.MustCompile(`Error retrieving VPC Peering`),
			},
		},
	})
}

const testAccCheckDataSourceAbrhaVPCPeeringConfig_Basic = `
resource "abrha_vpc" "vpc1" {
  name   = "%s"
  region = "nyc3"
}

resource "abrha_vpc" "vpc2" {
  name   = "%s"
  region = "nyc3"
}

resource "abrha_vpc_peering" "foobar" {
  name = "%s"
  vpc_ids = [
    abrha_vpc.vpc1.id,
    abrha_vpc.vpc2.id
  ]
  depends_on = [
    abrha_vpc.vpc1,
    abrha_vpc.vpc2
  ]
}
`

const testAccCheckDataSourceAbrhaVPCPeeringConfig_DoesNotExist = `
data "abrha_vpc_peering" "foobar" {
  id = "%s"
}
`
