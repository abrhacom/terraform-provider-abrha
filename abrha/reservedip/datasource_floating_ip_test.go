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

func TestAccDataSourceAbrhaFloatingIp_Basic(t *testing.T) {
	var floatingIp goApiAbrha.FloatingIP

	expectedURNRegEx, _ := regexp.Compile(`do:floatingip:(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceAbrhaFloatingIpConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceAbrhaFloatingIpExists("data.abrha_floating_ip.foobar", &floatingIp),
					resource.TestCheckResourceAttrSet(
						"data.abrha_floating_ip.foobar", "ip_address"),
					resource.TestCheckResourceAttr(
						"data.abrha_floating_ip.foobar", "region", "nyc3"),
					resource.TestMatchResourceAttr("data.abrha_floating_ip.foobar", "urn", expectedURNRegEx),
				),
			},
		},
	})
}

func testAccCheckDataSourceAbrhaFloatingIpExists(n string, floatingIp *goApiAbrha.FloatingIP) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No floating ip ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

		foundFloatingIp, _, err := client.FloatingIPs.Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundFloatingIp.IP != rs.Primary.ID {
			return fmt.Errorf("Floating ip not found")
		}

		*floatingIp = *foundFloatingIp

		return nil
	}
}

const testAccCheckDataSourceAbrhaFloatingIpConfig_basic = `
resource "abrha_floating_ip" "foo" {
  region = "nyc3"
}

data "abrha_floating_ip" "foobar" {
  ip_address = abrha_floating_ip.foo.ip_address
}`
