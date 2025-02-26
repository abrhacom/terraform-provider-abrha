package tag_test

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

func TestAccDataSourceAbrhaTag_Basic(t *testing.T) {
	var tag goApiAbrha.Tag
	tagName := acceptance.RandomTestName()
	resourceConfig := fmt.Sprintf(`
resource "abrha_tag" "foo" {
  name = "%s"
}`, tagName)
	dataSourceConfig := `
data "abrha_tag" "foobar" {
  name = abrha_tag.foo.name
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
					testAccCheckDataSourceAbrhaTagExists("data.abrha_tag.foobar", &tag),
					resource.TestCheckResourceAttr(
						"data.abrha_tag.foobar", "name", tagName),
					resource.TestCheckResourceAttrSet(
						"data.abrha_tag.foobar", "total_resource_count"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_tag.foobar", "vms_count"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_tag.foobar", "images_count"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_tag.foobar", "volumes_count"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_tag.foobar", "volume_snapshots_count"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_tag.foobar", "databases_count"),
				),
			},
		},
	})
}

func testAccCheckDataSourceAbrhaTagExists(n string, tag *goApiAbrha.Tag) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No tag ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

		foundTag, _, err := client.Tags.Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundTag.Name != rs.Primary.ID {
			return fmt.Errorf("Tag not found")
		}

		*tag = *foundTag

		return nil
	}
}
