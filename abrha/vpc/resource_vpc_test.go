package vpc_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAbrhaVPC_Basic(t *testing.T) {
	vpcName := acceptance.RandomTestName()
	vpcDesc := "A description for the VPC"
	vpcCreateConfig := fmt.Sprintf(testAccCheckAbrhaVPCConfig_Basic, vpcName, vpcDesc)

	updatedVPCName := acceptance.RandomTestName()
	updatedVPVDesc := "A brand new updated description for the VPC"
	vpcUpdateConfig := fmt.Sprintf(testAccCheckAbrhaVPCConfig_Basic, updatedVPCName, updatedVPVDesc)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaVPCDestroy,
		Steps: []resource.TestStep{
			{
				Config: vpcCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaVPCExists("abrha_vpc.foobar"),
					resource.TestCheckResourceAttr(
						"abrha_vpc.foobar", "name", vpcName),
					resource.TestCheckResourceAttr(
						"abrha_vpc.foobar", "default", "false"),
					resource.TestCheckResourceAttrSet(
						"abrha_vpc.foobar", "created_at"),
					resource.TestCheckResourceAttr(
						"abrha_vpc.foobar", "description", vpcDesc),
					resource.TestCheckResourceAttrSet(
						"abrha_vpc.foobar", "ip_range"),
					resource.TestCheckResourceAttrSet(
						"abrha_vpc.foobar", "urn"),
				),
			},
			{
				Config: vpcUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaVPCExists("abrha_vpc.foobar"),
					resource.TestCheckResourceAttr(
						"abrha_vpc.foobar", "name", updatedVPCName),
					resource.TestCheckResourceAttr(
						"abrha_vpc.foobar", "description", updatedVPVDesc),
					resource.TestCheckResourceAttr(
						"abrha_vpc.foobar", "default", "false"),
				),
			},
		},
	})
}

func TestAccAbrhaVPC_IPRange(t *testing.T) {
	vpcName := acceptance.RandomTestName()
	vpcCreateConfig := fmt.Sprintf(testAccCheckAbrhaVPCConfig_IPRange, vpcName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaVPCDestroy,
		Steps: []resource.TestStep{
			{
				Config: vpcCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaVPCExists("abrha_vpc.foobar"),
					resource.TestCheckResourceAttr(
						"abrha_vpc.foobar", "name", vpcName),
					resource.TestCheckResourceAttr(
						"abrha_vpc.foobar", "ip_range", "10.10.10.0/24"),
					resource.TestCheckResourceAttr(
						"abrha_vpc.foobar", "default", "false"),
				),
			},
		},
	})
}

// https://github.com/parspack/terraform-provider-parspack/issues/551
func TestAccAbrhaVPC_IPRangeRace(t *testing.T) {
	vpcNameOne := acceptance.RandomTestName()
	vpcNameTwo := acceptance.RandomTestName()
	vpcCreateConfig := fmt.Sprintf(testAccCheckAbrhaVPCConfig_IPRangeRace, vpcNameOne, vpcNameTwo)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaVPCDestroy,
		Steps: []resource.TestStep{
			{
				Config: vpcCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaVPCExists("abrha_vpc.foo"),
					testAccCheckAbrhaVPCExists("abrha_vpc.bar"),
					resource.TestCheckResourceAttr(
						"abrha_vpc.foo", "name", vpcNameOne),
					resource.TestCheckResourceAttrSet(
						"abrha_vpc.foo", "ip_range"),
					resource.TestCheckResourceAttr(
						"abrha_vpc.bar", "name", vpcNameTwo),
					resource.TestCheckResourceAttrSet(
						"abrha_vpc.bar", "ip_range"),
				),
			},
		},
	})
}

func testAccCheckAbrhaVPCDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "abrha_vpc" {
			continue
		}

		_, _, err := client.VPCs.Get(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("VPC resource still exists")
		}
	}

	return nil
}

func testAccCheckAbrhaVPCExists(resource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

		rs, ok := s.RootModule().Resources[resource]

		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID set for resource: %s", resource)
		}

		foundVPC, _, err := client.VPCs.Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundVPC.ID != rs.Primary.ID {
			return fmt.Errorf("Resource not found: %s : %s", resource, rs.Primary.ID)
		}

		return nil
	}
}

const testAccCheckAbrhaVPCConfig_Basic = `
resource "abrha_vpc" "foobar" {
  name        = "%s"
  description = "%s"
  region      = "nyc3"
}
`
const testAccCheckAbrhaVPCConfig_IPRange = `
resource "abrha_vpc" "foobar" {
  name     = "%s"
  region   = "nyc3"
  ip_range = "10.10.10.0/24"
}
`

const testAccCheckAbrhaVPCConfig_IPRangeRace = `
resource "abrha_vpc" "foo" {
  name   = "%s"
  region = "nyc3"
}

resource "abrha_vpc" "bar" {
  name   = "%s"
  region = "nyc3"
}
`
