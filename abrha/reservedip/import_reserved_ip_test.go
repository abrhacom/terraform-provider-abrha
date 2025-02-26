package reservedip_test

import (
	"testing"

	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAbrhaReservedIP_importBasicRegion(t *testing.T) {
	resourceName := "abrha_reserved_ip.foobar"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaReservedIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaReservedIPConfig_region,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAbrhaReservedIP_importBasicVm(t *testing.T) {
	resourceName := "abrha_reserved_ip.foobar"
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaReservedIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaReservedIPConfig_vm(name),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
