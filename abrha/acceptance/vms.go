package acceptance

import (
	"context"
	"fmt"
	"log"
	"strings"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/abrhacom/terraform-provider-abrha/abrha/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCheckAbrhaVmDestroy(s *terraform.State) error {
	client := TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "abrha_vm" {
			continue
		}

		id := rs.Primary.ID

		// Try to find the Vm
		_, _, err := client.Vms.Get(context.Background(), id)

		// Wait

		if err != nil && !strings.Contains(err.Error(), "404") {
			return fmt.Errorf(
				"Error waiting for vm (%s) to be destroyed: %s",
				rs.Primary.ID, err)
		}
	}

	return nil
}

func TestAccCheckAbrhaVmExists(n string, vm *goApiAbrha.Vm) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Vm ID is set")
		}

		client := TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

		id := rs.Primary.ID

		// Try to find the Vm
		retrieveVm, _, err := client.Vms.Get(context.Background(), id)

		if err != nil {
			return err
		}

		if retrieveVm.ID != rs.Primary.ID {
			return fmt.Errorf("Vm not found")
		}

		*vm = *retrieveVm

		return nil
	}
}

func TestAccCheckAbrhaVmConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "abrha_vm" "foobar" {
  name      = "%s"
  size      = "s-1vcpu-1gb"
  image     = "ubuntu-22-04-x64"
  region    = "nyc3"
  user_data = "foobar"
}`, name)
}

// TakeSnapshotsOfVm takes three snapshots of the given Vm. One will have the suffix -1 and two will have -0.
func TakeSnapshotsOfVm(snapName string, vm *goApiAbrha.Vm, snapshotsIDs *[]int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()
		for i := 0; i < 3; i++ {
			err := takeSnapshotOfVm(snapName, i%2, vm)
			if err != nil {
				return err
			}
		}
		retrieveVm, _, err := client.Vms.Get(context.Background(), (*vm).ID)
		if err != nil {
			return err
		}
		*snapshotsIDs = retrieveVm.SnapshotIDs
		return nil
	}
}

func takeSnapshotOfVm(snapName string, intSuffix int, vm *goApiAbrha.Vm) error {
	client := TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()
	action, _, err := client.VmActions.Snapshot(context.Background(), (*vm).ID, fmt.Sprintf("%s-%d", snapName, intSuffix))
	if err != nil {
		return err
	}
	err = util.WaitForAction(client, action)
	if err != nil {
		return err
	}
	return nil
}

func DeleteVmSnapshots(snapshotsId *[]int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		log.Printf("Deleting Vm snapshots")

		client := TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

		snapshots := *snapshotsId
		for _, value := range snapshots {
			log.Printf("Deleting %d", value)
			_, err := client.Images.Delete(context.Background(), value)
			if err != nil {
				return err
			}
		}
		return nil
	}
}
