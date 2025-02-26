package reservedipv6_test

import (
	"testing"

	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAbrhaReservedIPV6_importBasicRegion(t *testing.T) {
	resourceName := "abrha_reserved_ipv6.foobar"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaReservedIPV6Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaReservedIPV6Config_regionSlug,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
