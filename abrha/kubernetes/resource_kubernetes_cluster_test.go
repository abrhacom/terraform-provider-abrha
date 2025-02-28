package kubernetes_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"reflect"
	"regexp"
	"testing"
	"time"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/abrhacom/terraform-provider-abrha/abrha/kubernetes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	testClusterVersionPrevious = `data "abrha_kubernetes_versions" "latest" {
}

locals {
  previous_version = format("%s.",
    join(".", [
      split(".", data.abrha_kubernetes_versions.latest.latest_version)[0],
      tostring(parseint(split(".", data.abrha_kubernetes_versions.latest.latest_version)[1], 10) - 1)
    ])
  )
}

data "abrha_kubernetes_versions" "test" {
  version_prefix = local.previous_version
}`

	testClusterVersionLatest = `data "abrha_kubernetes_versions" "test" {
}`
)

func TestAccAbrhaKubernetesCluster_Basic(t *testing.T) {
	rName := acceptance.RandomTestName()
	var k8s goApiAbrha.KubernetesCluster
	expectedURNRegEx, _ := regexp.Compile(`do:kubernetes:[0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}`)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAbrhaKubernetesConfigBasic(testClusterVersionLatest, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaKubernetesClusterExists("abrha_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "region", "nyc1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "surge_upgrade", "true"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "ha", "false"),
					resource.TestCheckResourceAttrPair("abrha_kubernetes_cluster.foobar", "version", "data.abrha_kubernetes_versions.test", "latest_version"),
					resource.TestCheckResourceAttrSet("abrha_kubernetes_cluster.foobar", "cluster_subnet"),
					resource.TestCheckResourceAttrSet("abrha_kubernetes_cluster.foobar", "service_subnet"),
					resource.TestCheckResourceAttrSet("abrha_kubernetes_cluster.foobar", "endpoint"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "tags.#", "3"),
					resource.TestCheckTypeSetElemAttr("abrha_kubernetes_cluster.foobar", "tags.*", "foo"),
					resource.TestCheckTypeSetElemAttr("abrha_kubernetes_cluster.foobar", "tags.*", "foo"),
					resource.TestCheckResourceAttrSet("abrha_kubernetes_cluster.foobar", "status"),
					resource.TestCheckResourceAttrSet("abrha_kubernetes_cluster.foobar", "created_at"),
					resource.TestCheckResourceAttrSet("abrha_kubernetes_cluster.foobar", "updated_at"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.#", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.size", "s-1vcpu-2gb"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.node_count", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.actual_node_count", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.tags.*", "one"),
					resource.TestCheckTypeSetElemAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.tags.*", "two"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.labels.%", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.labels.priority", "high"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.nodes.#", "1"),
					resource.TestCheckResourceAttrSet("abrha_kubernetes_cluster.foobar", "node_pool.0.nodes.0.name"),
					resource.TestCheckResourceAttrSet("abrha_kubernetes_cluster.foobar", "node_pool.0.nodes.0.status"),
					resource.TestCheckResourceAttrSet("abrha_kubernetes_cluster.foobar", "node_pool.0.nodes.0.created_at"),
					resource.TestCheckResourceAttrSet("abrha_kubernetes_cluster.foobar", "node_pool.0.nodes.0.updated_at"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.taint.#", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.taint.0.effect", "PreferNoSchedule"),
					resource.TestCheckResourceAttrSet("abrha_kubernetes_cluster.foobar", "kube_config.0.raw_config"),
					resource.TestCheckResourceAttrSet("abrha_kubernetes_cluster.foobar", "kube_config.0.cluster_ca_certificate"),
					resource.TestCheckResourceAttrSet("abrha_kubernetes_cluster.foobar", "kube_config.0.host"),
					resource.TestCheckResourceAttrSet("abrha_kubernetes_cluster.foobar", "kube_config.0.token"),
					resource.TestCheckResourceAttrSet("abrha_kubernetes_cluster.foobar", "kube_config.0.expires_at"),
					resource.TestCheckResourceAttrSet("abrha_kubernetes_cluster.foobar", "vpc_uuid"),
					resource.TestCheckResourceAttrSet("abrha_kubernetes_cluster.foobar", "auto_upgrade"),
					resource.TestMatchResourceAttr("abrha_kubernetes_cluster.foobar", "urn", expectedURNRegEx),
					resource.TestCheckResourceAttrSet("abrha_kubernetes_cluster.foobar", "maintenance_policy.0.day"),
					resource.TestCheckResourceAttrSet("abrha_kubernetes_cluster.foobar", "maintenance_policy.0.start_time"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "registry_integration", "false"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "destroy_all_associated_resources", "false"),
				),
			},
			// Update: remove default node_pool taints
			{
				Config: fmt.Sprintf(`%s

resource "abrha_kubernetes_cluster" "foobar" {
  name          = "%s"
  region        = "lon1"
  version       = data.abrha_kubernetes_versions.test.latest_version
  surge_upgrade = true
  tags          = ["foo", "bar", "one"]

  node_pool {
    name       = "default"
    size       = "s-1vcpu-2gb"
    node_count = 1
    tags       = ["one", "two"]
    labels = {
      priority = "high"
    }
  }
}`, testClusterVersionLatest, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "tags.#", "3"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.#", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.tags.#", "2"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.nodes.#", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.taint.#", "0"),
				),
			},
		},
	})
}

