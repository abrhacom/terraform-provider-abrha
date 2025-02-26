package spaces_test

import (
	"fmt"
	"testing"

	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceAbrhaSpacesBucketObjects_basic(t *testing.T) {
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories:         acceptance.TestAccProviderFactories,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAbrhaSpacesObjectsConfigResources(name), // NOTE: contains no data source
				// Does not need Check
			},
			{
				Config: testAccDataSourceAbrhaSpacesObjectsConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaSpacesObjectsDataSourceExists("data.abrha_spaces_bucket_objects.yesh"),
					resource.TestCheckResourceAttr("data.abrha_spaces_bucket_objects.yesh", "keys.#", "2"),
					resource.TestCheckResourceAttr("data.abrha_spaces_bucket_objects.yesh", "keys.0", "arch/navajo/north_window"),
					resource.TestCheckResourceAttr("data.abrha_spaces_bucket_objects.yesh", "keys.1", "arch/navajo/sand_dune"),
				),
			},
		},
	})
}

func TestAccDataSourceAbrhaSpacesBucketObjects_all(t *testing.T) {
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories:         acceptance.TestAccProviderFactories,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAbrhaSpacesObjectsConfigResources(name), // NOTE: contains no data source
				// Does not need Check
			},
			{
				Config: testAccDataSourceAbrhaSpacesObjectsConfigAll(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaSpacesObjectsDataSourceExists("data.abrha_spaces_bucket_objects.yesh"),
					resource.TestCheckResourceAttr("data.abrha_spaces_bucket_objects.yesh", "keys.#", "7"),
					resource.TestCheckResourceAttr("data.abrha_spaces_bucket_objects.yesh", "keys.0", "arch/courthouse_towers/landscape"),
					resource.TestCheckResourceAttr("data.abrha_spaces_bucket_objects.yesh", "keys.1", "arch/navajo/north_window"),
					resource.TestCheckResourceAttr("data.abrha_spaces_bucket_objects.yesh", "keys.2", "arch/navajo/sand_dune"),
					resource.TestCheckResourceAttr("data.abrha_spaces_bucket_objects.yesh", "keys.3", "arch/partition/park_avenue"),
					resource.TestCheckResourceAttr("data.abrha_spaces_bucket_objects.yesh", "keys.4", "arch/rubicon"),
					resource.TestCheckResourceAttr("data.abrha_spaces_bucket_objects.yesh", "keys.5", "arch/three_gossips/broken"),
					resource.TestCheckResourceAttr("data.abrha_spaces_bucket_objects.yesh", "keys.6", "arch/three_gossips/turret"),
				),
			},
		},
	})
}

func TestAccDataSourceAbrhaSpacesBucketObjects_prefixes(t *testing.T) {
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories:         acceptance.TestAccProviderFactories,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAbrhaSpacesObjectsConfigResources(name), // NOTE: contains no data source
				// Does not need Check
			},
			{
				Config: testAccDataSourceAbrhaSpacesObjectsConfigPrefixes(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaSpacesObjectsDataSourceExists("data.abrha_spaces_bucket_objects.yesh"),
					resource.TestCheckResourceAttr("data.abrha_spaces_bucket_objects.yesh", "keys.#", "1"),
					resource.TestCheckResourceAttr("data.abrha_spaces_bucket_objects.yesh", "keys.0", "arch/rubicon"),
					resource.TestCheckResourceAttr("data.abrha_spaces_bucket_objects.yesh", "common_prefixes.#", "4"),
					resource.TestCheckResourceAttr("data.abrha_spaces_bucket_objects.yesh", "common_prefixes.0", "arch/courthouse_towers/"),
					resource.TestCheckResourceAttr("data.abrha_spaces_bucket_objects.yesh", "common_prefixes.1", "arch/navajo/"),
					resource.TestCheckResourceAttr("data.abrha_spaces_bucket_objects.yesh", "common_prefixes.2", "arch/partition/"),
					resource.TestCheckResourceAttr("data.abrha_spaces_bucket_objects.yesh", "common_prefixes.3", "arch/three_gossips/"),
				),
			},
		},
	})
}

func TestAccDataSourceAbrhaSpacesBucketObjects_encoded(t *testing.T) {
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories:         acceptance.TestAccProviderFactories,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAbrhaSpacesObjectsConfigExtraResource(name), // NOTE: contains no data source
				// Does not need Check
			},
			{
				Config: testAccDataSourceAbrhaSpacesObjectsConfigEncoded(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaSpacesObjectsDataSourceExists("data.abrha_spaces_bucket_objects.yesh"),
					resource.TestCheckResourceAttr("data.abrha_spaces_bucket_objects.yesh", "keys.#", "2"),
					resource.TestCheckResourceAttr("data.abrha_spaces_bucket_objects.yesh", "keys.0", "arch%2Fru%20b%20ic%20on"),
					resource.TestCheckResourceAttr("data.abrha_spaces_bucket_objects.yesh", "keys.1", "arch%2Frubicon"),
				),
			},
		},
	})
}

func TestAccDataSourceAbrhaSpacesBucketObjects_maxKeys(t *testing.T) {
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories:         acceptance.TestAccProviderFactories,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAbrhaSpacesObjectsConfigResources(name), // NOTE: contains no data source
				// Does not need Check
			},
			{
				Config: testAccDataSourceAbrhaSpacesObjectsConfigMaxKeys(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaSpacesObjectsDataSourceExists("data.abrha_spaces_bucket_objects.yesh"),
					resource.TestCheckResourceAttr("data.abrha_spaces_bucket_objects.yesh", "keys.#", "2"),
					resource.TestCheckResourceAttr("data.abrha_spaces_bucket_objects.yesh", "keys.0", "arch/courthouse_towers/landscape"),
					resource.TestCheckResourceAttr("data.abrha_spaces_bucket_objects.yesh", "keys.1", "arch/navajo/north_window"),
				),
			},
		},
	})
}

