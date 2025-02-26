package kubernetes_test

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

func TestAccDataSourceAbrhaKubernetesCluster_Basic(t *testing.T) {
	rName := acceptance.RandomTestName()
	var k8s goApiAbrha.KubernetesCluster
	expectedURNRegEx, _ := regexp.Compile(`do:kubernetes:[0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}`)
	resourceConfig := testAccAbrhaKubernetesConfigForDataSource(testClusterVersionLatest, rName)
	dataSourceConfig := `
data "abrha_kubernetes_cluster" "foobar" {
  name = abrha_kubernetes_cluster.foo.name
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"kubernetes": {
				Source:            "hashicorp/kubernetes",
				VersionConstraint: "1.13.2",
			},
		},
		CheckDestroy: testAccCheckAbrhaKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
			},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceAbrhaKubernetesClusterExists("data.abrha_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("data.abrha_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("data.abrha_kubernetes_cluster.foobar", "region", "lon1"),
					resource.TestCheckResourceAttrPair("data.abrha_kubernetes_cluster.foobar", "version", "data.abrha_kubernetes_versions.test", "latest_version"),
					resource.TestCheckResourceAttr("data.abrha_kubernetes_cluster.foobar", "node_pool.0.labels.priority", "high"),
					resource.TestCheckResourceAttrSet("data.abrha_kubernetes_cluster.foobar", "vpc_uuid"),
					resource.TestCheckResourceAttrSet("data.abrha_kubernetes_cluster.foobar", "auto_upgrade"),
					resource.TestMatchResourceAttr("data.abrha_kubernetes_cluster.foobar", "urn", expectedURNRegEx),
					resource.TestCheckResourceAttr("data.abrha_kubernetes_cluster.foobar", "maintenance_policy.0.day", "monday"),
					resource.TestCheckResourceAttr("data.abrha_kubernetes_cluster.foobar", "maintenance_policy.0.start_time", "00:00"),
					resource.TestCheckResourceAttrSet("data.abrha_kubernetes_cluster.foobar", "maintenance_policy.0.duration"),
				),
			},
		},
	})
}

func testAccAbrhaKubernetesConfigForDataSource(version string, rName string) string {
	return fmt.Sprintf(`%s

resource "abrha_kubernetes_cluster" "foo" {
  name         = "%s"
  region       = "lon1"
  version      = data.abrha_kubernetes_versions.test.latest_version
  tags         = ["foo", "bar"]
  auto_upgrade = true

  node_pool {
    name       = "default"
    size       = "s-1vcpu-2gb"
    node_count = 1
    tags       = ["one", "two"]
    labels = {
      priority = "high"
    }
  }
  maintenance_policy {
    day        = "monday"
    start_time = "00:00"
  }
}`, version, rName)
}

func testAccCheckDataSourceAbrhaKubernetesClusterExists(n string, cluster *goApiAbrha.KubernetesCluster) resource.TestCheckFunc {
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