func TestAccAbrhaKubernetesCluster_CreateWithHAControlPlane(t *testing.T) {
	rName := acceptance.RandomTestName()
	var k8s goApiAbrha.KubernetesCluster

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`%s

resource "abrha_kubernetes_cluster" "foobar" {
  name    = "%s"
  region  = "nyc1"
  ha      = true
  version = data.abrha_kubernetes_versions.test.latest_version

  node_pool {
    name       = "default"
    size       = "s-1vcpu-2gb"
    node_count = 1
  }
}
				`, testClusterVersionLatest, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaKubernetesClusterExists("abrha_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "region", "nyc1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "ha", "true"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "ipv4_address", ""),
					resource.TestCheckResourceAttrPair("abrha_kubernetes_cluster.foobar", "version", "data.abrha_kubernetes_versions.test", "latest_version"),
					resource.TestCheckResourceAttrSet("abrha_kubernetes_cluster.foobar", "status"),
					resource.TestCheckResourceAttrSet("abrha_kubernetes_cluster.foobar", "created_at"),
					resource.TestCheckResourceAttrSet("abrha_kubernetes_cluster.foobar", "updated_at"),
					resource.TestCheckResourceAttrSet("abrha_kubernetes_cluster.foobar", "endpoint"),
				),
			},
		},
	})
}

