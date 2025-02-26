package domain_test

import (
	"fmt"
	"testing"

	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAbrhaRecords_Basic(t *testing.T) {
	name1 := acceptance.RandomTestName("records") + ".com"

	resourcesConfig := fmt.Sprintf(`
resource "abrha_domain" "foo" {
  name = "%s"
}

resource "abrha_record" "mail" {
  name     = "mail"
  domain   = abrha_domain.foo.name
  type     = "MX"
  priority = 10
  value    = "mail.example.com."
}

resource "abrha_record" "www" {
  name   = "www"
  domain = abrha_domain.foo.name
  type   = "A"
  value  = "192.168.1.1"
}
`, name1)

	datasourceConfig := fmt.Sprintf(`
data "abrha_records" "result" {
  domain = "%s"
  filter {
    key    = "type"
    values = ["A"]
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
					resource.TestCheckResourceAttr("data.abrha_records.result", "records.#", "1"),
					resource.TestCheckResourceAttr("data.abrha_records.result", "records.0.domain", name1),
					resource.TestCheckResourceAttrPair("data.abrha_records.result", "records.0.id", "abrha_record.www", "id"),
					resource.TestCheckResourceAttrPair("data.abrha_records.result", "records.0.name", "abrha_record.www", "name"),
					resource.TestCheckResourceAttrPair("data.abrha_records.result", "records.0.type", "abrha_record.www", "type"),
					resource.TestCheckResourceAttrPair("data.abrha_records.result", "records.0.value", "abrha_record.www", "value"),
					resource.TestCheckResourceAttrPair("data.abrha_records.result", "records.0.priority", "abrha_record.www", "priority"),
					resource.TestCheckResourceAttrPair("data.abrha_records.result", "records.0.port", "abrha_record.www", "port"),
					resource.TestCheckResourceAttrPair("data.abrha_records.result", "records.0.ttl", "abrha_record.www", "ttl"),
					resource.TestCheckResourceAttrPair("data.abrha_records.result", "records.0.weight", "abrha_record.www", "weight"),
					resource.TestCheckResourceAttrPair("data.abrha_records.result", "records.0.flags", "abrha_record.www", "flags"),
					resource.TestCheckResourceAttrPair("data.abrha_records.result", "records.0.tag", "abrha_record.www", "tag"),
				),
			},
		},
	})
}
