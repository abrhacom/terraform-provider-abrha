package loadbalancer

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
	resource.AddTestSweepers("abrha_loadbalancer", &resource.Sweeper{
		Name: "abrha_loadbalancer",
		F:    sweepLoadbalancer,
	})

}

func sweepLoadbalancer(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	opt := &goApiAbrha.ListOptions{PerPage: 200}
	lbs, _, err := client.LoadBalancers.List(context.Background(), opt)
	if err != nil {
		return err
	}

	for _, l := range lbs {
		if strings.HasPrefix(l.Name, sweep.TestNamePrefix) {
			log.Printf("Destroying loadbalancer %s", l.Name)

			if _, err := client.LoadBalancers.Delete(context.Background(), l.ID); err != nil {
				return err
			}
		}
	}

	return nil
}
