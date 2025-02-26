package database_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAbrhaDatabaseKafkaTopic(t *testing.T) {
	name := acceptance.RandomTestName()
	dbConfig := fmt.Sprintf(testAccCheckAbrhaDatabaseClusterKafka, name, "3.5")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaDatabaseKafkaTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseKafkaTopicBasic, dbConfig, "topic-foobar"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"abrha_database_kafka_topic.foobar", "name", "topic-foobar"),
					resource.TestCheckResourceAttr(
						"abrha_database_kafka_topic.foobar", "state", "active"),
					resource.TestCheckResourceAttr(
						"abrha_database_kafka_topic.foobar", "replication_factor", "2"),
					resource.TestCheckResourceAttr(
						"abrha_database_kafka_topic.foobar", "partition_count", "3"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_kafka_topic.foobar", "config.0.cleanup_policy"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_kafka_topic.foobar", "config.0.compression_type"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_kafka_topic.foobar", "config.0.delete_retention_ms"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_kafka_topic.foobar", "config.0.file_delete_delay_ms"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_kafka_topic.foobar", "config.0.flush_messages"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_kafka_topic.foobar", "config.0.flush_ms"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_kafka_topic.foobar", "config.0.index_interval_bytes"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_kafka_topic.foobar", "config.0.max_compaction_lag_ms"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_kafka_topic.foobar", "config.0.message_down_conversion_enable"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_kafka_topic.foobar", "config.0.message_format_version"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_kafka_topic.foobar", "config.0.message_timestamp_difference_max_ms"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_kafka_topic.foobar", "config.0.message_timestamp_type"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_kafka_topic.foobar", "config.0.min_cleanable_dirty_ratio"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_kafka_topic.foobar", "config.0.min_compaction_lag_ms"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_kafka_topic.foobar", "config.0.min_insync_replicas"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_kafka_topic.foobar", "config.0.retention_bytes"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_kafka_topic.foobar", "config.0.retention_ms"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_kafka_topic.foobar", "config.0.segment_bytes"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_kafka_topic.foobar", "config.0.segment_index_bytes"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_kafka_topic.foobar", "config.0.segment_jitter_ms"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_kafka_topic.foobar", "config.0.segment_ms"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseKafkaTopicWithConfig, dbConfig, "topic-foobar", 5, 3, "compact", "snappy", 80000),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"abrha_database_kafka_topic.foobar", "name", "topic-foobar"),
					resource.TestCheckResourceAttr(
						"abrha_database_kafka_topic.foobar", "state", "active"),
					resource.TestCheckResourceAttr(
						"abrha_database_kafka_topic.foobar", "replication_factor", "3"),
					resource.TestCheckResourceAttr(
						"abrha_database_kafka_topic.foobar", "partition_count", "5"),
					resource.TestCheckResourceAttr(
						"abrha_database_kafka_topic.foobar", "config.0.cleanup_policy", "compact"),
					resource.TestCheckResourceAttr(
						"abrha_database_kafka_topic.foobar", "config.0.compression_type", "snappy"),
					resource.TestCheckResourceAttr(
						"abrha_database_kafka_topic.foobar", "config.0.delete_retention_ms", "80000"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_kafka_topic.foobar", "config.0.cleanup_policy"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_kafka_topic.foobar", "config.0.compression_type"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_kafka_topic.foobar", "config.0.delete_retention_ms"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_kafka_topic.foobar", "config.0.file_delete_delay_ms"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_kafka_topic.foobar", "config.0.flush_messages"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_kafka_topic.foobar", "config.0.flush_ms"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_kafka_topic.foobar", "config.0.index_interval_bytes"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_kafka_topic.foobar", "config.0.max_compaction_lag_ms"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_kafka_topic.foobar", "config.0.message_down_conversion_enable"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_kafka_topic.foobar", "config.0.message_format_version"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_kafka_topic.foobar", "config.0.message_timestamp_difference_max_ms"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_kafka_topic.foobar", "config.0.message_timestamp_type"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_kafka_topic.foobar", "config.0.min_cleanable_dirty_ratio"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_kafka_topic.foobar", "config.0.min_compaction_lag_ms"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_kafka_topic.foobar", "config.0.min_insync_replicas"),
					resource.TestCheckResourceAttr(
						"abrha_database_kafka_topic.foobar", "config.0.min_insync_replicas", "1"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_kafka_topic.foobar", "config.0.retention_bytes"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_kafka_topic.foobar", "config.0.retention_ms"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_kafka_topic.foobar", "config.0.segment_bytes"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_kafka_topic.foobar", "config.0.segment_index_bytes"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_kafka_topic.foobar", "config.0.segment_jitter_ms"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_kafka_topic.foobar", "config.0.segment_ms"),
				),
			},
		},
	})
}

func testAccCheckAbrhaDatabaseKafkaTopicDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "abrha_database_kafka_topic" {
			continue
		}
		clusterId := rs.Primary.Attributes["cluster_id"]
		name := rs.Primary.Attributes["name"]
		// Try to find the kafka topic
		_, _, err := client.Databases.GetTopic(context.Background(), clusterId, name)

		if err == nil {
			return fmt.Errorf("kafka topic still exists")
		}
	}

	return nil
}

const testAccCheckAbrhaDatabaseKafkaTopicBasic = `
%s

resource "abrha_database_kafka_topic" "foobar" {
  cluster_id = abrha_database_cluster.foobar.id
  name       = "%s"
}`

const testAccCheckAbrhaDatabaseKafkaTopicWithConfig = `
%s

resource "abrha_database_kafka_topic" "foobar" {
  cluster_id         = abrha_database_cluster.foobar.id
  name               = "%s"
  partition_count    = %d
  replication_factor = %d
  config {
    cleanup_policy      = "%s"
    compression_type    = "%s"
    delete_retention_ms = %d
  }
}`
