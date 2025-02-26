package database_test

import (
	"fmt"
	"testing"

	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAbrhaDatabaseKafkaConfig_Basic(t *testing.T) {
	name := acceptance.RandomTestName()
	dbConfig := fmt.Sprintf(testAccCheckAbrhaDatabaseClusterKafka, name, "3.7")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseKafkaConfigConfigBasic, dbConfig, 3000, true, "9223372036854776000"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("abrha_database_kafka_config.foobar", "group_initial_rebalance_delay_ms", "3000"),
					resource.TestCheckResourceAttr("abrha_database_kafka_config.foobar", "log_message_downconversion_enable", "true"),
					resource.TestCheckResourceAttr("abrha_database_kafka_config.foobar", "log_message_timestamp_difference_max_ms", "9223372036854776000"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseKafkaConfigConfigBasic, dbConfig, 300000, false, "0"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("abrha_database_kafka_config.foobar", "group_initial_rebalance_delay_ms", "300000"),
					resource.TestCheckResourceAttr("abrha_database_kafka_config.foobar", "log_message_downconversion_enable", "false"),
					resource.TestCheckResourceAttr("abrha_database_kafka_config.foobar", "log_message_timestamp_difference_max_ms", "0"),
				),
			},
		},
	})
}

const testAccCheckAbrhaDatabaseKafkaConfigConfigBasic = `
%s

resource "abrha_database_kafka_config" "foobar" {
  cluster_id                              = abrha_database_cluster.foobar.id
  group_initial_rebalance_delay_ms        = %d
  log_message_downconversion_enable       = %t
  log_message_timestamp_difference_max_ms = %s
}`
