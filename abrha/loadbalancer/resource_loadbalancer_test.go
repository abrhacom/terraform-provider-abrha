package loadbalancer_test

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/abrhacom/terraform-provider-abrha/abrha/loadbalancer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAbrhaLoadbalancer_Basic(t *testing.T) {
	var loadbalancer goApiAbrha.LoadBalancer
	name := acceptance.RandomTestName()

	expectedURNRegEx, _ := regexp.Compile(`do:loadbalancer:[0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}`)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaLoadbalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaLoadbalancerConfig_basic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaLoadbalancerExists("abrha_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "type", "REGIONAL"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "size_unit", "1"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "forwarding_rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"abrha_loadbalancer.foobar",
						"forwarding_rule.*",
						map[string]string{
							"entry_port":      "80",
							"entry_protocol":  "http",
							"target_port":     "80",
							"target_protocol": "http",
						},
					),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "healthcheck.#", "1"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "healthcheck.0.port", "22"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "healthcheck.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "vm_ids.#", "1"),
					resource.TestCheckResourceAttrSet(
						"abrha_loadbalancer.foobar", "vpc_uuid"),
					resource.TestMatchResourceAttr(
						"abrha_loadbalancer.foobar", "urn", expectedURNRegEx),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "enable_proxy_protocol", "true"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "enable_backend_keepalive", "true"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "disable_lets_encrypt_dns_records", "false"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "http_idle_timeout_seconds", "90"),
					resource.TestCheckResourceAttrSet(
						"abrha_loadbalancer.foobar", "project_id"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "network", "EXTERNAL"),
				),
			},
		},
	})
}

func TestAccAbrhaLoadbalancer_Updated(t *testing.T) {
	var loadbalancer goApiAbrha.LoadBalancer
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaLoadbalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaLoadbalancerConfig_basic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaLoadbalancerExists("abrha_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "type", "REGIONAL"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "size_unit", "1"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "forwarding_rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"abrha_loadbalancer.foobar",
						"forwarding_rule.*",
						map[string]string{
							"entry_port":      "80",
							"entry_protocol":  "http",
							"target_port":     "80",
							"target_protocol": "http",
							"tls_passthrough": "false",
						},
					),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "healthcheck.#", "1"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "healthcheck.0.port", "22"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "healthcheck.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "vm_ids.#", "1"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "enable_proxy_protocol", "true"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "enable_backend_keepalive", "true"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "disable_lets_encrypt_dns_records", "false"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "http_idle_timeout_seconds", "90"),
					resource.TestCheckResourceAttrSet(
						"abrha_loadbalancer.foobar", "project_id"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "network", "EXTERNAL"),
				),
			},
			{
				Config: testAccCheckAbrhaLoadbalancerConfig_updated(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaLoadbalancerExists("abrha_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "type", "REGIONAL"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "size_unit", "1"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "forwarding_rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"abrha_loadbalancer.foobar",
						"forwarding_rule.*",
						map[string]string{
							"entry_port":      "81",
							"entry_protocol":  "http",
							"target_port":     "81",
							"target_protocol": "http",
							"tls_passthrough": "false",
						},
					),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "healthcheck.#", "1"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "healthcheck.0.port", "22"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "healthcheck.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "vm_ids.#", "2"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "enable_proxy_protocol", "false"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "enable_backend_keepalive", "false"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "disable_lets_encrypt_dns_records", "true"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "http_idle_timeout_seconds", "120"),
					resource.TestCheckResourceAttrSet(
						"abrha_loadbalancer.foobar", "project_id"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "network", "EXTERNAL"),
				),
			},
		},
	})
}

func TestAccAbrhaLoadbalancer_vmTag(t *testing.T) {
	var loadbalancer goApiAbrha.LoadBalancer
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaLoadbalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaLoadbalancerConfig_vmTag(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaLoadbalancerExists("abrha_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "size_unit", "1"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "forwarding_rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"abrha_loadbalancer.foobar",
						"forwarding_rule.*",
						map[string]string{
							"entry_port":      "80",
							"entry_protocol":  "http",
							"target_port":     "80",
							"target_protocol": "http",
							"tls_passthrough": "false",
						},
					),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "healthcheck.#", "1"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "healthcheck.0.port", "22"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "healthcheck.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "vm_tag", "sample"),
				),
			},
		},
	})
}

