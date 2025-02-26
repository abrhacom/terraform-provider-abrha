package domain_test

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/abrhacom/terraform-provider-abrha/abrha/domain"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAbrhaRecordConstructFqdn(t *testing.T) {
	cases := []struct {
		Input, Output string
	}{
		{"www", "www.nonexample.com"},
		{"dev.www", "dev.www.nonexample.com"},
		{"*", "*.nonexample.com"},
		{"nonexample.com", "nonexample.com.nonexample.com"},
		{"test.nonexample.com", "test.nonexample.com.nonexample.com"},
		{"test.nonexample.com.", "test.nonexample.com"},
		{"@", "nonexample.com"},
	}

	domainName := "nonexample.com"
	for _, tc := range cases {
		actual := domain.ConstructFqdn(tc.Input, domainName)
		if actual != tc.Output {
			t.Fatalf("input: %s\noutput: %s", tc.Input, actual)
		}
	}
}

func TestAccAbrhaRecord_Basic(t *testing.T) {
	var record goApiAbrha.DomainRecord
	domain := acceptance.RandomTestName() + ".com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAbrhaRecordConfig_basic, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaRecordExists("abrha_record.foobar", &record),
					testAccCheckAbrhaRecordAttributes(&record),
					resource.TestCheckResourceAttr(
						"abrha_record.foobar", "name", "terraform"),
					resource.TestCheckResourceAttr(
						"abrha_record.foobar", "domain", domain),
					resource.TestCheckResourceAttr(
						"abrha_record.foobar", "value", "192.168.0.10"),
					resource.TestCheckResourceAttr(
						"abrha_record.foobar", "fqdn", strings.Join([]string{"terraform", domain}, ".")),
				),
			},
		},
	})
}

func TestAccAbrhaRecord_BasicFullName(t *testing.T) {
	var record goApiAbrha.DomainRecord
	domain := acceptance.RandomTestName("record") + ".com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAbrhaRecordConfig_basic_full_name, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaRecordExists("abrha_record.foobar", &record),
					testAccCheckAbrhaRecordAttributes(&record),
					resource.TestCheckResourceAttr(
						"abrha_record.foobar", "name", "terraform"),
					resource.TestCheckResourceAttr(
						"abrha_record.foobar", "domain", domain),
					resource.TestCheckResourceAttr(
						"abrha_record.foobar", "value", "192.168.0.10"),
					resource.TestCheckResourceAttr(
						"abrha_record.foobar", "fqdn", strings.Join([]string{"terraform", domain}, ".")),
				),
			},
		},
	})
}

func TestAccAbrhaRecord_Updated(t *testing.T) {
	var record goApiAbrha.DomainRecord
	domain := acceptance.RandomTestName("record") + ".com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAbrhaRecordConfig_basic, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaRecordExists("abrha_record.foobar", &record),
					testAccCheckAbrhaRecordAttributes(&record),
					resource.TestCheckResourceAttr(
						"abrha_record.foobar", "name", "terraform"),
					resource.TestCheckResourceAttr(
						"abrha_record.foobar", "domain", domain),
					resource.TestCheckResourceAttr(
						"abrha_record.foobar", "value", "192.168.0.10"),
					resource.TestCheckResourceAttr(
						"abrha_record.foobar", "type", "A"),
					resource.TestCheckResourceAttr(
						"abrha_record.foobar", "ttl", "1800"),
				),
			},
			{
				Config: fmt.Sprintf(
					testAccCheckAbrhaRecordConfig_new_value, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaRecordExists("abrha_record.foobar", &record),
					testAccCheckAbrhaRecordAttributesUpdated(&record),
					resource.TestCheckResourceAttr(
						"abrha_record.foobar", "name", "terraform"),
					resource.TestCheckResourceAttr(
						"abrha_record.foobar", "domain", domain),
					resource.TestCheckResourceAttr(
						"abrha_record.foobar", "value", "192.168.0.11"),
					resource.TestCheckResourceAttr(
						"abrha_record.foobar", "type", "A"),
					resource.TestCheckResourceAttr(
						"abrha_record.foobar", "ttl", "90"),
				),
			},
		},
	})
}

