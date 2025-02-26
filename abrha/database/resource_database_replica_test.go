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

func TestAccAbrhaDatabaseReplica_Basic(t *testing.T) {
	var databaseReplica goApiAbrha.DatabaseReplica
	var database goApiAbrha.Database

	databaseName := acceptance.RandomTestName()
	databaseReplicaName := acceptance.RandomTestName()

	databaseConfig := fmt.Sprintf(testAccCheckAbrhaDatabaseClusterConfigBasic, databaseName)
	replicaConfig := fmt.Sprintf(testAccCheckAbrhaDatabaseReplicaConfigBasic, databaseReplicaName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaDatabaseReplicaDestroy,
		Steps: []resource.TestStep{
			{
				Config: databaseConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseClusterExists("abrha_database_cluster.foobar", &database),
				),
			},
			{
				Config: databaseConfig + replicaConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseReplicaExists("abrha_database_replica.read-01", &databaseReplica),
					testAccCheckAbrhaDatabaseReplicaAttributes(&databaseReplica, databaseReplicaName),
					resource.TestCheckResourceAttr(
						"abrha_database_replica.read-01", "size", "db-s-1vcpu-2gb"),
					resource.TestCheckResourceAttr(
						"abrha_database_replica.read-01", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"abrha_database_replica.read-01", "name", databaseReplicaName),
					resource.TestCheckResourceAttrSet(
						"abrha_database_replica.read-01", "host"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_replica.read-01", "private_host"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_replica.read-01", "port"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_replica.read-01", "user"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_replica.read-01", "uri"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_replica.read-01", "private_uri"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_replica.read-01", "password"),
					resource.TestCheckResourceAttr(
						"abrha_database_replica.read-01", "tags.#", "1"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_replica.read-01", "private_network_uuid"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_replica.read-01", "uuid"),
				),
			},
		},
	})
}

func TestAccAbrhaDatabaseReplica_WithVPC(t *testing.T) {
	var database goApiAbrha.Database
	var databaseReplica goApiAbrha.DatabaseReplica

	vpcName := acceptance.RandomTestName()
	databaseName := acceptance.RandomTestName()
	databaseReplicaName := acceptance.RandomTestName()

	databaseConfig := fmt.Sprintf(testAccCheckAbrhaDatabaseClusterConfigWithVPC, vpcName, databaseName)
	replicaConfig := fmt.Sprintf(testAccCheckAbrhaDatabaseReplicaConfigWithVPC, databaseReplicaName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: databaseConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseClusterExists("abrha_database_cluster.foobar", &database),
				),
			},
			{
				Config: databaseConfig + replicaConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseReplicaExists("abrha_database_replica.read-01", &databaseReplica),
					testAccCheckAbrhaDatabaseClusterAttributes(&database, databaseName),
					resource.TestCheckResourceAttrPair(
						"abrha_database_replica.read-01", "private_network_uuid", "abrha_vpc.foobar", "id"),
				),
			},
		},
	})
}

