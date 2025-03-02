package volume_test

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

func TestAccAbrhaVolumeAttachment_Basic(t *testing.T) {
	var (
		volume = goApiAbrha.Volume{Name: acceptance.RandomTestName()}
		dName  = acceptance.RandomTestName()
		vm     goApiAbrha.Vm
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaVolumeAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaVolumeAttachmentConfig_basic(dName, volume.Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaVolumeExists("abrha_volume.foobar", &volume),
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar", &vm),
					testAccCheckAbrhaVolumeAttachmentExists("abrha_volume_attachment.foobar"),
					resource.TestCheckResourceAttrSet(
						"abrha_volume_attachment.foobar", "id"),
					resource.TestCheckResourceAttrSet(
						"abrha_volume_attachment.foobar", "vm_id"),
					resource.TestCheckResourceAttrSet(
						"abrha_volume_attachment.foobar", "volume_id"),
				),
			},
		},
	})
}

func TestAccAbrhaVolumeAttachment_Update(t *testing.T) {
	var (
		firstVolume  = goApiAbrha.Volume{Name: acceptance.RandomTestName()}
		secondVolume = goApiAbrha.Volume{Name: acceptance.RandomTestName()}
		dName        = acceptance.RandomTestName()
		vm           goApiAbrha.Vm
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaVolumeAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaVolumeAttachmentConfig_basic(dName, firstVolume.Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaVolumeExists("abrha_volume.foobar", &firstVolume),
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar", &vm),
					testAccCheckAbrhaVolumeAttachmentExists("abrha_volume_attachment.foobar"),
					resource.TestCheckResourceAttrSet(
						"abrha_volume_attachment.foobar", "id"),
					resource.TestCheckResourceAttrSet(
						"abrha_volume_attachment.foobar", "vm_id"),
					resource.TestCheckResourceAttrSet(
						"abrha_volume_attachment.foobar", "volume_id"),
				),
			},
			{
				Config: testAccCheckAbrhaVolumeAttachmentConfig_basic(dName, secondVolume.Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaVolumeExists("abrha_volume.foobar", &secondVolume),
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar", &vm),
					testAccCheckAbrhaVolumeAttachmentExists("abrha_volume_attachment.foobar"),
					resource.TestCheckResourceAttrSet(
						"abrha_volume_attachment.foobar", "id"),
					resource.TestCheckResourceAttrSet(
						"abrha_volume_attachment.foobar", "vm_id"),
					resource.TestCheckResourceAttrSet(
						"abrha_volume_attachment.foobar", "volume_id"),
				),
			},
		},
	})
}

func TestAccAbrhaVolumeAttachment_UpdateToSecondVolume(t *testing.T) {
	var (
		firstVolume  = goApiAbrha.Volume{Name: acceptance.RandomTestName()}
		secondVolume = goApiAbrha.Volume{Name: acceptance.RandomTestName()}
		dName        = acceptance.RandomTestName()
		vm           goApiAbrha.Vm
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaVolumeAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaVolumeAttachmentConfig_multiple_volumes(dName, firstVolume.Name, secondVolume.Name, "foobar"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaVolumeExists("abrha_volume.foobar", &firstVolume),
					testAccCheckAbrhaVolumeExists("abrha_volume.foobar_second", &secondVolume),
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar", &vm),
					testAccCheckAbrhaVolumeAttachmentExists("abrha_volume_attachment.foobar"),
					resource.TestCheckResourceAttrSet(
						"abrha_volume_attachment.foobar", "id"),
					resource.TestCheckResourceAttrSet(
						"abrha_volume_attachment.foobar", "vm_id"),
					resource.TestCheckResourceAttrSet(
						"abrha_volume_attachment.foobar", "volume_id"),
				),
			},
			{
				Config: testAccCheckAbrhaVolumeAttachmentConfig_multiple_volumes(dName, firstVolume.Name, secondVolume.Name, "foobar_second"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaVolumeExists("abrha_volume.foobar", &firstVolume),
					testAccCheckAbrhaVolumeExists("abrha_volume.foobar_second", &secondVolume),
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar", &vm),
					testAccCheckAbrhaVolumeAttachmentExists("abrha_volume_attachment.foobar"),
					resource.TestCheckResourceAttrSet(
						"abrha_volume_attachment.foobar", "id"),
					resource.TestCheckResourceAttrSet(
						"abrha_volume_attachment.foobar", "vm_id"),
					resource.TestCheckResourceAttrSet(
						"abrha_volume_attachment.foobar", "volume_id"),
				),
			},
		},
	})
}