func TestAccAbrhaRecord_HostnameValue(t *testing.T) {
	var record goApiAbrha.DomainRecord
	domain := acceptance.RandomTestName() + ".com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccCheckAbrhaRecordConfig_cname, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaRecordExists("abrha_record.foobar", &record),
					testAccCheckAbrhaRecordAttributesHostname("a.foobar-test-terraform.com", &record),
					resource.TestCheckResourceAttr(
						"abrha_record.foobar", "name", "terraform"),
					resource.TestCheckResourceAttr(
						"abrha_record.foobar", "domain", domain),
					resource.TestCheckResourceAttr(
						"abrha_record.foobar", "value", "a.foobar-test-terraform.com."),
					resource.TestCheckResourceAttr(
						"abrha_record.foobar", "type", "CNAME"),
				),
			},
		},
	})
}

func TestAccAbrhaRecord_ExternalHostnameValue(t *testing.T) {
	var record goApiAbrha.DomainRecord
	domain := acceptance.RandomTestName() + ".com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccCheckAbrhaRecordConfig_external_cname, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaRecordExists("abrha_record.foobar", &record),
					testAccCheckAbrhaRecordAttributesHostname("a.foobar-test-terraform.net", &record),
					resource.TestCheckResourceAttr(
						"abrha_record.foobar", "name", "terraform"),
					resource.TestCheckResourceAttr(
						"abrha_record.foobar", "domain", domain),
					resource.TestCheckResourceAttr(
						"abrha_record.foobar", "value", "a.foobar-test-terraform.net."),
					resource.TestCheckResourceAttr(
						"abrha_record.foobar", "type", "CNAME"),
				),
			},
		},
	})
}

func TestAccAbrhaRecord_FlagsAndTag(t *testing.T) {
	var record goApiAbrha.DomainRecord
	domain := acceptance.RandomTestName() + ".com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccCheckAbrhaRecordConfig_caa, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaRecordExists("abrha_record.foobar", &record),
					testAccCheckAbrhaRecordAttributesHostname("letsencrypt.org", &record),
					resource.TestCheckResourceAttr(
						"abrha_record.foobar", "name", "terraform"),
					resource.TestCheckResourceAttr(
						"abrha_record.foobar", "domain", domain),
					resource.TestCheckResourceAttr(
						"abrha_record.foobar", "value", "letsencrypt.org."),
					resource.TestCheckResourceAttr(
						"abrha_record.foobar", "type", "CAA"),
					resource.TestCheckResourceAttr(
						"abrha_record.foobar", "flags", "1"),
					resource.TestCheckResourceAttr(
						"abrha_record.foobar", "tag", "issue"),
				),
			},
		},
	})
}

func TestAccAbrhaRecord_MX(t *testing.T) {
	var record goApiAbrha.DomainRecord
	domain := acceptance.RandomTestName() + ".com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccCheckAbrhaRecordConfig_mx, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaRecordExists("abrha_record.foo_record", &record),
					testAccCheckAbrhaRecordAttributesHostname("foobar."+domain, &record),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "name", "terraform"),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "domain", domain),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "value", "foobar."+domain+"."),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "type", "MX"),
				),
			},
		},
	})
}

func TestAccAbrhaRecord_MX_at(t *testing.T) {
	var record goApiAbrha.DomainRecord
	domain := acceptance.RandomTestName() + ".com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccCheckAbrhaRecordConfig_mx_at, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaRecordExists("abrha_record.foo_record", &record),
					testAccCheckAbrhaRecordAttributesHostname("@", &record),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "name", "terraform"),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "domain", domain),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "value", "@"),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "type", "MX"),
				),
			},
		},
	})
}

func TestAccAbrhaRecord_SRV_zero_weight_port(t *testing.T) {
	var record goApiAbrha.DomainRecord
	domain := acceptance.RandomTestName() + ".com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccCheckAbrhaRecordConfig_srv_zero_weight_port, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaRecordExists("abrha_record.foo_record", &record),
					testAccCheckAbrhaRecordAttributesHostname("foobar."+domain, &record),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "name", "_service._protocol"),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "domain", domain),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "value", "foobar."+domain+"."),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "type", "SRV"),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "weight", "0"),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "port", "0"),
				),
			},
		},
	})
}