func TestAccAbrhaDatabaseReplica_Resize(t *testing.T) {
	var databaseReplica goApiAbrha.DatabaseReplica
	var database goApiAbrha.Database

	databaseName := acceptance.RandomTestName()
	databaseReplicaName := acceptance.RandomTestName()

	databaseConfig := fmt.Sprintf(testAccCheckAbrhaDatabaseClusterConfigBasic, databaseName)
	replicaConfig := fmt.Sprintf(testAccCheckAbrhaDatabaseReplicaConfigBasic, databaseReplicaName)
	resizedConfig := fmt.Sprintf(testAccCheckAbrhaDatabaseReplicaConfigResized, databaseReplicaName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaDatabaseReplicaDestroy,
		Steps: []resource.TestStep{
			{
				Config: databaseConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseClusterExists("abrha_database_cluster.foobar", &database),
				),
			},
			{
				Config: databaseConfig + replicaConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseReplicaExists("abrha_database_replica.read-01", &databaseReplica),
					testAccCheckAbrhaDatabaseReplicaAttributes(&databaseReplica, databaseReplicaName),
					resource.TestCheckResourceAttr(
						"abrha_database_replica.read-01", "size", "db-s-1vcpu-2gb"),
					resource.TestCheckResourceAttr(
						"abrha_database_replica.read-01", "storage_size_mib", "30720"),
					resource.TestCheckResourceAttr(
						"abrha_database_replica.read-01", "name", databaseReplicaName),
					resource.TestCheckResourceAttrSet(
						"abrha_database_replica.read-01", "uuid"),
				),
			},
			{
				Config: databaseConfig + resizedConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseReplicaExists("abrha_database_replica.read-01", &databaseReplica),
					testAccCheckAbrhaDatabaseReplicaAttributes(&databaseReplica, databaseReplicaName),
					resource.TestCheckResourceAttr(
						"abrha_database_replica.read-01", "size", "db-s-2vcpu-4gb"),
					resource.TestCheckResourceAttr(
						"abrha_database_replica.read-01", "storage_size_mib", "61440"),
					resource.TestCheckResourceAttr(
						"abrha_database_replica.read-01", "name", databaseReplicaName),
					resource.TestCheckResourceAttrSet(
						"abrha_database_replica.read-01", "uuid"),
				),
			},
		},
	})
}

func testAccCheckAbrhaDatabaseReplicaDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "abrha_database_replica" {
			continue
		}
		clusterId := rs.Primary.Attributes["cluster_id"]
		name := rs.Primary.Attributes["name"]
		// Try to find the database replica
		_, _, err := client.Databases.GetReplica(context.Background(), clusterId, name)

		if err == nil {
			return fmt.Errorf("DatabaseReplica still exists")
		}
	}

	return nil
}

func testAccCheckAbrhaDatabaseReplicaExists(n string, database *goApiAbrha.DatabaseReplica) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No DatabaseReplica cluster ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()
		clusterId := rs.Primary.Attributes["cluster_id"]
		name := rs.Primary.Attributes["name"]
		uuid := rs.Primary.Attributes["uuid"]

		foundDatabaseReplica, _, err := client.Databases.GetReplica(context.Background(), clusterId, name)

		if err != nil {
			return err
		}

		if foundDatabaseReplica.Name != name {
			return fmt.Errorf("DatabaseReplica not found")
		}

		if foundDatabaseReplica.ID != uuid {
			return fmt.Errorf("DatabaseReplica UUID not found")
		}

		*database = *foundDatabaseReplica

		return nil
	}
}

func testAccCheckAbrhaDatabaseReplicaAttributes(databaseReplica *goApiAbrha.DatabaseReplica, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if databaseReplica.Name != name {
			return fmt.Errorf("Bad name: %s", databaseReplica.Name)
		}

		return nil
	}
}

const testAccCheckAbrhaDatabaseReplicaConfigBasic = `
resource "abrha_database_replica" "read-01" {
  cluster_id = abrha_database_cluster.foobar.id
  name       = "%s"
  region     = "nyc3"
  size       = "db-s-1vcpu-2gb"
  tags       = ["staging"]
}`

const testAccCheckAbrhaDatabaseReplicaConfigResized = `
resource "abrha_database_replica" "read-01" {
  cluster_id       = abrha_database_cluster.foobar.id
  name             = "%s"
  region           = "nyc3"
  size             = "db-s-2vcpu-4gb"
  storage_size_mib = 61440
  tags             = ["staging"]
}`

const testAccCheckAbrhaDatabaseReplicaConfigWithVPC = `


resource "abrha_database_replica" "read-01" {
  cluster_id           = abrha_database_cluster.foobar.id
  name                 = "%s"
  region               = "nyc1"
  size                 = "db-s-1vcpu-2gb"
  tags                 = ["staging"]
  private_network_uuid = abrha_vpc.foobar.id
}`