func TestAccAbrhaKubernetesCluster_CreateWithRegistry(t *testing.T) {
	var (
		rName          = acceptance.RandomTestName()
		k8s            goApiAbrha.KubernetesCluster
		registryConfig = fmt.Sprintf(`
resource "abrha_container_registry" "foobar" {
  name                   = "%s"
  region                 = "nyc3"
  subscription_tier_slug = "starter"
}`, rName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			// Create container registry
			{
				Config: registryConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("abrha_container_registry.foobar", "name", rName),
					resource.TestCheckResourceAttr("abrha_container_registry.foobar", "endpoint", "registry.abrha.com/"+rName),
					resource.TestCheckResourceAttr("abrha_container_registry.foobar", "server_url", "registry.abrha.com"),
					resource.TestCheckResourceAttr("abrha_container_registry.foobar", "subscription_tier_slug", "starter"),
					resource.TestCheckResourceAttr("abrha_container_registry.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttrSet("abrha_container_registry.foobar", "created_at"),
					resource.TestCheckResourceAttrSet("abrha_container_registry.foobar", "storage_usage_bytes"),
				),
			},
			// Create cluster with registry integration enabled
			{
				Config: fmt.Sprintf(`%s

%s

resource "abrha_kubernetes_cluster" "foobar" {
  name                 = "%s"
  region               = "nyc3"
  registry_integration = true
  version              = data.abrha_kubernetes_versions.test.latest_version

  node_pool {
    name       = "default"
    size       = "s-1vcpu-2gb"
    node_count = 1
  }
}
				`, testClusterVersionLatest, registryConfig, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaKubernetesClusterExists("abrha_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "registry_integration", "true"),
					resource.TestCheckResourceAttrPair("abrha_kubernetes_cluster.foobar", "version", "data.abrha_kubernetes_versions.test", "latest_version"),
					resource.TestCheckResourceAttrSet("abrha_kubernetes_cluster.foobar", "status"),
					resource.TestCheckResourceAttrSet("abrha_kubernetes_cluster.foobar", "created_at"),
					resource.TestCheckResourceAttrSet("abrha_kubernetes_cluster.foobar", "updated_at"),
					resource.TestCheckResourceAttrSet("abrha_kubernetes_cluster.foobar", "endpoint"),
				),
			},
			// Disable registry integration
			{
				Config: fmt.Sprintf(`%s

%s

resource "abrha_kubernetes_cluster" "foobar" {
  name    = "%s"
  region  = "nyc3"
  version = data.abrha_kubernetes_versions.test.latest_version

  node_pool {
    name       = "default"
    size       = "s-1vcpu-2gb"
    node_count = 1
  }
}
				`, testClusterVersionLatest, registryConfig, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaKubernetesClusterExists("abrha_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "registry_integration", "false"),
				),
			},
			// Re-enable registry integration
			{
				Config: fmt.Sprintf(`%s

%s

resource "abrha_kubernetes_cluster" "foobar" {
  name                 = "%s"
  region               = "nyc3"
  version              = data.abrha_kubernetes_versions.test.latest_version
  registry_integration = true

  node_pool {
    name       = "default"
    size       = "s-1vcpu-2gb"
    node_count = 1
  }
}
				`, testClusterVersionLatest, registryConfig, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaKubernetesClusterExists("abrha_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "registry_integration", "true"),
				),
			},
		},
	})
}

func TestAccAbrhaKubernetesCluster_UpdateCluster(t *testing.T) {
	rName := acceptance.RandomTestName()
	var k8s goApiAbrha.KubernetesCluster

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAbrhaKubernetesConfigBasic(testClusterVersionLatest, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaKubernetesClusterExists("abrha_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "ha", "false"),
				),
			},
			{
				Config: testAccAbrhaKubernetesConfigBasic4(testClusterVersionLatest, rName+"-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaKubernetesClusterExists("abrha_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "name", rName+"-updated"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("abrha_kubernetes_cluster.foobar", "tags.*", "one"),
					resource.TestCheckTypeSetElemAttr("abrha_kubernetes_cluster.foobar", "tags.*", "two"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.labels.%", "0"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "surge_upgrade", "true"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "ha", "true"),
				),
			},
		},
	})
}

func TestAccAbrhaKubernetesCluster_MaintenancePolicy(t *testing.T) {
	rName := acceptance.RandomTestName()
	var k8s goApiAbrha.KubernetesCluster

	policy := `
	maintenance_policy {
		day = "monday"
		start_time = "00:00"
	}
`

	updatedPolicy := `
	maintenance_policy {
		day = "any"
		start_time = "04:00"
	}
`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAbrhaKubernetesConfigMaintenancePolicy(testClusterVersionLatest, rName, policy),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaKubernetesClusterExists("abrha_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "maintenance_policy.0.day", "monday"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "maintenance_policy.0.start_time", "00:00"),
					resource.TestCheckResourceAttrSet("abrha_kubernetes_cluster.foobar", "maintenance_policy.0.duration"),
				),
			},
			{
				Config: testAccAbrhaKubernetesConfigMaintenancePolicy(testClusterVersionLatest, rName, updatedPolicy),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaKubernetesClusterExists("abrha_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "maintenance_policy.0.day", "any"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "maintenance_policy.0.start_time", "04:00"),
					resource.TestCheckResourceAttrSet("abrha_kubernetes_cluster.foobar", "maintenance_policy.0.duration"),
				),
			},
		},
	})
}

func TestAccAbrhaKubernetesCluster_UpdatePoolDetails(t *testing.T) {
	rName := acceptance.RandomTestName()
	var k8s goApiAbrha.KubernetesCluster

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAbrhaKubernetesConfigBasic(testClusterVersionLatest, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaKubernetesClusterExists("abrha_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.#", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.node_count", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.actual_node_count", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.name", "default"),
				),
			},
			{
				Config: testAccAbrhaKubernetesConfigBasic2(testClusterVersionLatest, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaKubernetesClusterExists("abrha_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.name", "default-rename"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.#", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.node_count", "2"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.actual_node_count", "2"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.tags.#", "3"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.labels.%", "2"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.labels.priority", "high"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.labels.purpose", "awesome"),
				),
			},
		},
	})
}