func TestAccAbrhaRecord_UpdateBasic(t *testing.T) {
	var record goApiAbrha.DomainRecord
	domain := acceptance.RandomTestName() + ".com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccCheckAbrhaRecordConfig_updated_basic, domain, "terraform", "a.foobar-test-terraform.com.", "1800"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaRecordExists("abrha_record.foobar", &record),
					testAccCheckAbrhaRecordAttributesHostname("a.foobar-test-terraform.com", &record),
					resource.TestCheckResourceAttr(
						"abrha_record.foobar", "name", "terraform"),
					resource.TestCheckResourceAttr(
						"abrha_record.foobar", "domain", domain),
					resource.TestCheckResourceAttr(
						"abrha_record.foobar", "value", "a.foobar-test-terraform.com."),
					resource.TestCheckResourceAttr(
						"abrha_record.foobar", "type", "CNAME"),
					resource.TestCheckResourceAttr(
						"abrha_record.foobar", "ttl", "1800"),
				),
			},
			{
				Config: fmt.Sprintf(
					testAccCheckAbrhaRecordConfig_updated_basic, domain, "terraform-updated", "b.foobar-test-terraform.com.", "1000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaRecordExists("abrha_record.foobar", &record),
					testAccCheckAbrhaRecordAttributesHostname("b.foobar-test-terraform.com", &record),
					resource.TestCheckResourceAttr(
						"abrha_record.foobar", "name", "terraform-updated"),
					resource.TestCheckResourceAttr(
						"abrha_record.foobar", "domain", domain),
					resource.TestCheckResourceAttr(
						"abrha_record.foobar", "value", "b.foobar-test-terraform.com."),
					resource.TestCheckResourceAttr(
						"abrha_record.foobar", "type", "CNAME"),
					resource.TestCheckResourceAttr(
						"abrha_record.foobar", "ttl", "1000"),
				),
			},
		},
	})
}

func TestAccAbrhaRecord_MXUpdated(t *testing.T) {
	var record goApiAbrha.DomainRecord
	domain := acceptance.RandomTestName() + ".com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccCheckAbrhaRecordConfig_mx_updated, domain, "10"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaRecordExists("abrha_record.foo_record", &record),
					testAccCheckAbrhaRecordAttributesHostname("foobar."+domain, &record),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "name", "terraform"),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "domain", domain),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "value", "foobar."+domain+"."),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "type", "MX"),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "priority", "10"),
				),
			},
			{
				Config: fmt.Sprintf(
					testAccCheckAbrhaRecordConfig_mx_updated, domain, "20"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaRecordExists("abrha_record.foo_record", &record),
					testAccCheckAbrhaRecordAttributesHostname("foobar."+domain, &record),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "name", "terraform"),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "domain", domain),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "value", "foobar."+domain+"."),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "type", "MX"),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "priority", "20"),
				),
			},
		},
	})
}

func TestAccAbrhaRecord_SrvUpdated(t *testing.T) {
	var record goApiAbrha.DomainRecord
	domain := acceptance.RandomTestName() + ".com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccCheckAbrhaRecordConfig_srv_updated, domain, "5050", "100"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaRecordExists("abrha_record.foo_record", &record),
					testAccCheckAbrhaRecordAttributesHostname("foobar."+domain, &record),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "name", "_service._protocol"),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "domain", domain),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "value", "foobar."+domain+"."),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "type", "SRV"),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "port", "5050"),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "weight", "100"),
				),
			},
			{
				Config: fmt.Sprintf(
					testAccCheckAbrhaRecordConfig_srv_updated, domain, "6060", "150"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaRecordExists("abrha_record.foo_record", &record),
					testAccCheckAbrhaRecordAttributesHostname("foobar."+domain, &record),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "name", "_service._protocol"),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "domain", domain),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "value", "foobar."+domain+"."),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "type", "SRV"),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "port", "6060"),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "weight", "150"),
				),
			},
		},
	})
}

