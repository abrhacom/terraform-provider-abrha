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

func TestAccDataSourceAbrhaVmSnapshot_basic(t *testing.T) {
	var snapshot goApiAbrha.Snapshot
	vmName := acceptance.RandomTestName()
	snapName := acceptance.RandomTestName()
	resourceConfig := fmt.Sprintf(testAccCheckDataSourceAbrhaVmSnapshot_basic, vmName, snapName)
	dataSourceConfig := `
data "abrha_vm_snapshot" "foobar" {
  most_recent = true
  name        = abrha_vm_snapshot.foo.name
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
					testAccCheckDataSourceAbrhaVmSnapshotExists("data.abrha_vm_snapshot.foobar", &snapshot),
					resource.TestCheckResourceAttr("data.abrha_vm_snapshot.foobar", "name", snapName),
					resource.TestCheckResourceAttr("data.abrha_vm_snapshot.foobar", "min_disk_size", "25"),
					resource.TestCheckResourceAttr("data.abrha_vm_snapshot.foobar", "regions.#", "1"),
					resource.TestCheckResourceAttrSet("data.abrha_vm_snapshot.foobar", "vm_id"),
				),
			},
		},
	})
}

func TestAccDataSourceAbrhaVmSnapshot_regex(t *testing.T) {
	var snapshot goApiAbrha.Snapshot
	vmName := acceptance.RandomTestName()
	snapName := acceptance.RandomTestName()
	resourceConfig := fmt.Sprintf(testAccCheckDataSourceAbrhaVmSnapshot_basic, vmName, snapName)
	dataSourceConfig := fmt.Sprintf(`
data "abrha_vm_snapshot" "foobar" {
  most_recent = true
  name_regex  = "^%s"
}`, snapName)

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
					testAccCheckDataSourceAbrhaVmSnapshotExists("data.abrha_vm_snapshot.foobar", &snapshot),
					resource.TestCheckResourceAttr("data.abrha_vm_snapshot.foobar", "name", snapName),
					resource.TestCheckResourceAttr("data.abrha_vm_snapshot.foobar", "min_disk_size", "25"),
					resource.TestCheckResourceAttr("data.abrha_vm_snapshot.foobar", "regions.#", "1"),
					resource.TestCheckResourceAttrSet("data.abrha_vm_snapshot.foobar", "vm_id"),
				),
			},
		},
	})
}

func TestAccDataSourceAbrhaVmSnapshot_region(t *testing.T) {
	var snapshot goApiAbrha.Snapshot
	vmName := acceptance.RandomTestName()
	snapName := acceptance.RandomTestName()
	nycResourceConfig := fmt.Sprintf(testAccCheckDataSourceAbrhaVmSnapshot_basic, vmName, snapName)
	lonResourceConfig := fmt.Sprintf(`
resource "abrha_vm" "bar" {
  region = "lon1"
  name   = "%s"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
}

resource "abrha_vm_snapshot" "bar" {
  name  = "%s"
  vm_id = abrha_vm.bar.id
}`, vmName, snapName)
	dataSourceConfig := `
data "abrha_vm_snapshot" "foobar" {
  name   = abrha_vm_snapshot.bar.name
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
					testAccCheckDataSourceAbrhaVmSnapshotExists("data.abrha_vm_snapshot.foobar", &snapshot),
					resource.TestCheckResourceAttr("data.abrha_vm_snapshot.foobar", "name", snapName),
					resource.TestCheckResourceAttr("data.abrha_vm_snapshot.foobar", "min_disk_size", "25"),
					resource.TestCheckResourceAttr("data.abrha_vm_snapshot.foobar", "regions.#", "1"),
					resource.TestCheckResourceAttrSet("data.abrha_vm_snapshot.foobar", "vm_id"),
				),
			},
		},
	})
}

func testAccCheckDataSourceAbrhaVmSnapshotExists(n string, snapshot *goApiAbrha.Snapshot) resource.TestCheckFunc {
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

const testAccCheckDataSourceAbrhaVmSnapshot_basic = `
resource "abrha_vm" "foo" {
  name   = "%s"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = "nyc3"
}

resource "abrha_vm_snapshot" "foo" {
  name  = "%s"
  vm_id = abrha_vm.foo.id
}
`
