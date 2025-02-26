package database_test

import (
	"fmt"
	"testing"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAbrhaDatabaseUser_Basic(t *testing.T) {
	var user goApiAbrha.DatabaseUser

	databaseName := acceptance.RandomTestName()
	userName := acceptance.RandomTestName()

	resourceConfig := fmt.Sprintf(testAccCheckAbrhaDatabaseUserConfigBasic, databaseName, userName)
	datasourceConfig := fmt.Sprintf(testAccCheckAbrhaDatasourceDatabaseUserConfigBasic, userName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaDatabaseUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDatabaseUserExists("abrha_database_user.foobar_user", &user),
					testAccCheckAbrhaDatabaseUserAttributes(&user, userName),
					resource.TestCheckResourceAttr(
						"abrha_database_user.foobar_user", "name", userName),
					resource.TestCheckResourceAttrSet(
						"abrha_database_user.foobar_user", "role"),
					resource.TestCheckResourceAttrSet(
						"abrha_database_user.foobar_user", "password"),
					resource.TestCheckNoResourceAttr(
						"abrha_database_user.foobar_user", "access_cert"),
					resource.TestCheckNoResourceAttr(
						"abrha_database_user.foobar_user", "access_key"),
				),
			},
			{
				Config: resourceConfig + datasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("abrha_database_user.foobar_user", "name",
						"data.abrha_database_user.foobar_user", "name"),
					resource.TestCheckResourceAttrPair("abrha_database_user.foobar_user", "role",
						"data.abrha_database_user.foobar_user", "role"),
					resource.TestCheckResourceAttrPair("abrha_database_user.foobar_user", "password",
						"data.abrha_database_user.foobar_user", "password"),
					resource.TestCheckNoResourceAttr(
						"abrha_database_user.foobar_user", "access_cert"),
					resource.TestCheckNoResourceAttr(
						"abrha_database_user.foobar_user", "access_key"),
				),
			},
		},
	})
}

const testAccCheckAbrhaDatasourceDatabaseUserConfigBasic = `
data "abrha_database_user" "foobar_user" {
  cluster_id = abrha_database_cluster.foobar.id
  name       = "%s"
}`
