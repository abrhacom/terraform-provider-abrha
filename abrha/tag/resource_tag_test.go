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

func TestAccAbrhaTag_Basic(t *testing.T) {
	var tag goApiAbrha.Tag

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaTagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaTagConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaTagExists("abrha_tag.foobar", &tag),
					testAccCheckAbrhaTagAttributes(&tag),
					resource.TestCheckResourceAttr(
						"abrha_tag.foobar", "name", "foobar"),
					resource.TestCheckResourceAttrSet(
						"abrha_tag.foobar", "total_resource_count"),
					resource.TestCheckResourceAttrSet(
						"abrha_tag.foobar", "vms_count"),
					resource.TestCheckResourceAttrSet(
						"abrha_tag.foobar", "images_count"),
					resource.TestCheckResourceAttrSet(
						"abrha_tag.foobar", "volumes_count"),
					resource.TestCheckResourceAttrSet(
						"abrha_tag.foobar", "volume_snapshots_count"),
					resource.TestCheckResourceAttrSet(
						"abrha_tag.foobar", "databases_count"),
				),
			},
		},
	})
}

func testAccCheckAbrhaTagDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "abrha_tag" {
			continue
		}

		// Try to find the key
		_, _, err := client.Tags.Get(context.Background(), rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Tag still exists")
		}
	}

	return nil
}

func testAccCheckAbrhaTagAttributes(tag *goApiAbrha.Tag) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if tag.Name != "foobar" {
			return fmt.Errorf("Bad name: %s", tag.Name)
		}

		return nil
	}
}

func testAccCheckAbrhaTagExists(n string, tag *goApiAbrha.Tag) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

		// Try to find the tag
		foundTag, _, err := client.Tags.Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}

		*tag = *foundTag

		return nil
	}
}

var testAccCheckAbrhaTagConfig_basic = `
resource "abrha_tag" "foobar" {
  name = "foobar"
}`
