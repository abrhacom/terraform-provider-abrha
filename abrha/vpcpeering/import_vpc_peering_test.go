package vpcpeering_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAbrhaVPCPeering_importBasic(t *testing.T) {
	resourceName := "abrha_vpc_peering.foobar"
	vpcPeeringName := acceptance.RandomTestName()
	vpcPeeringCreateConfig := fmt.Sprintf(testAccCheckAbrhaVPCPeeringConfig_Basic, vpcPeeringName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaVPCPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: vpcPeeringCreateConfig,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Test importing non-existent resource provides expected error.
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateId:     "123abc",
				ExpectError:       regexp.MustCompile(`(Please verify the ID is correct|Cannot import non-existent remote object)`),
			},
		},
	})
}
