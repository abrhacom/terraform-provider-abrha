package project

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
	resource.AddTestSweepers("abrha_project", &resource.Sweeper{
		Name: "abrha_project",
		F:    sweepProjects,
		Dependencies: []string{
			//	"abrha_spaces_bucket", TODO: Add when Spaces sweeper exists.
			"abrha_vm",
			"abrha_domain",
		},
	})
}

func sweepProjects(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	opt := &goApiAbrha.ListOptions{PerPage: 200}
	projects, _, err := client.Projects.List(context.Background(), opt)
	if err != nil {
		return err
	}

	for _, p := range projects {
		if strings.HasPrefix(p.Name, sweep.TestNamePrefix) {
			log.Printf("[DEBUG] Destroying project %s", p.Name)

			resp, err := client.Projects.Delete(context.Background(), p.ID)
			if err != nil {
				// Projects with resources can not be deleted.
				if resp.StatusCode == http.StatusPreconditionFailed {
					log.Printf("[DEBUG] Skipping project %s: %s", p.Name, err.Error())
				} else {
					return err
				}
			}
		}
	}

	return nil
}
