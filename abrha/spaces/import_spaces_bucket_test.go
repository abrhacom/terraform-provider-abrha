package spaces_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAbrhaBucket_importBasic(t *testing.T) {
	resourceName := "abrha_spaces_bucket.bucket"
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAbrhaBucketConfigImport(name),
			},

			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdPrefix:     fmt.Sprintf("%s,", "sfo3"),
				ImportStateVerifyIgnore: []string{"acl"},
			},
			// Test importing non-existent resource provides expected error.
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateId:     "sfo3,nonexistent-bucket",
				ExpectError:       regexp.MustCompile(`(Please verify the ID is correct|Cannot import non-existent remote object)`),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateId:     "bucket",
				ExpectError:       regexp.MustCompile(`importing a Spaces bucket requires the format: <region>,<name>`),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateId:     "nyc2,",
				ExpectError:       regexp.MustCompile(`importing a Spaces bucket requires the format: <region>,<name>`),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateId:     ",bucket",
				ExpectError:       regexp.MustCompile(`importing a Spaces bucket requires the format: <region>,<name>`),
			},
		},
	})
}
