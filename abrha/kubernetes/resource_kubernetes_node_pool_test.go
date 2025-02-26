package kubernetes_test

import (
	"context"
	"fmt"
	"testing"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAbrhaKubernetesNodePool_Basic(t *testing.T) {
	rName := acceptance.RandomTestName()
	var k8s goApiAbrha.KubernetesCluster
	var k8sPool goApiAbrha.KubernetesNodePool

	clusterConfig := fmt.Sprintf(`%s
resource "abrha_kubernetes_cluster" "foobar" {
  name    = "%s"
  region  = "lon1"
  version = data.abrha_kubernetes_versions.test.latest_version
  tags    = ["foo", "bar"]

  node_pool {
    name       = "default"
    size       = "s-1vcpu-2gb"
    node_count = 1
    tags       = ["one", "two"]
  }
}
`, testClusterVersionLatest, rName)

	nodePoolConfig := fmt.Sprintf(`resource abrha_kubernetes_node_pool "barfoo" {
	cluster_id = abrha_kubernetes_cluster.foobar.id

	name    = "%s"
	size  = "s-1vcpu-2gb"
	node_count = 1
	tags  = ["three","four"]
}
`, rName)

	nodePoolAddTaintConfig := fmt.Sprintf(`resource abrha_kubernetes_node_pool "barfoo" {
	cluster_id = abrha_kubernetes_cluster.foobar.id

	name       = "%s"
	size       = "s-1vcpu-2gb"
	node_count = 1
	tags       = ["three","four"]
	taint {
		key    = "k1"
		value  = "v1"
		effect = "NoExecute"
	}
}
`, rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: clusterConfig + nodePoolConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaKubernetesClusterExists("abrha_kubernetes_cluster.foobar", &k8s),
					testAccCheckAbrhaKubernetesNodePoolExists("abrha_kubernetes_node_pool.barfoo", &k8s, &k8sPool),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "name", rName),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "taint.#", "0"),
				),
			},
			// Update: add taint
			{
				Config: clusterConfig + nodePoolAddTaintConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "taint.#", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "taint.0.effect", "NoExecute"),
				),
			},
			// Update: remove all taints (ensure all taints are removed from resource)
			{
				Config: clusterConfig + nodePoolConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "taint.#", "0"),
				),
			},
		},
	})
}

func TestAccAbrhaKubernetesNodePool_Update(t *testing.T) {
	rName := acceptance.RandomTestName()
	var k8s goApiAbrha.KubernetesCluster
	var k8sPool goApiAbrha.KubernetesNodePool

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAbrhaKubernetesConfigBasicWithNodePool(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaKubernetesClusterExists("abrha_kubernetes_cluster.foobar", &k8s),
					testAccCheckAbrhaKubernetesNodePoolExists("abrha_kubernetes_node_pool.barfoo", &k8s, &k8sPool),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "name", rName),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "tags.#", "2"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "labels.%", "0"),
					resource.TestCheckNoResourceAttr("abrha_kubernetes_node_pool.barfoo", "labels.priority"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "node_count", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "actual_node_count", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "nodes.#", "1"),
				),
			},
			{
				Config: testAccAbrhaKubernetesConfigBasicWithNodePool2(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaKubernetesClusterExists("abrha_kubernetes_cluster.foobar", &k8s),
					testAccCheckAbrhaKubernetesNodePoolExists("abrha_kubernetes_node_pool.barfoo", &k8s, &k8sPool),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "name", rName+"-updated"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "tags.#", "3"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "labels.%", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "labels.priority", "high"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "node_count", "2"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "actual_node_count", "2"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "nodes.#", "2"),
				),
			},
			// Update NodePool Taint
			{
				Config: testAccAbrhaKubernetesConfigBasicWithNodePoolTaint(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaKubernetesNodePoolExists("abrha_kubernetes_node_pool.barfoo", &k8s, &k8sPool),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "name", rName+"-tainted"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "nodes.#", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "taint.#", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "taint.0.effect", "NoSchedule"),
				),
			},
			// Add second NodePool Taint
			{
				Config: testAccAbrhaKubernetesConfigBasicWithNodePoolTaint2(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaKubernetesNodePoolExists("abrha_kubernetes_node_pool.barfoo", &k8s, &k8sPool),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "name", rName+"-tainted"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "nodes.#", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "taint.#", "2"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "taint.0.effect", "NoSchedule"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "taint.1.effect", "PreferNoSchedule"),
				),
			},
		},
	})
}