func TestAccAbrhaKubernetesCluster_UpdatePoolSize(t *testing.T) {
	rName := acceptance.RandomTestName()
	var k8s goApiAbrha.KubernetesCluster

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAbrhaKubernetesConfigBasic(testClusterVersionLatest, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaKubernetesClusterExists("abrha_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.#", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.node_count", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.actual_node_count", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.size", "s-1vcpu-2gb"),
				),
			},
			{
				Config: testAccAbrhaKubernetesConfigBasic3(testClusterVersionLatest, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaKubernetesClusterExists("abrha_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.#", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.node_count", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.actual_node_count", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.size", "s-2vcpu-4gb"),
				),
			},
		},
	})
}

func TestAccAbrhaKubernetesCluster_CreatePoolWithAutoScale(t *testing.T) {
	rName := acceptance.RandomTestName()
	var k8s goApiAbrha.KubernetesCluster

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			// Create with auto-scaling and explicit node_count.
			{
				Config: fmt.Sprintf(`%s

resource "abrha_kubernetes_cluster" "foobar" {
  name         = "%s"
  region       = "lon1"
  version      = data.abrha_kubernetes_versions.test.latest_version
  auto_upgrade = true

  node_pool {
    name       = "default"
    size       = "s-1vcpu-2gb"
    node_count = 1
    auto_scale = true
    min_nodes  = 1
    max_nodes  = 3
  }
  maintenance_policy {
    start_time = "05:00"
    day        = "sunday"
  }
}
				`, testClusterVersionLatest, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaKubernetesClusterExists("abrha_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.#", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.node_count", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.actual_node_count", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.size", "s-1vcpu-2gb"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.auto_scale", "true"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.min_nodes", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.max_nodes", "3"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "auto_upgrade", "true"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "maintenance_policy.0.day", "sunday"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "maintenance_policy.0.start_time", "05:00"),
				),
			},
			// Remove node_count, keep auto-scaling.
			{
				Config: fmt.Sprintf(`%s

resource "abrha_kubernetes_cluster" "foobar" {
  name    = "%s"
  region  = "lon1"
  version = data.abrha_kubernetes_versions.test.latest_version

  node_pool {
    name       = "default"
    size       = "s-1vcpu-2gb"
    auto_scale = true
    min_nodes  = 1
    max_nodes  = 3
  }
}
				`, testClusterVersionLatest, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaKubernetesClusterExists("abrha_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.#", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.node_count", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.actual_node_count", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.size", "s-1vcpu-2gb"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.auto_scale", "true"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.min_nodes", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.max_nodes", "3"),
				),
			},
			// Update node_count, keep auto-scaling.
			{
				Config: fmt.Sprintf(`%s

resource "abrha_kubernetes_cluster" "foobar" {
  name    = "%s"
  region  = "lon1"
  version = data.abrha_kubernetes_versions.test.latest_version

  node_pool {
    name       = "default"
    size       = "s-1vcpu-2gb"
    node_count = 2
    auto_scale = true
    min_nodes  = 1
    max_nodes  = 3
  }
}
				`, testClusterVersionLatest, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaKubernetesClusterExists("abrha_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.#", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.node_count", "2"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.actual_node_count", "2"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.size", "s-1vcpu-2gb"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.auto_scale", "true"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.min_nodes", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.max_nodes", "3"),
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
    node_count = 2
  }
}
				`, testClusterVersionLatest, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaKubernetesClusterExists("abrha_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.#", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.node_count", "2"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.actual_node_count", "2"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.size", "s-1vcpu-2gb"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.auto_scale", "false"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.min_nodes", "0"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.max_nodes", "0"),
				),
			},
		},
	})
}

func TestAccAbrhaKubernetesCluster_UpdatePoolWithAutoScale(t *testing.T) {
	rName := acceptance.RandomTestName()
	var k8s goApiAbrha.KubernetesCluster

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			// Create with auto-scaling disabled.
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
			`, testClusterVersionLatest, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaKubernetesClusterExists("abrha_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.#", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.node_count", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.actual_node_count", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.size", "s-1vcpu-2gb"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.auto_scale", "false"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.min_nodes", "0"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.max_nodes", "0"),
				),
			},
			// Enable auto-scaling with explicit node_count.
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
    auto_scale = true
    min_nodes  = 1
    max_nodes  = 3
  }
}
				`, testClusterVersionLatest, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaKubernetesClusterExists("abrha_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.#", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.node_count", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.actual_node_count", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.size", "s-1vcpu-2gb"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.auto_scale", "true"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.min_nodes", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.max_nodes", "3"),
				),
			},
			// Remove node_count, keep auto-scaling.
			{
				Config: fmt.Sprintf(`%s

