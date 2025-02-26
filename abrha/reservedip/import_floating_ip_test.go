package reservedip_test

import (
	"testing"

	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAbrhaFloatingIP_importBasicRegion(t *testing.T) {
	resourceName := "abrha_floating_ip.foobar"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaFloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaFloatingIPConfig_region,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAbrhaFloatingIP_importBasicVm(t *testing.T) {
	resourceName := "abrha_floating_ip.foobar"
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaFloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaFloatingIPConfig_vm(name),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
