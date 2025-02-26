package domain

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
	resource.AddTestSweepers("abrha_domain", &resource.Sweeper{
		Name: "abrha_domain",
		F:    sweepDomain,
	})

}

func sweepDomain(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	opt := &goApiAbrha.ListOptions{PerPage: 200}
	domains, _, err := client.Domains.List(context.Background(), opt)
	if err != nil {
		return err
	}

	for _, d := range domains {
		if strings.HasPrefix(d.Name, sweep.TestNamePrefix) {
			log.Printf("Destroying domain %s", d.Name)

			if _, err := client.Domains.Delete(context.Background(), d.Name); err != nil {
				return err
			}
		}
	}

	return nil
}