resource "abrha_kubernetes_cluster" "foobar" {
  name    = "%s"
  region  = "lon1"
  version = data.abrha_kubernetes_versions.test.latest_version

  node_pool {
    name       = "default"
    size       = "s-1vcpu-2gb"
    auto_scale = true
    min_nodes  = 1
    max_nodes  = 3
  }
}
				`, testClusterVersionLatest, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaKubernetesClusterExists("abrha_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.#", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.node_count", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.actual_node_count", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.size", "s-1vcpu-2gb"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.auto_scale", "true"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.min_nodes", "1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "node_pool.0.max_nodes", "3"),
				),
			},
		},
	})
}

func TestAccAbrhaKubernetesCluster_KubernetesProviderInteroperability(t *testing.T) {
	rName := acceptance.RandomTestName()
	var k8s goApiAbrha.KubernetesCluster

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"kubernetes": {
				Source:            "hashicorp/kubernetes",
				VersionConstraint: "2.0.1",
			},
		},
		CheckDestroy: testAccCheckAbrhaKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAbrhaKubernetesConfig_KubernetesProviderInteroperability(testClusterVersionLatest, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaKubernetesClusterExists("abrha_kubernetes_cluster.foobar", &k8s), resource.TestCheckResourceAttrSet("abrha_kubernetes_cluster.foobar", "kube_config.0.raw_config"),
					resource.TestCheckResourceAttrSet("abrha_kubernetes_cluster.foobar", "kube_config.0.cluster_ca_certificate"),
					resource.TestCheckResourceAttrSet("abrha_kubernetes_cluster.foobar", "kube_config.0.host"),
					resource.TestCheckResourceAttrSet("abrha_kubernetes_cluster.foobar", "kube_config.0.token"),
				),
			},
		},
	})
}

func TestAccAbrhaKubernetesCluster_UpgradeVersion(t *testing.T) {
	rName := acceptance.RandomTestName()
	var k8s goApiAbrha.KubernetesCluster

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAbrhaKubernetesConfigBasic(testClusterVersionPrevious, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaKubernetesClusterExists("abrha_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttrPair("abrha_kubernetes_cluster.foobar", "version", "data.abrha_kubernetes_versions.test", "latest_version"),
				),
			},
			{
				Config: testAccAbrhaKubernetesConfigBasic(testClusterVersionLatest, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPtr("abrha_kubernetes_cluster.foobar", "id", &k8s.ID),
					testAccCheckAbrhaKubernetesClusterExists("abrha_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttrPair("abrha_kubernetes_cluster.foobar", "version", "data.abrha_kubernetes_versions.test", "latest_version"),
				),
			},
		},
	})
}

func TestAccAbrhaKubernetesCluster_DestroyAssociated(t *testing.T) {
	rName := acceptance.RandomTestName()
	var k8s goApiAbrha.KubernetesCluster

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAbrhaKubernetesConfigDestroyAssociated(testClusterVersionPrevious, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaKubernetesClusterExists("abrha_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttrPair("abrha_kubernetes_cluster.foobar", "version", "data.abrha_kubernetes_versions.test", "latest_version"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "destroy_all_associated_resources", "true"),
				),
			},
		},
	})
}

func TestAccAbrhaKubernetesCluster_VPCNative(t *testing.T) {
	rName := acceptance.RandomTestName()
	var k8s goApiAbrha.KubernetesCluster

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAbrhaKubernetesConfigVPCNative(testClusterVersionLatest, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaKubernetesClusterExists("abrha_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "region", "nyc1"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "cluster_subnet", "192.168.0.0/20"),
					resource.TestCheckResourceAttr("abrha_kubernetes_cluster.foobar", "service_subnet", "192.168.16.0/22"),
				),
			},
		},
	})
}

func testAccAbrhaKubernetesConfigBasic(testClusterVersion string, rName string) string {
	return fmt.Sprintf(`%s

