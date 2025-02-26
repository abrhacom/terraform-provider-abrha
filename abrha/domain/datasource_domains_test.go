package domain_test

import (
	"fmt"
	"testing"

	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAbrhaDomains_Basic(t *testing.T) {
	name1 := acceptance.RandomTestName() + ".com"
	name2 := acceptance.RandomTestName() + ".com"

	resourcesConfig := fmt.Sprintf(`
resource "abrha_domain" "foo" {
  name = "%s"
}

resource "abrha_domain" "bar" {
  name = "%s"
}
`, name1, name2)

	datasourceConfig := fmt.Sprintf(`
data "abrha_domains" "result" {
  filter {
    key    = "name"
    values = ["%s"]
  }
}
`, name1)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourcesConfig,
			},
			{
				Config: resourcesConfig + datasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.abrha_domains.result", "domains.#", "1"),
					resource.TestCheckResourceAttrPair("data.abrha_domains.result", "domains.0.name", "abrha_domain.foo", "name"),
					resource.TestCheckResourceAttrPair("data.abrha_domains.result", "domains.0.urn", "abrha_domain.foo", "urn"),
					resource.TestCheckResourceAttrPair("data.abrha_domains.result", "domains.0.ttl", "abrha_domain.foo", "ttl"),
				),
			},
			{
				Config: resourcesConfig,
			},
		},
	})
}
