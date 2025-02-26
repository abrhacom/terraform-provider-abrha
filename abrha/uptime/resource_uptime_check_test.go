package uptime_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testAccCheckAbrhaUptimeCheckConfig_Basic = `
resource "abrha_uptime_check" "foobar" {
  name    = "%s"
  target  = "%s"
  regions = ["%s"]
}
`

func TestAccAbrhaUptimeCheck_Basic(t *testing.T) {
	checkName := acceptance.RandomTestName()
	checkTarget := "https://www.landingpage.com"
	checkRegions := "eu_west"
	checkCreateConfig := fmt.Sprintf(testAccCheckAbrhaUptimeCheckConfig_Basic, checkName, checkTarget, checkRegions)

	updatedCheckName := acceptance.RandomTestName()
	updatedCheckRegions := "us_east"
	checkUpdateConfig := fmt.Sprintf(testAccCheckAbrhaUptimeCheckConfig_Basic, updatedCheckName, checkTarget, updatedCheckRegions)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaUptimeCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: checkCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaUptimeCheckExists("abrha_uptime_check.foobar"),
					resource.TestCheckResourceAttr(
						"abrha_uptime_check.foobar", "name", checkName),
					resource.TestCheckResourceAttr(
						"abrha_uptime_check.foobar", "target", checkTarget),
					resource.TestCheckResourceAttr("abrha_uptime_check.foobar", "regions.#", "1"),
					resource.TestCheckTypeSetElemAttr("abrha_uptime_check.foobar", "regions.*", "eu_west"),
				),
			},
			{
				Config: checkUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaUptimeCheckExists("abrha_uptime_check.foobar"),
					resource.TestCheckResourceAttr(
						"abrha_uptime_check.foobar", "name", updatedCheckName),
					resource.TestCheckResourceAttr(
						"abrha_uptime_check.foobar", "target", checkTarget),
					resource.TestCheckResourceAttr("abrha_uptime_check.foobar", "regions.#", "1"),
					resource.TestCheckTypeSetElemAttr("abrha_uptime_check.foobar", "regions.*", "us_east"),
				),
			},
		},
	})
}

func testAccCheckAbrhaUptimeCheckDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "abrha_uptime_check" {
			continue
		}

		_, _, err := client.UptimeChecks.Get(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Uptime Check resource still exists")
		}
	}

	return nil
}

func testAccCheckAbrhaUptimeCheckExists(resource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

		rs, ok := s.RootModule().Resources[resource]

		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID set for resource: %s", resource)
		}

		foundUptimeCheck, _, err := client.UptimeChecks.Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundUptimeCheck.ID != rs.Primary.ID {
			return fmt.Errorf("Resource not found: %s : %s", resource, rs.Primary.ID)
		}

		return nil
	}
}
