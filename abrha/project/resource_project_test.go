package project_test

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

func TestAccAbrhaProject_CreateWithDefaults(t *testing.T) {

	expectedName := generateProjectName()
	createConfig := fixtureCreateWithDefaults(expectedName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: createConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaProjectExists("abrha_project.myproj"),
					resource.TestCheckResourceAttr(
						"abrha_project.myproj", "name", expectedName),
					resource.TestCheckResourceAttr(
						"abrha_project.myproj", "description", ""),
					resource.TestCheckResourceAttr(
						"abrha_project.myproj", "purpose", "Web Application"),
					resource.TestCheckResourceAttr(
						"abrha_project.myproj", "environment", ""),
					resource.TestCheckResourceAttrSet("abrha_project.myproj", "id"),
					resource.TestCheckResourceAttrSet("abrha_project.myproj", "owner_uuid"),
					resource.TestCheckResourceAttrSet("abrha_project.myproj", "owner_id"),
					resource.TestCheckResourceAttrSet("abrha_project.myproj", "created_at"),
					resource.TestCheckResourceAttrSet("abrha_project.myproj", "updated_at"),
				),
			},
		},
	})
}

func TestAccAbrhaProject_CreateWithIsDefault(t *testing.T) {
	expectedName := generateProjectName()
	expectedIsDefault := "true"
	createConfig := fixtureCreateWithIsDefault(expectedName, expectedIsDefault)

	var (
		originalDefaultProject = &goApiAbrha.Project{}
		client                 = &goApiAbrha.Client{}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)

			// Get an store original default project ID
			client = acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()
			defaultProject, _, defaultProjErr := client.Projects.GetDefault(context.Background())
			if defaultProjErr != nil {
				t.Errorf("Error locating default project %s", defaultProjErr)
			}
			originalDefaultProject = defaultProject
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config:             createConfig,
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaProjectExists("abrha_project.myproj"),
					// Restore original default project. This must happen here
					// to ensure it runs even if the tests fails.
					func(*terraform.State) error {
						t.Logf("Restoring original default project: %s (%s)", originalDefaultProject.Name, originalDefaultProject.ID)
						originalDefaultProject.IsDefault = true
						updateReq := &goApiAbrha.UpdateProjectRequest{
							Name:        originalDefaultProject.Name,
							Description: originalDefaultProject.Description,
							Purpose:     originalDefaultProject.Purpose,
							Environment: originalDefaultProject.Environment,
							IsDefault:   true,
						}
						_, _, err := client.Projects.Update(context.Background(), originalDefaultProject.ID, updateReq)
						if err != nil {
							return fmt.Errorf("Error restoring default project %s", err)
						}
						return nil
					},
					resource.TestCheckResourceAttr(
						"abrha_project.myproj", "name", expectedName),
					resource.TestCheckResourceAttr(
						"abrha_project.myproj", "description", ""),
					resource.TestCheckResourceAttr(
						"abrha_project.myproj", "purpose", "Web Application"),
					resource.TestCheckResourceAttr(
						"abrha_project.myproj", "environment", ""),
					resource.TestCheckResourceAttr(
						"abrha_project.myproj", "is_default", expectedIsDefault),
					resource.TestCheckResourceAttrSet("abrha_project.myproj", "id"),
					resource.TestCheckResourceAttrSet("abrha_project.myproj", "owner_uuid"),
					resource.TestCheckResourceAttrSet("abrha_project.myproj", "owner_id"),
					resource.TestCheckResourceAttrSet("abrha_project.myproj", "created_at"),
					resource.TestCheckResourceAttrSet("abrha_project.myproj", "updated_at"),
				),
			},
		},
	})
}

func TestAccAbrhaProject_CreateWithInitialValues(t *testing.T) {

	expectedName := generateProjectName()
	expectedDescription := "A simple project for a web app."
	expectedPurpose := "My Basic Web App"
	expectedEnvironment := "Production"

	createConfig := fixtureCreateWithInitialValues(expectedName, expectedDescription,
		expectedPurpose, expectedEnvironment)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: createConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaProjectExists("abrha_project.myproj"),
					resource.TestCheckResourceAttr(
						"abrha_project.myproj", "name", expectedName),
					resource.TestCheckResourceAttr(
						"abrha_project.myproj", "description", expectedDescription),
					resource.TestCheckResourceAttr(
						"abrha_project.myproj", "purpose", expectedPurpose),
					resource.TestCheckResourceAttr(
						"abrha_project.myproj", "environment", expectedEnvironment),
				),
			},
		},
	})
}

