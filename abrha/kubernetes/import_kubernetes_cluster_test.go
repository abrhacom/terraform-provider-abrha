package kubernetes_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/abrhacom/terraform-provider-abrha/abrha/kubernetes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var (
	clusterStateIgnore = []string{
		"kube_config",            // because kube_config was completely different for imported state
		"node_pool.0.node_count", // because import test failed before DO had started the node in pool
		"updated_at",             // because removing default tag updates the resource outside of Terraform
		"registry_integration",   // registry_integration state can not be known via the API
	}
)

func TestAccAbrhaKubernetesCluster_ImportBasic(t *testing.T) {
	clusterName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAbrhaKubernetesConfigBasic(testClusterVersionLatest, clusterName),
				// Remove the default node pool tag so that the import code which infers
				// the need to add the tag gets triggered.
				Check: testAccAbrhaKubernetesRemoveDefaultNodePoolTag(clusterName),
			},
			{
				ResourceName:            "abrha_kubernetes_cluster.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: clusterStateIgnore,
			},
		},
	})
}

func TestAccAbrhaKubernetesCluster_ImportErrorNonDefaultNodePool(t *testing.T) {
	testName1 := acceptance.RandomTestName()
	testName2 := acceptance.RandomTestName()

	config := fmt.Sprintf(testAccAbrhaKubernetesCusterWithMultipleNodePools, testClusterVersionLatest, testName1, testName2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				// Remove the default node pool tag before importing in order to
				// trigger the multiple node pool import error.
				Check: testAccAbrhaKubernetesRemoveDefaultNodePoolTag(testName1),
			},
			{
				ResourceName:      "abrha_kubernetes_cluster.foobar",
				ImportState:       true,
				ImportStateVerify: false,
				ExpectError:       regexp.MustCompile(kubernetes.MultipleNodePoolImportError.Error()),
			},
		},
	})
}

func TestAccAbrhaKubernetesCluster_ImportNonDefaultNodePool(t *testing.T) {
	testName1 := acceptance.RandomTestName()
	testName2 := acceptance.RandomTestName()

	config := fmt.Sprintf(testAccAbrhaKubernetesCusterWithMultipleNodePools, testClusterVersionLatest, testName1, testName2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:            "abrha_kubernetes_cluster.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: clusterStateIgnore,
			},
			// Import the non-default node pool as a separate abrha_kubernetes_node_pool resource.
			{
				ResourceName:            "abrha_kubernetes_node_pool.barfoo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: clusterStateIgnore,
			},
		},
	})
}

func testAccAbrhaKubernetesRemoveDefaultNodePoolTag(clusterName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

		clusters, resp, err := client.Kubernetes.List(context.Background(), &goApiAbrha.ListOptions{})
		if err != nil {
			if resp != nil && resp.StatusCode == 404 {
				return fmt.Errorf("No clusters found")
			}

			return fmt.Errorf("Error listing Kubernetes clusters: %s", err)
		}

		var cluster *goApiAbrha.KubernetesCluster
		for _, c := range clusters {
			if c.Name == clusterName {
				cluster = c
				break
			}
		}
		if cluster == nil {
			return fmt.Errorf("Unable to find Kubernetes cluster with name: %s", clusterName)
		}

		for _, nodePool := range cluster.NodePools {
			tags := make([]string, 0)
			for _, tag := range nodePool.Tags {
				if tag != kubernetes.ParspackKubernetesDefaultNodePoolTag {
					tags = append(tags, tag)
				}
			}

			if len(tags) != len(nodePool.Tags) {
				nodePoolUpdateRequest := &goApiAbrha.KubernetesNodePoolUpdateRequest{
					Tags: tags,
				}

				_, _, err := client.Kubernetes.UpdateNodePool(context.Background(), cluster.ID, nodePool.ID, nodePoolUpdateRequest)
				if err != nil {
					return err
				}
			}
		}

		return nil
	}
}

const testAccAbrhaKubernetesCusterWithMultipleNodePools = `%s

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
`
