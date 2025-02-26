package reservedip_test

import (
	"context"
	"testing"

	"fmt"
	"regexp"

	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAbrhaFloatingIPAssignment_importBasic(t *testing.T) {
	resourceName := "abrha_floating_ip_assignment.foobar"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaFloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaFloatingIPAssignmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaFloatingIPAttachmentExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccFLIPAssignmentImportID(resourceName),
				// floating_ip_assignments are "virtual" resources that have unique, timestamped IDs.
				// As the imported one will have a different ID that the initial one, the states will not match.
				// Verify the attachment is correct using an ImportStateCheck function instead.
				ImportStateVerify: false,
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					if len(s) != 1 {
						return fmt.Errorf("expected 1 state: %+v", s)
					}

					rs := s[0]
					flipID := rs.Attributes["ip_address"]
					vmID := rs.Attributes["vm_id"]

					client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()
					foundFloatingIP, _, err := client.FloatingIPs.Get(context.Background(), flipID)
					if err != nil {
						return err
					}

					if foundFloatingIP.IP != flipID || foundFloatingIP.Vm.ID != vmID {
						return fmt.Errorf("wrong floating IP attachment found")
					}

					return nil
				},
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateId:     "192.0.2.1",
				ExpectError:       regexp.MustCompile("joined with a comma"),
			},
		},
	})
}

func testAccFLIPAssignmentImportID(n string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return "", fmt.Errorf("Not found: %s", n)
		}

		flip := rs.Primary.Attributes["ip_address"]
		vm := rs.Primary.Attributes["vm_id"]

		return fmt.Sprintf("%s,%s", flip, vm), nil
	}
}