func TestAccAbrhaLoadbalancer_minimal(t *testing.T) {
	var loadbalancer goApiAbrha.LoadBalancer
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaLoadbalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaLoadbalancerConfig_minimal(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaLoadbalancerExists("abrha_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "size_unit", "1"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "forwarding_rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"abrha_loadbalancer.foobar",
						"forwarding_rule.*",
						map[string]string{
							"entry_port":      "80",
							"entry_protocol":  "http",
							"target_port":     "80",
							"target_protocol": "http",
							"tls_passthrough": "false",
						},
					),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "healthcheck.#", "1"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "healthcheck.0.port", "80"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "healthcheck.0.protocol", "http"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "sticky_sessions.#", "1"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "sticky_sessions.0.type", "none"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "vm_ids.#", "1"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "enable_proxy_protocol", "false"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "enable_backend_keepalive", "false"),
					resource.TestCheckResourceAttrSet(
						"abrha_loadbalancer.foobar", "project_id"),
				),
			},
		},
	})
}

func TestAccAbrhaLoadbalancer_NonDefaultProject(t *testing.T) {
	var loadbalancer goApiAbrha.LoadBalancer
	lbName := acceptance.RandomTestName()
	projectName := acceptance.RandomTestName()

	projectConfig := `


resource "abrha_project" "test" {
  name = "%s"
}
`
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaLoadbalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(projectConfig, projectName) + testAccCheckAbrhaLoadbalancerConfig_NonDefaultProject(projectName, lbName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaLoadbalancerExists("abrha_loadbalancer.test", &loadbalancer),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.test", "name", lbName),
					resource.TestCheckResourceAttrPair(
						"abrha_loadbalancer.test", "project_id", "abrha_project.test", "id"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.test", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.test", "size_unit", "1"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.test", "forwarding_rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"abrha_loadbalancer.test",
						"forwarding_rule.*",
						map[string]string{
							"entry_port":      "80",
							"entry_protocol":  "http",
							"target_port":     "80",
							"target_protocol": "http",
							"tls_passthrough": "false",
						},
					),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.test", "healthcheck.#", "1"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.test", "healthcheck.0.port", "80"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.test", "healthcheck.0.protocol", "http"),
				),
			},
			{
				// The load balancer must be destroyed before the project which
				// discovers that asynchronously.
				Config: fmt.Sprintf(projectConfig, projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckFunc(
						func(s *terraform.State) error {
							time.Sleep(10 * time.Second)
							return nil
						},
					),
				),
			},
		},
	})
}

func TestAccAbrhaLoadbalancer_minimalUDP(t *testing.T) {
	var loadbalancer goApiAbrha.LoadBalancer
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaLoadbalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaLoadbalancerConfig_minimalUDP(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaLoadbalancerExists("abrha_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "size_unit", "1"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "forwarding_rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"abrha_loadbalancer.foobar",
						"forwarding_rule.*",
						map[string]string{
							"entry_port":      "80",
							"entry_protocol":  "udp",
							"target_port":     "80",
							"target_protocol": "udp",
							"tls_passthrough": "false",
						},
					),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "healthcheck.#", "1"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "healthcheck.0.port", "80"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "healthcheck.0.protocol", "http"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "sticky_sessions.#", "1"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "sticky_sessions.0.type", "none"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "vm_ids.#", "1"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "enable_proxy_protocol", "false"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "enable_backend_keepalive", "false"),
				),
			},
		},
	})
}

func TestAccAbrhaLoadbalancer_stickySessions(t *testing.T) {
	var loadbalancer goApiAbrha.LoadBalancer
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaLoadbalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaLoadbalancerConfig_stickySessions(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaLoadbalancerExists("abrha_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "size_unit", "1"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "forwarding_rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"abrha_loadbalancer.foobar",
						"forwarding_rule.*",
						map[string]string{
							"entry_port":      "80",
							"entry_protocol":  "http",
							"target_port":     "80",
							"target_protocol": "http",
							"tls_passthrough": "false",
						},
					),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "healthcheck.#", "1"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "healthcheck.0.port", "80"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "healthcheck.0.protocol", "http"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "sticky_sessions.#", "1"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "sticky_sessions.0.type", "cookies"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "sticky_sessions.0.cookie_name", "sessioncookie"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "sticky_sessions.0.cookie_ttl_seconds", "1800"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "vm_ids.#", "1"),
				),
			},
		},
	})
}