resource "abrha_kubernetes_cluster" "foobar" {
  name          = "%s"
  region        = "nyc1"
  version       = data.abrha_kubernetes_versions.test.latest_version
  surge_upgrade = true
  tags          = ["foo", "bar", "one"]

  node_pool {
    name       = "default"
    size       = "s-1vcpu-2gb"
    node_count = 1
    tags       = ["one", "two"]
    labels = {
      priority = "high"
    }
    taint {
      key    = "key1"
      value  = "val1"
      effect = "PreferNoSchedule"
    }
  }
}
`, testClusterVersion, rName)
}

func testAccAbrhaKubernetesConfigMaintenancePolicy(testClusterVersion string, rName string, policy string) string {
	return fmt.Sprintf(`%s

resource "abrha_kubernetes_cluster" "foobar" {
  name          = "%s"
  region        = "lon1"
  version       = data.abrha_kubernetes_versions.test.latest_version
  surge_upgrade = true
  tags          = ["foo", "bar", "one"]

%s

  node_pool {
    name       = "default"
    size       = "s-1vcpu-2gb"
    node_count = 1
    tags       = ["one", "two"]
    labels = {
      priority = "high"
    }
    taint {
      key    = "key1"
      value  = "val1"
      effect = "PreferNoSchedule"
    }
  }
}
`, testClusterVersion, rName, policy)
}

func testAccAbrhaKubernetesConfigBasic2(testClusterVersion string, rName string) string {
	return fmt.Sprintf(`%s

resource "abrha_kubernetes_cluster" "foobar" {
  name          = "%s"
  region        = "lon1"
  version       = data.abrha_kubernetes_versions.test.latest_version
  surge_upgrade = true
  tags          = ["foo", "bar"]

  node_pool {
    name       = "default-rename"
    size       = "s-1vcpu-2gb"
    node_count = 2
    tags       = ["one", "two", "three"]
    labels = {
      priority = "high"
      purpose  = "awesome"
    }
  }
}
`, testClusterVersion, rName)
}

func testAccAbrhaKubernetesConfigBasic3(testClusterVersion string, rName string) string {
	return fmt.Sprintf(`%s

resource "abrha_kubernetes_cluster" "foobar" {
  name    = "%s"
  region  = "lon1"
  version = data.abrha_kubernetes_versions.test.latest_version
  tags    = ["foo", "bar"]

  node_pool {
    name       = "default"
    size       = "s-2vcpu-4gb"
    node_count = 1
    tags       = ["one", "two"]
  }
}
`, testClusterVersion, rName)
}

func testAccAbrhaKubernetesConfigBasic4(testClusterVersion string, rName string) string {
	return fmt.Sprintf(`%s

resource "abrha_kubernetes_cluster" "foobar" {
  name          = "%s"
  region        = "lon1"
  surge_upgrade = true
  ha            = true
  version       = data.abrha_kubernetes_versions.test.latest_version
  tags          = ["one", "two"]

  node_pool {
    name       = "default"
    size       = "s-2vcpu-4gb"
    node_count = 1
    tags       = ["foo", "bar"]
  }
}
`, testClusterVersion, rName)
}

func testAccAbrhaKubernetesConfig_KubernetesProviderInteroperability(testClusterVersion string, rName string) string {
	return fmt.Sprintf(`%s

resource "abrha_kubernetes_cluster" "foobar" {
  name    = "%s"
  region  = "lon1"
  version = data.abrha_kubernetes_versions.test.latest_version

  node_pool {
    name       = "default"
    size       = "s-2vcpu-4gb"
    node_count = 1
  }
}

provider "kubernetes" {
  host = abrha_kubernetes_cluster.foobar.endpoint
  cluster_ca_certificate = base64decode(
    abrha_kubernetes_cluster.foobar.kube_config[0].cluster_ca_certificate
  )
  token = abrha_kubernetes_cluster.foobar.kube_config[0].token
}

