package kubernetes_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAbrhaKubernetesNodePool_Import(t *testing.T) {
	testName1 := acceptance.RandomTestName()
	testName2 := acceptance.RandomTestName()

	config := fmt.Sprintf(`%s

resource "abrha_kubernetes_cluster" "foobar" {
  name    = "%s"
  region  = "lon1"
  version = data.abrha_kubernetes_versions.test.latest_version

  node_pool {
    name       = "default"
    size       = "s-1vcpu-2gb"
    node_count = 1
  }
}

resource "abrha_kubernetes_node_pool" "barfoo" {
  cluster_id = abrha_kubernetes_cluster.foobar.id
  name       = "%s"
  size       = "s-1vcpu-2gb"
  node_count = 1
}
`, testClusterVersionLatest, testName1, testName2)
	resourceName := "abrha_kubernetes_node_pool.barfoo"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateId:     "this-is-not-a-valid-ID",
				ExpectError:       regexp.MustCompile("Did not find the cluster owning the node pool"),
			},
		},
	})
}
