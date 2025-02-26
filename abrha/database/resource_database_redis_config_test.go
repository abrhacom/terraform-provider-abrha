package database_test

import (
	"fmt"
	"testing"

	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAbrhaDatabaseRedisConfig_Basic(t *testing.T) {
	name := acceptance.RandomTestName()
	dbConfig := fmt.Sprintf(testAccCheckAbrhaDatabaseClusterRedis, name, "7")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseRedisConfigConfigBasic, dbConfig, "noeviction", 3600, "KA"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"abrha_database_redis_config.foobar", "maxmemory_policy", "noeviction"),
					resource.TestCheckResourceAttr(
						"abrha_database_redis_config.foobar", "timeout", "3600"),
					resource.TestCheckResourceAttr(
						"abrha_database_redis_config.foobar", "notify_keyspace_events", "KA"),
					resource.TestCheckResourceAttr(
						"abrha_database_redis_config.foobar", "ssl", "true"),
					resource.TestCheckResourceAttr(
						"abrha_database_redis_config.foobar", "persistence", "rdb"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseRedisConfigConfigBasic, dbConfig, "allkeys-lru", 0, "KEA"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"abrha_database_redis_config.foobar", "maxmemory_policy", "allkeys-lru"),
					resource.TestCheckResourceAttr(
						"abrha_database_redis_config.foobar", "timeout", "0"),
					resource.TestCheckResourceAttr(
						"abrha_database_redis_config.foobar", "notify_keyspace_events", "KEA"),
					resource.TestCheckResourceAttr(
						"abrha_database_redis_config.foobar", "ssl", "true"),
					resource.TestCheckResourceAttr(
						"abrha_database_redis_config.foobar", "persistence", "rdb"),
				),
			},
		},
	})
}

const testAccCheckAbrhaDatabaseRedisConfigConfigBasic = `
%s

resource "abrha_database_redis_config" "foobar" {
  cluster_id             = abrha_database_cluster.foobar.id
  maxmemory_policy       = "%s"
  timeout                = %d
  notify_keyspace_events = "%s"
}`
