package domain_test

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

func TestAccDataSourceAbrhaDomain_Basic(t *testing.T) {
	var domain goApiAbrha.Domain
	domainName := acceptance.RandomTestName() + ".com"
	expectedURN := fmt.Sprintf("do:domain:%s", domainName)

	resourceConfig := fmt.Sprintf(`
resource "abrha_domain" "foo" {
  name       = "%s"
  ip_address = "192.168.0.10"
}
`, domainName)

	dataSourceConfig := `
data "abrha_domain" "foobar" {
  name = abrha_domain.foo.name
}`

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
					testAccCheckDataSourceAbrhaDomainExists("data.abrha_domain.foobar", &domain),
					testAccCheckDataSourceAbrhaDomainAttributes(&domain, domainName),
					resource.TestCheckResourceAttr(
						"data.abrha_domain.foobar", "name", domainName),
					resource.TestCheckResourceAttr(
						"data.abrha_domain.foobar", "urn", expectedURN),
				),
			},
		},
	})
}

func testAccCheckDataSourceAbrhaDomainAttributes(domain *goApiAbrha.Domain, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if domain.Name != name {
			return fmt.Errorf("Bad name: %s", domain.Name)
		}

		return nil
	}
}

func testAccCheckDataSourceAbrhaDomainExists(n string, domain *goApiAbrha.Domain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

		foundDomain, _, err := client.Domains.Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundDomain.Name != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		*domain = *foundDomain

		return nil
	}
}
