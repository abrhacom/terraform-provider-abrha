package database_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAbrhaDatabaseConnectionPool_Basic(t *testing.T) {
	var databaseConnectionPool goApiAbrha.DatabasePool
	databaseName := acceptance.RandomTestName()
	databaseConnectionPoolName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaDatabaseConnectionPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseConnectionPoolConfigBasic, databaseName, databaseConnectionPoolName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseConnectionPoolExists("abrha_database_connection_pool.pool-01", &databaseConnectionPool),
					testAccCheckAbrhaDatabaseConnectionPoolAttributes(&databaseConnectionPool, databaseConnectionPoolName),
					resource.TestCheckResourceAttr(
						"abrha_database_connection_pool.pool-01", "name", databaseConnectionPoolName),
					resource.TestCheckResourceAttr(
						"abrha_database_connection_pool.pool-01", "size", "10"),
					resource.TestCheckResourceAttr(
						"abrha_database_connection_pool.pool-01", "mode", "transaction"),
					resource.TestCheckResourceAttr(
						"abrha_database_connection_pool.pool-01", "db_name", "defaultdb"),
					resource.TestCheckResourceAttr(
						"abrha_database_connection_pool.pool-01", "user", "doadmin"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_connection_pool.pool-01", "host"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_connection_pool.pool-01", "private_host"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_connection_pool.pool-01", "port"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_connection_pool.pool-01", "uri"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_connection_pool.pool-01", "private_uri"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_connection_pool.pool-01", "password"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseConnectionPoolConfigUpdated, databaseName, databaseConnectionPoolName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseConnectionPoolExists("abrha_database_connection_pool.pool-01", &databaseConnectionPool),
					testAccCheckAbrhaDatabaseConnectionPoolAttributes(&databaseConnectionPool, databaseConnectionPoolName),
					resource.TestCheckResourceAttr(
						"abrha_database_connection_pool.pool-01", "name", databaseConnectionPoolName),
					resource.TestCheckResourceAttr(
						"abrha_database_connection_pool.pool-01", "mode", "session"),
				),
			},
		},
	})
}

func TestAccAbrhaDatabaseConnectionPool_InboundUser(t *testing.T) {

	var databaseConnectionPool goApiAbrha.DatabasePool
	databaseName := acceptance.RandomTestName()
	databaseConnectionPoolName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaDatabaseConnectionPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseConnectionPoolConfigInboundUser, databaseName, databaseConnectionPoolName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseConnectionPoolExists("abrha_database_connection_pool.pool-01", &databaseConnectionPool),
					testAccCheckAbrhaDatabaseConnectionPoolAttributes(&databaseConnectionPool, databaseConnectionPoolName),
					resource.TestCheckResourceAttr(
						"abrha_database_connection_pool.pool-01", "name", databaseConnectionPoolName),
					resource.TestCheckResourceAttr(
						"abrha_database_connection_pool.pool-01", "size", "10"),
					resource.TestCheckResourceAttr(
						"abrha_database_connection_pool.pool-01", "mode", "transaction"),
					resource.TestCheckResourceAttr(
						"abrha_database_connection_pool.pool-01", "db_name", "defaultdb"),
					resource.TestCheckResourceAttr(
						"abrha_database_connection_pool.pool-01", "user", ""),
					resource.TestCheckResourceAttrSet(
						"abrha_database_connection_pool.pool-01", "host"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_connection_pool.pool-01", "private_host"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_connection_pool.pool-01", "port"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_connection_pool.pool-01", "uri"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_connection_pool.pool-01", "private_uri"),
				),
			},
		},
	})
}

func TestAccAbrhaDatabaseConnectionPool_BadModeName(t *testing.T) {
	databaseName := acceptance.RandomTestName()
	databaseConnectionPoolName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaDatabaseConnectionPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config:      fmt.Sprintf(testAccCheckAbrhaDatabaseConnectionPoolConfigBad, databaseName, databaseConnectionPoolName),
				ExpectError: regexp.MustCompile(`expected mode to be one of`),
			},
		},
	})
}

func testAccCheckAbrhaDatabaseConnectionPoolDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "abrha_database_connection_pool" {
			continue
		}
		clusterId := rs.Primary.Attributes["cluster_id"]
		name := rs.Primary.Attributes["name"]
		// Try to find the database connection_pool
		_, _, err := client.Databases.GetPool(context.Background(), clusterId, name)

		if err == nil {
			return fmt.Errorf("DatabaseConnectionPool still exists")
		}
	}

	return nil
}

func testAccCheckAbrhaDatabaseConnectionPoolExists(n string, database *goApiAbrha.DatabasePool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No DatabaseConnectionPool ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()
		clusterId := rs.Primary.Attributes["cluster_id"]
		name := rs.Primary.Attributes["name"]

		foundDatabaseConnectionPool, _, err := client.Databases.GetPool(context.Background(), clusterId, name)

		if err != nil {
			return err
		}

		if foundDatabaseConnectionPool.Name != name {
			return fmt.Errorf("DatabaseConnectionPool not found")
		}

		*database = *foundDatabaseConnectionPool

		return nil
	}
}

func testAccCheckAbrhaDatabaseConnectionPoolAttributes(databaseConnectionPool *goApiAbrha.DatabasePool, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if databaseConnectionPool.Name != name {
			return fmt.Errorf("Bad name: %s", databaseConnectionPool.Name)
		}

		return nil
	}
}

const testAccCheckAbrhaDatabaseConnectionPoolConfigBasic = `
resource "abrha_database_cluster" "foobar" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

resource "abrha_database_connection_pool" "pool-01" {
  cluster_id = abrha_database_cluster.foobar.id
  name       = "%s"
  mode       = "transaction"
  size       = 10
  db_name    = "defaultdb"
  user       = "doadmin"
}`

const testAccCheckAbrhaDatabaseConnectionPoolConfigUpdated = `
resource "abrha_database_cluster" "foobar" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

resource "abrha_database_connection_pool" "pool-01" {
  cluster_id = abrha_database_cluster.foobar.id
  name       = "%s"
  mode       = "session"
  size       = 10
  db_name    = "defaultdb"
}`

const testAccCheckAbrhaDatabaseConnectionPoolConfigBad = `
resource "abrha_database_cluster" "foobar" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

resource "abrha_database_connection_pool" "pool-01" {
  cluster_id = abrha_database_cluster.foobar.id
  name       = "%s"
  mode       = "transactional"
  size       = 10
  db_name    = "defaultdb"
  user       = "doadmin"
}`

const testAccCheckAbrhaDatabaseConnectionPoolConfigInboundUser = `
resource "abrha_database_cluster" "foobar" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

resource "abrha_database_connection_pool" "pool-01" {
  cluster_id = abrha_database_cluster.foobar.id
  name       = "%s"
  mode       = "transaction"
  size       = 10
  db_name    = "defaultdb"
}`
