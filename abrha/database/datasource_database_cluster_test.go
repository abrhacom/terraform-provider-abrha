package database_test

import (
	"context"
	"fmt"
	"testing"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceAbrhaDatabaseCluster_Basic(t *testing.T) {
	var database goApiAbrha.Database
	databaseName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceAbrhaDatabaseClusterConfigBasic, databaseName),
			},
			{
				Config: fmt.Sprintf(testAccCheckDataSourceAbrhaDatabaseClusterConfigWithDatasource, databaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceAbrhaDatabaseClusterExists("data.abrha_database_cluster.foobar", &database),
					resource.TestCheckResourceAttr(
						"data.abrha_database_cluster.foobar", "name", databaseName),
					resource.TestCheckResourceAttr(
						"data.abrha_database_cluster.foobar", "engine", "pg"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_database_cluster.foobar", "host"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_database_cluster.foobar", "private_host"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_database_cluster.foobar", "port"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_database_cluster.foobar", "user"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_database_cluster.foobar", "password"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_database_cluster.foobar", "private_network_uuid"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_database_cluster.foobar", "project_id"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_database_cluster.foobar", "storage_size_mib"),
					testAccCheckAbrhaDatabaseClusterURIPassword(
						"abrha_database_cluster.foobar", "uri"),
					testAccCheckAbrhaDatabaseClusterURIPassword(
						"abrha_database_cluster.foobar", "private_uri"),
				),
			},
		},
	})
}

func testAccCheckDataSourceAbrhaDatabaseClusterExists(n string, databaseCluster *goApiAbrha.Database) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

		foundCluster, _, err := client.Databases.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}

		if foundCluster.ID != rs.Primary.ID {
			return fmt.Errorf("DatabaseCluster not found")
		}

		*databaseCluster = *foundCluster

		return nil
	}
}

const testAccCheckDataSourceAbrhaDatabaseClusterConfigBasic = `
resource "abrha_database_cluster" "foobar" {
  name             = "%s"
  engine           = "pg"
  version          = "15"
  size             = "db-s-1vcpu-1gb"
  region           = "nyc1"
  node_count       = 1
  tags             = ["production"]
  storage_size_mib = 10240
}
`

const testAccCheckDataSourceAbrhaDatabaseClusterConfigWithDatasource = `
resource "abrha_database_cluster" "foobar" {
  name             = "%s"
  engine           = "pg"
  version          = "15"
  size             = "db-s-1vcpu-1gb"
  region           = "nyc1"
  node_count       = 1
  tags             = ["production"]
  storage_size_mib = 10240
}

data "abrha_database_cluster" "foobar" {
  name = abrha_database_cluster.foobar.name
}
`
