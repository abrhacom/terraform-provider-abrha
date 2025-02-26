package loadbalancer_test

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

func TestAccDataSourceAbrhaLoadBalancer_BasicByName(t *testing.T) {
	var loadbalancer goApiAbrha.LoadBalancer
	testName := acceptance.RandomTestName()
	resourceConfig := testAccCheckDataSourceAbrhaLoadBalancerConfig(testName, "lb-small")
	dataSourceConfig := `
data "abrha_loadbalancer" "foobar" {
  name = abrha_loadbalancer.foo.name
}`

	expectedURNRegEx, _ := regexp.Compile(`do:loadbalancer:[0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}`)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
			},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceAbrhaLoadBalancerExists("data.abrha_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "name", testName),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "size_unit", "1"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "type", ""),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "forwarding_rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"data.abrha_loadbalancer.foobar",
						"forwarding_rule.*",
						map[string]string{
							"entry_port":      "80",
							"entry_protocol":  "http",
							"target_port":     "80",
							"target_protocol": "http",
						},
					),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "healthcheck.#", "1"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "healthcheck.0.port", "22"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "healthcheck.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "vm_ids.#", "2"),
					resource.TestMatchResourceAttr(
						"data.abrha_loadbalancer.foobar", "urn", expectedURNRegEx),
					resource.TestCheckResourceAttrSet(
						"data.abrha_loadbalancer.foobar", "vpc_uuid"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "enable_proxy_protocol", "false"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "enable_backend_keepalive", "false"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "disable_lets_encrypt_dns_records", "false"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_loadbalancer.foobar", "project_id"),
					resource.TestCheckResourceAttrSet(
						"data.abrha_loadbalancer.foobar", "http_idle_timeout_seconds"),
				),
			},
		},
	})
}

func TestAccDataSourceAbrhaLoadBalancer_BasicById(t *testing.T) {
	var loadbalancer goApiAbrha.LoadBalancer
	testName := acceptance.RandomTestName()
	resourceConfig := testAccCheckDataSourceAbrhaLoadBalancerConfig(testName, "lb-small")
	dataSourceConfig := `
data "abrha_loadbalancer" "foobar" {
  id = abrha_loadbalancer.foo.id
}`

	expectedURNRegEx, _ := regexp.Compile(`do:loadbalancer:[0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}`)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
			},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceAbrhaLoadBalancerExists("data.abrha_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "name", testName),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "size_unit", "1"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "type", ""),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "forwarding_rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"data.abrha_loadbalancer.foobar",
						"forwarding_rule.*",
						map[string]string{
							"entry_port":      "80",
							"entry_protocol":  "http",
							"target_port":     "80",
							"target_protocol": "http",
						},
					),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "healthcheck.#", "1"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "healthcheck.0.port", "22"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "healthcheck.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "vm_ids.#", "2"),
					resource.TestMatchResourceAttr(
						"data.abrha_loadbalancer.foobar", "urn", expectedURNRegEx),
					resource.TestCheckResourceAttrSet(
						"data.abrha_loadbalancer.foobar", "vpc_uuid"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "enable_proxy_protocol", "false"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "enable_backend_keepalive", "false"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "disable_lets_encrypt_dns_records", "false"),
				),
			},
		},
	})
}

func TestAccDataSourceAbrhaLoadBalancer_LargeByName(t *testing.T) {
	var loadbalancer goApiAbrha.LoadBalancer
	testName := acceptance.RandomTestName()
	resourceConfig := testAccCheckDataSourceAbrhaLoadBalancerConfig(testName, "lb-large")
	dataSourceConfig := `
data "abrha_loadbalancer" "foobar" {
  name = abrha_loadbalancer.foo.name
}`

	expectedURNRegEx, _ := regexp.Compile(`do:loadbalancer:[0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}`)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
			},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceAbrhaLoadBalancerExists("data.abrha_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "name", testName),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "size_unit", "6"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "type", ""),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "forwarding_rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"data.abrha_loadbalancer.foobar",
						"forwarding_rule.*",
						map[string]string{
							"entry_port":      "80",
							"entry_protocol":  "http",
							"target_port":     "80",
							"target_protocol": "http",
						},
					),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "healthcheck.#", "1"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "healthcheck.0.port", "22"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "healthcheck.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "vm_ids.#", "2"),
					resource.TestMatchResourceAttr(
						"data.abrha_loadbalancer.foobar", "urn", expectedURNRegEx),
					resource.TestCheckResourceAttrSet(
						"data.abrha_loadbalancer.foobar", "vpc_uuid"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "enable_proxy_protocol", "false"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "enable_backend_keepalive", "false"),
				),
			},
		},
	})
}

