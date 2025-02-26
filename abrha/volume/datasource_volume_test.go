package volume_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceAbrhaVolume_Basic(t *testing.T) {
	var volume goApiAbrha.Volume
	testName := acceptance.RandomTestName("volume")
	resourceConfig := testAccCheckDataSourceAbrhaVolumeConfig_basic(testName)
	dataSourceConfig := `
data "abrha_volume" "foobar" {
  name = abrha_volume.foo.name
}`

	expectedURNRegEx, _ := regexp.Compile(`do:volume:[0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}`)

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
					testAccCheckDataSourceAbrhaVolumeExists("data.abrha_volume.foobar", &volume),
					resource.TestCheckResourceAttr(
						"data.abrha_volume.foobar", "name", testName),
					resource.TestCheckResourceAttr(
						"data.abrha_volume.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"data.abrha_volume.foobar", "size", "10"),
					resource.TestCheckResourceAttr(
						"data.abrha_volume.foobar", "vm_ids.#", "0"),
					resource.TestCheckResourceAttr(
						"data.abrha_volume.foobar", "tags.#", "2"),
					resource.TestMatchResourceAttr("data.abrha_volume.foobar", "urn", expectedURNRegEx),
				),
			},
		},
	})
}

func TestAccDataSourceAbrhaVolume_RegionScoped(t *testing.T) {
	var volume goApiAbrha.Volume
	name := acceptance.RandomTestName()
	resourceConfig := testAccCheckDataSourceAbrhaVolumeConfig_region_scoped(name)
	dataSourceConfig := `
data "abrha_volume" "foobar" {
  name   = abrha_volume.foo.name
  region = "lon1"
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
					testAccCheckDataSourceAbrhaVolumeExists("data.abrha_volume.foobar", &volume),
					resource.TestCheckResourceAttr(
						"data.abrha_volume.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"data.abrha_volume.foobar", "region", "lon1"),
					resource.TestCheckResourceAttr(
						"data.abrha_volume.foobar", "size", "20"),
					resource.TestCheckResourceAttr(
						"data.abrha_volume.foobar", "vm_ids.#", "0"),
					resource.TestCheckResourceAttr(
						"data.abrha_volume.foobar", "tags.#", "0"),
				),
			},
		},
	})
}

func testAccCheckDataSourceAbrhaVolumeExists(n string, volume *goApiAbrha.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Volume ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

		foundVolume, _, err := client.Storage.GetVolume(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundVolume.ID != rs.Primary.ID {
			return fmt.Errorf("Volume not found")
		}

		*volume = *foundVolume

		return nil
	}
}

func testAccCheckDataSourceAbrhaVolumeConfig_basic(testName string) string {
	return fmt.Sprintf(`
resource "abrha_volume" "foo" {
  region = "nyc3"
  name   = "%s"
  size   = 10
  tags   = ["foo", "bar"]
}`, testName)
}

func testAccCheckDataSourceAbrhaVolumeConfig_region_scoped(name string) string {
	return fmt.Sprintf(`
resource "abrha_volume" "foo" {
  region = "nyc3"
  name   = "%s"
  size   = 10
  tags   = ["foo", "bar"]
}

resource "abrha_volume" "bar" {
  region = "lon1"
  name   = "%s"
  size   = 20
}`, name, name)
}
