package firewall

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
	resource.AddTestSweepers("abrha_firewall", &resource.Sweeper{
		Name: "abrha_firewall",
		F:    sweepFirewall,
	})

}

func sweepFirewall(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	opt := &goApiAbrha.ListOptions{PerPage: 200}
	fws, _, err := client.Firewalls.List(context.Background(), opt)
	if err != nil {
		return err
	}

	for _, f := range fws {
		if strings.HasPrefix(f.Name, sweep.TestNamePrefix) {
			log.Printf("Destroying firewall %s", f.Name)

			if _, err := client.Firewalls.Delete(context.Background(), f.ID); err != nil {
				return err
			}
		}
	}

	return nil
}