func TestAccDataSourceAbrhaLoadBalancer_LargeById(t *testing.T) {
	var loadbalancer goApiAbrha.LoadBalancer
	testName := acceptance.RandomTestName()
	resourceConfig := testAccCheckDataSourceAbrhaLoadBalancerConfig(testName, "lb-large")
	dataSourceConfig := `
data "abrha_loadbalancer" "foobar" {
  id = abrha_loadbalancer.foo.id
}`

	expectedURNRegEx, _ := regexp.Compile(`do:loadbalancer:[0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}`)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
			},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceAbrhaLoadBalancerExists("data.abrha_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "name", testName),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "size_unit", "6"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "type", ""),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "forwarding_rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"data.abrha_loadbalancer.foobar",
						"forwarding_rule.*",
						map[string]string{
							"entry_port":      "80",
							"entry_protocol":  "http",
							"target_port":     "80",
							"target_protocol": "http",
						},
					),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "healthcheck.#", "1"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "healthcheck.0.port", "22"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "healthcheck.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "vm_ids.#", "2"),
					resource.TestMatchResourceAttr(
						"data.abrha_loadbalancer.foobar", "urn", expectedURNRegEx),
					resource.TestCheckResourceAttrSet(
						"data.abrha_loadbalancer.foobar", "vpc_uuid"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "enable_proxy_protocol", "false"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "enable_backend_keepalive", "false"),
				),
			},
		},
	})
}

func TestAccDataSourceAbrhaLoadBalancer_Size2ByName(t *testing.T) {
	var loadbalancer goApiAbrha.LoadBalancer
	testName := acceptance.RandomTestName()
	resourceConfig := testAccCheckDataSourceAbrhaLoadBalancerConfigSizeUnit(testName, 2)
	dataSourceConfig := `
data "abrha_loadbalancer" "foobar" {
  name = abrha_loadbalancer.foo.name
}`

	expectedURNRegEx, _ := regexp.Compile(`do:loadbalancer:[0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}`)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
			},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceAbrhaLoadBalancerExists("data.abrha_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "name", testName),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "size_unit", "2"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "forwarding_rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"data.abrha_loadbalancer.foobar",
						"forwarding_rule.*",
						map[string]string{
							"entry_port":      "80",
							"entry_protocol":  "http",
							"target_port":     "80",
							"target_protocol": "http",
						},
					),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "healthcheck.#", "1"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "healthcheck.0.port", "22"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "healthcheck.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "vm_ids.#", "2"),
					resource.TestMatchResourceAttr(
						"data.abrha_loadbalancer.foobar", "urn", expectedURNRegEx),
					resource.TestCheckResourceAttrSet(
						"data.abrha_loadbalancer.foobar", "vpc_uuid"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "enable_proxy_protocol", "false"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "enable_backend_keepalive", "false"),
				),
			},
		},
	})
}

func TestAccDataSourceAbrhaLoadBalancer_Size2ById(t *testing.T) {
	var loadbalancer goApiAbrha.LoadBalancer
	testName := acceptance.RandomTestName()
	resourceConfig := testAccCheckDataSourceAbrhaLoadBalancerConfigSizeUnit(testName, 2)
	dataSourceConfig := `
data "abrha_loadbalancer" "foobar" {
  id = abrha_loadbalancer.foo.id
}`

	expectedURNRegEx, _ := regexp.Compile(`do:loadbalancer:[0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}`)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
			},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceAbrhaLoadBalancerExists("data.abrha_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "name", testName),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "size_unit", "2"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "forwarding_rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"data.abrha_loadbalancer.foobar",
						"forwarding_rule.*",
						map[string]string{
							"entry_port":      "80",
							"entry_protocol":  "http",
							"target_port":     "80",
							"target_protocol": "http",
						},
					),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "healthcheck.#", "1"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "healthcheck.0.port", "22"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "healthcheck.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "vm_ids.#", "2"),
					resource.TestMatchResourceAttr(
						"data.abrha_loadbalancer.foobar", "urn", expectedURNRegEx),
					resource.TestCheckResourceAttrSet(
						"data.abrha_loadbalancer.foobar", "vpc_uuid"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "enable_proxy_protocol", "false"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "enable_backend_keepalive", "false"),
				),
			},
		},
	})
}

