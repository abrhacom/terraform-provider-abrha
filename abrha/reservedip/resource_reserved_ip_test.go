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

func TestAccAbrhaReservedIP_Region(t *testing.T) {
	var reservedIP goApiAbrha.ReservedIP

	expectedURNRegEx, _ := regexp.Compile(`do:reservedip:(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaReservedIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaReservedIPConfig_region,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaReservedIPExists("abrha_reserved_ip.foobar", &reservedIP),
					resource.TestCheckResourceAttr(
						"abrha_reserved_ip.foobar", "region", "nyc3"),
					resource.TestMatchResourceAttr("abrha_reserved_ip.foobar", "urn", expectedURNRegEx),
				),
			},
		},
	})
}

func TestAccAbrhaReservedIP_Vm(t *testing.T) {
	var reservedIP goApiAbrha.ReservedIP
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaReservedIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaReservedIPConfig_vm(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaReservedIPExists("abrha_reserved_ip.foobar", &reservedIP),
					resource.TestCheckResourceAttr(
						"abrha_reserved_ip.foobar", "region", "nyc3"),
				),
			},
			{
				Config: testAccCheckAbrhaReservedIPConfig_Reassign(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaReservedIPExists("abrha_reserved_ip.foobar", &reservedIP),
					resource.TestCheckResourceAttr(
						"abrha_reserved_ip.foobar", "region", "nyc3"),
				),
			},
			{
				Config: testAccCheckAbrhaReservedIPConfig_Unassign(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaReservedIPExists("abrha_reserved_ip.foobar", &reservedIP),
					resource.TestCheckResourceAttr(
						"abrha_reserved_ip.foobar", "region", "nyc3"),
				),
			},
		},
	})
}

func testAccCheckAbrhaReservedIPDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "abrha_reserved_ip" {
			continue
		}

		// Try to find the key
		_, _, err := client.ReservedIPs.Get(context.Background(), rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Reserved IP still exists")
		}
	}

	return nil
}

func testAccCheckAbrhaReservedIPExists(n string, reservedIP *goApiAbrha.ReservedIP) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

		// Try to find the ReservedIP
		foundReservedIP, _, err := client.ReservedIPs.Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundReservedIP.IP != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		*reservedIP = *foundReservedIP

		return nil
	}
}

var testAccCheckAbrhaReservedIPConfig_region = `
resource "abrha_reserved_ip" "foobar" {
  region = "nyc3"
}`

func testAccCheckAbrhaReservedIPConfig_vm(name string) string {
	return fmt.Sprintf(`
resource "abrha_vm" "foobar" {
  name               = "%s"
  size               = "s-1vcpu-1gb"
  image              = "ubuntu-22-04-x64"
  region             = "nyc3"
  ipv6               = true
  private_networking = true
}

resource "abrha_reserved_ip" "foobar" {
  vm_id = abrha_vm.foobar.id
  region     = abrha_vm.foobar.region
}`, name)
}

func testAccCheckAbrhaReservedIPConfig_Reassign(name string) string {
	return fmt.Sprintf(`
resource "abrha_vm" "baz" {
  name               = "%s"
  size               = "s-1vcpu-1gb"
  image              = "ubuntu-22-04-x64"
  region             = "nyc3"
  ipv6               = true
  private_networking = true
}

resource "abrha_reserved_ip" "foobar" {
  vm_id = abrha_vm.baz.id
  region     = abrha_vm.baz.region
}`, name)
}

func testAccCheckAbrhaReservedIPConfig_Unassign(name string) string {
	return fmt.Sprintf(`
resource "abrha_vm" "baz" {
  name               = "%s"
  size               = "s-1vcpu-1gb"
  image              = "ubuntu-22-04-x64"
  region             = "nyc3"
  ipv6               = true
  private_networking = true
}

resource "abrha_reserved_ip" "foobar" {
  region = "nyc3"
}`, name)
}