func TestAccAbrhaRecord_CaaUpdated(t *testing.T) {
	var record goApiAbrha.DomainRecord
	domain := acceptance.RandomTestName() + ".com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccCheckAbrhaRecordConfig_caa_updated, domain, "20", "issue"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaRecordExists("abrha_record.foo_record", &record),
					testAccCheckAbrhaRecordAttributesHostname("letsencrypt.org", &record),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "name", "terraform"),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "domain", domain),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "value", "letsencrypt.org."),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "type", "CAA"),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "flags", "20"),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "tag", "issue"),
				),
			},
			{
				Config: fmt.Sprintf(
					testAccCheckAbrhaRecordConfig_caa_updated, domain, "50", "issuewild"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaRecordExists("abrha_record.foo_record", &record),
					testAccCheckAbrhaRecordAttributesHostname("letsencrypt.org", &record),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "name", "terraform"),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "domain", domain),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "value", "letsencrypt.org."),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "type", "CAA"),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "flags", "50"),
					resource.TestCheckResourceAttr(
						"abrha_record.foo_record", "tag", "issuewild"),
				),
			},
		},
	})
}

func TestAccAbrhaRecord_iodefCAA(t *testing.T) {
	var record goApiAbrha.DomainRecord
	domain := acceptance.RandomTestName() + ".com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccCheckAbrhaRecordConfig_iodef, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaRecordExists("abrha_record.CAA_iodef", &record),
					resource.TestCheckResourceAttr(
						"abrha_record.CAA_iodef", "name", "@"),
					resource.TestCheckResourceAttr(
						"abrha_record.CAA_iodef", "domain", domain),
					resource.TestCheckResourceAttr(
						"abrha_record.CAA_iodef", "value", "mailto:caa-failures@example.com"),
					resource.TestCheckResourceAttr(
						"abrha_record.CAA_iodef", "type", "CAA"),
					resource.TestCheckResourceAttr(
						"abrha_record.CAA_iodef", "flags", "0"),
					resource.TestCheckResourceAttr(
						"abrha_record.CAA_iodef", "tag", "iodef"),
				),
			},
		},
	})
}

func TestAccAbrhaRecord_TXT(t *testing.T) {
	var record goApiAbrha.DomainRecord
	domain := acceptance.RandomTestName() + ".com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAbrhaRecordTXT, domain, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaRecordExists("abrha_record.txt", &record),
					resource.TestCheckResourceAttr(
						"abrha_record.txt", "type", "TXT"),
					resource.TestCheckResourceAttr(
						"abrha_record.txt", "domain", domain),
					resource.TestCheckResourceAttr(
						"abrha_record.txt", "value", "v=spf1 a:smtp01.example.com a:mail.example.com -all"),
					resource.TestCheckResourceAttr(
						"abrha_record.txt", "fqdn", domain),
				),
			},
		},
	})
}

func TestAccAbrhaRecord_ExpectedErrors(t *testing.T) {
	var (
		srvNoPort = `resource "abrha_record" "pgsql_default_pub_srv" {
  domain = "example.com"

  type = "SRV"
  name = "_postgresql_.tcp.example.com"

  // priority can be 0, but must be set.
  priority = 0
  weight   = 0
  value    = "srv.example.com"
}`
		srvNoPriority = `resource "abrha_record" "pgsql_default_pub_srv" {
  domain = "example.com"

  type = "SRV"
  name = "_postgresql_.tcp.example.com"

  port   = 3600
  weight = 0
  value  = "srv.example.com"
}`
		srvNoWeight = `resource "abrha_record" "pgsql_default_pub_srv" {
  domain = "example.com"

  type = "SRV"
  name = "_postgresql._tcp.example.com"

  port     = 3600
  priority = 10
  value    = "srv.example.com"
}`
		mxNoPriority = `resource "abrha_record" "foo_record" {
  domain = "example.com"

  name  = "terraform"
  value = "mail."
  type  = "MX"
}`
		caaNoFlags = `resource "abrha_record" "foo_record" {
  domain = "example.com"

  name  = "cert"
  type  = "CAA"
  value = "letsencrypt.org."
  tag   = "issue"
}`
		caaNoTag = `resource "abrha_record" "foo_record" {
  domain = "example.com"

  name  = "cert"
  type  = "CAA"
  value = "letsencrypt.org."
  flags = 1
}`
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config:      srvNoPort,
				ExpectError: regexp.MustCompile("`port` is required for when type is `SRV`"),
			},
			{
				Config:      srvNoPriority,
				ExpectError: regexp.MustCompile("`priority` is required for when type is `SRV`"),
			},
			{
				Config:      srvNoWeight,
				ExpectError: regexp.MustCompile("`weight` is required for when type is `SRV`"),
			},
			{
				Config:      mxNoPriority,
				ExpectError: regexp.MustCompile("`priority` is required for when type is `MX`"),
			},
			{
				Config:      caaNoFlags,
				ExpectError: regexp.MustCompile("`flags` is required for when type is `CAA`"),
			},
			{
				Config:      caaNoTag,
				ExpectError: regexp.MustCompile("`tag` is required for when type is `CAA`"),
			},
		},
	})
}

func testAccCheckAbrhaRecordDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "abrha_record" {
			continue
		}
		domain := rs.Primary.Attributes["domain"]
		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, _, err = client.Domains.Record(context.Background(), domain, id)

		if err == nil {
			return fmt.Errorf("Record still exists")
		}
	}

	return nil
}

func testAccCheckAbrhaRecordAttributes(record *goApiAbrha.DomainRecord) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if record.Data != "192.168.0.10" {
			return fmt.Errorf("Bad value: %s", record.Data)
		}

		return nil
	}
}

func testAccCheckAbrhaRecordAttributesUpdated(record *goApiAbrha.DomainRecord) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if record.Data != "192.168.0.11" {
			return fmt.Errorf("Bad value: %s", record.Data)
		}

		return nil
	}
}

func testAccCheckAbrhaRecordExists(n string, record *goApiAbrha.DomainRecord) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

		domain := rs.Primary.Attributes["domain"]
		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		foundRecord, _, err := client.Domains.Record(context.Background(), domain, id)

		if err != nil {
			return err
		}

		if strconv.Itoa(foundRecord.ID) != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		*record = *foundRecord

		return nil
	}
}

func testAccCheckAbrhaRecordAttributesHostname(data string, record *goApiAbrha.DomainRecord) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if record.Data != data {
			return fmt.Errorf("Bad value: expected %s, got %s", data, record.Data)
		}

		return nil
	}
}

const testAccCheckAbrhaRecordConfig_basic = `
resource "abrha_domain" "foobar" {
  name       = "%s"
  ip_address = "192.168.0.10"
}

resource "abrha_record" "foobar" {
  domain = abrha_domain.foobar.name

  name  = "terraform"
  value = "192.168.0.10"
  type  = "A"
}`

const testAccCheckAbrhaRecordConfig_basic_full_name = `
resource "abrha_domain" "foobar" {
  name       = "%s"
  ip_address = "192.168.0.10"
}

resource "abrha_record" "foobar" {
  domain = abrha_domain.foobar.name

  name  = "terraform.${abrha_domain.foobar.name}."
  value = "192.168.0.10"
  type  = "A"
}`

const testAccCheckAbrhaRecordConfig_new_value = `
resource "abrha_domain" "foobar" {
  name       = "%s"
  ip_address = "192.168.0.10"
}

resource "abrha_record" "foobar" {
  domain = abrha_domain.foobar.name

  name  = "terraform"
  value = "192.168.0.11"
  type  = "A"
  ttl   = 90
}`

const testAccCheckAbrhaRecordConfig_cname = `
resource "abrha_domain" "foobar" {
  name       = "%s"
  ip_address = "192.168.0.10"
}

resource "abrha_record" "foobar" {
  domain = abrha_domain.foobar.name

  name  = "terraform"
  value = "a.foobar-test-terraform.com."
  type  = "CNAME"
}`

const testAccCheckAbrhaRecordConfig_mx_at = `
resource "abrha_domain" "foobar" {
  name       = "%s"
  ip_address = "192.168.0.10"
}

resource "abrha_record" "foo_record" {
  domain = abrha_domain.foobar.name

  name     = "terraform"
  value    = "${abrha_domain.foobar.name}."
  type     = "MX"
  priority = "10"
}`