func TestAccAbrhaLoadbalancer_sslTermination(t *testing.T) {
	var loadbalancer goApiAbrha.LoadBalancer
	name := acceptance.RandomTestName()
	privateKeyMaterial, leafCertMaterial, certChainMaterial := acceptance.GenerateTestCertMaterial(t)
	certName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaLoadbalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaLoadbalancerConfig_sslTermination(
					certName, name, privateKeyMaterial, leafCertMaterial, certChainMaterial, "certificate_id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaLoadbalancerExists("abrha_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "size_unit", "1"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "forwarding_rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"abrha_loadbalancer.foobar",
						"forwarding_rule.*",
						map[string]string{
							"entry_port":       "443",
							"entry_protocol":   "https",
							"target_port":      "80",
							"target_protocol":  "http",
							"certificate_name": certName,
							"tls_passthrough":  "false",
						},
					),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "redirect_http_to_https", "true"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "enable_proxy_protocol", "true"),
				),
			},
		},
	})
}

func TestAccAbrhaLoadbalancer_sslCertByName(t *testing.T) {
	var loadbalancer goApiAbrha.LoadBalancer
	name := acceptance.RandomTestName()
	privateKeyMaterial, leafCertMaterial, certChainMaterial := acceptance.GenerateTestCertMaterial(t)
	certName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: testAccCheckAbrhaLoadbalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaLoadbalancerConfig_sslTermination(
					certName, name, privateKeyMaterial, leafCertMaterial, certChainMaterial, "certificate_name"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaLoadbalancerExists("abrha_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "forwarding_rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"abrha_loadbalancer.foobar",
						"forwarding_rule.*",
						map[string]string{
							"entry_port":       "443",
							"entry_protocol":   "https",
							"target_port":      "80",
							"target_protocol":  "http",
							"certificate_name": certName,
							"tls_passthrough":  "false",
						},
					),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "size_unit", "1"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "redirect_http_to_https", "true"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "enable_proxy_protocol", "true"),
				),
			},
		},
	})
}

