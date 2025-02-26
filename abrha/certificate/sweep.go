package certificate

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
	resource.AddTestSweepers("abrha_certificate", &resource.Sweeper{
		Name: "abrha_certificate",
		F:    sweepCertificate,
	})

}

func sweepCertificate(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	opt := &goApiAbrha.ListOptions{PerPage: 200}
	certs, _, err := client.Certificates.List(context.Background(), opt)
	if err != nil {
		return err
	}

	for _, c := range certs {
		if strings.HasPrefix(c.Name, sweep.TestNamePrefix) {
			log.Printf("Destroying certificate %s", c.Name)

			if _, err := client.Certificates.Delete(context.Background(), c.ID); err != nil {
				return err
			}
		}
	}

	return nil
}