func TestAccAbrhaKubernetesNodePool_CreateWithAutoScale(t *testing.T) {
	rName := acceptance.RandomTestName()
	var k8s goApiAbrha.KubernetesCluster
	var k8sPool goApiAbrha.KubernetesNodePool

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			// Create without auto-scaling.
			{
				Config: fmt.Sprintf(`%s

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

resource abrha_kubernetes_node_pool "barfoo" {
  cluster_id = "${abrha_kubernetes_cluster.foobar.id}"
  name       = "%s"
  size       = "s-1vcpu-2gb"
  node_count = 1
  auto_scale = true
  min_nodes  = 1
  max_nodes  = 5
}
				`, testClusterVersionLatest, rName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaKubernetesClusterExists("abrha_kubernetes_cluster.foobar", &k8s),
					testAccCheckAbrhaKubernetesNodePoolExists("abrha_kubernetes_node_pool.barfoo", &k8s, &k8sPool),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "name", rName),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "node_count", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "actual_node_count", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "nodes.#", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "auto_scale", "true"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "min_nodes", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "max_nodes", "5"),
				),
			},
			// Remove node count, keep auto-scaling.
			{
				Config: fmt.Sprintf(`%s

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

resource abrha_kubernetes_node_pool "barfoo" {
  cluster_id = "${abrha_kubernetes_cluster.foobar.id}"
  name       = "%s"
  size       = "s-1vcpu-2gb"
  auto_scale = true
  min_nodes  = 1
  max_nodes  = 3
}
				`, testClusterVersionLatest, rName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaKubernetesClusterExists("abrha_kubernetes_cluster.foobar", &k8s),
					testAccCheckAbrhaKubernetesNodePoolExists("abrha_kubernetes_node_pool.barfoo", &k8s, &k8sPool),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "name", rName),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "node_count", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "actual_node_count", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "nodes.#", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "auto_scale", "true"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "min_nodes", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "max_nodes", "3"),
				),
			},
			// Update node count, keep auto-scaling.
			{
				Config: fmt.Sprintf(`%s

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

resource abrha_kubernetes_node_pool "barfoo" {
  cluster_id = "${abrha_kubernetes_cluster.foobar.id}"
  name       = "%s"
  size       = "s-1vcpu-2gb"
  node_count = 2
  auto_scale = true
  min_nodes  = 1
  max_nodes  = 3
}
				`, testClusterVersionLatest, rName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaKubernetesClusterExists("abrha_kubernetes_cluster.foobar", &k8s),
					testAccCheckAbrhaKubernetesNodePoolExists("abrha_kubernetes_node_pool.barfoo", &k8s, &k8sPool),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "name", rName),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "node_count", "2"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "actual_node_count", "2"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "nodes.#", "2"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "auto_scale", "true"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "min_nodes", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "max_nodes", "3"),
				),
			},
			// Disable auto-scaling.
			{
				Config: fmt.Sprintf(`%s

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

resource abrha_kubernetes_node_pool "barfoo" {
  cluster_id = "${abrha_kubernetes_cluster.foobar.id}"
  name       = "%s"
  size       = "s-1vcpu-2gb"
  node_count = 2
}
				`, testClusterVersionLatest, rName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaKubernetesClusterExists("abrha_kubernetes_cluster.foobar", &k8s),
					testAccCheckAbrhaKubernetesNodePoolExists("abrha_kubernetes_node_pool.barfoo", &k8s, &k8sPool),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "name", rName),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "node_count", "2"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "actual_node_count", "2"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "nodes.#", "2"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "auto_scale", "false"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "min_nodes", "0"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "max_nodes", "0"),
				),
			},
		},
	})
}

func TestAccAbrhaKubernetesNodePool_UpdateWithAutoScale(t *testing.T) {
	rName := acceptance.RandomTestName()
	var k8s goApiAbrha.KubernetesCluster
	var k8sPool goApiAbrha.KubernetesNodePool

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			// Create without auto-scaling.
			{
				Config: fmt.Sprintf(`%s

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

resource abrha_kubernetes_node_pool "barfoo" {
  cluster_id = "${abrha_kubernetes_cluster.foobar.id}"
  name       = "%s"
  size       = "s-1vcpu-2gb"
  node_count = 1
}
				`, testClusterVersionLatest, rName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaKubernetesClusterExists("abrha_kubernetes_cluster.foobar", &k8s),
					testAccCheckAbrhaKubernetesNodePoolExists("abrha_kubernetes_node_pool.barfoo", &k8s, &k8sPool),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "name", rName),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "node_count", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "actual_node_count", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "nodes.#", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "auto_scale", "false"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "min_nodes", "0"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "max_nodes", "0"),
				),
			},
			// Update to enable auto-scaling.
			{
				Config: fmt.Sprintf(`%s

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

resource abrha_kubernetes_node_pool "barfoo" {
  cluster_id = "${abrha_kubernetes_cluster.foobar.id}"
  name       = "%s"
  size       = "s-1vcpu-2gb"
  node_count = 1
  auto_scale = true
  min_nodes  = 1
  max_nodes  = 3
}
				`, testClusterVersionLatest, rName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaKubernetesClusterExists("abrha_kubernetes_cluster.foobar", &k8s),
					testAccCheckAbrhaKubernetesNodePoolExists("abrha_kubernetes_node_pool.barfoo", &k8s, &k8sPool),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "name", rName),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "node_count", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "actual_node_count", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "nodes.#", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "auto_scale", "true"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "min_nodes", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "max_nodes", "3"),
				),
			},
			// Remove node count, keep auto-scaling.
			{
				Config: fmt.Sprintf(`%s

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

resource abrha_kubernetes_node_pool "barfoo" {
  cluster_id = "${abrha_kubernetes_cluster.foobar.id}"
  name       = "%s"
  size       = "s-1vcpu-2gb"
  auto_scale = true
  min_nodes  = 1
  max_nodes  = 3
}
				`, testClusterVersionLatest, rName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaKubernetesClusterExists("abrha_kubernetes_cluster.foobar", &k8s),
					testAccCheckAbrhaKubernetesNodePoolExists("abrha_kubernetes_node_pool.barfoo", &k8s, &k8sPool),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "name", rName),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "node_count", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "actual_node_count", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "nodes.#", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "auto_scale", "true"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "min_nodes", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_node_pool.barfoo", "max_nodes", "3"),
				),
			},
		},
	})
}

