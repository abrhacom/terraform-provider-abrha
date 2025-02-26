package spaces_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAbrhaBucketPolicy_importBasic(t *testing.T) {
	resourceName := "abrha_spaces_bucket_policy.policy"
	name := acceptance.RandomTestName()

	bucketPolicy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":"*","Action":"s3:*","Resource":"*"}]}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaSpacesBucketPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAbrhaSpacesBucketPolicy(name, bucketPolicy),
			},

			{
				ResourceName:        resourceName,
				ImportState:         true,
				ImportStateVerify:   true,
				ImportStateIdPrefix: fmt.Sprintf("%s,", testAccAbrhaSpacesBucketPolicy_TestRegion),
			},
			// Test importing non-existent resource provides expected error.
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateId:     "policy",
				ExpectError:       regexp.MustCompile(`importing a Spaces bucket policy requires the format: <region>,<bucket>`),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateId:     "nyc2,",
				ExpectError:       regexp.MustCompile(`importing a Spaces bucket policy requires the format: <region>,<bucket>`),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateId:     ",policy",
				ExpectError:       regexp.MustCompile(`importing a Spaces bucket policy requires the format: <region>,<bucket>`),
			},
		},
	})
}
