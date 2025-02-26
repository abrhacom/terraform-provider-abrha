package database_test

import (
	"fmt"
	"regexp"
	"testing"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAbrhaDatabaseReplica_importBasic(t *testing.T) {
	var database goApiAbrha.Database
	resourceName := "abrha_database_replica.read-01"
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
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// Requires passing both the cluster ID and replica name
				ImportStateIdFunc: testAccDatabaseReplicaImportID(resourceName),
			},
			// Test importing non-existent resource provides expected error.
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateId:     fmt.Sprintf("%s,%s", "this-cluster-id-does-not-exist", databaseReplicaName),
				ExpectError:       regexp.MustCompile(`(Please verify the ID is correct|Cannot import non-existent remote object)`),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateId:     "replica",
				ExpectError:       regexp.MustCompile("joined with a comma"),
			},
		},
	})
}

func testAccDatabaseReplicaImportID(n string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return "", fmt.Errorf("Not found: %s", n)
		}

		clusterId := rs.Primary.Attributes["cluster_id"]
		name := rs.Primary.Attributes["name"]

		return fmt.Sprintf("%s,%s", clusterId, name), nil
	}
}
