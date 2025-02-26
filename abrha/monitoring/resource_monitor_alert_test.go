package monitoring_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	multipleSlackChannel = `
slack {
	channel = "production-alerts"
	url		= "https://hooks.slack.com/services/T1234567/AAAAAAAA/ZZZZZZ"
}

slack {
	channel = "abrha-cloud-alerts"
	url		= "https://hooks.slack.com/services/T2345678/BBBBBBBB/XXXXXX"
}
`

	testAccAlertPolicy = `
resource "abrha_vm" "web" {
  image  = "ubuntu-20-04-x64"
  name   = "%s"
  region = "fra1"
  size   = "s-1vcpu-1gb"
}

resource "abrha_monitor_alert" "%s" {
  alerts {
    email = ["benny@abrha.com"]
      %s
  }
  window      = "%s"
  type        = "%s"
  compare     = "GreaterThan"
  value       = 95
  entities    = [abrha_vm.web.id]
  description = "%s"
}
`

	testAccAlertPolicySlackEmailAlerts = `
resource "abrha_vm" "web" {
  image  = "ubuntu-20-04-x64"
  name   = "%s"
  region = "fra1"
  size   = "s-1vcpu-1gb"
}

resource "abrha_monitor_alert" "%s" {
  alerts {
    email = ["benny@abrha.com"]
    slack {
      channel = "production-alerts"
      url     = "https://hooks.slack.com/services/T1234567/AAAAAAAA/ZZZZZZ"
    }
  }
  window      = "5m"
  type        = "api/public/v1/insights/vm/cpu"
  compare     = "GreaterThan"
  value       = 95
  entities    = [abrha_vm.web.id]
  description = "%s"
}
`

	testAccAlertPolicyWithTag = `
resource "abrha_tag" "test" {
  name = "%s"
}

resource "abrha_vm" "web" {
  image  = "ubuntu-20-04-x64"
  name   = "%s"
  region = "fra1"
  size   = "s-1vcpu-1gb"
  tags   = [abrha_tag.test.name]
}

resource "abrha_monitor_alert" "%s" {
  alerts {
    email = ["benny@abrha.com"]
  }
  window      = "%s"
  type        = "%s"
  compare     = "GreaterThan"
  value       = 95
  tags        = [abrha_tag.test.name]
  description = "%s"
}
`

	testAccAlertPolicyAddVm = `
resource "abrha_vm" "web" {
  image  = "ubuntu-20-04-x64"
  name   = "%s"
  region = "fra1"
  size   = "s-1vcpu-1gb"
}

resource "abrha_vm" "web2" {
  image  = "ubuntu-20-04-x64"
  name   = "%s"
  region = "fra1"
  size   = "s-1vcpu-1gb"
}


resource "abrha_monitor_alert" "%s" {
  alerts {
    email = ["benny@abrha.com"]
      %s
  }
  window      = "%s"
  type        = "%s"
  compare     = "GreaterThan"
  value       = 95
  entities    = [abrha_vm.web.id, abrha_vm.web2.id]
  description = "%s"
}
`
)

func TestAccAbrhaMonitorAlert(t *testing.T) {
	var randName = acceptance.RandomTestName()
	resourceName := fmt.Sprintf("abrha_monitor_alert.%s", randName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories:         acceptance.TestAccProviderFactories,
		CheckDestroy:              testAccCheckAbrhaMonitorAlertDestroy,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccAlertPolicy, randName, randName, "", "5m", "api/public/v1/insights/vm/cpu", randName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "api/public/v1/insights/vm/cpu"),
					resource.TestCheckResourceAttr(resourceName, "compare", "GreaterThan"),
					resource.TestCheckResourceAttr(resourceName, "alerts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "alerts.0.email.0", "benny@abrha.com"),
				),
			},
		},
	})
}

