package database_test

import (
	"fmt"
	"testing"

	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAbrhaDatabaseMongoDBConfig_Basic(t *testing.T) {
	name := acceptance.RandomTestName()
	dbConfig := fmt.Sprintf(testAccCheckAbrhaDatabaseClusterMongoDB, name, "7")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseMongoDBConfigBasic, dbConfig, "available", 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("abrha_database_mongodb_config.foobar", "default_read_concern", "available"),
					resource.TestCheckResourceAttr("abrha_database_mongodb_config.foobar", "transaction_lifetime_limit_seconds", "1"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseMongoDBConfigBasic, dbConfig, "majority", 100),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("abrha_database_mongodb_config.foobar", "default_read_concern", "majority"),
					resource.TestCheckResourceAttr("abrha_database_mongodb_config.foobar", "transaction_lifetime_limit_seconds", "100"),
				),
			},
		},
	})
}

const testAccCheckAbrhaDatabaseMongoDBConfigBasic = `
%s

resource "abrha_database_mongodb_config" "foobar" {
  cluster_id                         = abrha_database_cluster.foobar.id
  default_read_concern               = "%s"
  transaction_lifetime_limit_seconds = %d
}`
