package database_test

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"testing"
	"time"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceAbrhaDatabaseCA(t *testing.T) {
	var database goApiAbrha.Database
	databaseName := acceptance.RandomTestName()
	databaseConfig := fmt.Sprintf(testAccCheckAbrhaDatabaseClusterConfigBasic, databaseName)

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
				Config: databaseConfig + testAccCheckAbrhaDatasourceCAConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.abrha_database_ca.ca", "certificate"),
					resource.TestCheckFunc(
						// Do some basic validation by parsing the certificate.
						func(s *terraform.State) error {
							rs, ok := s.RootModule().Resources["data.abrha_database_ca.ca"]
							if !ok {
								return fmt.Errorf("Not found: %s", "data.abrha_database_ca.ca")
							}

							certString := rs.Primary.Attributes["certificate"]
							block, _ := pem.Decode([]byte(certString))
							if block == nil {
								return fmt.Errorf("failed to parse certificate PEM")
							}
							cert, err := x509.ParseCertificate(block.Bytes)
							if err != nil {
								return fmt.Errorf("failed to parse certificate: " + err.Error())
							}

							if !cert.IsCA {
								return fmt.Errorf("not a CA cert")
							}

							return nil
						},
					),
				),
			},
		},
	})
}

const (
	testAccCheckAbrhaDatasourceCAConfig = `

data "abrha_database_ca" "ca" {
  cluster_id = abrha_database_cluster.foobar.id
}`
)
