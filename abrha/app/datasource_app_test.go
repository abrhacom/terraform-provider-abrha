package app_test

import (
	"fmt"
	"testing"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAbrhaApp_Basic(t *testing.T) {
	var app goApiAbrha.App
	appName := acceptance.RandomTestName()
	appCreateConfig := fmt.Sprintf(testAccCheckAbrhaAppConfig_basic, appName)
	appDataConfig := fmt.Sprintf(testAccCheckDataSourceAbrhaAppConfig, appCreateConfig)

	updatedAppCreateConfig := fmt.Sprintf(testAccCheckAbrhaAppConfig_addService, appName)
	updatedAppDataConfig := fmt.Sprintf(testAccCheckDataSourceAbrhaAppConfig, updatedAppCreateConfig)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.TestAccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: appCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaAppExists("abrha_app.foobar", &app),
				),
			},
			{
				Config: appDataConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.abrha_app.foobar", "spec.0.name", appName),
					resource.TestCheckResourceAttrPair("abrha_app.foobar", "default_ingress",
						"data.abrha_app.foobar", "default_ingress"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_app.foobar", "project_id"),
					resource.TestCheckResourceAttrPair("abrha_app.foobar", "live_url",
						"data.abrha_app.foobar", "live_url"),
					resource.TestCheckResourceAttrPair("abrha_app.foobar", "active_deployment_id",
						"data.abrha_app.foobar", "active_deployment_id"),
					resource.TestCheckResourceAttrPair("abrha_app.foobar", "urn",
						"data.abrha_app.foobar", "urn"),
					resource.TestCheckResourceAttrPair("abrha_app.foobar", "updated_at",
						"data.abrha_app.foobar", "updated_at"),
					resource.TestCheckResourceAttrPair("abrha_app.foobar", "created_at",
						"data.abrha_app.foobar", "created_at"),
					resource.TestCheckResourceAttr(
						"abrha_app.foobar", "spec.0.alert.0.rule", "DEPLOYMENT_FAILED"),
					resource.TestCheckResourceAttr(
						"data.abrha_app.foobar", "spec.0.service.0.instance_count", "1"),
					resource.TestCheckResourceAttr(
						"data.abrha_app.foobar", "spec.0.service.0.instance_size_slug", "basic-xxs"),
					resource.TestCheckResourceAttr(
						"data.abrha_app.foobar", "spec.0.ingress.0.rule.0.match.0.path.0.prefix", "/"),
					resource.TestCheckResourceAttr(
						"data.abrha_app.foobar", "spec.0.ingress.0.rule.0.component.0.name", "go-service"),
					resource.TestCheckResourceAttr(
						"data.abrha_app.foobar", "spec.0.service.0.git.0.repo_clone_url",
						"https://github.com/parspack/sample-golang.git"),
					resource.TestCheckResourceAttr(
						"data.abrha_app.foobar", "spec.0.service.0.git.0.branch", "main"),
					resource.TestCheckResourceAttr(
						"abrha_app.foobar", "spec.0.service.0.alert.0.value", "75"),
					resource.TestCheckResourceAttr(
						"abrha_app.foobar", "spec.0.service.0.alert.0.operator", "GREATER_THAN"),
					resource.TestCheckResourceAttr(
						"abrha_app.foobar", "spec.0.service.0.alert.0.window", "TEN_MINUTES"),
					resource.TestCheckResourceAttr(
						"abrha_app.foobar", "spec.0.service.0.alert.0.rule", "CPU_UTILIZATION"),
					resource.TestCheckResourceAttr(
						"abrha_app.foobar", "spec.0.service.0.log_destination.0.name", "ServiceLogs"),
					resource.TestCheckResourceAttr(
						"abrha_app.foobar", "spec.0.service.0.log_destination.0.papertrail.0.endpoint", "syslog+tls://example.com:12345"),
				),
			},
			{
				Config: updatedAppDataConfig,
			},
			{
				Config: updatedAppDataConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.abrha_app.foobar", "spec.0.service.0.name", "go-service"),
					resource.TestCheckResourceAttr(
						"data.abrha_app.foobar", "spec.0.ingress.0.rule.0.match.0.path.0.prefix", "/go"),
					resource.TestCheckResourceAttr(
						"data.abrha_app.foobar", "spec.0.ingress.0.rule.0.component.0.preserve_path_prefix", "false"),
					resource.TestCheckResourceAttr(
						"data.abrha_app.foobar", "spec.0.ingress.0.rule.0.component.0.name", "go-service"),
					resource.TestCheckResourceAttr(
						"data.abrha_app.foobar", "spec.0.service.1.name", "python-service"),
					resource.TestCheckResourceAttr(
						"data.abrha_app.foobar", "spec.0.ingress.0.rule.1.match.0.path.0.prefix", "/python"),
					resource.TestCheckResourceAttr(
						"data.abrha_app.foobar", "spec.0.ingress.0.rule.1.component.0.preserve_path_prefix", "true"),
					resource.TestCheckResourceAttr(
						"data.abrha_app.foobar", "spec.0.ingress.0.rule.1.component.0.name", "python-service"),
				),
			},
		},
	})
}

const testAccCheckDataSourceAbrhaAppConfig = `
%s

data "abrha_app" "foobar" {
  app_id = abrha_app.foobar.id
}`