func TestAccAbrhaLoadbalancer_resize(t *testing.T) {
	var loadbalancer goApiAbrha.LoadBalancer
	name := acceptance.RandomTestName()

	lbConfig := `resource "abrha_loadbalancer" "foobar" {
  name      = "%s"
  region    = "nyc3"
  size_unit = %d

  forwarding_rule {
    entry_port     = 80
    entry_protocol = "http"

    target_port     = 80
    target_protocol = "http"
  }

  healthcheck {
    port     = 22
    protocol = "tcp"
  }
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: testAccCheckAbrhaLoadbalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(lbConfig, name, 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaLoadbalancerExists("abrha_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "size_unit", "1"),
				),
			},
			{
				Config: fmt.Sprintf(lbConfig, name, 2),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaLoadbalancerExists("abrha_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "size_unit", "2"),
				),
			},
		},
	})
}

func TestAccAbrhaLoadbalancer_multipleRules(t *testing.T) {
	var loadbalancer goApiAbrha.LoadBalancer
	rName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaLoadbalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaLoadbalancerConfig_multipleRules(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaLoadbalancerExists("abrha_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "name", rName),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "forwarding_rule.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"abrha_loadbalancer.foobar",
						"forwarding_rule.*",
						map[string]string{
							"entry_port":      "443",
							"entry_protocol":  "https",
							"target_port":     "443",
							"target_protocol": "https",
							"tls_passthrough": "true",
						},
					),
					resource.TestCheckTypeSetElemNestedAttrs(
						"abrha_loadbalancer.foobar",
						"forwarding_rule.*",
						map[string]string{
							"entry_port":      "80",
							"entry_protocol":  "http",
							"target_port":     "80",
							"target_protocol": "http",
							"tls_passthrough": "false",
						},
					),
				),
			},
			{
				Config: testAccCheckAbrhaLoadbalancerConfig_multipleRulesUDP(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaLoadbalancerExists("abrha_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "name", rName),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "forwarding_rule.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"abrha_loadbalancer.foobar",
						"forwarding_rule.*",
						map[string]string{
							"entry_port":      "443",
							"entry_protocol":  "udp",
							"target_port":     "443",
							"target_protocol": "udp",
							"tls_passthrough": "false",
						},
					),
					resource.TestCheckTypeSetElemNestedAttrs(
						"abrha_loadbalancer.foobar",
						"forwarding_rule.*",
						map[string]string{
							"entry_port":      "444",
							"entry_protocol":  "udp",
							"target_port":     "444",
							"target_protocol": "udp",
							"tls_passthrough": "false",
						},
					),
				),
			},
		},
	})
}

func TestAccAbrhaLoadbalancer_WithVPC(t *testing.T) {
	var loadbalancer goApiAbrha.LoadBalancer
	lbName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaLoadbalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaLoadbalancerConfig_WithVPC(lbName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaLoadbalancerExists("abrha_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "name", lbName),
					resource.TestCheckResourceAttrSet(
						"abrha_loadbalancer.foobar", "vpc_uuid"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "vm_ids.#", "1"),
				),
			},
		},
	})
}

func TestAccAbrhaLoadbalancer_Firewall(t *testing.T) {
	var loadbalancer goApiAbrha.LoadBalancer
	lbName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaLoadbalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaLoadbalancerConfig_Firewall(lbName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaLoadbalancerExists("abrha_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "name", lbName),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "firewall.#", "1"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "firewall.0.deny.0", "cidr:1.2.0.0/16"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "firewall.0.deny.1", "ip:2.3.4.5"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "firewall.0.allow.0", "ip:1.2.3.4"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.foobar", "firewall.0.allow.1", "cidr:2.3.4.0/24"),
				),
			},
		},
	})
}

func TestLoadbalancerDiffCheck(t *testing.T) {
	cases := []struct {
		name         string
		attrs        map[string]interface{}
		expectError  bool
		errorMessage string
	}{
		{
			name:         "Missing type and region",
			expectError:  true,
			errorMessage: "missing 'region' value",
		},
		{
			name:        "Missing type",
			expectError: false,
			attrs: map[string]interface{}{
				"region": "nyc3",
			},
		},
		{
			name:         "Empty region",
			expectError:  true,
			errorMessage: "missing 'region' value",
			attrs: map[string]interface{}{
				"region": "",
			},
		},
		{
			name:         "Regional type without region",
			expectError:  true,
			errorMessage: "'region' must be set and not be empty when 'type' is 'REGIONAL'",
			attrs: map[string]interface{}{
				"type": "REGIONAL",
			},
		},
		{
			name:         "Regional type with empty region",
			expectError:  true,
			errorMessage: "'region' must be set and not be empty when 'type' is 'REGIONAL'",
			attrs: map[string]interface{}{
				"type":   "REGIONAL",
				"region": "",
			},
		},
		{
			name:        "Regional type with region",
			expectError: false,
			attrs: map[string]interface{}{
				"type":   "REGIONAL",
				"region": "nyc3",
			},
		},
		{
			name:        "Global type without region",
			expectError: false,
			attrs: map[string]interface{}{
				"type": "GLOBAL",
			},
		},
		{
			name:        "Global type with empty region",
			expectError: false,
			attrs: map[string]interface{}{
				"type":   "GLOBAL",
				"region": "",
			},
		},
		{
			name:         "Global type with region",
			expectError:  true,
			errorMessage: "'region' must be empty or not set when 'type' is 'GLOBAL'",
			attrs: map[string]interface{}{
				"type":   "GLOBAL",
				"region": "nyc3",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var s *terraform.InstanceState
			conf := terraform.NewResourceConfigRaw(c.attrs)

			r := loadbalancer.ResourceAbrhaLoadbalancer()
			_, err := r.Diff(context.Background(), s, conf, nil)

			if c.expectError {
				if err.Error() != c.errorMessage {
					t.Fatalf("Expected %s, got %s", c.errorMessage, err)
				}
			}
		})
	}
}

func TestAccAbrhaGlobalLoadbalancer(t *testing.T) {
	var loadbalancer goApiAbrha.LoadBalancer
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaLoadbalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaGlobalLoadbalancerConfig_basic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaLoadbalancerExists("abrha_loadbalancer.lorem", &loadbalancer),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.lorem", "name", "global-"+name),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.lorem", "type", "GLOBAL"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.lorem", "glb_settings.0.target_protocol", "http"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.lorem", "glb_settings.0.target_port", "80"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.lorem", "glb_settings.0.cdn.0.is_enabled", "true"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.lorem", "domains.#", "2"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.lorem", "domains.1.name", "test.github.io"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.lorem", "domains.0.name", "test-2.github.io"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.lorem", "vm_ids.#", "1"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.lorem", "target_load_balancer_ids.#", "1"),
				),
			},
			{
				Config: testAccCheckAbrhaGlobalLoadbalancerConfig_updated(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaLoadbalancerExists("abrha_loadbalancer.lorem", &loadbalancer),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.lorem", "name", "global-"+name),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.lorem", "type", "GLOBAL"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.lorem", "glb_settings.0.target_protocol", "http"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.lorem", "glb_settings.0.target_port", "80"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.lorem", "glb_settings.0.cdn.0.is_enabled", "false"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "glb_settings.0.region_priorities.%", "2"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "glb_settings.0.region_priorities.nyc1", "1"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "glb_settings.0.region_priorities.nyc2", "2"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.lorem", "domains.#", "2"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.lorem", "domains.1.name", "test-updated.github.io"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.lorem", "domains.0.name", "test-updated-2.github.io"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.lorem", "vm_ids.#", "1"),
					resource.TestCheckResourceAttr(
						"abrha_loadbalancer.lorem", "target_load_balancer_ids.#", "1"),
				),
			},
		},
	})
}

func testAccCheckAbrhaLoadbalancerDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "abrha_loadbalancer" {
			continue
		}

		_, _, err := client.LoadBalancers.Get(context.Background(), rs.Primary.ID)

		if err != nil && !strings.Contains(err.Error(), "404") {
			return fmt.Errorf(
				"Error waiting for loadbalancer (%s) to be destroyed: %s",
				rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckAbrhaLoadbalancerExists(n string, loadbalancer *goApiAbrha.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Loadbalancer ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

		lb, _, err := client.LoadBalancers.Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}

		if lb.ID != rs.Primary.ID {
			return fmt.Errorf("Loabalancer not found")
		}

		*loadbalancer = *lb

		return nil
	}
}

func testAccCheckAbrhaLoadbalancerConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "abrha_vm" "foobar" {
  name   = "%s"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = "nyc3"
}

resource "abrha_loadbalancer" "foobar" {
  name    = "%s"
  region  = "nyc3"
  type    = "REGIONAL"
  network = "EXTERNAL"

  forwarding_rule {
    entry_port     = 80
    entry_protocol = "http"

    target_port     = 80
    target_protocol = "http"
  }

  healthcheck {
    port     = 22
    protocol = "tcp"
  }

  enable_proxy_protocol     = true
  enable_backend_keepalive  = true
  http_idle_timeout_seconds = 90

  vm_ids = [abrha_vm.foobar.id]
}`, name, name)
}

