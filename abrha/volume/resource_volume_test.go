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

func TestAccAbrhaVolume_Basic(t *testing.T) {
	name := acceptance.RandomTestName("volume")

	expectedURNRegEx, _ := regexp.Compile(`do:volume:[0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}`)

	volume := goApiAbrha.Volume{
		Name: name,
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAbrhaVolumeConfig_basic, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaVolumeExists("abrha_volume.foobar", &volume),
					resource.TestCheckResourceAttr(
						"abrha_volume.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"abrha_volume.foobar", "size", "100"),
					resource.TestCheckResourceAttr(
						"abrha_volume.foobar", "region", "nyc1"),
					resource.TestCheckResourceAttr(
						"abrha_volume.foobar", "description", "peace makes plenty"),
					resource.TestCheckResourceAttr(
						"abrha_volume.foobar", "tags.#", "2"),
					resource.TestMatchResourceAttr("abrha_volume.foobar", "urn", expectedURNRegEx),
				),
			},
		},
	})
}

const testAccCheckAbrhaVolumeConfig_basic = `
resource "abrha_volume" "foobar" {
  region      = "nyc1"
  name        = "%s"
  size        = 100
  description = "peace makes plenty"
  tags        = ["foo", "bar"]
}`

func testAccCheckAbrhaVolumeExists(rn string, volume *goApiAbrha.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("not found: %s", rn)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no volume ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

		got, _, err := client.Storage.GetVolume(context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}
		if got.Name != volume.Name {
			return fmt.Errorf("wrong volume found, want %q got %q", volume.Name, got.Name)
		}
		// get the computed volume details
		*volume = *got
		return nil
	}
}

func testAccCheckAbrhaVolumeDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "abrha_volume" {
			continue
		}

		// Try to find the volume
		_, _, err := client.Storage.GetVolume(context.Background(), rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Volume still exists")
		}
	}

	return nil
}

func TestAccAbrhaVolume_Vm(t *testing.T) {
	var (
		volume = goApiAbrha.Volume{Name: acceptance.RandomTestName()}
		dName  = acceptance.RandomTestName()
		vm     goApiAbrha.Vm
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaVolumeConfig_vm(dName, volume.Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaVolumeExists("abrha_volume.foobar", &volume),
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar", &vm),
					// the vm should see an attached volume
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "volume_ids.#", "1"),
				),
			},
		},
	})
}

func testAccCheckAbrhaVolumeConfig_vm(dName, vName string) string {
	return fmt.Sprintf(`
resource "abrha_volume" "foobar" {
  region      = "nyc1"
  name        = "%s"
  size        = 100
  description = "peace makes plenty"
}

resource "abrha_vm" "foobar" {
  name               = "%s"
  size               = "s-1vcpu-1gb"
  image              = "ubuntu-22-04-x64"
  region             = "nyc1"
  ipv6               = true
  private_networking = true
  volume_ids         = [abrha_volume.foobar.id]
}`, vName, dName)
}

func TestAccAbrhaVolume_LegacyFilesystemType(t *testing.T) {
	name := acceptance.RandomTestName()

	volume := goApiAbrha.Volume{
		Name: name,
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAbrhaVolumeConfig_legacy_filesystem_type, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaVolumeExists("abrha_volume.foobar", &volume),
					resource.TestCheckResourceAttr(
						"abrha_volume.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"abrha_volume.foobar", "size", "100"),
					resource.TestCheckResourceAttr(
						"abrha_volume.foobar", "region", "nyc1"),
					resource.TestCheckResourceAttr(
						"abrha_volume.foobar", "description", "peace makes plenty"),
					resource.TestCheckResourceAttr(
						"abrha_volume.foobar", "filesystem_type", "xfs"),
				),
			},
		},
	})
}

const testAccCheckAbrhaVolumeConfig_legacy_filesystem_type = `
resource "abrha_volume" "foobar" {
  region          = "nyc1"
  name            = "%s"
  size            = 100
  description     = "peace makes plenty"
  filesystem_type = "xfs"
}`

func TestAccAbrhaVolume_FilesystemType(t *testing.T) {
	name := acceptance.RandomTestName()

	volume := goApiAbrha.Volume{
		Name: name,
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAbrhaVolumeConfig_filesystem_type, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaVolumeExists("abrha_volume.foobar", &volume),
					resource.TestCheckResourceAttr(
						"abrha_volume.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"abrha_volume.foobar", "size", "100"),
					resource.TestCheckResourceAttr(
						"abrha_volume.foobar", "region", "nyc1"),
					resource.TestCheckResourceAttr(
						"abrha_volume.foobar", "description", "peace makes plenty"),
					resource.TestCheckResourceAttr(
						"abrha_volume.foobar", "initial_filesystem_type", "xfs"),
					resource.TestCheckResourceAttr(
						"abrha_volume.foobar", "initial_filesystem_label", "label"),
					resource.TestCheckResourceAttr(
						"abrha_volume.foobar", "filesystem_type", "xfs"),
					resource.TestCheckResourceAttr(
						"abrha_volume.foobar", "filesystem_label", "label"),
				),
			},
		},
	})
}

const testAccCheckAbrhaVolumeConfig_filesystem_type = `
resource "abrha_volume" "foobar" {
  region                   = "nyc1"
  name                     = "%s"
  size                     = 100
  description              = "peace makes plenty"
  initial_filesystem_type  = "xfs"
  initial_filesystem_label = "label"
}`

