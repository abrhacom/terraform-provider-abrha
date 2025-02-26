package project_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/abrhacom/terraform-provider-abrha/abrha/project"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAbrhaProjectResources_Basic(t *testing.T) {
	projectName := generateProjectName()
	vmName := generateVmName()

	baseConfig := fmt.Sprintf(`
resource "abrha_project" "foo" {
  name = "%s"
}

resource "abrha_vm" "foobar" {
  name      = "%s"
  size      = "s-1vcpu-1gb"
  image     = "ubuntu-22-04-x64"
  region    = "nyc3"
  user_data = "foobar"
}
`, projectName, vmName)

	projectResourcesConfigEmpty := `
resource "abrha_project_resources" "barfoo" {
  project   = abrha_project.foo.id
  resources = []
}
`

	projectResourcesConfigWithVm := `
resource "abrha_project_resources" "barfoo" {
  project   = abrha_project.foo.id
  resources = [abrha_vm.foobar.urn]
}
`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaProjectResourcesDestroy,
		Steps: []resource.TestStep{
			{
				Config: baseConfig + projectResourcesConfigEmpty,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("abrha_project_resources.barfoo", "project"),
					resource.TestCheckResourceAttr("abrha_project_resources.barfoo", "resources.#", "0"),
					testProjectMembershipCount("abrha_project_resources.barfoo", 0),
				),
			},
			{
				// Add a resource to the abrha_project_resources.
				Config: baseConfig + projectResourcesConfigWithVm,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("abrha_project_resources.barfoo", "project"),
					resource.TestCheckResourceAttr("abrha_project_resources.barfoo", "resources.#", "1"),
					testProjectMembershipCount("abrha_project_resources.barfoo", 1),
				),
			},
			{
				// Remove the resource that was added.
				Config: baseConfig + projectResourcesConfigEmpty,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("abrha_project_resources.barfoo", "project"),
					resource.TestCheckResourceAttr("abrha_project_resources.barfoo", "resources.#", "0"),
					testProjectMembershipCount("abrha_project_resources.barfoo", 0),
				),
			},
		},
	})
}

func testProjectMembershipCount(name string, expectedCount int) resource.TestCheckFunc {
	return acceptance.TestResourceInstanceState(name, func(is *terraform.InstanceState) error {
		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

		projectId, ok := is.Attributes["project"]
		if !ok {
			return fmt.Errorf("project attribute not set")
		}

		resources, err := project.LoadResourceURNs(client, projectId)
		if err != nil {
			return fmt.Errorf("Error retrieving project resources: %s", err)
		}

		actualCount := len(*resources)

		if actualCount != expectedCount {
			return fmt.Errorf("project membership count mismatch: expected=%d, actual=%d",
				expectedCount, actualCount)
		}

		return nil
	})
}

func testAccCheckAbrhaProjectResourcesDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "abrha_project":
			_, _, err := client.Projects.Get(context.Background(), rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("Project resource still exists")
			}

		case "abrha_vm":
			id := rs.Primary.ID

			_, _, err := client.Vms.Get(context.Background(), id)
			if err == nil {
				return fmt.Errorf("Vm resource still exists")
			}
		}
	}

	return nil
}
