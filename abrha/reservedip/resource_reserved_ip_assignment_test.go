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

func TestAccAbrhaReservedIPAssignment(t *testing.T) {
	var reservedIP goApiAbrha.ReservedIP

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaReservedIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaReservedIPAssignmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaReservedIPAttachmentExists("abrha_reserved_ip_assignment.foobar"),
					resource.TestMatchResourceAttr(
						"abrha_reserved_ip_assignment.foobar", "id", regexp.MustCompile("[0-9.]+")),
					resource.TestMatchResourceAttr(
						"abrha_reserved_ip_assignment.foobar", "vm_id", regexp.MustCompile("[0-9]+")),
				),
			},
			{
				Config: testAccCheckAbrhaReservedIPAssignmentReassign,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaReservedIPAttachmentExists("abrha_reserved_ip_assignment.foobar"),
					resource.TestMatchResourceAttr(
						"abrha_reserved_ip_assignment.foobar", "id", regexp.MustCompile("[0-9.]+")),
					resource.TestMatchResourceAttr(
						"abrha_reserved_ip_assignment.foobar", "vm_id", regexp.MustCompile("[0-9]+")),
				),
			},
			{
				Config: testAccCheckAbrhaReservedIPAssignmentDeleteAssignment,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaReservedIPExists("abrha_reserved_ip.foobar", &reservedIP),
					resource.TestMatchResourceAttr(
						"abrha_reserved_ip.foobar", "ip_address", regexp.MustCompile("[0-9.]+")),
				),
			},
		},
	})
}

func TestAccAbrhaReservedIPAssignment_createBeforeDestroy(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaReservedIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaReservedIPAssignmentConfig_createBeforeDestroy,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaReservedIPAttachmentExists("abrha_reserved_ip_assignment.foobar"),
					resource.TestMatchResourceAttr(
						"abrha_reserved_ip_assignment.foobar", "id", regexp.MustCompile("[0-9.]+")),
					resource.TestMatchResourceAttr(
						"abrha_reserved_ip_assignment.foobar", "vm_id", regexp.MustCompile("[0-9]+")),
				),
			},
			{
				Config: testAccCheckAbrhaReservedIPAssignmentConfig_createBeforeDestroyReassign,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaReservedIPAttachmentExists("abrha_reserved_ip_assignment.foobar"),
					resource.TestMatchResourceAttr(
						"abrha_reserved_ip_assignment.foobar", "id", regexp.MustCompile("[0-9.]+")),
					resource.TestMatchResourceAttr(
						"abrha_reserved_ip_assignment.foobar", "vm_id", regexp.MustCompile("[0-9]+")),
				),
			},
		},
	})
}

func testAccCheckAbrhaReservedIPAttachmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.Attributes["ip_address"] == "" {
			return fmt.Errorf("No floating IP is set")
		}
		fipID := rs.Primary.Attributes["ip_address"]
		vmID := rs.Primary.Attributes["vm_id"]

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

		// Try to find the ReservedIP
		foundReservedIP, _, err := client.ReservedIPs.Get(context.Background(), fipID)
		if err != nil {
			return err
		}

		if foundReservedIP.IP != fipID || foundReservedIP.Vm.ID != vmID {
			return fmt.Errorf("wrong floating IP attachment found")
		}

		return nil
	}
}

var testAccCheckAbrhaReservedIPAssignmentConfig = `
resource "abrha_reserved_ip" "foobar" {
  region = "nyc3"
}

resource "abrha_vm" "foobar" {
  count              = 2
  name               = "tf-acc-test-${count.index}"
  size               = "s-1vcpu-1gb"
  image              = "ubuntu-22-04-x64"
  region             = "nyc3"
  ipv6               = true
  private_networking = true
}

resource "abrha_reserved_ip_assignment" "foobar" {
  ip_address = abrha_reserved_ip.foobar.ip_address
  vm_id      = abrha_vm.foobar.0.id
}
`

var testAccCheckAbrhaReservedIPAssignmentReassign = `
resource "abrha_reserved_ip" "foobar" {
  region = "nyc3"
}

resource "abrha_vm" "foobar" {
  count              = 2
  name               = "tf-acc-test-${count.index}"
  size               = "s-1vcpu-1gb"
  image              = "ubuntu-22-04-x64"
  region             = "nyc3"
  ipv6               = true
  private_networking = true
}

resource "abrha_reserved_ip_assignment" "foobar" {
  ip_address = abrha_reserved_ip.foobar.ip_address
  vm_id      = abrha_vm.foobar.1.id
}
`

var testAccCheckAbrhaReservedIPAssignmentDeleteAssignment = `
resource "abrha_reserved_ip" "foobar" {
  region = "nyc3"
}

resource "abrha_vm" "foobar" {
  count              = 2
  name               = "tf-acc-test-${count.index}"
  size               = "s-1vcpu-1gb"
  image              = "ubuntu-22-04-x64"
  region             = "nyc3"
  ipv6               = true
  private_networking = true
}
`

var testAccCheckAbrhaReservedIPAssignmentConfig_createBeforeDestroy = `
resource "abrha_vm" "foobar" {
  image  = "ubuntu-22-04-x64"
  name   = "tf-acc-test"
  region = "nyc3"
  size   = "s-1vcpu-1gb"

  lifecycle {
    create_before_destroy = true
  }
}

resource "abrha_reserved_ip" "foobar" {
  region = "nyc3"
}

resource "abrha_reserved_ip_assignment" "foobar" {
  ip_address = abrha_reserved_ip.foobar.id
  vm_id      = abrha_vm.foobar.id

  lifecycle {
    create_before_destroy = true
  }
}
`

var testAccCheckAbrhaReservedIPAssignmentConfig_createBeforeDestroyReassign = `
resource "abrha_vm" "foobar" {
  image  = "ubuntu-22-04-x64"
  name   = "tf-acc-test-01"
  region = "nyc3"
  size   = "s-1vcpu-1gb"

  lifecycle {
    create_before_destroy = true
  }
}

resource "abrha_reserved_ip" "foobar" {
  region = "nyc3"
}

resource "abrha_reserved_ip_assignment" "foobar" {
  ip_address = abrha_reserved_ip.foobar.id
  vm_id      = abrha_vm.foobar.id

  lifecycle {
    create_before_destroy = true
  }
}
`
