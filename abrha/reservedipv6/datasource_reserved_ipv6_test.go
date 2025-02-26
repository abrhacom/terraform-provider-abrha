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

func TestAccDataSourceAbrhaReservedIPV6_Basic(t *testing.T) {
	var reservedIPv6 goApiAbrha.ReservedIPV6

	expectedURNRegex, _ := regexp.Compile(`do:reservedipv6:` + ipv6Regex)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceAbrhaReservedIPConfig_Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceAbrhaReservedIPV6Exists("data.abrha_reserved_ipv6.foobar", &reservedIPv6),
					resource.TestCheckResourceAttrSet(
						"data.abrha_reserved_ipv6.foobar", "ip"),
					resource.TestCheckResourceAttr(
						"data.abrha_reserved_ipv6.foobar", "region_slug", "nyc3"),
					resource.TestMatchResourceAttr("data.abrha_reserved_ipv6.foobar", "urn", expectedURNRegex),
				),
			},
		},
	})
}

func TestAccDataSourceAbrhaReservedIPV6_FindsReservedIP(t *testing.T) {
	var reservedIPv6 goApiAbrha.ReservedIPV6

	expectedURNRegex, _ := regexp.Compile(`do:reservedipv6:` + ipv6Regex)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceAbrhaReservedIPConfig_FindsFloatingIP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceAbrhaReservedIPV6Exists("data.abrha_reserved_ipv6.foobar", &reservedIPv6),
					resource.TestCheckResourceAttrSet(
						"data.abrha_reserved_ipv6.foobar", "ip"),
					resource.TestCheckResourceAttr(
						"data.abrha_reserved_ipv6.foobar", "region_slug", "nyc3"),
					resource.TestMatchResourceAttr("data.abrha_reserved_ipv6.foobar", "urn", expectedURNRegex),
				),
			},
		},
	})
}

func testAccCheckDataSourceAbrhaReservedIPV6Exists(n string, reservedIPv6 *goApiAbrha.ReservedIPV6) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No reserved IPv6 ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

		foundReservedIP, _, err := client.ReservedIPV6s.Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundReservedIP.IP != rs.Primary.ID {
			return fmt.Errorf("reserved IPv6 not found")
		}

		*reservedIPv6 = *foundReservedIP

		return nil
	}
}

const testAccCheckDataSourceAbrhaReservedIPConfig_FindsFloatingIP = `
resource "abrha_reserved_ipv6" "foo" {
  region_slug = "nyc3"
}

data "abrha_reserved_ipv6" "foobar" {
  ip = abrha_reserved_ipv6.foo.ip
}`

const testAccCheckDataSourceAbrhaReservedIPConfig_Basic = `
resource "abrha_reserved_ipv6" "foo" {
  region_slug = "nyc3"
}

data "abrha_reserved_ipv6" "foobar" {
  ip = abrha_reserved_ipv6.foo.ip
}`
