package registry_test

import (
	"context"
	"fmt"
	"testing"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceAbrhaContainerRegistry_Basic(t *testing.T) {
	var reg goApiAbrha.Registry
	regName := acceptance.RandomTestName()

	resourceConfig := fmt.Sprintf(`
resource "abrha_container_registry" "foo" {
  name                   = "%s"
  subscription_tier_slug = "basic"
}
`, regName)

	dataSourceConfig := `
data "abrha_container_registry" "foobar" {
  name = abrha_container_registry.foo.name
}
`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
			},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceAbrhaContainerRegistryExists("data.abrha_container_registry.foobar", &reg),
					resource.TestCheckResourceAttr(
						"data.abrha_container_registry.foobar", "name", regName),
					resource.TestCheckResourceAttr(
						"data.abrha_container_registry.foobar", "subscription_tier_slug", "basic"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_container_registry.foobar", "region"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_container_registry.foobar", "created_at"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_container_registry.foobar", "storage_usage_bytes"),
				),
			},
		},
	})
}

func testAccCheckDataSourceAbrhaContainerRegistryExists(n string, reg *goApiAbrha.Registry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No registry ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

		foundReg, _, err := client.Registry.Get(context.Background())

		if err != nil {
			return err
		}

		if foundReg.Name != rs.Primary.ID {
			return fmt.Errorf("Registry not found")
		}

		*reg = *foundReg

		return nil
	}
}
