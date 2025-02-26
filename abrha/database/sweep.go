package database

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
	resource.AddTestSweepers("abrha_database_cluster", &resource.Sweeper{
		Name: "abrha_database_cluster",
		F:    testSweepDatabaseCluster,
	})

}

func testSweepDatabaseCluster(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	opt := &goApiAbrha.ListOptions{PerPage: 200}
	databases, _, err := client.Databases.List(context.Background(), opt)
	if err != nil {
		return err
	}

	for _, db := range databases {
		if strings.HasPrefix(db.Name, sweep.TestNamePrefix) {
			log.Printf("Destroying database cluster %s", db.Name)

			if _, err := client.Databases.Delete(context.Background(), db.ID); err != nil {
				return err
			}
		}
	}

	return nil
}
