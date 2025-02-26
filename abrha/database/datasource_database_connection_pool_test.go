package database_test

import (
	"fmt"
	"testing"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAbrhaDatabaseConnectionPool_Basic(t *testing.T) {
	var pool goApiAbrha.DatabasePool

	databaseName := acceptance.RandomTestName()
	poolName := acceptance.RandomTestName()

	resourceConfig := fmt.Sprintf(testAccCheckAbrhaDatabaseConnectionPoolConfigBasic, databaseName, poolName)
	datasourceConfig := fmt.Sprintf(testAccCheckAbrhaDatasourceDatabaseConnectionPoolConfigBasic, poolName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaDatabaseConnectionPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseConnectionPoolExists("abrha_database_connection_pool.pool-01", &pool),
					testAccCheckAbrhaDatabaseConnectionPoolAttributes(&pool, poolName),
					resource.TestCheckResourceAttr(
						"abrha_database_connection_pool.pool-01", "name", poolName),
					resource.TestCheckResourceAttrSet(
						"abrha_database_connection_pool.pool-01", "cluster_id"),
					resource.TestCheckResourceAttr(
						"abrha_database_connection_pool.pool-01", "size", "10"),
					resource.TestCheckResourceAttr(
						"abrha_database_connection_pool.pool-01", "mode", "transaction"),
					resource.TestCheckResourceAttr(
						"abrha_database_connection_pool.pool-01", "db_name", "defaultdb"),
					resource.TestCheckResourceAttr(
						"abrha_database_connection_pool.pool-01", "user", "doadmin"),
				),
			},
			{
				Config: resourceConfig + datasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("abrha_database_connection_pool.pool-01", "name",
						"data.abrha_database_connection_pool.pool-01", "name"),
					resource.TestCheckResourceAttrPair("abrha_database_connection_pool.pool-01", "mode",
						"data.abrha_database_connection_pool.pool-01", "mode"),
					resource.TestCheckResourceAttrPair("abrha_database_connection_pool.pool-01", "size",
						"data.abrha_database_connection_pool.pool-01", "size"),
					resource.TestCheckResourceAttrPair("abrha_database_connection_pool.pool-01", "db_name",
						"data.abrha_database_connection_pool.pool-01", "db_name"),
					resource.TestCheckResourceAttrPair("abrha_database_connection_pool.pool-01", "user",
						"data.abrha_database_connection_pool.pool-01", "user"),
				),
			},
		},
	})
}

const testAccCheckAbrhaDatasourceDatabaseConnectionPoolConfigBasic = `
data "abrha_database_connection_pool" "pool-01" {
  cluster_id = abrha_database_cluster.foobar.id
  name       = "%s"
}`
