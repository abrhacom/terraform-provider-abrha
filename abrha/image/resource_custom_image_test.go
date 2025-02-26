package image_test

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/abrhacom/terraform-provider-abrha/abrha/image"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAbrhaCustomImageFull(t *testing.T) {
	rString := acceptance.RandomTestName()
	name := fmt.Sprintf("abrha_custom_image.%s", rString)
	regions := `["nyc3"]`
	updatedString := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaCustomImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaCustomImageConfig(rString, rString, regions, "Unknown OS"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", rString),
					resource.TestCheckResourceAttr(name, "description", fmt.Sprintf("%s-description", rString)),
					resource.TestCheckResourceAttr(name, "distribution", "Unknown OS"),
					resource.TestCheckResourceAttr(name, "public", "false"),
					resource.TestCheckResourceAttr(name, "regions.0", "nyc3"),
					resource.TestCheckResourceAttr(name, "status", "available"),
					resource.TestCheckResourceAttr(name, "tags.0", "flatcar"),
					resource.TestCheckResourceAttr(name, "type", "custom"),
					resource.TestCheckResourceAttr(name, "slug", ""),
					resource.TestCheckResourceAttrSet(name, "created_at"),
					resource.TestCheckResourceAttrSet(name, "image_id"),
					resource.TestCheckResourceAttrSet(name, "min_disk_size"),
					resource.TestCheckResourceAttrSet(name, "size_gigabytes"),
				),
			},
			{
				Config: testAccCheckAbrhaCustomImageConfig(rString, updatedString, regions, "CoreOS"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", updatedString),
					resource.TestCheckResourceAttr(name, "description", fmt.Sprintf("%s-description", updatedString)),
					resource.TestCheckResourceAttr(name, "distribution", "CoreOS"),
				),
			},
		},
	})
}

func TestAccAbrhaCustomImageMultiRegion(t *testing.T) {
	rString := acceptance.RandomTestName()
	name := fmt.Sprintf("abrha_custom_image.%s", rString)
	regions := `["nyc3", "nyc2"]`
	regionsUpdated := `["nyc3", "nyc2", "tor1"]`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaCustomImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaCustomImageConfig(rString, rString, regions, "Unknown OS"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", rString),
					resource.TestCheckTypeSetElemAttr(name, "regions.*", "nyc2"),
					resource.TestCheckTypeSetElemAttr(name, "regions.*", "nyc3"),
				),
			},
			{
				Config: testAccCheckAbrhaCustomImageConfig(rString, rString, regionsUpdated, "Unknown OS"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", rString),
					resource.TestCheckTypeSetElemAttr(name, "regions.*", "nyc2"),
					resource.TestCheckTypeSetElemAttr(name, "regions.*", "nyc3"),
					resource.TestCheckTypeSetElemAttr(name, "regions.*", "tor1"),
				),
			},
			{
				Config: testAccCheckAbrhaCustomImageConfig(rString, rString, regions, "Unknown OS"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", rString),
					resource.TestCheckTypeSetElemAttr(name, "regions.*", "nyc2"),
					resource.TestCheckTypeSetElemAttr(name, "regions.*", "nyc3"),
				),
			},
		},
	})
}

func testAccCheckAbrhaCustomImageConfig(rName string, name string, regions string, distro string) string {
	return fmt.Sprintf(`
resource "abrha_custom_image" "%s" {
  name         = "%s"
  url          = "https://stable.release.flatcar-linux.net/amd64-usr/2605.7.0/flatcar_production_abrha_image.bin.bz2"
  regions      = %s
  description  = "%s-description"
  distribution = "%s"
  tags = [
    "flatcar"
  ]
}
`, rName, name, regions, name, distro)
}

func testAccCheckAbrhaCustomImageDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "abrha_custom_image" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		// Try to find the Image by ID
		i, resp, err := client.Images.GetByID(context.Background(), id)
		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}

			return err
		}

		if i.Status != image.ImageDeletedStatus {
			return fmt.Errorf("Image %d not destroyed", id)
		}
	}

	return nil
}
