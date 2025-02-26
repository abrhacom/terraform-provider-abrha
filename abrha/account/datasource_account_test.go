package account_test

import (
	"testing"

	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAbrhaAccount_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceAbrhaAccountConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.abrha_account.foobar", "uuid"),
				),
			},
		},
	})
}

const testAccCheckDataSourceAbrhaAccountConfig_basic = `
data "abrha_account" "foobar" {
}`