func TestAccAbrhaProject_UpdateWithInitialValues(t *testing.T) {

	expectedName := generateProjectName()
	expectedDesc := "A simple project for a web app."
	expectedPurpose := "My Basic Web App"
	expectedEnv := "Production"

	createConfig := fixtureCreateWithInitialValues(expectedName, expectedDesc,
		expectedPurpose, expectedEnv)

	expectedUpdateName := generateProjectName()
	expectedUpdateDesc := "A simple project for Beta testing."
	expectedUpdatePurpose := "MyWeb App, (Beta)"
	expectedUpdateEnv := "Staging"

	updateConfig := fixtureUpdateWithValues(expectedUpdateName, expectedUpdateDesc,
		expectedUpdatePurpose, expectedUpdateEnv)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: createConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaProjectExists("abrha_project.myproj"),
					resource.TestCheckResourceAttr(
						"abrha_project.myproj", "name", expectedName),
					resource.TestCheckResourceAttr(
						"abrha_project.myproj", "description", expectedDesc),
					resource.TestCheckResourceAttr(
						"abrha_project.myproj", "purpose", expectedPurpose),
					resource.TestCheckResourceAttr(
						"abrha_project.myproj", "environment", expectedEnv),
				),
			},
			{
				Config: updateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaProjectExists("abrha_project.myproj"),
					resource.TestCheckResourceAttr(
						"abrha_project.myproj", "name", expectedUpdateName),
					resource.TestCheckResourceAttr(
						"abrha_project.myproj", "description", expectedUpdateDesc),
					resource.TestCheckResourceAttr(
						"abrha_project.myproj", "purpose", expectedUpdatePurpose),
					resource.TestCheckResourceAttr(
						"abrha_project.myproj", "environment", expectedUpdateEnv),
				),
			},
		},
	})
}

func TestAccAbrhaProject_CreateWithVmResource(t *testing.T) {

	expectedName := generateProjectName()
	expectedVmName := generateVmName()

	createConfig := fixtureCreateWithVmResource(expectedVmName, expectedName)
	destroyConfig := fixtureCreateWithDefaults(expectedName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: createConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaProjectExists("abrha_project.myproj"),
					resource.TestCheckResourceAttr(
						"abrha_project.myproj", "name", expectedName),
					resource.TestCheckResourceAttr("abrha_project.myproj", "resources.#", "1"),
				),
			},
			{
				Config: destroyConfig,
			},
			{
				Config: destroyConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaProjectExists("abrha_project.myproj"),
					resource.TestCheckResourceAttr(
						"abrha_project.myproj", "name", expectedName),
					resource.TestCheckResourceAttr("abrha_project.myproj", "resources.#", "0"),
				),
			},
		},
	})
}

func TestAccAbrhaProject_CreateWithUnacceptedResourceExpectError(t *testing.T) {

	expectedName := generateProjectName()
	vpcName := acceptance.RandomTestName()
	vpcDesc := "A description for the VPC"

	createConfig := fixtureCreateWithUnacceptedResource(vpcName, vpcDesc, expectedName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config:      createConfig,
				ExpectError: regexp.MustCompile(`Error creating Project: Error assigning resources`),
			},
		},
	})
}

func TestAccAbrhaProject_UpdateWithVmResource(t *testing.T) {

	expectedName := generateProjectName()
	expectedVmName := generateVmName()

	createConfig := fixtureCreateWithVmResource(expectedVmName, expectedName)

	updateConfig := fixtureCreateWithDefaults(expectedName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: createConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaProjectExists("abrha_project.myproj"),
					resource.TestCheckResourceAttr(
						"abrha_project.myproj", "name", expectedName),
					resource.TestCheckResourceAttr("abrha_project.myproj", "resources.#", "1"),
				),
			},
			{
				Config: updateConfig,
			},
			{
				Config: updateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaProjectExists("abrha_project.myproj"),
					resource.TestCheckResourceAttr(
						"abrha_project.myproj", "name", expectedName),
					resource.TestCheckResourceAttr("abrha_project.myproj", "resources.#", "0"),
				),
			},
		},
	})
}

func TestAccAbrhaProject_UpdateFromVmToSpacesResource(t *testing.T) {
	expectedName := generateProjectName()
	expectedVmName := generateVmName()
	expectedSpacesName := generateSpacesName()

	createConfig := fixtureCreateWithVmResource(expectedVmName, expectedName)
	updateConfig := fixtureCreateWithSpacesResource(expectedSpacesName, expectedName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: createConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaProjectExists("abrha_project.myproj"),
					resource.TestCheckResourceAttr(
						"abrha_project.myproj", "name", expectedName),
					resource.TestCheckResourceAttr("abrha_project.myproj", "resources.#", "1"),
					resource.TestCheckResourceAttrSet("abrha_vm.foobar", "urn"),
				),
			},
			{
				Config: updateConfig,
			},
			{
				Config: updateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaProjectExists("abrha_project.myproj"),
					testAccCheckAbrhaProjectResourceURNIsPresent("abrha_project.myproj", "do:spaces:"+generateSpacesName()),
					resource.TestCheckResourceAttr(
						"abrha_project.myproj", "name", expectedName),
					resource.TestCheckResourceAttr("abrha_project.myproj", "resources.#", "1"),
					resource.TestCheckResourceAttrSet("abrha_spaces_bucket.foobar", "urn"),
				),
			},
		},
	})
}