func TestAccDataSourceAbrhaLoadBalancer_multipleRulesByName(t *testing.T) {
	var loadbalancer goApiAbrha.LoadBalancer
	testName := acceptance.RandomTestName()
	resourceConfig := testAccCheckAbrhaLoadbalancerConfig_multipleRules(testName)
	dataSourceConfig := `
data "abrha_loadbalancer" "foobar" {
  name = abrha_loadbalancer.foobar.name
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaLoadbalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
			},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaLoadbalancerExists("abrha_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "name", testName),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "size_unit", "1"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "forwarding_rule.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"data.abrha_loadbalancer.foobar",
						"forwarding_rule.*",
						map[string]string{
							"entry_port":      "443",
							"entry_protocol":  "https",
							"target_port":     "443",
							"target_protocol": "https",
						},
					),
					resource.TestCheckTypeSetElemNestedAttrs(
						"data.abrha_loadbalancer.foobar",
						"forwarding_rule.*",
						map[string]string{
							"entry_port":      "80",
							"entry_protocol":  "http",
							"target_port":     "80",
							"target_protocol": "http",
						},
					),
				),
			},
		},
	})
}

func TestAccDataSourceAbrhaLoadBalancer_multipleRulesById(t *testing.T) {
	var loadbalancer goApiAbrha.LoadBalancer
	testName := acceptance.RandomTestName()
	resourceConfig := testAccCheckAbrhaLoadbalancerConfig_multipleRules(testName)
	dataSourceConfig := `
data "abrha_loadbalancer" "foobar" {
  id = abrha_loadbalancer.foobar.id
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaLoadbalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
			},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaLoadbalancerExists("abrha_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "name", testName),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "size_unit", "1"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "forwarding_rule.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"data.abrha_loadbalancer.foobar",
						"forwarding_rule.*",
						map[string]string{
							"entry_port":      "443",
							"entry_protocol":  "https",
							"target_port":     "443",
							"target_protocol": "https",
						},
					),
					resource.TestCheckTypeSetElemNestedAttrs(
						"data.abrha_loadbalancer.foobar",
						"forwarding_rule.*",
						map[string]string{
							"entry_port":      "80",
							"entry_protocol":  "http",
							"target_port":     "80",
							"target_protocol": "http",
						},
					),
				),
			},
		},
	})
}

