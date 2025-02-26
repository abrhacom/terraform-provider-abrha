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

func TestAccAbrhaVolumeSnapshot_Basic(t *testing.T) {
	var snapshot goApiAbrha.Snapshot
	volName := acceptance.RandomTestName("volume")
	snapName := acceptance.RandomTestName("snapshot")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaVolumeSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAbrhaVolumeSnapshotConfig_basic, volName, snapName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaVolumeSnapshotExists("abrha_volume_snapshot.foobar", &snapshot),
					resource.TestCheckResourceAttr(
						"abrha_volume_snapshot.foobar", "name", snapName),
					resource.TestCheckResourceAttr(
						"abrha_volume_snapshot.foobar", "size", "0"),
					resource.TestCheckResourceAttr(
						"abrha_volume_snapshot.foobar", "regions.#", "1"),
					resource.TestCheckResourceAttr(
						"abrha_volume_snapshot.foobar", "min_disk_size", "100"),
					resource.TestCheckResourceAttr(
						"abrha_volume_snapshot.foobar", "tags.#", "2"),
					resource.TestCheckResourceAttrSet(
						"abrha_volume_snapshot.foobar", "volume_id"),
				),
			},
		},
	})
}

func testAccCheckAbrhaVolumeSnapshotExists(n string, snapshot *goApiAbrha.Snapshot) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Volume Snapshot ID is set")
		}

		foundSnapshot, _, err := client.Snapshots.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}

		if foundSnapshot.ID != rs.Primary.ID {
			return fmt.Errorf("Volume Snapshot not found")
		}

		*snapshot = *foundSnapshot

		return nil
	}
}

func testAccCheckAbrhaVolumeSnapshotDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "abrha_volume_snapshot" {
			continue
		}

		// Try to find the snapshot
		_, _, err := client.Snapshots.Get(context.Background(), rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Volume Snapshot still exists")
		}
	}

	return nil
}

const testAccCheckAbrhaVolumeSnapshotConfig_basic = `
resource "abrha_volume" "foo" {
  region      = "nyc1"
  name        = "%s"
  size        = 100
  description = "peace makes plenty"
}

resource "abrha_volume_snapshot" "foobar" {
  name      = "%s"
  volume_id = abrha_volume.foo.id
  tags      = ["foo", "bar"]
}`

func TestAccAbrhaVolumeSnapshot_UpdateTags(t *testing.T) {
	var snapshot goApiAbrha.Snapshot
	volName := acceptance.RandomTestName("volume")
	snapName := acceptance.RandomTestName("snapshot")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaVolumeSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAbrhaVolumeSnapshotConfig_basic, volName, snapName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaVolumeSnapshotExists("abrha_volume_snapshot.foobar", &snapshot),
					resource.TestCheckResourceAttr("abrha_volume_snapshot.foobar", "tags.#", "2"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckAbrhaVolumeSnapshotConfig_basic_tag_update, volName, snapName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaVolumeSnapshotExists("abrha_volume_snapshot.foobar", &snapshot),
					resource.TestCheckResourceAttr("abrha_volume_snapshot.foobar", "tags.#", "3"),
				),
			},
		},
	})
}

const testAccCheckAbrhaVolumeSnapshotConfig_basic_tag_update = `
resource "abrha_volume" "foo" {
  region      = "nyc1"
  name        = "%s"
  size        = 100
  description = "peace makes plenty"
}

resource "abrha_volume_snapshot" "foobar" {
  name      = "%s"
  volume_id = abrha_volume.foo.id
  tags      = ["foo", "bar", "baz"]
}`