func TestAccAbrhaVolume_Resize(t *testing.T) {
	var (
		volume = goApiAbrha.Volume{Name: acceptance.RandomTestName()}
		dName  = acceptance.RandomTestName()
		vm     goApiAbrha.Vm
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaVolumeConfig_resize(dName, volume.Name, 20),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaVolumeExists("abrha_volume.foobar", &volume),
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar", &vm),
					// the vm should see an attached volume
					resource.TestCheckResourceAttr("abrha_vm.foobar", "volume_ids.#", "1"),
					resource.TestCheckResourceAttr("abrha_volume.foobar", "size", "20"),
				),
			},
			{
				Config: testAccCheckAbrhaVolumeConfig_resize(dName, volume.Name, 50),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaVolumeExists("abrha_volume.foobar", &volume),
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar", &vm),
					// the vm should see an attached volume
					resource.TestCheckResourceAttr("abrha_vm.foobar", "volume_ids.#", "1"),
					resource.TestCheckResourceAttr("abrha_volume.foobar", "size", "50"),
				),
			},
		},
	})
}

func testAccCheckAbrhaVolumeConfig_resize(dName, vName string, vSize int) string {
	return fmt.Sprintf(`
resource "abrha_volume" "foobar" {
  region      = "nyc1"
  name        = "%s"
  size        = %d
  description = "peace makes plenty"
}

resource "abrha_vm" "foobar" {
  name               = "%s"
  size               = "s-1vcpu-1gb"
  image              = "ubuntu-22-04-x64"
  region             = "nyc1"
  ipv6               = true
  private_networking = true
  volume_ids         = [abrha_volume.foobar.id]
}`, vName, vSize, dName)
}

func TestAccAbrhaVolume_CreateFromSnapshot(t *testing.T) {
	volName := acceptance.RandomTestName()
	snapName := acceptance.RandomTestName()
	restoredName := acceptance.RandomTestName()

	volume := goApiAbrha.Volume{
		Name: restoredName,
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaVolumeConfig_create_from_snapshot(volName, snapName, restoredName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaVolumeExists("abrha_volume.foobar", &volume),
					// the vm should see an attached volume
					resource.TestCheckResourceAttr("abrha_volume.foobar", "region", "nyc1"),
					resource.TestCheckResourceAttr("abrha_volume.foobar", "size", "100"),
				),
			},
		},
	})
}

func testAccCheckAbrhaVolumeConfig_create_from_snapshot(volume, snapshot, restored string) string {
	return fmt.Sprintf(`
resource "abrha_volume" "foo" {
  region      = "nyc1"
  name        = "%s"
  size        = 100
  description = "peace makes plenty"
}

resource "abrha_volume_snapshot" "foo" {
  name      = "%s"
  volume_id = abrha_volume.foo.id
}

resource "abrha_volume" "foobar" {
  region      = "nyc1"
  name        = "%s"
  size        = abrha_volume_snapshot.foo.min_disk_size
  snapshot_id = abrha_volume_snapshot.foo.id
}`, volume, snapshot, restored)
}

func TestAccAbrhaVolume_UpdateTags(t *testing.T) {
	name := acceptance.RandomTestName()

	volume := goApiAbrha.Volume{
		Name: name,
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAbrhaVolumeConfig_basic, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaVolumeExists("abrha_volume.foobar", &volume),
					resource.TestCheckResourceAttr(
						"abrha_volume.foobar", "tags.#", "2"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckAbrhaVolumeConfig_basic_tag_update, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaVolumeExists("abrha_volume.foobar", &volume),
					resource.TestCheckResourceAttr(
						"abrha_volume.foobar", "tags.#", "3"),
				),
			},
		},
	})
}

const testAccCheckAbrhaVolumeConfig_basic_tag_update = `
resource "abrha_volume" "foobar" {
  region      = "nyc1"
  name        = "%s"
  size        = 100
  description = "peace makes plenty"
  tags        = ["foo", "bar", "baz"]
}`

func TestAccAbrhaVolume_createWithRegionSlugUpperCase(t *testing.T) {
	name := acceptance.RandomTestName("volume")

	expectedURNRegEx, _ := regexp.Compile(`do:volume:[0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}`)

	volume := goApiAbrha.Volume{
		Name: name,
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAbrhaVolumeConfig_createWithRegionSlugUpperCase, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaVolumeExists("abrha_volume.foobar", &volume),
					resource.TestCheckResourceAttr(
						"abrha_volume.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"abrha_volume.foobar", "size", "100"),
					resource.TestCheckResourceAttr(
						"abrha_volume.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"abrha_volume.foobar", "description", "peace makes plenty"),
					resource.TestCheckResourceAttr(
						"abrha_volume.foobar", "tags.#", "2"),
					resource.TestMatchResourceAttr("abrha_volume.foobar", "urn", expectedURNRegEx),
				),
			},
		},
	})
}

const testAccCheckAbrhaVolumeConfig_createWithRegionSlugUpperCase = `
resource "abrha_volume" "foobar" {
  region      = "NYC3"
  name        = "%s"
  size        = 100
  description = "peace makes plenty"
  tags        = ["foo", "bar"]
}`
