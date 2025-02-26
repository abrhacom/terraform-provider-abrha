package database_test

import (
	"fmt"
	"testing"

	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAbrhaDatabasePostgreSQLConfig_Basic(t *testing.T) {
	name := acceptance.RandomTestName()
	dbConfig := fmt.Sprintf(testAccCheckAbrhaDatabaseClusterPostgreSQL, name, "15")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabasePostgreSQLConfigConfigBasic, dbConfig, "UTC", 30.5, 32, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("abrha_database_postgresql_config.foobar", "jit", "false"),
					resource.TestCheckResourceAttr("abrha_database_postgresql_config.foobar", "timezone", "UTC"),
					resource.TestCheckResourceAttr("abrha_database_postgresql_config.foobar", "shared_buffers_percentage", "30.5"),
					resource.TestCheckResourceAttr("abrha_database_postgresql_config.foobar", "work_mem", "32"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabasePostgreSQLConfigConfigBasic, dbConfig, "UTC", 20.0, 16, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("abrha_database_postgresql_config.foobar", "jit", "true"),
					resource.TestCheckResourceAttr("abrha_database_postgresql_config.foobar", "timezone", "UTC"),
					resource.TestCheckResourceAttr("abrha_database_postgresql_config.foobar", "shared_buffers_percentage", "20"),
					resource.TestCheckResourceAttr("abrha_database_postgresql_config.foobar", "work_mem", "16"),
				),
			},
		},
	})
}

const testAccCheckAbrhaDatabasePostgreSQLConfigConfigBasic = `
%s

resource "abrha_database_postgresql_config" "foobar" {
  cluster_id                = abrha_database_cluster.foobar.id
  timezone                  = "%s"
  shared_buffers_percentage = %f
  work_mem                  = %d
  jit                       = %t
  timescaledb {
    max_background_workers = 1
  }
}`