func testAccAbrhaKubernetesConfigBasicWithNodePool(rName string) string {
	return fmt.Sprintf(`%s

resource "abrha_kubernetes_cluster" "foobar" {
  name    = "%s"
  region  = "lon1"
  version = data.abrha_kubernetes_versions.test.latest_version
  tags    = ["foo", "bar"]

  node_pool {
    name       = "default"
    size       = "s-1vcpu-2gb"
    node_count = 1
    tags       = ["one", "two"]
  }
}

resource abrha_kubernetes_node_pool "barfoo" {
  cluster_id = "${abrha_kubernetes_cluster.foobar.id}"

  name       = "%s"
  size       = "s-1vcpu-2gb"
  node_count = 1
  tags       = ["three", "four"]
}
`, testClusterVersionLatest, rName, rName)
}

func testAccAbrhaKubernetesConfigBasicWithNodePoolTaint(rName string) string {
	return fmt.Sprintf(`%s

resource "abrha_kubernetes_cluster" "foobar" {
  name    = "%s"
  region  = "lon1"
  version = data.abrha_kubernetes_versions.test.latest_version
  tags    = ["foo", "bar"]

  node_pool {
    name       = "default"
    size       = "s-1vcpu-2gb"
    node_count = 1
    tags       = ["one", "two"]
  }
}

resource abrha_kubernetes_node_pool "barfoo" {
  cluster_id = abrha_kubernetes_cluster.foobar.id

  name       = "%s-tainted"
  size       = "s-1vcpu-2gb"
  node_count = 1
  tags       = ["three", "four"]
  labels = {
    priority = "high"
  }
  taint {
    key    = "key1"
    value  = "val1"
    effect = "NoSchedule"
  }
}
`, testClusterVersionLatest, rName, rName)
}

func testAccAbrhaKubernetesConfigBasicWithNodePoolTaint2(rName string) string {
	return fmt.Sprintf(`%s

resource "abrha_kubernetes_cluster" "foobar" {
  name    = "%s"
  region  = "lon1"
  version = data.abrha_kubernetes_versions.test.latest_version
  tags    = ["foo", "bar"]

  node_pool {
    name       = "default"
    size       = "s-1vcpu-2gb"
    node_count = 1
    tags       = ["one", "two"]
  }
}

resource abrha_kubernetes_node_pool "barfoo" {
  cluster_id = abrha_kubernetes_cluster.foobar.id

  name       = "%s-tainted"
  size       = "s-1vcpu-2gb"
  node_count = 1
  tags       = ["three", "four"]
  labels = {
    priority = "high"
  }
  taint {
    key    = "key1"
    value  = "val1"
    effect = "NoSchedule"
  }
  taint {
    key    = "key2"
    value  = "val2"
    effect = "PreferNoSchedule"
  }
}
`, testClusterVersionLatest, rName, rName)
}

func testAccAbrhaKubernetesConfigBasicWithNodePool2(rName string) string {
	return fmt.Sprintf(`%s

resource "abrha_kubernetes_cluster" "foobar" {
  name    = "%s"
  region  = "lon1"
  version = data.abrha_kubernetes_versions.test.latest_version
  tags    = ["foo", "bar"]

  node_pool {
    name       = "default"
    size       = "s-1vcpu-2gb"
    node_count = 1
    tags       = ["one", "two"]
  }
}

resource abrha_kubernetes_node_pool "barfoo" {
  cluster_id = "${abrha_kubernetes_cluster.foobar.id}"

  name       = "%s-updated"
  size       = "s-1vcpu-2gb"
  node_count = 2
  tags       = ["one", "two", "three"]
  labels = {
    priority = "high"
  }
}
`, testClusterVersionLatest, rName, rName)
}

func testAccCheckAbrhaKubernetesNodePoolExists(n string, cluster *goApiAbrha.KubernetesCluster, pool *goApiAbrha.KubernetesNodePool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

		foundPool, _, err := client.Kubernetes.GetNodePool(context.Background(), cluster.ID, rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundPool.ID != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		*pool = *foundPool

		return nil
	}
}