func testAccCheckAbrhaLoadbalancerConfig_updated(name string) string {
	return fmt.Sprintf(`
resource "abrha_vm" "foobar" {
  name   = "%s-01"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = "nyc3"
}

resource "abrha_vm" "foo" {
  name   = "%s-02"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = "nyc3"
}

resource "abrha_loadbalancer" "foobar" {
  name    = "%s"
  region  = "nyc3"
  type    = "REGIONAL"
  network = "EXTERNAL"

  forwarding_rule {
    entry_port     = 81
    entry_protocol = "http"

    target_port     = 81
    target_protocol = "http"
  }

  healthcheck {
    port     = 22
    protocol = "tcp"
  }

  enable_proxy_protocol            = false
  enable_backend_keepalive         = false
  disable_lets_encrypt_dns_records = true
  http_idle_timeout_seconds        = 120

  vm_ids = [abrha_vm.foobar.id, abrha_vm.foo.id]
}`, name, name, name)
}

func testAccCheckAbrhaLoadbalancerConfig_vmTag(name string) string {
	return fmt.Sprintf(`
resource "abrha_tag" "barbaz" {
  name = "sample"
}

resource "abrha_vm" "foobar" {
  name   = "%s"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = "nyc3"
  tags   = [abrha_tag.barbaz.id]
}

resource "abrha_loadbalancer" "foobar" {
  name   = "%s"
  region = "nyc3"

  forwarding_rule {
    entry_port     = 80
    entry_protocol = "http"

    target_port     = 80
    target_protocol = "http"
  }

  healthcheck {
    port     = 22
    protocol = "tcp"
  }

  vm_tag = abrha_tag.barbaz.name

  depends_on = [abrha_vm.foobar]
}`, name, name)
}

