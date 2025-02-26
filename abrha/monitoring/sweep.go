package monitoring

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
	resource.AddTestSweepers("abrha_monitor_alert", &resource.Sweeper{
		Name: "abrha_monitor_alert",
		F:    sweepMonitoringAlerts,
	})

}

func sweepMonitoringAlerts(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	opt := &goApiAbrha.ListOptions{PerPage: 200}
	alerts, _, err := client.Monitoring.ListAlertPolicies(context.Background(), opt)
	if err != nil {
		return err
	}

	for _, a := range alerts {
		if strings.HasPrefix(a.Description, sweep.TestNamePrefix) {
			log.Printf("[DEBUG] Destroying alert %s", a.Description)

			if _, err := client.Monitoring.DeleteAlertPolicy(context.Background(), a.UUID); err != nil {
				return err
			}
		}
	}

	return nil
}
