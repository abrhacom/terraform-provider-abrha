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

func TestAccAbrhaFloatingIPAssignment(t *testing.T) {
	var floatingIP goApiAbrha.FloatingIP

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaFloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaFloatingIPAssignmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaFloatingIPAttachmentExists("abrha_floating_ip_assignment.foobar"),
					resource.TestMatchResourceAttr(
						"abrha_floating_ip_assignment.foobar", "id", regexp.MustCompile("[0-9.]+")),
					resource.TestMatchResourceAttr(
						"abrha_floating_ip_assignment.foobar", "vm_id", regexp.MustCompile("[0-9]+")),
				),
			},
			{
				Config: testAccCheckAbrhaFloatingIPAssignmentReassign,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaFloatingIPAttachmentExists("abrha_floating_ip_assignment.foobar"),
					resource.TestMatchResourceAttr(
						"abrha_floating_ip_assignment.foobar", "id", regexp.MustCompile("[0-9.]+")),
					resource.TestMatchResourceAttr(
						"abrha_floating_ip_assignment.foobar", "vm_id", regexp.MustCompile("[0-9]+")),
				),
			},
			{
				Config: testAccCheckAbrhaFloatingIPAssignmentDeleteAssignment,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaFloatingIPExists("abrha_floating_ip.foobar", &floatingIP),
					resource.TestMatchResourceAttr(
						"abrha_floating_ip.foobar", "ip_address", regexp.MustCompile("[0-9.]+")),
				),
			},
		},
	})
}

func TestAccAbrhaFloatingIPAssignment_createBeforeDestroy(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaFloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaFloatingIPAssignmentConfig_createBeforeDestroy,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaFloatingIPAttachmentExists("abrha_floating_ip_assignment.foobar"),
					resource.TestMatchResourceAttr(
						"abrha_floating_ip_assignment.foobar", "id", regexp.MustCompile("[0-9.]+")),
					resource.TestMatchResourceAttr(
						"abrha_floating_ip_assignment.foobar", "vm_id", regexp.MustCompile("[0-9]+")),
				),
			},
			{
				Config: testAccCheckAbrhaFloatingIPAssignmentConfig_createBeforeDestroyReassign,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaFloatingIPAttachmentExists("abrha_floating_ip_assignment.foobar"),
					resource.TestMatchResourceAttr(
						"abrha_floating_ip_assignment.foobar", "id", regexp.MustCompile("[0-9.]+")),
					resource.TestMatchResourceAttr(
						"abrha_floating_ip_assignment.foobar", "vm_id", regexp.MustCompile("[0-9]+")),
				),
			},
		},
	})
}

func testAccCheckAbrhaFloatingIPAttachmentExists(n string) resource.TestCheckFunc {
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

		// Try to find the FloatingIP
		foundFloatingIP, _, err := client.FloatingIPs.Get(context.Background(), fipID)
		if err != nil {
			return err
		}

		if foundFloatingIP.IP != fipID || foundFloatingIP.Vm.ID != vmID {
			return fmt.Errorf("wrong floating IP attachment found")
		}

		return nil
	}
}

var testAccCheckAbrhaFloatingIPAssignmentConfig = `
resource "abrha_floating_ip" "foobar" {
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

resource "abrha_floating_ip_assignment" "foobar" {
  ip_address = abrha_floating_ip.foobar.ip_address
  vm_id = abrha_vm.foobar[0].id
}
`

var testAccCheckAbrhaFloatingIPAssignmentReassign = `
resource "abrha_floating_ip" "foobar" {
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

resource "abrha_floating_ip_assignment" "foobar" {
  ip_address = abrha_floating_ip.foobar.ip_address
  vm_id = abrha_vm.foobar[1].id
}
`

var testAccCheckAbrhaFloatingIPAssignmentDeleteAssignment = `
resource "abrha_floating_ip" "foobar" {
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

var testAccCheckAbrhaFloatingIPAssignmentConfig_createBeforeDestroy = `
resource "abrha_vm" "foobar" {
  image  = "ubuntu-22-04-x64"
  name   = "tf-acc-test-01"
  region = "nyc3"
  size   = "s-1vcpu-1gb"

  lifecycle {
    create_before_destroy = true
  }
}

resource "abrha_floating_ip" "foobar" {
  region = "nyc3"
}

resource "abrha_floating_ip_assignment" "foobar" {
  ip_address = abrha_floating_ip.foobar.id
  vm_id = abrha_vm.foobar.id

  lifecycle {
    create_before_destroy = true
  }
}
`

var testAccCheckAbrhaFloatingIPAssignmentConfig_createBeforeDestroyReassign = `
resource "abrha_vm" "foobar" {
  image  = "ubuntu-18-04-x64"
  name   = "tf-acc-test-01"
  region = "nyc3"
  size   = "s-1vcpu-1gb"

  lifecycle {
    create_before_destroy = true
  }
}

resource "abrha_floating_ip" "foobar" {
  region = "nyc3"
}

resource "abrha_floating_ip_assignment" "foobar" {
  ip_address = abrha_floating_ip.foobar.id
  vm_id = abrha_vm.foobar.id

  lifecycle {
    create_before_destroy = true
  }
}
`
