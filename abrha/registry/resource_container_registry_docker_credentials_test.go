package registry_test

import (
	"context"
	"fmt"
	"testing"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAbrhaContainerRegistryDockerCredentials_Basic(t *testing.T) {
	var reg goApiAbrha.Registry
	name := acceptance.RandomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaContainerRegistryDockerCredentialsDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAbrhaContainerRegistryDockerCredentialsConfig_basic, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaContainerRegistryDockerCredentialsExists("abrha_container_registry.foobar", &reg),
					testAccCheckAbrhaContainerRegistryDockerCredentialsAttributes(&reg, name),
					resource.TestCheckResourceAttr(
						"abrha_container_registry_docker_credentials.foobar", "registry_name", name),
					resource.TestCheckResourceAttr(
						"abrha_container_registry_docker_credentials.foobar", "write", "true"),
					resource.TestCheckResourceAttrSet(
						"abrha_container_registry_docker_credentials.foobar", "docker_credentials"),
					resource.TestCheckResourceAttrSet(
						"abrha_container_registry_docker_credentials.foobar", "credential_expiration_time"),
				),
			},
		},
	})
}

func TestAccAbrhaContainerRegistryDockerCredentials_withExpiry(t *testing.T) {
	var reg goApiAbrha.Registry
	name := acceptance.RandomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaContainerRegistryDockerCredentialsDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAbrhaContainerRegistryDockerCredentialsConfig_withExpiry, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaContainerRegistryDockerCredentialsExists("abrha_container_registry.foobar", &reg),
					testAccCheckAbrhaContainerRegistryDockerCredentialsAttributes(&reg, name),
					resource.TestCheckResourceAttr(
						"abrha_container_registry_docker_credentials.foobar", "registry_name", name),
					resource.TestCheckResourceAttr(
						"abrha_container_registry_docker_credentials.foobar", "write", "true"),
					resource.TestCheckResourceAttr(
						"abrha_container_registry_docker_credentials.foobar", "expiry_seconds", "3600"),
					resource.TestCheckResourceAttrSet(
						"abrha_container_registry_docker_credentials.foobar", "docker_credentials"),
					resource.TestCheckResourceAttrSet(
						"abrha_container_registry_docker_credentials.foobar", "credential_expiration_time"),
				),
			},
		},
	})
}

func testAccCheckAbrhaContainerRegistryDockerCredentialsDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "abrha_container_registry_docker_credentials" {
			continue
		}

		// Try to find the key
		_, _, err := client.Registry.Get(context.Background())

		if err == nil {
			return fmt.Errorf("Container Registry still exists")
		}
	}

	return nil
}

func testAccCheckAbrhaContainerRegistryDockerCredentialsAttributes(reg *goApiAbrha.Registry, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if reg.Name != name {
			return fmt.Errorf("Bad name: %s", reg.Name)
		}

		return nil
	}
}

func testAccCheckAbrhaContainerRegistryDockerCredentialsExists(n string, reg *goApiAbrha.Registry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

		// Try to find the registry
		foundReg, _, err := client.Registry.Get(context.Background())

		if err != nil {
			return err
		}

		*reg = *foundReg

		return nil
	}
}

var testAccCheckAbrhaContainerRegistryDockerCredentialsConfig_basic = `
resource "abrha_container_registry" "foobar" {
  name                   = "%s"
  subscription_tier_slug = "basic"
}

resource "abrha_container_registry_docker_credentials" "foobar" {
  registry_name = abrha_container_registry.foobar.name
  write         = true
}`

var testAccCheckAbrhaContainerRegistryDockerCredentialsConfig_withExpiry = `
resource "abrha_container_registry" "foobar" {
  name                   = "%s"
  subscription_tier_slug = "basic"
}

resource "abrha_container_registry_docker_credentials" "foobar" {
  registry_name  = abrha_container_registry.foobar.name
  write          = true
  expiry_seconds = 3600
}`
