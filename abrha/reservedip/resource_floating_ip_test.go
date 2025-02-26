package reservedip_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAbrhaFloatingIP_Region(t *testing.T) {
	var floatingIP goApiAbrha.FloatingIP

	expectedURNRegEx, _ := regexp.Compile(`do:floatingip:(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaFloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaFloatingIPConfig_region,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaFloatingIPExists("abrha_floating_ip.foobar", &floatingIP),
					resource.TestCheckResourceAttr(
						"abrha_floating_ip.foobar", "region", "nyc3"),
					resource.TestMatchResourceAttr("abrha_floating_ip.foobar", "urn", expectedURNRegEx),
				),
			},
		},
	})
}

func TestAccAbrhaFloatingIP_Vm(t *testing.T) {
	var floatingIP goApiAbrha.FloatingIP
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaFloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaFloatingIPConfig_vm(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaFloatingIPExists("abrha_floating_ip.foobar", &floatingIP),
					resource.TestCheckResourceAttr(
						"abrha_floating_ip.foobar", "region", "nyc3"),
				),
			},
			{
				Config: testAccCheckAbrhaFloatingIPConfig_Reassign(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaFloatingIPExists("abrha_floating_ip.foobar", &floatingIP),
					resource.TestCheckResourceAttr(
						"abrha_floating_ip.foobar", "region", "nyc3"),
				),
			},
			{
				Config: testAccCheckAbrhaFloatingIPConfig_Unassign(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaFloatingIPExists("abrha_floating_ip.foobar", &floatingIP),
					resource.TestCheckResourceAttr(
						"abrha_floating_ip.foobar", "region", "nyc3"),
				),
			},
		},
	})
}

func testAccCheckAbrhaFloatingIPDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "abrha_floating_ip" {
			continue
		}

		// Try to find the key
		_, _, err := client.FloatingIPs.Get(context.Background(), rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Floating IP still exists")
		}
	}

	return nil
}

func testAccCheckAbrhaFloatingIPExists(n string, floatingIP *goApiAbrha.FloatingIP) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

		// Try to find the FloatingIP
		foundFloatingIP, _, err := client.FloatingIPs.Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundFloatingIP.IP != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		*floatingIP = *foundFloatingIP

		return nil
	}
}

var testAccCheckAbrhaFloatingIPConfig_region = `
resource "abrha_floating_ip" "foobar" {
  region = "nyc3"
}`

func testAccCheckAbrhaFloatingIPConfig_vm(name string) string {
	return fmt.Sprintf(`
resource "abrha_vm" "foobar" {
  name               = "%s"
  size               = "s-1vcpu-1gb"
  image              = "ubuntu-22-04-x64"
  region             = "nyc3"
  ipv6               = true
  private_networking = true
}

resource "abrha_floating_ip" "foobar" {
  vm_id = abrha_vm.foobar.id
  region     = abrha_vm.foobar.region
}`, name)
}

func testAccCheckAbrhaFloatingIPConfig_Reassign(name string) string {
	return fmt.Sprintf(`
resource "abrha_vm" "baz" {
  name               = "%s"
  size               = "s-1vcpu-1gb"
  image              = "ubuntu-22-04-x64"
  region             = "nyc3"
  ipv6               = true
  private_networking = true
}

resource "abrha_floating_ip" "foobar" {
  vm_id = abrha_vm.baz.id
  region     = abrha_vm.baz.region
}`, name)
}

func testAccCheckAbrhaFloatingIPConfig_Unassign(name string) string {
	return fmt.Sprintf(`
resource "abrha_vm" "baz" {
  name               = "%s"
  size               = "s-1vcpu-1gb"
  image              = "ubuntu-22-04-x64"
  region             = "nyc3"
  ipv6               = true
  private_networking = true
}

resource "abrha_floating_ip" "foobar" {
  region = "nyc3"
}`, name)
}
