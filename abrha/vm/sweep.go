package vm

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
	resource.AddTestSweepers("abrha_vm", &resource.Sweeper{
		Name: "abrha_vm",
		F:    sweepVms,
	})
}

func sweepVms(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	opt := &goApiAbrha.ListOptions{PerPage: 200}
	vms, _, err := client.Vms.List(context.Background(), opt)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Found %d vms to sweep", len(vms))

	var swept int
	for _, d := range vms {
		if strings.HasPrefix(d.Name, sweep.TestNamePrefix) {
			log.Printf("Destroying Vm %s", d.Name)

			if _, err := client.Vms.Delete(context.Background(), d.ID); err != nil {
				return err
			}
			swept++
		}
	}
	log.Printf("[DEBUG] Deleted %d of %d vms", swept, len(vms))

	return nil
}
