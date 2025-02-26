package app

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
	resource.AddTestSweepers("abrha_app", &resource.Sweeper{
		Name: "abrha_app",
		F:    sweepApp,
	})

}

func sweepApp(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	opt := &goApiAbrha.ListOptions{PerPage: 200}
	apps, _, err := client.Apps.List(context.Background(), opt)
	if err != nil {
		return err
	}

	for _, app := range apps {
		if strings.HasPrefix(app.Spec.Name, sweep.TestNamePrefix) {
			log.Printf("Destroying app %s", app.Spec.Name)

			if _, err := client.Apps.Delete(context.Background(), app.ID); err != nil {
				return err
			}
		}
	}

	return nil
}
