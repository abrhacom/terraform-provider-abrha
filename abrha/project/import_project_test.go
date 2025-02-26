package project_test

import (
	"testing"

	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAbrhaProject_importBasic(t *testing.T) {
	name := generateProjectName()
	resourceName := "abrha_project.myproj"
	createConfig := fixtureCreateWithDefaults(name)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: createConfig,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
