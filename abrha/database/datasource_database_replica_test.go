package database_test

import (
	"fmt"
	"testing"
	"time"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceAbrhaDatabaseReplica_Basic(t *testing.T) {
	var databaseReplica goApiAbrha.DatabaseReplica
	var database goApiAbrha.Database

	databaseName := acceptance.RandomTestName()
	databaseReplicaName := acceptance.RandomTestName()

	databaseConfig := fmt.Sprintf(testAccCheckAbrhaDatabaseClusterConfigBasic, databaseName)
	replicaConfig := fmt.Sprintf(testAccCheckAbrhaDatabaseReplicaConfigBasic, databaseReplicaName)
	datasourceReplicaConfig := fmt.Sprintf(testAccCheckAbrhaDatasourceDatabaseReplicaConfigBasic, databaseReplicaName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaDatabaseReplicaDestroy,
		Steps: []resource.TestStep{
			{
				Config: databaseConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseClusterExists("abrha_database_cluster.foobar", &database),
					resource.TestCheckFunc(
						func(s *terraform.State) error {
							time.Sleep(30 * time.Second)
							return nil
						},
					),
				),
			},
			{
				Config: databaseConfig + replicaConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseReplicaExists("abrha_database_replica.read-01", &databaseReplica),
					testAccCheckAbrhaDatabaseReplicaAttributes(&databaseReplica, databaseReplicaName),
				),
			},
			{
				Config: databaseConfig + replicaConfig + datasourceReplicaConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("abrha_database_replica.read-01", "cluster_id",
						"data.abrha_database_replica.my_db_replica", "cluster_id"),
					resource.TestCheckResourceAttrPair("abrha_database_replica.read-01", "name",
						"data.abrha_database_replica.my_db_replica", "name"),
					resource.TestCheckResourceAttrPair("abrha_database_replica.read-01", "uuid",
						"data.abrha_database_replica.my_db_replica", "uuid"),
					resource.TestCheckResourceAttr(
						"data.abrha_database_replica.my_db_replica", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"data.abrha_database_replica.my_db_replica", "name", databaseReplicaName),
					resource.TestCheckResourceAttrSet(
						"data.abrha_database_replica.my_db_replica", "host"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_database_replica.my_db_replica", "private_host"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_database_replica.my_db_replica", "port"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_database_replica.my_db_replica", "user"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_database_replica.my_db_replica", "uri"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_database_replica.my_db_replica", "private_uri"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_database_replica.my_db_replica", "password"),
					resource.TestCheckResourceAttr(
						"data.abrha_database_replica.my_db_replica", "tags.#", "1"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_database_replica.my_db_replica", "private_network_uuid"),
					resource.TestCheckResourceAttr(
						"data.abrha_database_replica.my_db_replica", "storage_size_mib", "30720"),
				),
			},
		},
	})
}

const (
	testAccCheckAbrhaDatasourceDatabaseReplicaConfigBasic = `
data "abrha_database_replica" "my_db_replica" {
  cluster_id = abrha_database_cluster.foobar.id
  name       = "%s"
}`
)
