package vpcpeering_test

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

func TestAccAbrhaVPCPeering_Basic(t *testing.T) {
	var vpcPeering goApiAbrha.VPCPeering
	vpcPeeringName := acceptance.RandomTestName()
	vpcPeeringCreateConfig := fmt.Sprintf(testAccCheckAbrhaVPCPeeringConfig_Basic, vpcPeeringName)

	updateVPCPeeringName := acceptance.RandomTestName()
	vpcPeeringUpdateConfig := fmt.Sprintf(testAccCheckAbrhaVPCPeeringConfig_Basic, updateVPCPeeringName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaVPCPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: vpcPeeringCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaVPCPeeringExists("abrha_vpc_peering.foobar", &vpcPeering),
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
			{
				Config: vpcPeeringUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaVPCPeeringExists("abrha_vpc_peering.foobar", &vpcPeering),
					resource.TestCheckResourceAttr(
						"abrha_vpc_peering.foobar", "name", updateVPCPeeringName),
				),
			},
		},
	})
}

func testAccCheckAbrhaVPCPeeringDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "abrha_vpc_peering" {
			continue
		}

		_, _, err := client.VPCs.GetVPCPeering(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("VPC Peering resource still exists")
		}
	}

	return nil
}

func testAccCheckAbrhaVPCPeeringExists(resource string, vpcPeering *goApiAbrha.VPCPeering) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

		rs, ok := s.RootModule().Resources[resource]

		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID set for resource: %s", resource)
		}

		foundVPCPeering, _, err := client.VPCs.GetVPCPeering(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundVPCPeering.ID != rs.Primary.ID {
			return fmt.Errorf("Resource not found: %s : %s", resource, rs.Primary.ID)
		}

		*vpcPeering = *foundVPCPeering

		return nil
	}
}

const testAccCheckAbrhaVPCPeeringConfig_Basic = `
resource "abrha_vpc" "vpc1" {
  name   = "vpc1"
  region = "nyc3"
}

resource "abrha_vpc" "vpc2" {
  name   = "vpc2"
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
