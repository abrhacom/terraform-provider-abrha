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

func TestAccAbrhaSSHKey_Basic(t *testing.T) {
	var key goApiAbrha.Key
	name := acceptance.RandomTestName()
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("abrha@ssh-acceptance-test")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaSSHKeyConfig_basic(name, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaSSHKeyExists("abrha_ssh_key.foobar", &key),
					resource.TestCheckResourceAttr(
						"abrha_ssh_key.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"abrha_ssh_key.foobar", "public_key", publicKeyMaterial),
				),
			},
		},
	})
}

func testAccCheckAbrhaSSHKeyDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "abrha_ssh_key" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		// Try to find the key
		_, _, err = client.Keys.GetByID(context.Background(), id)

		if err == nil {
			return fmt.Errorf("SSH key still exists")
		}
	}

	return nil
}

func testAccCheckAbrhaSSHKeyExists(n string, key *goApiAbrha.Key) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		// Try to find the key
		foundKey, _, err := client.Keys.GetByID(context.Background(), id)

		if err != nil {
			return err
		}

		if strconv.Itoa(foundKey.ID) != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		*key = *foundKey

		return nil
	}
}

func testAccCheckAbrhaSSHKeyConfig_basic(name, key string) string {
	return fmt.Sprintf(`
resource "abrha_ssh_key" "foobar" {
  name       = "%s"
  public_key = "%s"
}`, name, key)
}
