package volume

import (
	"context"
	"fmt"
	"log"
	"strings"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/abrhacom/terraform-provider-abrha/abrha/sweep"
	"github.com/abrhacom/terraform-provider-abrha/abrha/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("abrha_volume", &resource.Sweeper{
		Name:         "abrha_volume",
		F:            testSweepVolumes,
		Dependencies: []string{"abrha_vm"},
	})
}

func testSweepVolumes(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	opt := &goApiAbrha.ListVolumeParams{
		ListOptions: &goApiAbrha.ListOptions{PerPage: 200},
	}
	volumes, _, err := client.Storage.ListVolumes(context.Background(), opt)
	if err != nil {
		return err
	}

	for _, v := range volumes {
		if strings.HasPrefix(v.Name, sweep.TestNamePrefix) {

			if len(v.VmIDs) > 0 {
				log.Printf("Detaching volume %v from Vm %v", v.ID, v.VmIDs[0])

				action, _, err := client.StorageActions.DetachByVmID(context.Background(), v.ID, v.VmIDs[0])
				if err != nil {
					return fmt.Errorf("Error resizing volume (%s): %s", v.ID, err)
				}

				if err = util.WaitForAction(client, action); err != nil {
					return fmt.Errorf(
						"Error waiting for volume (%s): %s", v.ID, err)
				}
			}

			log.Printf("Destroying Volume %s", v.Name)

			if _, err := client.Storage.DeleteVolume(context.Background(), v.ID); err != nil {
				return err
			}
		}
	}

	return nil
}
