package reservedipv6_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const ipv6Regex = "(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$"

func TestAccAbrhaReservedIPV6_RegionSlug(t *testing.T) {
	var reservedIPv6 goApiAbrha.ReservedIPV6

	expectedURNRegex, _ := regexp.Compile(`do:reservedipv6:` + ipv6Regex)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaReservedIPV6Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaReservedIPV6Config_regionSlug,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaReservedIPV6Exists("abrha_reserved_ipv6.foobar", &reservedIPv6),
					resource.TestCheckResourceAttr(
						"abrha_reserved_ipv6.foobar", "region_slug", "nyc3"),
					resource.TestMatchResourceAttr("abrha_reserved_ipv6.foobar", "urn", expectedURNRegex),
				),
			},
		},
	})
}

func testAccCheckAbrhaReservedIPV6Destroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "abrha_reserved_ipv6" {
			continue
		}

		// Try to find the key
		_, _, err := client.ReservedIPV6s.Get(context.Background(), rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Reserved IPv6 still exists")
		}
	}

	return nil
}

func testAccCheckAbrhaReservedIPV6Exists(n string, reservedIPv6 *goApiAbrha.ReservedIPV6) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

		// Try to find the ReservedIPv6
		foundReservedIP, _, err := client.ReservedIPV6s.Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundReservedIP.IP != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		*reservedIPv6 = *foundReservedIP

		return nil
	}
}

var testAccCheckAbrhaReservedIPV6Config_regionSlug = `
resource "abrha_reserved_ipv6" "foobar" {
  region_slug = "nyc3"
}`
