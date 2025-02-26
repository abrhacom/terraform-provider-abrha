package domain_test

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

func TestAccAbrhaDomain_Basic(t *testing.T) {
	var domain goApiAbrha.Domain
	domainName := acceptance.RandomTestName() + ".com"

	expectedURN := fmt.Sprintf("do:domain:%s", domainName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDomainConfig_basic, domainName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDomainExists("abrha_domain.foobar", &domain),
					testAccCheckAbrhaDomainAttributes(&domain, domainName),
					resource.TestCheckResourceAttr(
						"abrha_domain.foobar", "name", domainName),
					resource.TestCheckResourceAttr(
						"abrha_domain.foobar", "ip_address", "192.168.0.10"),
					resource.TestCheckResourceAttr(
						"abrha_domain.foobar", "urn", expectedURN),
				),
			},
		},
	})
}

func TestAccAbrhaDomain_WithoutIp(t *testing.T) {
	var domain goApiAbrha.Domain
	domainName := acceptance.RandomTestName() + ".com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAbrhaDomainConfig_withoutIp, domainName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaDomainExists("abrha_domain.foobar", &domain),
					testAccCheckAbrhaDomainAttributes(&domain, domainName),
					resource.TestCheckResourceAttr(
						"abrha_domain.foobar", "name", domainName),
					resource.TestCheckNoResourceAttr(
						"abrha_domain.foobar", "ip_address"),
				),
			},
		},
	})
}

func testAccCheckAbrhaDomainDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "abrha_domain" {
			continue
		}

		// Try to find the domain
		_, _, err := client.Domains.Get(context.Background(), rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Domain still exists")
		}
	}

	return nil
}

func testAccCheckAbrhaDomainAttributes(domain *goApiAbrha.Domain, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if domain.Name != name {
			return fmt.Errorf("Bad name: %s", domain.Name)
		}

		return nil
	}
}

func testAccCheckAbrhaDomainExists(n string, domain *goApiAbrha.Domain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

		foundDomain, _, err := client.Domains.Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundDomain.Name != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		*domain = *foundDomain

		return nil
	}
}

const testAccCheckAbrhaDomainConfig_basic = `
resource "abrha_domain" "foobar" {
  name       = "%s"
  ip_address = "192.168.0.10"
}`

const testAccCheckAbrhaDomainConfig_withoutIp = `
resource "abrha_domain" "foobar" {
  name = "%s"
}`
