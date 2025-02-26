package reservedipv6_test

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

func TestAccAbrhaReservedIPV6Assignment(t *testing.T) {
	var reservedIPv6 goApiAbrha.ReservedIPV6

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaReservedIPV6Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaReservedIPV6AssignmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaReservedIPV6AttachmentExists("abrha_reserved_ipv6_assignment.foobar"),
					resource.TestMatchResourceAttr(
						"abrha_reserved_ipv6_assignment.foobar", "ip", regexp.MustCompile(ipv6Regex)),
					resource.TestMatchResourceAttr(
						"abrha_reserved_ipv6_assignment.foobar", "vm_id", regexp.MustCompile("[0-9]+")),
				),
			},

			{
				Config: testAccCheckAbrhaReservedIPV6AssignmentDeleteAssignment,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaReservedIPV6Exists("abrha_reserved_ipv6.foobar", &reservedIPv6),
					resource.TestMatchResourceAttr(
						"abrha_reserved_ipv6.foobar", "ip", regexp.MustCompile(ipv6Regex)),
				),
			},
		},
	})
}

func TestAccAbrhaReservedIPV6Assignment_createBeforeDestroy(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaReservedIPV6Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaReservedIPV6AssignmentConfig_createBeforeDestroy,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaReservedIPV6AttachmentExists("abrha_reserved_ipv6_assignment.foobar"),
					resource.TestMatchResourceAttr(
						"abrha_reserved_ipv6_assignment.foobar", "ip", regexp.MustCompile(ipv6Regex)),
					resource.TestMatchResourceAttr(
						"abrha_reserved_ipv6_assignment.foobar", "vm_id", regexp.MustCompile("[0-9]+")),
				),
			},
			{
				Config: testAccCheckAbrhaReservedIPV6AssignmentConfig_createBeforeDestroyReassign,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaReservedIPV6AttachmentExists("abrha_reserved_ipv6_assignment.foobar"),
					resource.TestMatchResourceAttr(
						"abrha_reserved_ipv6_assignment.foobar", "ip", regexp.MustCompile(ipv6Regex)),
					resource.TestMatchResourceAttr(
						"abrha_reserved_ipv6_assignment.foobar", "vm_id", regexp.MustCompile("[0-9]+")),
				),
			},
		},
	})
}

func TestAccAbrhaReservedIPV6Assignment_unassignAndAssignToNewVm(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaReservedIPV6Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaReservedIPV6AssignmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaReservedIPV6AttachmentExists("abrha_reserved_ipv6_assignment.foobar"),
					resource.TestMatchResourceAttr(
						"abrha_reserved_ipv6_assignment.foobar", "ip", regexp.MustCompile(ipv6Regex)),
					resource.TestMatchResourceAttr(
						"abrha_reserved_ipv6_assignment.foobar", "vm_id", regexp.MustCompile("[0-9]+")),
				),
			},
			{
				Config: testAccCheckAbrhaReservedIPV6AssignmentUnAssign,
			},
			{
				Config: testAccCheckAbrhaReservedIPV6ReAssignmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaReservedIPV6AttachmentExists("abrha_reserved_ipv6_assignment.foobar"),
					resource.TestMatchResourceAttr(
						"abrha_reserved_ipv6_assignment.foobar", "ip", regexp.MustCompile(ipv6Regex)),
					resource.TestMatchResourceAttr(
						"abrha_reserved_ipv6_assignment.foobar", "vm_id", regexp.MustCompile("[0-9]+")),
				),
			},
		},
	})
}

func testAccCheckAbrhaReservedIPV6AttachmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.Attributes["ip"] == "" {
			return fmt.Errorf("No reserved IPv6 is set")
		}
		fipID := rs.Primary.Attributes["ip"]
		vmID := rs.Primary.Attributes["vm_id"]

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

		// Try to find the ReservedIPv6
		foundReservedIP, _, err := client.ReservedIPV6s.Get(context.Background(), fipID)
		if err != nil {
			return err
		}

		if foundReservedIP.IP != fipID || foundReservedIP.Vm.ID != vmID {
			return fmt.Errorf("wrong floating IP attachment found")
		}

		return nil
	}
}

var testAccCheckAbrhaReservedIPV6AssignmentConfig = `
resource "abrha_reserved_ipv6" "foobar" {
  region_slug = "nyc3"
}

resource "abrha_vm" "foobar" {
  count  = 1
  name   = "tf-acc-test-assign"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = "nyc3"
  ipv6   = true
}

resource "abrha_reserved_ipv6_assignment" "foobar" {
  ip         = abrha_reserved_ipv6.foobar.ip
  vm_id = abrha_vm.foobar.0.id
}
`

var testAccCheckAbrhaReservedIPV6AssignmentUnAssign = `
resource "abrha_reserved_ipv6" "foobar" {
  region_slug = "nyc3"
}

resource "abrha_vm" "foobar" {
  count  = 1
  image  = "ubuntu-22-04-x64"
  name   = "tf-acc-test-assign"
  region = "nyc3"
  size   = "s-1vcpu-1gb"
  ipv6   = true
}

`

var testAccCheckAbrhaReservedIPV6ReAssignmentConfig = `
resource "abrha_reserved_ipv6" "foobar" {
  region_slug = "nyc3"
}

resource "abrha_vm" "foobar1" {
  count  = 1
  name   = "tf-acc-test-reassign"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = "nyc3"
  ipv6   = true
}
resource "abrha_vm" "foobar" {
  count  = 1
  name   = "tf-acc-test-assign"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = "nyc3"
  ipv6   = true
}

resource "abrha_reserved_ipv6_assignment" "foobar" {
  ip         = abrha_reserved_ipv6.foobar.ip
  vm_id = abrha_vm.foobar1.0.id
}
`

var testAccCheckAbrhaReservedIPV6AssignmentDeleteAssignment = `
resource "abrha_reserved_ipv6" "foobar" {
  region_slug = "nyc3"
}

resource "abrha_vm" "foobar" {
  count  = 1
  name   = "tf-acc-test-${count.index}"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = "nyc3"
  ipv6   = true
}
`

var testAccCheckAbrhaReservedIPV6AssignmentConfig_createBeforeDestroy = `
resource "abrha_vm" "foobar" {
  image  = "ubuntu-22-04-x64"
  name   = "tf-acc-test"
  region = "nyc3"
  size   = "s-1vcpu-1gb"
  ipv6   = true

  lifecycle {
    create_before_destroy = true
  }
}

resource "abrha_reserved_ipv6" "foobar" {
  region_slug = "nyc3"
}

resource "abrha_reserved_ipv6_assignment" "foobar" {
  ip         = abrha_reserved_ipv6.foobar.ip
  vm_id = abrha_vm.foobar.id

  lifecycle {
    create_before_destroy = true
  }
}
`
var testAccCheckAbrhaReservedIPV6AssignmentConfig_createBeforeDestroyReassign = `
resource "abrha_vm" "foobar" {
  image  = "ubuntu-22-04-x64"
  name   = "tf-acc-test-01"
  region = "nyc3"
  size   = "s-1vcpu-1gb"
  ipv6   = true

  lifecycle {
    create_before_destroy = true
  }
}

resource "abrha_reserved_ipv6" "foobar" {
  region_slug = "nyc3"
}

resource "abrha_reserved_ipv6_assignment" "foobar" {
  ip         = abrha_reserved_ipv6.foobar.ip
  vm_id = abrha_vm.foobar.id

  lifecycle {
    create_before_destroy = true
  }
}
`
