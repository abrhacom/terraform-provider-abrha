package sshkey_test

import (
	"fmt"
	"testing"

	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAbrhaSSHKeys_Basic(t *testing.T) {
	keyName1 := acceptance.RandomTestName("datasource1")
	pubKey1, _, err := acctest.RandSSHKeyPair("abrha@ssh-acceptance-test")
	if err != nil {
		t.Fatalf("Unable to generate public key: %v", err)
		return
	}
	keyName2 := acceptance.RandomTestName("datasource2")
	pubKey2, _, err := acctest.RandSSHKeyPair("abrha@ssh-acceptance-test")
	if err != nil {
		t.Fatalf("Unable to generate public key: %v", err)
		return
	}

	resourcesConfig := fmt.Sprintf(`
resource "abrha_ssh_key" "foo" {
  name       = "%s"
  public_key = "%s"
}

resource "abrha_ssh_key" "bar" {
  name       = "%s"
  public_key = "%s"
}
`, keyName1, pubKey1, keyName2, pubKey2)

	datasourceConfig := fmt.Sprintf(`
data "abrha_ssh_keys" "result" {
  sort {
    key       = "name"
    direction = "asc"
  }
  filter {
    key    = "name"
    values = ["%s", "%s"]
  }
}
`, keyName1, keyName2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourcesConfig,
			},
			{
				Config: resourcesConfig + datasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.abrha_ssh_keys.result", "ssh_keys.#", "2"),
					resource.TestCheckResourceAttr("data.abrha_ssh_keys.result", "ssh_keys.0.name", keyName1),
					resource.TestCheckResourceAttr("data.abrha_ssh_keys.result", "ssh_keys.1.name", keyName2),
				),
			},
			{
				Config: resourcesConfig,
			},
		},
	})
}
