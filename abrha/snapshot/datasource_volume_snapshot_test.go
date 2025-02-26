package snapshot_test

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

func TestAccDataSourceAbrhaVolumeSnapshot_basic(t *testing.T) {
	var snapshot goApiAbrha.Snapshot
	volName := acceptance.RandomTestName("volume")
	snapName := acceptance.RandomTestName("snapshot")
	resourceConfig := fmt.Sprintf(testAccCheckDataSourceAbrhaVolumeSnapshot_basic, volName, snapName)
	dataSourceConfig := `
data "abrha_volume_snapshot" "foobar" {
  most_recent = true
  name        = abrha_volume_snapshot.foo.name
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
					testAccCheckDataSourceAbrhaVolumeSnapshotExists("data.abrha_volume_snapshot.foobar", &snapshot),
					resource.TestCheckResourceAttr("data.abrha_volume_snapshot.foobar", "name", snapName),
					resource.TestCheckResourceAttr("data.abrha_volume_snapshot.foobar", "size", "0"),
					resource.TestCheckResourceAttr("data.abrha_volume_snapshot.foobar", "min_disk_size", "100"),
					resource.TestCheckResourceAttr("data.abrha_volume_snapshot.foobar", "regions.#", "1"),
					resource.TestCheckResourceAttr("data.abrha_volume_snapshot.foobar", "tags.#", "2"),
					resource.TestCheckResourceAttrSet("data.abrha_volume_snapshot.foobar", "volume_id"),
				),
			},
		},
	})
}

func TestAccDataSourceAbrhaVolumeSnapshot_regex(t *testing.T) {
	var snapshot goApiAbrha.Snapshot
	volName := acceptance.RandomTestName("volume")
	snapName := acceptance.RandomTestName("snapshot")
	resourceConfig := fmt.Sprintf(testAccCheckDataSourceAbrhaVolumeSnapshot_basic, volName, snapName)
	dataSourceConfig := `
data "abrha_volume_snapshot" "foobar" {
  most_recent = true
  name_regex  = "^${abrha_volume_snapshot.foo.name}"
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
					testAccCheckDataSourceAbrhaVolumeSnapshotExists("data.abrha_volume_snapshot.foobar", &snapshot),
					resource.TestCheckResourceAttr("data.abrha_volume_snapshot.foobar", "name", snapName),
					resource.TestCheckResourceAttr("data.abrha_volume_snapshot.foobar", "size", "0"),
					resource.TestCheckResourceAttr("data.abrha_volume_snapshot.foobar", "min_disk_size", "100"),
					resource.TestCheckResourceAttr("data.abrha_volume_snapshot.foobar", "regions.#", "1"),
					resource.TestCheckResourceAttr("data.abrha_volume_snapshot.foobar", "tags.#", "2"),
					resource.TestCheckResourceAttrSet("data.abrha_volume_snapshot.foobar", "volume_id"),
				),
			},
		},
	})
}

func TestAccDataSourceAbrhaVolumeSnapshot_region(t *testing.T) {
	var snapshot goApiAbrha.Snapshot
	vmName := acceptance.RandomTestName()
	snapName := acceptance.RandomTestName()
	nycResourceConfig := fmt.Sprintf(testAccCheckDataSourceAbrhaVolumeSnapshot_basic, vmName, snapName)
	lonResourceConfig := fmt.Sprintf(`
resource "abrha_volume" "bar" {
  region      = "lon1"
  name        = "%s"
  size        = 100
  description = "peace makes plenty"
}

resource "abrha_volume_snapshot" "bar" {
  name      = "%s"
  volume_id = abrha_volume.bar.id
}`, vmName, snapName)
	dataSourceConfig := `
data "abrha_volume_snapshot" "foobar" {
  name   = abrha_volume_snapshot.bar.name
  region = "lon1"
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: nycResourceConfig + lonResourceConfig,
			},
			{
				Config: nycResourceConfig + lonResourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceAbrhaVolumeSnapshotExists("data.abrha_volume_snapshot.foobar", &snapshot),
					resource.TestCheckResourceAttr("data.abrha_volume_snapshot.foobar", "name", snapName),
					resource.TestCheckResourceAttr("data.abrha_volume_snapshot.foobar", "size", "0"),
					resource.TestCheckResourceAttr("data.abrha_volume_snapshot.foobar", "min_disk_size", "100"),
					resource.TestCheckResourceAttr("data.abrha_volume_snapshot.foobar", "regions.#", "1"),
					resource.TestCheckResourceAttr("data.abrha_volume_snapshot.foobar", "tags.#", "0"),
					resource.TestCheckResourceAttrSet("data.abrha_volume_snapshot.foobar", "volume_id"),
				),
			},
		},
	})
}

func testAccCheckDataSourceAbrhaVolumeSnapshotExists(n string, snapshot *goApiAbrha.Snapshot) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No snapshot ID is set")
		}

		foundSnapshot, _, err := client.Snapshots.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}

		if foundSnapshot.ID != rs.Primary.ID {
			return fmt.Errorf("Snapshot not found")
		}

		*snapshot = *foundSnapshot

		return nil
	}
}

const testAccCheckDataSourceAbrhaVolumeSnapshot_basic = `
resource "abrha_volume" "foo" {
  region      = "nyc1"
  name        = "%s"
  size        = 100
  description = "peace makes plenty"
}

resource "abrha_volume_snapshot" "foo" {
  name      = "%s"
  volume_id = abrha_volume.foo.id
  tags      = ["foo", "bar"]
}
`
