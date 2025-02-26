package kubernetes_test

import (
	"fmt"
	"testing"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAbrhaKubernetesVersions_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceAbrhaKubernetesVersionsConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.abrha_kubernetes_versions.foobar", "latest_version"),
				),
			},
		},
	})
}

func TestAccDataSourceAbrhaKubernetesVersions_Filtered(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceAbrhaKubernetesVersionsConfig_filtered,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.abrha_kubernetes_versions.foobar", "valid_versions.#", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceAbrhaKubernetesVersions_CreateCluster(t *testing.T) {
	rName := acceptance.RandomTestName()
	var k8s goApiAbrha.KubernetesCluster

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceAbrhaKubernetesVersionsConfig_create, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.abrha_kubernetes_versions.foobar", "latest_version"),
					testAccCheckAbrhaKubernetesClusterExists(
						"abrha_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr(
						"abrha_kubernetes_cluster.foobar", "name", rName),
				),
			},
		},
	})
}

const testAccCheckDataSourceAbrhaKubernetesVersionsConfig_basic = `
data "abrha_kubernetes_versions" "foobar" {}`

const testAccCheckDataSourceAbrhaKubernetesVersionsConfig_filtered = `
data "abrha_kubernetes_versions" "foobar" {
  version_prefix = "1.12." # No longer supported, should be empty
}`

const testAccCheckDataSourceAbrhaKubernetesVersionsConfig_create = `
data "abrha_kubernetes_versions" "foobar" {
}

resource "abrha_kubernetes_cluster" "foobar" {
  name    = "%s"
  region  = "lon1"
  version = data.abrha_kubernetes_versions.foobar.latest_version

  node_pool {
    name       = "default"
    size       = "s-1vcpu-2gb"
    node_count = 1
  }
}`
