package image_test

import (
	"fmt"
	"regexp"
	"testing"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAbrhaImage_Basic(t *testing.T) {
	var vm goApiAbrha.Vm
	var snapshotsId []int
	snapName := acceptance.RandomTestName()
	vmName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckAbrhaVmDestroy,
		Steps: []resource.TestStep{
			// Creates a Vm and takes multiple snapshots of it.
			// One will have the suffix -1 and two will have -0
			{
				Config: acceptance.TestAccCheckAbrhaVmConfig_basic(vmName),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar", &vm),
					acceptance.TakeSnapshotsOfVm(snapName, &vm, &snapshotsId),
				),
			},
			// Find snapshot with suffix -1
			{
				Config: testAccCheckAbrhaImageConfig_basic(snapName, 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.abrha_image.foobar", "name", snapName+"-1"),
					resource.TestCheckResourceAttr(
						"data.abrha_image.foobar", "min_disk_size", "25"),
					resource.TestCheckResourceAttr(
						"data.abrha_image.foobar", "private", "true"),
					resource.TestCheckResourceAttr(
						"data.abrha_image.foobar", "type", "snapshot"),
				),
			},
			// Expected error with  suffix -0 as multiple exist
			{
				Config:      testAccCheckAbrhaImageConfig_basic(snapName, 0),
				ExpectError: regexp.MustCompile(`.*too many images found with name tf-acc-test-.*\ .found 2, expected 1.`),
			},
			{
				Config:      testAccCheckAbrhaImageConfig_nonexisting(snapName),
				Destroy:     false,
				ExpectError: regexp.MustCompile(`.*no image found with name tf-acc-test-.*-nonexisting`),
			},
			{
				Config: " ",
				Check: resource.ComposeTestCheckFunc(
					acceptance.DeleteVmSnapshots(&snapshotsId),
				),
			},
		},
	})
}

func TestAccAbrhaImage_PublicSlug(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckAbrhaVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaImageConfig_slug("ubuntu-22-04-x64"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.abrha_image.foobar", "slug", "ubuntu-22-04-x64"),
					resource.TestCheckResourceAttr(
						"data.abrha_image.foobar", "min_disk_size", "7"),
					resource.TestCheckResourceAttr(
						"data.abrha_image.foobar", "private", "false"),
					resource.TestCheckResourceAttr(
						"data.abrha_image.foobar", "type", "snapshot"),
					resource.TestCheckResourceAttr(
						"data.abrha_image.foobar", "distribution", "Ubuntu"),
				),
			},
		},
	})
}

func testAccCheckAbrhaImageConfig_basic(name string, sInt int) string {
	return fmt.Sprintf(`
data "abrha_image" "foobar" {
  name = "%s-%d"
}
`, name, sInt)
}

func testAccCheckAbrhaImageConfig_nonexisting(name string) string {
	return fmt.Sprintf(`
data "abrha_image" "foobar" {
  name = "%s-nonexisting"
}
`, name)
}

func testAccCheckAbrhaImageConfig_slug(slug string) string {
	return fmt.Sprintf(`
data "abrha_image" "foobar" {
  slug = "%s"
}
`, slug)
}