func TestAccAbrhaVolumeAttachment_Multiple(t *testing.T) {
	var (
		firstVolume  = goApiAbrha.Volume{Name: acceptance.RandomTestName()}
		secondVolume = goApiAbrha.Volume{Name: acceptance.RandomTestName()}
		dName        = acceptance.RandomTestName()
		vm           goApiAbrha.Vm
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaVolumeAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaVolumeAttachmentConfig_multiple(dName, firstVolume.Name, secondVolume.Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaVolumeExists("abrha_volume.foobar", &firstVolume),
					testAccCheckAbrhaVolumeExists("abrha_volume.barfoo", &secondVolume),
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar", &vm),
					testAccCheckAbrhaVolumeAttachmentExists("abrha_volume_attachment.foobar"),
					testAccCheckAbrhaVolumeAttachmentExists("abrha_volume_attachment.barfoo"),
					resource.TestCheckResourceAttrSet(
						"abrha_volume_attachment.foobar", "id"),
					resource.TestCheckResourceAttrSet(
						"abrha_volume_attachment.foobar", "vm_id"),
					resource.TestCheckResourceAttrSet(
						"abrha_volume_attachment.foobar", "volume_id"),
					resource.TestCheckResourceAttrSet(
						"abrha_volume_attachment.barfoo", "id"),
					resource.TestCheckResourceAttrSet(
						"abrha_volume_attachment.barfoo", "vm_id"),
					resource.TestCheckResourceAttrSet(
						"abrha_volume_attachment.barfoo", "volume_id"),
				),
			},
		},
	})
}

func testAccCheckAbrhaVolumeAttachmentExists(rn string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("not found: %s", rn)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no volume ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

		volumeId := rs.Primary.Attributes["volume_id"]
		vmId := rs.Primary.Attributes["vm_id"]

		got, _, err := client.Storage.GetVolume(context.Background(), volumeId)
		if err != nil {
			return err
		}

		if len(got.VmIDs) == 0 || got.VmIDs[0] != vmId {
			return fmt.Errorf("wrong volume attachment found for volume %s, got %q wanted %q", volumeId, got.VmIDs[0], vmId)
		}

		return nil
	}
}

func testAccCheckAbrhaVolumeAttachmentDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "abrha_volume_attachment" {
			continue
		}
	}

	return nil
}

func testAccCheckAbrhaVolumeAttachmentConfig_basic(dName, vName string) string {
	return fmt.Sprintf(`
resource "abrha_volume" "foobar" {
  region      = "nyc1"
  name        = "%s"
  size        = 5
  description = "peace makes plenty"
}

resource "abrha_vm" "foobar" {
  name   = "%s"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = "nyc1"
}

resource "abrha_volume_attachment" "foobar" {
  vm_id     = abrha_vm.foobar.id
  volume_id = abrha_volume.foobar.id
}`, vName, dName)
}

func testAccCheckAbrhaVolumeAttachmentConfig_multiple(dName, vName, vSecondName string) string {
	return fmt.Sprintf(`
resource "abrha_volume" "foobar" {
  region      = "nyc1"
  name        = "%s"
  size        = 5
  description = "peace makes plenty"
}

resource "abrha_volume" "barfoo" {
  region      = "nyc1"
  name        = "%s"
  size        = 5
  description = "peace makes plenty"
}

resource "abrha_vm" "foobar" {
  name   = "%s"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = "nyc1"
}

resource "abrha_volume_attachment" "foobar" {
  vm_id     = abrha_vm.foobar.id
  volume_id = abrha_volume.foobar.id
}

resource "abrha_volume_attachment" "barfoo" {
  vm_id     = abrha_vm.foobar.id
  volume_id = abrha_volume.barfoo.id
}`, vName, vSecondName, dName)
}

func testAccCheckAbrhaVolumeAttachmentConfig_multiple_volumes(dName, vName, vSecondName, activeVolume string) string {
	return fmt.Sprintf(`
resource "abrha_volume" "foobar" {
  region      = "nyc1"
  name        = "%s"
  size        = 5
  description = "peace makes plenty"
}

resource "abrha_volume" "foobar_second" {
  region      = "nyc1"
  name        = "%s"
  size        = 5
  description = "peace makes plenty"
}

resource "abrha_vm" "foobar" {
  name   = "%s"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = "nyc1"
}

resource "abrha_volume_attachment" "foobar" {
  vm_id     = abrha_vm.foobar.id
  volume_id = abrha_volume.%s.id
}`, vName, vSecondName, dName, activeVolume)
}
