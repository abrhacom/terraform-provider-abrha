package project_test

import (
	"fmt"
	"testing"

	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAbrhaProject_DefaultProject(t *testing.T) {
	config := `
data "abrha_project" "default" {
}
`
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.abrha_project.default", "id"),
					resource.TestCheckResourceAttrSet("data.abrha_project.default", "name"),
					resource.TestCheckResourceAttr("data.abrha_project.default", "is_default", "true"),
				),
			},
		},
	})
}

func TestAccDataSourceAbrhaProject_NonDefaultProject(t *testing.T) {
	nonDefaultProjectName := acceptance.RandomTestName("project")
	resourceConfig := fmt.Sprintf(`
resource "abrha_project" "foo" {
  name = "%s"
}`, nonDefaultProjectName)
	dataSourceConfig := `
data "abrha_project" "bar" {
  id = abrha_project.foo.id
}

data "abrha_project" "barfoo" {
  name = abrha_project.foo.name
}
`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
			},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.abrha_project.bar", "id"),
					resource.TestCheckResourceAttr("data.abrha_project.bar", "is_default", "false"),
					resource.TestCheckResourceAttr("data.abrha_project.bar", "name", nonDefaultProjectName),
					resource.TestCheckResourceAttr("data.abrha_project.barfoo", "is_default", "false"),
					resource.TestCheckResourceAttr("data.abrha_project.barfoo", "name", nonDefaultProjectName),
				),
			},
		},
	})
}
