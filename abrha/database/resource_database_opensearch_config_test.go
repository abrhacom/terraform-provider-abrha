package database_test

import (
	"fmt"
	"testing"

	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAbrhaDatabaseOpensearchConfig_Basic(t *testing.T) {
	name := acceptance.RandomTestName()
	dbConfig := fmt.Sprintf(testAccCheckAbrhaDatabaseClusterOpensearch, name, "2")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseOpensearchConfigConfigBasic, dbConfig, true, 10, "1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("abrha_database_opensearch_config.foobar", "enable_security_audit", "true"),
					resource.TestCheckResourceAttr("abrha_database_opensearch_config.foobar", "ism_enabled", "true"),
					resource.TestCheckResourceAttr("abrha_database_opensearch_config.foobar", "ism_history_enabled", "true"),
					resource.TestCheckResourceAttr("abrha_database_opensearch_config.foobar", "ism_history_max_age_hours", "10"),
					resource.TestCheckResourceAttr("abrha_database_opensearch_config.foobar", "ism_history_max_docs", "1"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDatabaseOpensearchConfigConfigBasic, dbConfig, false, 1, "1000000000000000000"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("abrha_database_opensearch_config.foobar", "enable_security_audit", "false"),
					resource.TestCheckResourceAttr("abrha_database_opensearch_config.foobar", "ism_enabled", "true"),
					resource.TestCheckResourceAttr("abrha_database_opensearch_config.foobar", "ism_history_enabled", "true"),
					resource.TestCheckResourceAttr("abrha_database_opensearch_config.foobar", "ism_history_max_age_hours", "1"),
					resource.TestCheckResourceAttr("abrha_database_opensearch_config.foobar", "ism_history_max_docs", "1000000000000000000"),
				),
			},
		},
	})
}

const testAccCheckAbrhaDatabaseOpensearchConfigConfigBasic = `
%s

resource "abrha_database_opensearch_config" "foobar" {
  cluster_id                = abrha_database_cluster.foobar.id
  enable_security_audit     = %t
  ism_enabled               = true
  ism_history_enabled       = true
  ism_history_max_age_hours = %d
  ism_history_max_docs      = %s
}`