func testAccCheckAbrhaLoadbalancerConfig_minimal(name string) string {
	return fmt.Sprintf(`
resource "abrha_vm" "foobar" {
  name   = "%s"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = "nyc3"
}

resource "abrha_loadbalancer" "foobar" {
  name      = "%s"
  region    = "nyc3"
  size_unit = 1

  forwarding_rule {
    entry_port     = 80
    entry_protocol = "http"

    target_port     = 80
    target_protocol = "http"
  }

  vm_ids = [abrha_vm.foobar.id]
}`, name, name)
}

func testAccCheckAbrhaLoadbalancerConfig_NonDefaultProject(projectName, lbName string) string {
	return fmt.Sprintf(`
resource "abrha_tag" "test" {
  name = "%s"
}

resource "abrha_loadbalancer" "test" {
  name       = "%s"
  region     = "nyc3"
  size       = "lb-small"
  project_id = abrha_project.test.id

  forwarding_rule {
    entry_port     = 80
    entry_protocol = "http"

    target_port     = 80
    target_protocol = "http"
  }

  vm_tag = abrha_tag.test.name
}`, projectName, lbName)
}

func testAccCheckAbrhaLoadbalancerConfig_minimalUDP(name string) string {
	return fmt.Sprintf(`
resource "abrha_vm" "foobar" {
  name   = "%s"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = "nyc3"
}

resource "abrha_loadbalancer" "foobar" {
  name   = "%s"
  region = "nyc3"
  size   = "lb-small"

  forwarding_rule {
    entry_port     = 80
    entry_protocol = "udp"

    target_port     = 80
    target_protocol = "udp"
  }

  vm_ids = [abrha_vm.foobar.id]
}`, name, name)
}

func testAccCheckAbrhaLoadbalancerConfig_stickySessions(name string) string {
	return fmt.Sprintf(`
resource "abrha_vm" "foobar" {
  name   = "%s"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = "nyc3"
}

resource "abrha_loadbalancer" "foobar" {
  name   = "%s"
  region = "nyc3"
  size   = "lb-small"

  forwarding_rule {
    entry_port     = 80
    entry_protocol = "http"

    target_port     = 80
    target_protocol = "http"
  }

  sticky_sessions {
    type               = "cookies"
    cookie_name        = "sessioncookie"
    cookie_ttl_seconds = 1800
  }

  vm_ids = [abrha_vm.foobar.id]
}`, name, name)
}

func testAccCheckAbrhaLoadbalancerConfig_sslTermination(certName string, name string, privateKeyMaterial, leafCert, certChain, certAttribute string) string {
	return fmt.Sprintf(`
resource "abrha_certificate" "foobar" {
  name              = "%s"
  private_key       = <<EOF
%s
EOF
  leaf_certificate  = <<EOF
%s
EOF
  certificate_chain = <<EOF
%s
EOF
}

resource "abrha_loadbalancer" "foobar" {
  name                   = "%s"
  region                 = "nyc3"
  size                   = "lb-small"
  redirect_http_to_https = true
  enable_proxy_protocol  = true

  forwarding_rule {
    entry_port     = 443
    entry_protocol = "https"

    target_port     = 80
    target_protocol = "http"

    %s = abrha_certificate.foobar.id
  }
}`, certName, privateKeyMaterial, leafCert, certChain, name, certAttribute)
}

func testAccCheckAbrhaLoadbalancerConfig_multipleRules(rName string) string {
	return fmt.Sprintf(`
resource "abrha_loadbalancer" "foobar" {
  name   = "%s"
  region = "nyc3"
  size   = "lb-small"

  forwarding_rule {
    entry_port     = 443
    entry_protocol = "https"

    target_port     = 443
    target_protocol = "https"

    tls_passthrough = true
  }

  forwarding_rule {
    entry_port      = 80
    target_protocol = "http"
    entry_protocol  = "http"
    target_port     = 80
  }
}`, rName)
}

func testAccCheckAbrhaLoadbalancerConfig_multipleRulesUDP(rName string) string {
	return fmt.Sprintf(`
resource "abrha_loadbalancer" "foobar" {
  name   = "%s"
  region = "nyc3"
  size   = "lb-small"

  forwarding_rule {
    entry_port     = 443
    entry_protocol = "udp"

    target_port     = 443
    target_protocol = "udp"
  }

  forwarding_rule {
    entry_port      = 444
    target_protocol = "udp"
    entry_protocol  = "udp"
    target_port     = 444
  }
}`, rName)
}

