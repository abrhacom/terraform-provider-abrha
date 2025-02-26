package vmautoscale

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
	resource.AddTestSweepers("abrha_vm_autoscale", &resource.Sweeper{
		Name: "abrha_vm_autoscale",
		F:    sweepVmAutoscale,
	})
}

func sweepVmAutoscale(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()
	pools, _, err := client.VmAutoscale.List(context.Background(), &goApiAbrha.ListOptions{PerPage: 200})
	if err != nil {
		return err
	}
	for _, pool := range pools {
		if strings.HasPrefix(pool.Name, sweep.TestNamePrefix) {
			log.Printf("Destroying vm autoscale pool %s", pool.Name)
			if _, err = client.VmAutoscale.DeleteDangerous(context.Background(), pool.ID); err != nil {
				return err
			}
		}
	}
	return nil
}