func TestAccDataSourceAbrhaLoadBalancer_tlsCertByName(t *testing.T) {
	var loadbalancer goApiAbrha.LoadBalancer
	testName := acceptance.RandomTestName()
	name := acceptance.RandomTestName()
	privateKeyMaterial, leafCertMaterial, certChainMaterial := acceptance.GenerateTestCertMaterial(t)
	resourceConfig := testAccCheckAbrhaLoadbalancerConfig_sslTermination(
		testName+"-cert", name, privateKeyMaterial, leafCertMaterial, certChainMaterial, "certificate_name",
	)
	dataSourceConfig := `
data "abrha_loadbalancer" "foobar" {
  name = abrha_loadbalancer.foobar.name
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaLoadbalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
			},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaLoadbalancerExists("data.abrha_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "forwarding_rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"data.abrha_loadbalancer.foobar",
						"forwarding_rule.*",
						map[string]string{
							"entry_port":       "443",
							"entry_protocol":   "https",
							"target_port":      "80",
							"target_protocol":  "http",
							"certificate_name": testName + "-cert",
							"tls_passthrough":  "false",
						},
					),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "size_unit", "1"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "redirect_http_to_https", "true"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "enable_proxy_protocol", "true"),
				),
			},
		},
	})
}

func TestAccDataSourceAbrhaLoadBalancer_tlsCertById(t *testing.T) {
	var loadbalancer goApiAbrha.LoadBalancer
	testName := acceptance.RandomTestName()
	name := acceptance.RandomTestName()
	privateKeyMaterial, leafCertMaterial, certChainMaterial := acceptance.GenerateTestCertMaterial(t)
	resourceConfig := testAccCheckAbrhaLoadbalancerConfig_sslTermination(
		testName+"-cert", name, privateKeyMaterial, leafCertMaterial, certChainMaterial, "certificate_name",
	)
	dataSourceConfig := `
data "abrha_loadbalancer" "foobar" {
  id = abrha_loadbalancer.foobar.id
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaLoadbalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
			},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAbrhaLoadbalancerExists("data.abrha_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "forwarding_rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"data.abrha_loadbalancer.foobar",
						"forwarding_rule.*",
						map[string]string{
							"entry_port":       "443",
							"entry_protocol":   "https",
							"target_port":      "80",
							"target_protocol":  "http",
							"certificate_name": testName + "-cert",
							"tls_passthrough":  "false",
						},
					),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "size_unit", "1"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "redirect_http_to_https", "true"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "enable_proxy_protocol", "true"),
				),
			},
		},
	})
}
func TestAccDataSourceAbrhaGlobalLoadBalancer(t *testing.T) {
	var loadbalancer goApiAbrha.LoadBalancer
	testName := acceptance.RandomTestName()
	resourceConfig := testAccCheckDataSourceAbrhaGlobalLoadBalancerConfig(testName)
	dataSourceConfig := `
data "abrha_loadbalancer" "foobar" {
  name = abrha_loadbalancer.lorem.name
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
			},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceAbrhaLoadBalancerExists("data.abrha_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "name", testName),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "type", "GLOBAL"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "glb_settings.0.target_protocol", "HTTP"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "glb_settings.0.target_port", "80"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "glb_settings.0.cdn.0.is_enabled", "true"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "glb_settings.0.region_priorities.%", "2"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "glb_settings.0.region_priorities.nyc1", "1"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "glb_settings.0.region_priorities.nyc2", "2"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "domains.#", "2"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "domains.1.name", "test.github.io"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "domains.0.name", "test-2.github.io"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "vm_ids.#", "1"),
					resource.TestCheckResourceAttr(
						"data.abrha_loadbalancer.foobar", "network", ""),
				),
			},
		},
	})
}

func testAccCheckDataSourceAbrhaLoadBalancerExists(n string, loadbalancer *goApiAbrha.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Load Balancer ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

		foundLoadbalancer, _, err := client.LoadBalancers.Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundLoadbalancer.ID != rs.Primary.ID {
			return fmt.Errorf("Load Balancer not found")
		}

		*loadbalancer = *foundLoadbalancer

		return nil
	}
}

func testAccCheckDataSourceAbrhaLoadBalancerConfig(testName string, sizeSlug string) string {
	return fmt.Sprintf(`
resource "abrha_tag" "foo" {
  name = "%s"
}

resource "abrha_vm" "foo" {
  count              = 2
  image              = "ubuntu-22-04-x64"
  name               = "%s-${count.index}"
  region             = "nyc3"
  size               = "s-1vcpu-1gb"
  private_networking = true
  tags               = [abrha_tag.foo.id]
}

resource "abrha_loadbalancer" "foo" {
  name   = "%s"
  region = "nyc3"
  size   = "%s"
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

  vm_tag = abrha_tag.foo.id
  depends_on  = ["abrha_vm.foo"]
}`, testName, testName, testName, sizeSlug)
}

func testAccCheckDataSourceAbrhaLoadBalancerConfigSizeUnit(testName string, sizeUnit uint32) string {
	return fmt.Sprintf(`
resource "abrha_tag" "foo" {
  name = "%s"
}

resource "abrha_vm" "foo" {
  count              = 2
  image              = "ubuntu-22-04-x64"
  name               = "%s-${count.index}"
  region             = "nyc3"
  size               = "s-1vcpu-1gb"
  private_networking = true
  tags               = [abrha_tag.foo.id]
}

resource "abrha_loadbalancer" "foo" {
  name      = "%s"
  region    = "nyc3"
  size_unit = "%d"

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

  vm_tag = abrha_tag.foo.id
  depends_on  = ["abrha_vm.foo"]
}`, testName, testName, testName, sizeUnit)
}

func testAccCheckDataSourceAbrhaGlobalLoadBalancerConfig(name string) string {
	return fmt.Sprintf(`
resource "abrha_vm" "foobar" {
  name   = "%s"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = "blr1"
}

resource "abrha_loadbalancer" "lorem" {
  name    = "%s"
  type    = "GLOBAL"
  network = "EXTERNAL"

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
    region_priorities = {
      nyc1 = 1
      nyc2 = 2
    }
    failover_threshold = 10
  }

  domains {
    name       = "test.github.io"
    is_managed = false
  }

  domains {
    name       = "test-2.github.io"
    is_managed = false
  }

  vm_ids = [abrha_vm.foobar.id]
}`, name, name)
}
