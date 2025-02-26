package database_test

import (
	"fmt"
	"testing"

	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAbrhaDatabaseMySQLConfig_Basic(t *testing.T) {
	name := acceptance.RandomTestName()
	dbConfig := fmt.Sprintf(testAccCheckAbrhaDatabaseClusterMySQL, name, "8")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseMySQLConfigConfigBasic, dbConfig, 10, "UTC", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("abrha_database_mysql_config.foobar", "connect_timeout", "10"),
					resource.TestCheckResourceAttr("abrha_database_mysql_config.foobar", "default_time_zone", "UTC"),
					resource.TestCheckResourceAttr("abrha_database_mysql_config.foobar", "sql_require_primary_key", "false"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseMySQLConfigConfigBasic, dbConfig, 15, "SYSTEM", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("abrha_database_mysql_config.foobar", "connect_timeout", "15"),
					resource.TestCheckResourceAttr("abrha_database_mysql_config.foobar", "default_time_zone", "SYSTEM"),
					resource.TestCheckResourceAttr("abrha_database_mysql_config.foobar", "sql_require_primary_key", "false"),
				),
			},
		},
	})
}

const testAccCheckAbrhaDatabaseMySQLConfigConfigBasic = `
%s

resource "abrha_database_mysql_config" "foobar" {
  cluster_id              = abrha_database_cluster.foobar.id
  connect_timeout         = %d
  default_time_zone       = "%s"
  sql_require_primary_key = "%t"
}`
