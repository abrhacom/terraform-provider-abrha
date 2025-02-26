package vm_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/abrhacom/terraform-provider-abrha/abrha/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAbrhaVm_importBasic(t *testing.T) {
	resourceName := "abrha_vm.foobar"
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckAbrhaVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: acceptance.TestAccCheckAbrhaVmConfig_basic(name),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"ssh_keys", "user_data", "resize_disk", "graceful_shutdown"}, //we ignore these attributes as we do not set to state
			},
			// Test importing non-existent resource provides expected error.
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateId:     "123",
				ExpectError:       regexp.MustCompile(`The resource you were accessing could not be found.`),
			},
		},
	})
}

func TestAccAbrhaVm_ImportWithNoImageSlug(t *testing.T) {
	var (
		vm           goApiAbrha.Vm
		restoredVm   goApiAbrha.Vm
		snapshotID   = goApiAbrha.PtrTo(0)
		name         = acceptance.RandomTestName()
		restoredName = acceptance.RandomTestName("restored")
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckAbrhaVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: acceptance.TestAccCheckAbrhaVmConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar", &vm),
					takeVmSnapshot(t, name, &vm, snapshotID),
				),
			},
		},
	})

	importConfig := testAccCheckAbrhaVmConfig_fromSnapshot(t, restoredName, *snapshotID)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckAbrhaVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: importConfig,
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.from-snapshot", &restoredVm),
				),
			},
			{
				ResourceName:      "abrha_vm.from-snapshot",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"ssh_keys", "user_data", "resize_disk", "graceful_shutdown"}, //we ignore the ssh_keys, resize_disk and user_data as we do not set to state
			},
			{
				Config: " ",
				Check: resource.ComposeTestCheckFunc(
					acceptance.DeleteVmSnapshots(&[]int{*snapshotID}),
				),
			},
		},
	})
}

func takeVmSnapshot(t *testing.T, name string, vm *goApiAbrha.Vm, snapshotID *int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

		action, _, err := client.VmActions.Snapshot(context.Background(), (*vm).ID, name)
		if err != nil {
			return err
		}
		err = util.WaitForAction(client, action)
		if err != nil {
			return err
		}

		retrieveVm, _, err := client.Vms.Get(context.Background(), (*vm).ID)
		if err != nil {
			return err
		}

		*snapshotID = retrieveVm.SnapshotIDs[0]
		return nil
	}
}

func testAccCheckAbrhaVmConfig_fromSnapshot(t *testing.T, name string, snapshotID int) string {
	return fmt.Sprintf(`
resource "abrha_vm" "from-snapshot" {
  name   = "%s"
  size   = "%s"
  image  = "%d"
  region = "nyc3"
}`, name, defaultSize, snapshotID)
}