func TestAccAbrhaProject_WithManyResources(t *testing.T) {
	projectName := generateProjectName()
	domainBase := acceptance.RandomTestName("project")

	createConfig := fixtureCreateDomainResources(domainBase)
	updateConfig := fixtureWithManyResources(domainBase, projectName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: createConfig,
			},
			{
				Config: updateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaProjectExists("abrha_project.myproj"),
					resource.TestCheckResourceAttr(
						"abrha_project.myproj", "name", projectName),
					resource.TestCheckResourceAttr("abrha_project.myproj", "resources.#", "10"),
				),
			},
		},
	})
}

func testAccCheckAbrhaProjectResourceURNIsPresent(resource, expectedURN string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

		rs, ok := s.RootModule().Resources[resource]

		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID set for resource: %s", resource)
		}

		projectResources, _, err := client.Projects.ListResources(context.Background(), rs.Primary.ID, nil)
		if err != nil {
			return fmt.Errorf("Error Retrieving project resources to confirm.")
		}

		for _, v := range projectResources {

			if v.URN == expectedURN {
				return nil
			}

		}

		return nil
	}
}

func testAccCheckAbrhaProjectDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "abrha_project" {
			continue
		}

		_, _, err := client.Projects.Get(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Project resource still exists")
		}
	}

	return nil
}

func testAccCheckAbrhaProjectExists(resource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

		rs, ok := s.RootModule().Resources[resource]

		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID set for resource: %s", resource)
		}

		foundProject, _, err := client.Projects.Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundProject.ID != rs.Primary.ID {
			return fmt.Errorf("Resource not found: %s : %s", resource, rs.Primary.ID)
		}

		return nil
	}
}

func generateProjectName() string {
	return acceptance.RandomTestName("project")
}

func generateVmName() string {
	return acceptance.RandomTestName("vm")
}

func generateSpacesName() string {
	return acceptance.RandomTestName("space")
}

func fixtureCreateWithDefaults(name string) string {
	return fmt.Sprintf(`
resource "abrha_project" "myproj" {
  name = "%s"
}`, name)
}

func fixtureUpdateWithValues(name, description, purpose, environment string) string {
	return fixtureCreateWithInitialValues(name, description, purpose, environment)
}

func fixtureCreateWithInitialValues(name, description, purpose, environment string) string {
	return fmt.Sprintf(`
resource "abrha_project" "myproj" {
  name        = "%s"
  description = "%s"
  purpose     = "%s"
  environment = "%s"
}`, name, description, purpose, environment)
}

func fixtureCreateWithVmResource(vmName, name string) string {
	return fmt.Sprintf(`
resource "abrha_vm" "foobar" {
  name      = "%s"
  size      = "s-1vcpu-1gb"
  image     = "ubuntu-22-04-x64"
  region    = "nyc3"
  user_data = "foobar"
}

resource "abrha_project" "myproj" {
  name      = "%s"
  resources = [abrha_vm.foobar.urn]
}`, vmName, name)

}

func fixtureCreateWithUnacceptedResource(vpcName, vpcDesc, name string) string {
	return fmt.Sprintf(`
resource "abrha_vpc" "foobar" {
  name        = "%s"
  description = "%s"
  region      = "nyc3"
}

resource "abrha_project" "myproj" {
  name      = "%s"
  resources = [abrha_vpc.foobar.urn]
}`, vpcName, vpcDesc, name)

}

func fixtureCreateWithSpacesResource(spacesBucketName, name string) string {
	return fmt.Sprintf(`
resource "abrha_spaces_bucket" "foobar" {
  name   = "%s"
  acl    = "public-read"
  region = "ams3"
}

resource "abrha_project" "myproj" {
  name      = "%s"
  resources = [abrha_spaces_bucket.foobar.urn]
}`, spacesBucketName, name)

}

func fixtureCreateDomainResources(domainBase string) string {
	return fmt.Sprintf(`
resource "abrha_domain" "foobar" {
  count = 10
  name  = "%s-${count.index}.com"
}`, domainBase)
}

func fixtureWithManyResources(domainBase string, name string) string {
	return fmt.Sprintf(`
resource "abrha_domain" "foobar" {
  count = 10
  name  = "%s-${count.index}.com"
}

resource "abrha_project" "myproj" {
  name      = "%s"
  resources = abrha_domain.foobar[*].urn
}`, domainBase, name)
}

func fixtureCreateWithIsDefault(name string, is_default string) string {
	return fmt.Sprintf(`
resource "abrha_project" "myproj" {
  name       = "%s"
  is_default = "%s"
}`, name, is_default)
}