func testAccCheckAbrhaLoadbalancerConfig_WithVPC(name string) string {
	return fmt.Sprintf(`
resource "abrha_vpc" "foobar" {
  name   = "%s"
  region = "nyc3"
}

resource "abrha_vm" "foobar" {
  name     = "%s"
  size     = "s-1vcpu-1gb"
  image    = "ubuntu-22-04-x64"
  region   = "nyc3"
  vpc_uuid = abrha_vpc.foobar.id
}

resource "abrha_loadbalancer" "foobar" {
  name   = "%s"
  region = "nyc3"
  size   = "lb-small"

  forwarding_rule {
    entry_port     = 80
    entry_protocol = "http"

    target_port     = 80
    target_protocol = "http"
  }

  vpc_uuid    = abrha_vpc.foobar.id
  vm_ids = [abrha_vm.foobar.id]
}`, acceptance.RandomTestName(), acceptance.RandomTestName(), name)
}

func testAccCheckAbrhaLoadbalancerConfig_Firewall(name string) string {
	return fmt.Sprintf(`
resource "abrha_vm" "foobar" {
  name   = "%s"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = "nyc3"
}

resource "abrha_loadbalancer" "foobar" {
  name   = "%s"
  region = "nyc3"
  size   = "lb-small"

  forwarding_rule {
    entry_port     = 80
    entry_protocol = "http"

    target_port     = 80
    target_protocol = "http"
  }

  firewall {
    deny  = ["cidr:1.2.0.0/16", "ip:2.3.4.5"]
    allow = ["ip:1.2.3.4", "cidr:2.3.4.0/24"]
  }

  vm_ids = [abrha_vm.foobar.id]
}`, acceptance.RandomTestName(), name)
}

func testAccCheckAbrhaGlobalLoadbalancerConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "abrha_vm" "foobar" {
  name   = "%s"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = "blr1"
}

resource "abrha_loadbalancer" "foobar" {
  name   = "%s"
  region = "blr1"
  type   = "REGIONAL"

  forwarding_rule {
    entry_port     = 80
    entry_protocol = "http"

    target_port     = 80
    target_protocol = "http"
  }

  healthcheck {
    port     = 22
    protocol = "tcp"
  }

  enable_proxy_protocol     = true
  enable_backend_keepalive  = true
  http_idle_timeout_seconds = 90

  vm_ids = [abrha_vm.foobar.id]
}

resource "abrha_loadbalancer" "lorem" {
  name = "global-%s"
  type = "GLOBAL"

  healthcheck {
    port     = 80
    protocol = "http"
    path     = "/"
  }

  glb_settings {
    target_protocol = "http"
    target_port     = "80"
    cdn {
      is_enabled = true
    }
  }

  domains {
    name       = "test.github.io"
    is_managed = false
  }

  domains {
    name       = "test-2.github.io"
    is_managed = false
  }

  vm_ids              = [abrha_vm.foobar.id]
  target_load_balancer_ids = [abrha_loadbalancer.foobar.id]
}`, name, name, name)
}

func testAccCheckAbrhaGlobalLoadbalancerConfig_updated(name string) string {
	return fmt.Sprintf(`
resource "abrha_vm" "foobar" {
  name   = "%s"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = "blr1"
}

resource "abrha_loadbalancer" "foobar" {
  name   = "%s"
  region = "blr1"
  type   = "REGIONAL"

  forwarding_rule {
    entry_port     = 80
    entry_protocol = "http"

    target_port     = 80
    target_protocol = "http"
  }

  healthcheck {
    port     = 22
    protocol = "tcp"
  }

  enable_proxy_protocol     = true
  enable_backend_keepalive  = true
  http_idle_timeout_seconds = 90

  vm_ids = [abrha_vm.foobar.id]
}

resource "abrha_loadbalancer" "lorem" {
  name = "global-%s"
  type = "GLOBAL"

  healthcheck {
    port     = 80
    protocol = "http"
    path     = "/"
  }

  glb_settings {
    target_protocol = "http"
    target_port     = "80"
    cdn {
      is_enabled = false
    }
    region_priorities = {
      nyc1 = 1
      nyc2 = 2
    }
    failover_threshold = 10
  }

  domains {
    name       = "test-updated.github.io"
    is_managed = false
  }

  domains {
    name       = "test-updated-2.github.io"
    is_managed = false
  }

  vm_ids              = [abrha_vm.foobar.id]
  target_load_balancer_ids = [abrha_loadbalancer.foobar.id]
}`, name, name, name)
}
