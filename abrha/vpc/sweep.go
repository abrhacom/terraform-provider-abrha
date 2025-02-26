package vpc

import (
	"context"
	"log"
	"net/http"
	"strings"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/abrhacom/terraform-provider-abrha/abrha/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("abrha_vpc", &resource.Sweeper{
		Name: "abrha_vpc",
		F:    sweepVPC,
		Dependencies: []string{
			"abrha_vm",
			"abrha_vpc_peering",
		},
	})
}

func sweepVPC(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	opt := &goApiAbrha.ListOptions{PerPage: 200}
	vpcs, _, err := client.VPCs.List(context.Background(), opt)
	if err != nil {
		return err
	}

	for _, v := range vpcs {
		if strings.HasPrefix(v.Name, sweep.TestNamePrefix) {
			log.Printf("[DEBUG] Destroying VPC %s", v.Name)
			resp, err := client.VPCs.Delete(context.Background(), v.ID)
			if err != nil {
				if resp.StatusCode == http.StatusForbidden {
					log.Printf("[DEBUG] Skipping VPC %s; still contains resources", v.Name)
				} else {
					return err
				}
			}
		}
	}

	return nil
}
