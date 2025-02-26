package cdn

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
	resource.AddTestSweepers("abrha_cdn", &resource.Sweeper{
		Name: "abrha_cdn",
		F:    sweepCDN,
	})

}

func sweepCDN(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	opt := &goApiAbrha.ListOptions{PerPage: 200}
	cdns, _, err := client.CDNs.List(context.Background(), opt)
	if err != nil {
		return err
	}

	for _, c := range cdns {
		if strings.HasPrefix(c.Origin, sweep.TestNamePrefix) {
			log.Printf("Destroying CDN %s", c.Origin)

			if _, err := client.CDNs.Delete(context.Background(), c.ID); err != nil {
				return err
			}
		}
	}

	return nil
}
