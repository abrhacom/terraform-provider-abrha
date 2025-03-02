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

func TestAccAbrhaVmSnapshot_Basic(t *testing.T) {
	var snapshot goApiAbrha.Snapshot
	dName := acceptance.RandomTestName()
	snapName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaVmSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAbrhaVmSnapshotConfig_basic, dName, snapName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaVmSnapshotExists("abrha_vm_snapshot.foobar", &snapshot),
					resource.TestCheckResourceAttr(
						"abrha_vm_snapshot.foobar", "name", snapName),
				),
			},
		},
	})
}

func testAccCheckAbrhaVmSnapshotExists(n string, snapshot *goApiAbrha.Snapshot) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Vm Snapshot ID is set")
		}

		foundSnapshot, _, err := client.Snapshots.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}

		if foundSnapshot.ID != rs.Primary.ID {
			return fmt.Errorf("Vm Snapshot not found")
		}

		*snapshot = *foundSnapshot

		return nil
	}
}

func testAccCheckAbrhaVmSnapshotDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "abrha_vm_snapshot" {
			continue
		}

		// Try to find the snapshot
		_, _, err := client.Snapshots.Get(context.Background(), rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Vm Snapshot still exists")
		}
	}

	return nil
}

const testAccCheckAbrhaVmSnapshotConfig_basic = `
resource "abrha_vm" "foo" {
  name      = "%s"
  size      = "s-1vcpu-1gb"
  image     = "ubuntu-22-04-x64"
  region    = "nyc3"
  user_data = "foobar"
}

resource "abrha_vm_snapshot" "foobar" {
  vm_id = abrha_vm.foo.id
  name  = "%s"
}`
