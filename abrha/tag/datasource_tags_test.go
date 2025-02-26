package tag_test

import (
	"fmt"
	"testing"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAbrhaTags_Basic(t *testing.T) {
	var tag goApiAbrha.Tag
	tagName := acceptance.RandomTestName()
	resourceConfig := fmt.Sprintf(`
resource "abrha_tag" "foo" {
  name = "%s"
}`, tagName)
	dataSourceConfig := `
data "abrha_tags" "foobar" {
  filter {
    key    = "name"
    values = [abrha_tag.foo.name]
  }
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
					testAccCheckDataSourceAbrhaTagExists("abrha_tag.foo", &tag),
					resource.TestCheckResourceAttr(
						"data.abrha_tags.foobar", "tags.0.name", tagName),
					resource.TestCheckResourceAttrSet(
						"data.abrha_tags.foobar", "tags.0.total_resource_count"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_tags.foobar", "tags.0.vms_count"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_tags.foobar", "tags.0.images_count"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_tags.foobar", "tags.0.volumes_count"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_tags.foobar", "tags.0.volume_snapshots_count"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_tags.foobar", "tags.0.databases_count"),
				),
			},
		},
	})
}