const testAccCheckAbrhaRecordConfig_mx = `
resource "abrha_domain" "foobar" {
  name       = "%s"
  ip_address = "192.168.0.10"
}

resource "abrha_record" "foo_record" {
  domain = abrha_domain.foobar.name

  name     = "terraform"
  value    = "foobar.${abrha_domain.foobar.name}."
  type     = "MX"
  priority = "10"
}`

const testAccCheckAbrhaRecordConfig_external_cname = `
resource "abrha_domain" "foobar" {
  name       = "%s"
  ip_address = "192.168.0.10"
}

resource "abrha_record" "foobar" {
  domain = abrha_domain.foobar.name

  name  = "terraform"
  value = "a.foobar-test-terraform.net."
  type  = "CNAME"
}`

const testAccCheckAbrhaRecordConfig_caa = `
resource "abrha_domain" "foobar" {
  name       = "%s"
  ip_address = "192.168.0.10"
}

resource "abrha_record" "foobar" {
  domain = abrha_domain.foobar.name

  name  = "terraform"
  type  = "CAA"
  value = "letsencrypt.org."
  flags = 1
  tag   = "issue"
}`

const testAccCheckAbrhaRecordConfig_srv_zero_weight_port = `
resource "abrha_domain" "foobar" {
  name       = "%s"
  ip_address = "192.168.0.10"
}

resource "abrha_record" "foo_record" {
  domain = abrha_domain.foobar.name

  name     = "_service._protocol"
  value    = "foobar.${abrha_domain.foobar.name}."
  type     = "SRV"
  priority = 10
  port     = 0
  weight   = 0
}`

const testAccCheckAbrhaRecordConfig_updated_basic = `
resource "abrha_domain" "foobar" {
  name       = "%s"
  ip_address = "192.168.0.10"
}

resource "abrha_record" "foobar" {
  domain = abrha_domain.foobar.name

  name  = "%s"
  value = "%s"
  type  = "CNAME"
  ttl   = "%s"
}`

const testAccCheckAbrhaRecordConfig_mx_updated = `
resource "abrha_domain" "foobar" {
  name       = "%s"
  ip_address = "192.168.0.10"
}

resource "abrha_record" "foo_record" {
  domain = abrha_domain.foobar.name

  name     = "terraform"
  value    = "foobar.${abrha_domain.foobar.name}."
  type     = "MX"
  priority = "%s"
}`

const testAccCheckAbrhaRecordConfig_srv_updated = `
resource "abrha_domain" "foobar" {
  name       = "%s"
  ip_address = "192.168.0.10"
}

resource "abrha_record" "foo_record" {
  domain = abrha_domain.foobar.name

  name     = "_service._protocol"
  value    = "foobar.${abrha_domain.foobar.name}."
  type     = "SRV"
  priority = "10"
  port     = "%s"
  weight   = "%s"
}`

const testAccCheckAbrhaRecordConfig_caa_updated = `
resource "abrha_domain" "foobar" {
  name       = "%s"
  ip_address = "192.168.0.10"
}

resource "abrha_record" "foo_record" {
  domain = abrha_domain.foobar.name

  name  = "terraform"
  type  = "CAA"
  value = "letsencrypt.org."
  flags = "%s"
  tag   = "%s"
}`

const testAccCheckAbrhaRecordConfig_iodef = `
resource "abrha_domain" "foobar" {
  name = "%s"
}
resource "abrha_record" "CAA_iodef" {
  domain = abrha_domain.foobar.name
  type   = "CAA"
  tag    = "iodef"
  flags  = "0"
  name   = "@"
  value  = "mailto:caa-failures@example.com"
}`

const testAccCheckAbrhaRecordTXT = `
resource "abrha_domain" "foobar" {
  name = "%s"
}
resource "abrha_record" "txt" {
  domain = abrha_domain.foobar.name
  type   = "TXT"
  name   = "%s."
  value  = "v=spf1 a:smtp01.example.com a:mail.example.com -all"
}`