func testAccCheckAbrhaSpacesObjectsDataSourceExists(addr string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[addr]
		if !ok {
			return fmt.Errorf("Can't find Spaces objects data source: %s", addr)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Spaces objects data source ID not set")
		}

		return nil
	}
}

func testAccDataSourceAbrhaSpacesObjectsConfigResources(name string) string {
	return fmt.Sprintf(`
resource "abrha_spaces_bucket" "objects_bucket" {
  name          = "%s"
  region        = "nyc3"
  force_destroy = true
}

resource "abrha_spaces_bucket_object" "object1" {
  bucket  = abrha_spaces_bucket.objects_bucket.name
  region  = abrha_spaces_bucket.objects_bucket.region
  key     = "arch/three_gossips/turret"
  content = "Delicate"
}

resource "abrha_spaces_bucket_object" "object2" {
  bucket  = abrha_spaces_bucket.objects_bucket.name
  region  = abrha_spaces_bucket.objects_bucket.region
  key     = "arch/three_gossips/broken"
  content = "Dark Angel"
}

resource "abrha_spaces_bucket_object" "object3" {
  bucket  = abrha_spaces_bucket.objects_bucket.name
  region  = abrha_spaces_bucket.objects_bucket.region
  key     = "arch/navajo/north_window"
  content = "Balanced Rock"
}

resource "abrha_spaces_bucket_object" "object4" {
  bucket  = abrha_spaces_bucket.objects_bucket.name
  region  = abrha_spaces_bucket.objects_bucket.region
  key     = "arch/navajo/sand_dune"
  content = "Queen Victoria Rock"
}

resource "abrha_spaces_bucket_object" "object5" {
  bucket  = abrha_spaces_bucket.objects_bucket.name
  region  = abrha_spaces_bucket.objects_bucket.region
  key     = "arch/partition/park_avenue"
  content = "Double-O"
}

resource "abrha_spaces_bucket_object" "object6" {
  bucket  = abrha_spaces_bucket.objects_bucket.name
  region  = abrha_spaces_bucket.objects_bucket.region
  key     = "arch/courthouse_towers/landscape"
  content = "Fiery Furnace"
}

resource "abrha_spaces_bucket_object" "object7" {
  bucket  = abrha_spaces_bucket.objects_bucket.name
  region  = abrha_spaces_bucket.objects_bucket.region
  key     = "arch/rubicon"
  content = "Devils Garden"
}
`, name)
}

func testAccDataSourceAbrhaSpacesObjectsConfigBasic(name string) string {
	return fmt.Sprintf(`
%s

data "abrha_spaces_bucket_objects" "yesh" {
  bucket    = abrha_spaces_bucket.objects_bucket.name
  region    = abrha_spaces_bucket.objects_bucket.region
  prefix    = "arch/navajo/"
  delimiter = "/"
}
`, testAccDataSourceAbrhaSpacesObjectsConfigResources(name))
}

func testAccDataSourceAbrhaSpacesObjectsConfigAll(name string) string {
	return fmt.Sprintf(`
%s

data "abrha_spaces_bucket_objects" "yesh" {
  bucket = abrha_spaces_bucket.objects_bucket.name
  region = abrha_spaces_bucket.objects_bucket.region
}
`, testAccDataSourceAbrhaSpacesObjectsConfigResources(name))
}

func testAccDataSourceAbrhaSpacesObjectsConfigPrefixes(name string) string {
	return fmt.Sprintf(`
%s

data "abrha_spaces_bucket_objects" "yesh" {
  bucket    = abrha_spaces_bucket.objects_bucket.name
  region    = abrha_spaces_bucket.objects_bucket.region
  prefix    = "arch/"
  delimiter = "/"
}
`, testAccDataSourceAbrhaSpacesObjectsConfigResources(name))
}

func testAccDataSourceAbrhaSpacesObjectsConfigExtraResource(name string) string {
	return fmt.Sprintf(`
%s

resource "abrha_spaces_bucket_object" "object8" {
  bucket  = abrha_spaces_bucket.objects_bucket.name
  region  = abrha_spaces_bucket.objects_bucket.region
  key     = "arch/ru b ic on"
  content = "Goose Island"
}
`, testAccDataSourceAbrhaSpacesObjectsConfigResources(name))
}

func testAccDataSourceAbrhaSpacesObjectsConfigEncoded(name string) string {
	return fmt.Sprintf(`
%s

data "abrha_spaces_bucket_objects" "yesh" {
  bucket        = abrha_spaces_bucket.objects_bucket.name
  region        = abrha_spaces_bucket.objects_bucket.region
  encoding_type = "url"
  prefix        = "arch/ru"
}
`, testAccDataSourceAbrhaSpacesObjectsConfigExtraResource(name))
}

func testAccDataSourceAbrhaSpacesObjectsConfigMaxKeys(name string) string {
	return fmt.Sprintf(`
%s

data "abrha_spaces_bucket_objects" "yesh" {
  bucket   = abrha_spaces_bucket.objects_bucket.name
  region   = abrha_spaces_bucket.objects_bucket.region
  max_keys = 2
}
`, testAccDataSourceAbrhaSpacesObjectsConfigResources(name))
}
