package database_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAbrhaDatabaseDB_Basic(t *testing.T) {
	var databaseDB goApiAbrha.DatabaseDB
	databaseClusterName := fmt.Sprintf("foobar-test-terraform-%s", acctest.RandString(10))
	databaseDBName := fmt.Sprintf("foobar-test-db-terraform-%s", acctest.RandString(10))
	databaseDBNameUpdated := databaseDBName + "-up"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaDatabaseDBDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseDBConfigBasic, databaseClusterName, databaseDBName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseDBExists("abrha_database_db.foobar_db", &databaseDB),
					testAccCheckAbrhaDatabaseDBAttributes(&databaseDB, databaseDBName),
					resource.TestCheckResourceAttr(
						"abrha_database_db.foobar_db", "name", databaseDBName),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseDBConfigBasic, databaseClusterName, databaseDBNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseDBExists("abrha_database_db.foobar_db", &databaseDB),
					testAccCheckAbrhaDatabaseDBNotExists("abrha_database_db.foobar_db", databaseDBName),
					testAccCheckAbrhaDatabaseDBAttributes(&databaseDB, databaseDBNameUpdated),
					resource.TestCheckResourceAttr(
						"abrha_database_db.foobar_db", "name", databaseDBNameUpdated),
				),
			},
		},
	})
}

func testAccCheckAbrhaDatabaseDBDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "abrha_database_db" {
			continue
		}
		clusterID := rs.Primary.Attributes["cluster_id"]
		name := rs.Primary.Attributes["name"]

		// Try to find the database DB
		_, _, err := client.Databases.GetDB(context.Background(), clusterID, name)

		if err == nil {
			return fmt.Errorf("Database DB still exists")
		}
	}

	return nil
}

func testAccCheckAbrhaDatabaseDBExists(n string, databaseDB *goApiAbrha.DatabaseDB) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Database DB ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()
		clusterID := rs.Primary.Attributes["cluster_id"]
		name := rs.Primary.Attributes["name"]

		foundDatabaseDB, _, err := client.Databases.GetDB(context.Background(), clusterID, name)

		if err != nil {
			return err
		}

		if foundDatabaseDB.Name != name {
			return fmt.Errorf("Database DB not found")
		}

		*databaseDB = *foundDatabaseDB

		return nil
	}
}

func testAccCheckAbrhaDatabaseDBNotExists(n string, databaseDBName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Database DB ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()
		clusterID := rs.Primary.Attributes["cluster_id"]

		_, resp, err := client.Databases.GetDB(context.Background(), clusterID, databaseDBName)

		if err != nil && resp.StatusCode != http.StatusNotFound {
			return err
		}

		if err == nil {
			return fmt.Errorf("Database DB %s still exists", databaseDBName)
		}

		return nil
	}
}

func testAccCheckAbrhaDatabaseDBAttributes(databaseDB *goApiAbrha.DatabaseDB, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if databaseDB.Name != name {
			return fmt.Errorf("Bad name: %s", databaseDB.Name)
		}

		return nil
	}
}

const testAccCheckAbrhaDatabaseDBConfigBasic = `
resource "abrha_database_cluster" "foobar" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1

  maintenance_window {
    day  = "friday"
    hour = "13:00:00"
  }
}

resource "abrha_database_db" "foobar_db" {
  cluster_id = abrha_database_cluster.foobar.id
  name       = "%s"
}`
