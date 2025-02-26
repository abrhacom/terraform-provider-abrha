package snapshot

import (
	"context"
	"log"
	"strings"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/abrhacom/terraform-provider-abrha/abrha/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("abrha_vm_snapshot", &resource.Sweeper{
		Name:         "abrha_vm_snapshot",
		F:            testSweepVmSnapshots,
		Dependencies: []string{"abrha_vm"},
	})

	resource.AddTestSweepers("abrha_volume_snapshot", &resource.Sweeper{
		Name:         "abrha_volume_snapshot",
		F:            testSweepVolumeSnapshots,
		Dependencies: []string{"abrha_volume"},
	})
}

func testSweepVmSnapshots(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	opt := &goApiAbrha.ListOptions{PerPage: 200}
	snapshots, _, err := client.Snapshots.ListVm(context.Background(), opt)
	if err != nil {
		return err
	}

	for _, s := range snapshots {
		if strings.HasPrefix(s.Name, sweep.TestNamePrefix) {
			log.Printf("Destroying Vm Snapshot %s", s.Name)

			if _, err := client.Snapshots.Delete(context.Background(), s.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func testSweepVolumeSnapshots(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	opt := &goApiAbrha.ListOptions{PerPage: 200}
	snapshots, _, err := client.Snapshots.ListVolume(context.Background(), opt)
	if err != nil {
		return err
	}

	for _, s := range snapshots {
		if strings.HasPrefix(s.Name, sweep.TestNamePrefix) {
			log.Printf("Destroying Volume Snapshot %s", s.Name)

			if _, err := client.Snapshots.Delete(context.Background(), s.ID); err != nil {
				return err
			}
		}
	}

	return nil
}
