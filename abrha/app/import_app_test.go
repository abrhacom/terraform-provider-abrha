package app_test

import (
	"fmt"
	"testing"

	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAbrhaApp_importBasic(t *testing.T) {
	resourceName := "abrha_app.foobar"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: testAccCheckAbrhaAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAbrhaAppConfig_basic, acceptance.RandomTestName()),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
