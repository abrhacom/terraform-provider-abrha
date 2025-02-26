package project_test

import (
	"fmt"
	"testing"

	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAbrhaProjects_Basic(t *testing.T) {
	prodProjectName := acceptance.RandomTestName("project")
	stagingProjectName := acceptance.RandomTestName("project")

	resourcesConfig := fmt.Sprintf(`
resource "abrha_project" "prod" {
  name        = "%s"
  environment = "Production"
}

resource "abrha_project" "staging" {
  name        = "%s"
  environment = "Staging"
}
`, prodProjectName, stagingProjectName)

	datasourcesConfig := fmt.Sprintf(`
data "abrha_projects" "prod" {
  filter {
    key    = "environment"
    values = ["Production"]
  }
  filter {
    key    = "is_default"
    values = ["false"]
  }
}

data "abrha_projects" "staging" {
  filter {
    key    = "name"
    values = ["%s"]
  }
  filter {
    key    = "is_default"
    values = ["false"]
  }
}

data "abrha_projects" "both" {
  filter {
    key    = "environment"
    values = ["Production"]
  }
  filter {
    key    = "name"
    values = ["%s"]
  }
}
`, stagingProjectName, stagingProjectName)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourcesConfig,
			},
			{
				Config: resourcesConfig + datasourcesConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.abrha_projects.prod", "projects.#", "1"),
					resource.TestCheckResourceAttr("data.abrha_projects.prod", "projects.0.name", prodProjectName),
					resource.TestCheckResourceAttr("data.abrha_projects.prod", "projects.0.environment", "Production"),
					resource.TestCheckResourceAttr("data.abrha_projects.staging", "projects.#", "1"),
					resource.TestCheckResourceAttr("data.abrha_projects.staging", "projects.0.name", stagingProjectName),
					resource.TestCheckResourceAttr("data.abrha_projects.staging", "projects.0.environment", "Staging"),
					resource.TestCheckResourceAttr("data.abrha_projects.both", "projects.#", "0"),
				),
			},
		},
	})
}
