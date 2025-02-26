package sshkey_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceAbrhaSSHKey_Basic(t *testing.T) {
	var key goApiAbrha.Key
	keyName := acceptance.RandomTestName()

	pubKey, _, err := acctest.RandSSHKeyPair("abrha@ssh-acceptance-test")
	if err != nil {
		t.Fatalf("Unable to generate public key: %v", err)
		return
	}

	resourceConfig := fmt.Sprintf(`
resource "abrha_ssh_key" "foo" {
  name       = "%s"
  public_key = "%s"
}`, keyName, pubKey)

	dataSourceConfig := `
data "abrha_ssh_key" "foobar" {
  name = abrha_ssh_key.foo.name
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
			},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceAbrhaSSHKeyExists("data.abrha_ssh_key.foobar", &key),
					resource.TestCheckResourceAttr(
						"data.abrha_ssh_key.foobar", "name", keyName),
					resource.TestCheckResourceAttr(
						"data.abrha_ssh_key.foobar", "public_key", pubKey),
				),
			},
		},
	})
}

func testAccCheckDataSourceAbrhaSSHKeyExists(n string, key *goApiAbrha.Key) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ssh key ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		foundKey, _, err := client.Keys.GetByID(context.Background(), id)

		if err != nil {
			return err
		}

		if foundKey.ID != id {
			return fmt.Errorf("Key not found")
		}

		*key = *foundKey

		return nil
	}
}
