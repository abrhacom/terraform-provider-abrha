package uptime

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
	resource.AddTestSweepers("abrha_uptime_check", &resource.Sweeper{
		Name: "abrha_uptime_check",
		F:    sweepUptimeCheck,
	})

	// Note: Deleting the check will delete associated alerts. So no sweeper is
	// needed for abrha_uptime_alert
}

func sweepUptimeCheck(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	opt := &goApiAbrha.ListOptions{PerPage: 200}
	checks, _, err := client.UptimeChecks.List(context.Background(), opt)
	if err != nil {
		return err
	}

	for _, c := range checks {
		if strings.HasPrefix(c.Name, sweep.TestNamePrefix) {
			log.Printf("[DEBUG] Deleting uptime check %s", c.Name)

			if _, err := client.UptimeChecks.Delete(context.Background(), c.ID); err != nil {
				return err
			}
		}
	}

	return nil
}
