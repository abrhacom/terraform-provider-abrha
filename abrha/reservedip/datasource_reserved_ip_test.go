package reservedip_test

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

func TestAccDataSourceAbrhaReservedIP_Basic(t *testing.T) {
	var reservedIP goApiAbrha.ReservedIP

	expectedURNRegEx, _ := regexp.Compile(`do:reservedip:(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceAbrhaReservedIPConfig_Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceAbrhaReservedIPExists("data.abrha_reserved_ip.foobar", &reservedIP),
					resource.TestCheckResourceAttrSet(
						"data.abrha_reserved_ip.foobar", "ip_address"),
					resource.TestCheckResourceAttr(
						"data.abrha_reserved_ip.foobar", "region", "nyc3"),
					resource.TestMatchResourceAttr("data.abrha_reserved_ip.foobar", "urn", expectedURNRegEx),
				),
			},
		},
	})
}

func TestAccDataSourceAbrhaReservedIP_FindsFloatingIP(t *testing.T) {
	var reservedIP goApiAbrha.ReservedIP

	expectedURNRegEx, _ := regexp.Compile(`do:reservedip:(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceAbrhaReservedIPConfig_FindsFloatingIP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceAbrhaReservedIPExists("data.abrha_reserved_ip.foobar", &reservedIP),
					resource.TestCheckResourceAttrSet(
						"data.abrha_reserved_ip.foobar", "ip_address"),
					resource.TestCheckResourceAttr(
						"data.abrha_reserved_ip.foobar", "region", "nyc3"),
					resource.TestMatchResourceAttr("data.abrha_reserved_ip.foobar", "urn", expectedURNRegEx),
				),
			},
		},
	})
}

func testAccCheckDataSourceAbrhaReservedIPExists(n string, reservedIP *goApiAbrha.ReservedIP) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No reserved IP ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

		foundReservedIP, _, err := client.ReservedIPs.Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundReservedIP.IP != rs.Primary.ID {
			return fmt.Errorf("reserved IP not found")
		}

		*reservedIP = *foundReservedIP

		return nil
	}
}

const testAccCheckDataSourceAbrhaReservedIPConfig_FindsFloatingIP = `
resource "abrha_floating_ip" "foo" {
  region = "nyc3"
}

data "abrha_reserved_ip" "foobar" {
  ip_address = abrha_floating_ip.foo.ip_address
}`

const testAccCheckDataSourceAbrhaReservedIPConfig_Basic = `
resource "abrha_reserved_ip" "foo" {
  region = "nyc3"
}

data "abrha_reserved_ip" "foobar" {
  ip_address = abrha_reserved_ip.foo.ip_address
}`