func TestAccAbrhaMonitorAlertSlackEmailAlerts(t *testing.T) {
	var randName = acceptance.RandomTestName()
	resourceName := fmt.Sprintf("abrha_monitor_alert.%s", randName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories:         acceptance.TestAccProviderFactories,
		CheckDestroy:              testAccCheckAbrhaMonitorAlertDestroy,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccAlertPolicySlackEmailAlerts, randName, randName, randName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "api/public/v1/insights/vm/cpu"),
					resource.TestCheckResourceAttr(resourceName, "compare", "GreaterThan"),
					resource.TestCheckResourceAttr(resourceName, "alerts.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "alerts.0.email.0"),
					resource.TestCheckResourceAttr(resourceName, "alerts.0.slack.0.channel", "production-alerts"),
					resource.TestCheckResourceAttr(resourceName, "alerts.0.slack.0.url", "https://hooks.slack.com/services/T1234567/AAAAAAAA/ZZZZZZ"),
				),
			},
		},
	})
}

func TestAccAbrhaMonitorAlertUpdate(t *testing.T) {
	var randName = acceptance.RandomTestName()
	resourceName := fmt.Sprintf("abrha_monitor_alert.%s", randName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories:         acceptance.TestAccProviderFactories,
		CheckDestroy:              testAccCheckAbrhaMonitorAlertDestroy,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccAlertPolicy, randName, randName, "", "10m", "api/public/v1/insights/vm/cpu", randName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", randName),
					resource.TestCheckResourceAttr(resourceName, "type", "api/public/v1/insights/vm/cpu"),
					resource.TestCheckResourceAttr(resourceName, "compare", "GreaterThan"),
					resource.TestCheckResourceAttr(resourceName, "alerts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "entities.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "alerts.0.email.0", "benny@abrha.com"),
					resource.TestCheckResourceAttr(resourceName, "alerts.0.slack.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "window", "10m"),
				),
			},
			{
				Config: fmt.Sprintf(testAccAlertPolicyAddVm, randName, randName, randName, multipleSlackChannel, "5m", "api/public/v1/insights/vm/memory_utilization_percent", "Alert about memory usage"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "Alert about memory usage"),
					resource.TestCheckResourceAttr(resourceName, "type", "api/public/v1/insights/vm/memory_utilization_percent"),
					resource.TestCheckResourceAttr(resourceName, "compare", "GreaterThan"),
					resource.TestCheckResourceAttr(resourceName, "alerts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "entities.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "alerts.0.email.0", "benny@abrha.com"),
					resource.TestCheckResourceAttr(resourceName, "alerts.0.slack.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "alerts.0.slack.0.channel", "production-alerts"),
					resource.TestCheckResourceAttr(resourceName, "alerts.0.slack.0.url", "https://hooks.slack.com/services/T1234567/AAAAAAAA/ZZZZZZ"),
					resource.TestCheckResourceAttr(resourceName, "alerts.0.slack.1.channel", "abrha-cloud-alerts"),
					resource.TestCheckResourceAttr(resourceName, "alerts.0.slack.1.url", "https://hooks.slack.com/services/T2345678/BBBBBBBB/XXXXXX"),
					resource.TestCheckResourceAttr(resourceName, "window", "5m"),
				),
			},
		},
	})
}

func TestAccAbrhaMonitorAlertWithTag(t *testing.T) {
	var (
		randName = acceptance.RandomTestName()
		tagName  = acceptance.RandomTestName()
	)
	resourceName := fmt.Sprintf("abrha_monitor_alert.%s", randName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories:         acceptance.TestAccProviderFactories,
		CheckDestroy:              testAccCheckAbrhaMonitorAlertDestroy,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccAlertPolicyWithTag, tagName, randName, randName, "5m", "api/public/v1/insights/vm/cpu", randName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "api/public/v1/insights/vm/cpu"),
					resource.TestCheckResourceAttr(resourceName, "compare", "GreaterThan"),
					resource.TestCheckResourceAttr(resourceName, "alerts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "alerts.0.email.0", "benny@abrha.com"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", tagName),
				),
			},
		},
	})
}

func testAccCheckAbrhaMonitorAlertDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "abrha_monitor_alert" {
			continue
		}
		uuid := rs.Primary.Attributes["uuid"]

		// Try to find the monitor alert
		_, _, err := client.Monitoring.GetAlertPolicy(context.Background(), uuid)

		if err == nil {
			return fmt.Errorf("Monitor alert still exists")
		}
	}

	return nil
}