resource "kubernetes_namespace" "example" {
  metadata {
    name = "example-namespace"
  }
}
`, testClusterVersion, rName)
}

func testAccAbrhaKubernetesConfigDestroyAssociated(testClusterVersion string, rName string) string {
	return fmt.Sprintf(`%s

resource "abrha_kubernetes_cluster" "foobar" {
  name                             = "%s"
  region                           = "nyc1"
  version                          = data.abrha_kubernetes_versions.test.latest_version
  destroy_all_associated_resources = true

  node_pool {
    name       = "default"
    size       = "s-1vcpu-2gb"
    node_count = 1
  }
}
`, testClusterVersion, rName)
}

func testAccAbrhaKubernetesConfigVPCNative(testClusterVersion string, rName string) string {
	return fmt.Sprintf(`%s

resource "abrha_kubernetes_cluster" "foobar" {
  name           = "%s"
  region         = "nyc1"
  version        = data.abrha_kubernetes_versions.test.latest_version
  cluster_subnet = "192.168.0.0/20"
  service_subnet = "192.168.16.0/22"
  node_pool {
    name       = "default"
    size       = "s-1vcpu-2gb"
    node_count = 1
  }
}
`, testClusterVersion, rName)
}

func testAccCheckAbrhaKubernetesClusterDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "abrha_kubernetes_cluster" {
			continue
		}

		// Try to find the cluster
		_, _, err := client.Kubernetes.Get(context.Background(), rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("K8s Cluster still exists")
		}
	}

	return nil
}

func testAccCheckAbrhaKubernetesClusterExists(n string, cluster *goApiAbrha.KubernetesCluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

		foundCluster, _, err := client.Kubernetes.Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundCluster.ID != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		*cluster = *foundCluster

		return nil
	}
}

func Test_filterTags(t *testing.T) {
	tests := []struct {
		have []string
		want []string
	}{
		{
			have: []string{"k8s", "foo"},
			want: []string{"foo"},
		},
		{
			have: []string{"k8s", "k8s:looks-like-a-uuid", "bar"},
			want: []string{"bar"},
		},
		{
			have: []string{"k8s", "k8s:looks-like-a-uuid", "bar", "k8s-this-is-ok"},
			want: []string{"bar", "k8s-this-is-ok"},
		},
		{
			have: []string{"k8s", "k8s:looks-like-a-uuid", "terraform:default-node-pool", "baz"},
			want: []string{"baz"},
		},
	}

	for _, tt := range tests {
		filteredTags := kubernetes.FilterTags(tt.have)
		if !reflect.DeepEqual(filteredTags, tt.want) {
			t.Errorf("filterTags returned %+v, expected %+v", filteredTags, tt.want)
		}
	}
}

func Test_renderKubeconfig(t *testing.T) {
	certAuth := []byte("LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURKekNDQWWlOQT09Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K")
	expected := fmt.Sprintf(`apiVersion: v1
kind: Config
clusters:
- cluster:
    certificate-authority-data: %v
    server: https://6a37a0f6-c355-4527-b54d-521beffd9817.k8s.onparspack.com
  name: do-lon1-test-cluster
contexts:
- context:
    cluster: do-lon1-test-cluster
    user: do-lon1-test-cluster-admin
  name: do-lon1-test-cluster
current-context: do-lon1-test-cluster
users:
- name: do-lon1-test-cluster-admin
  user:
    token: 97ae2bbcfd85c34155a56b822ffa73909d6770b28eb7e5dfa78fa83e02ffc60f
`, base64.StdEncoding.EncodeToString(certAuth))

	creds := goApiAbrha.KubernetesClusterCredentials{
		Server:                   "https://6a37a0f6-c355-4527-b54d-521beffd9817.k8s.onparspack.com",
		CertificateAuthorityData: certAuth,
		Token:                    "97ae2bbcfd85c34155a56b822ffa73909d6770b28eb7e5dfa78fa83e02ffc60f",
		ExpiresAt:                time.Now(),
	}
	kubeConfigRendered, err := kubernetes.RenderKubeconfig("test-cluster", "lon1", &creds)
	if err != nil {
		t.Errorf("error calling renderKubeconfig: %s", err)

	}
	got := string(kubeConfigRendered)

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("renderKubeconfig returned %+v\n, expected %+v\n", got, expected)
	}
}
