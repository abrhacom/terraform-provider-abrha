package domain_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceAbrhaRecord_Basic(t *testing.T) {
	var record goApiAbrha.DomainRecord
	recordDomain := fmt.Sprintf("%s.com", acceptance.RandomTestName())
	recordName := acceptance.RandomTestName("record")
	resourceConfig := fmt.Sprintf(`
resource "abrha_domain" "foo" {
  name       = "%s"
  ip_address = "192.168.0.10"
}

resource "abrha_record" "foo" {
  domain = abrha_domain.foo.name
  type   = "A"
  name   = "%s"
  value  = "192.168.0.10"
}`, recordDomain, recordName)
	dataSourceConfig := `
data "abrha_record" "foobar" {
  name   = abrha_record.foo.name
  domain = abrha_domain.foo.name
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
					testAccCheckDataSourceAbrhaRecordExists("data.abrha_record.foobar", &record),
					testAccCheckDataSourceAbrhaRecordAttributes(&record, recordName, "A"),
					resource.TestCheckResourceAttr(
						"data.abrha_record.foobar", "name", recordName),
					resource.TestCheckResourceAttr(
						"data.abrha_record.foobar", "type", "A"),
				),
			},
		},
	})
}

func testAccCheckDataSourceAbrhaRecordAttributes(record *goApiAbrha.DomainRecord, name string, r_type string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if record.Name != name {
			return fmt.Errorf("Bad name: %s", record.Name)
		}

		if record.Type != r_type {
			return fmt.Errorf("Bad type: %s", record.Type)
		}

		return nil
	}
}

func testAccCheckDataSourceAbrhaRecordExists(n string, record *goApiAbrha.DomainRecord) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		domain := rs.Primary.Attributes["domain"]
		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		foundRecord, _, err := client.Domains.Record(context.Background(), domain, id)
		if err != nil {
			return err
		}

		if foundRecord.Name != rs.Primary.Attributes["name"] {
			return fmt.Errorf("Record not found")
		}

		*record = *foundRecord

		return nil
	}
}
