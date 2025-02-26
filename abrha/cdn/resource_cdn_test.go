package cdn_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const originSuffix = ".ams3.parspackspaces.com"

func TestAccAbrhaCDN_Create(t *testing.T) {

	bucketName := generateBucketName()
	cdnCreateConfig := fmt.Sprintf(testAccCheckAbrhaCDNConfig_Create, bucketName)

	expectedOrigin := bucketName + originSuffix
	expectedTTL := "3600"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaCDNDestroy,
		Steps: []resource.TestStep{
			{
				Config: cdnCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaCDNExists("abrha_cdn.foobar"),
					resource.TestCheckResourceAttr(
						"abrha_cdn.foobar", "origin", expectedOrigin),
					resource.TestCheckResourceAttr("abrha_cdn.foobar", "ttl", expectedTTL),
				),
			},
		},
	})
}

func TestAccAbrhaCDN_CreateWithNeedCloudflareCert(t *testing.T) {
	domain := os.Getenv("DO_TEST_SUBDOMAIN")
	if domain == "" {
		t.Skip("Test requires an active DO manage sub domain. Set DO_TEST_SUBDOMAIN")
	}

	bucketName := generateBucketName()
	cdnCreateConfig := fmt.Sprintf(testAccCheckAbrhaCDNConfig_CreateWithNeedCloudflareCert, bucketName, domain)

	expectedOrigin := bucketName + originSuffix
	expectedTTL := "3600"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaCDNDestroy,
		Steps: []resource.TestStep{
			{
				Config: cdnCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaCDNExists("abrha_cdn.foobar"),
					resource.TestCheckResourceAttr(
						"abrha_cdn.foobar", "origin", expectedOrigin),
					resource.TestCheckResourceAttr("abrha_cdn.foobar", "ttl", expectedTTL),
					resource.TestCheckResourceAttr("abrha_cdn.foobar", "certificate_name", "needs-cloudflare-cert"),
				),
			},
		},
	})
}

func TestAccAbrhaCDN_Create_with_TTL(t *testing.T) {

	bucketName := generateBucketName()
	ttl := 600
	cdnCreateConfig := fmt.Sprintf(testAccCheckAbrhaCDNConfig_Create_with_TTL, bucketName, ttl)

	expectedOrigin := bucketName + originSuffix
	expectedTTL := "600"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaCDNDestroy,
		Steps: []resource.TestStep{
			{
				Config: cdnCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaCDNExists("abrha_cdn.foobar"),
					resource.TestCheckResourceAttr(
						"abrha_cdn.foobar", "origin", expectedOrigin),
					resource.TestCheckResourceAttr("abrha_cdn.foobar", "ttl", expectedTTL),
				),
			},
		},
	})
}

func TestAccAbrhaCDN_Create_and_Update(t *testing.T) {

	bucketName := generateBucketName()
	ttl := 600

	cdnCreateConfig := fmt.Sprintf(testAccCheckAbrhaCDNConfig_Create, bucketName)
	cdnUpdateConfig := fmt.Sprintf(testAccCheckAbrhaCDNConfig_Create_with_TTL, bucketName, ttl)

	expectedOrigin := bucketName + originSuffix
	expectedTTL := "600"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaCDNDestroy,
		Steps: []resource.TestStep{
			{
				Config: cdnCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaCDNExists("abrha_cdn.foobar"),
					resource.TestCheckResourceAttr(
						"abrha_cdn.foobar", "origin", expectedOrigin),
					resource.TestCheckResourceAttr("abrha_cdn.foobar", "ttl", "3600"),
				),
			},
			{
				Config: cdnUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaCDNExists("abrha_cdn.foobar"),
					resource.TestCheckResourceAttr(
						"abrha_cdn.foobar", "origin", expectedOrigin),
					resource.TestCheckResourceAttr("abrha_cdn.foobar", "ttl", expectedTTL),
				),
			},
		},
	})
}

func TestAccAbrhaCDN_CustomDomain(t *testing.T) {
	spaceName := generateBucketName()
	certName := acceptance.RandomTestName()
	updatedCertName := generateBucketName()
	domain := acceptance.RandomTestName() + ".com"
	config := testAccCheckAbrhaCDNConfig_CustomDomain(domain, spaceName, certName)
	updatedConfig := testAccCheckAbrhaCDNConfig_CustomDomain(domain, spaceName, updatedCertName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaCDNDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"tls": {
				Source:            "hashicorp/tls",
				VersionConstraint: "3.0.0",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaCDNExists("abrha_cdn.space_cdn"),
					resource.TestCheckResourceAttr(
						"abrha_cdn.space_cdn", "certificate_name", certName),
					resource.TestCheckResourceAttr(
						"abrha_cdn.space_cdn", "custom_domain", "foo."+domain),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaCDNExists("abrha_cdn.space_cdn"),
					resource.TestCheckResourceAttr(
						"abrha_cdn.space_cdn", "certificate_name", updatedCertName),
					resource.TestCheckResourceAttr(
						"abrha_cdn.space_cdn", "custom_domain", "foo."+domain),
				),
			},
		},
	})
}

func testAccCheckAbrhaCDNDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "abrha_cdn" {
			continue
		}

		_, _, err := client.CDNs.Get(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("CDN resource still exists")
		}
	}

	return nil
}

func testAccCheckAbrhaCDNExists(resource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

		rs, ok := s.RootModule().Resources[resource]

		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID set for resource: %s", resource)
		}

		foundCDN, _, err := client.CDNs.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}

		if foundCDN.ID != rs.Primary.ID {
			return fmt.Errorf("Resource not found: %s : %s", resource, rs.Primary.ID)
		}

		return nil
	}
}

func generateBucketName() string {
	return acceptance.RandomTestName("cdn")
}

const testAccCheckAbrhaCDNConfig_Create = `
resource "abrha_spaces_bucket" "bucket" {
  name   = "%s"
  region = "ams3"
  acl    = "public-read"
}

resource "abrha_cdn" "foobar" {
  origin = abrha_spaces_bucket.bucket.bucket_domain_name
}`

const testAccCheckAbrhaCDNConfig_CreateWithNeedCloudflareCert = `
resource "abrha_spaces_bucket" "bucket" {
  name   = "%s"
  region = "ams3"
  acl    = "public-read"
}

resource "abrha_cdn" "foobar" {
  origin           = abrha_spaces_bucket.bucket.bucket_domain_name
  certificate_name = "needs-cloudflare-cert"
  custom_domain    = "%s"
}`

const testAccCheckAbrhaCDNConfig_Create_with_TTL = `
resource "abrha_spaces_bucket" "bucket" {
  name   = "%s"
  region = "ams3"
  acl    = "public-read"
}

resource "abrha_cdn" "foobar" {
  origin = abrha_spaces_bucket.bucket.bucket_domain_name
  ttl    = %d
}`

func testAccCheckAbrhaCDNConfig_CustomDomain(domain string, spaceName string, certName string) string {
	return fmt.Sprintf(`
resource "tls_private_key" "example" {
  algorithm = "RSA"
}

resource "tls_self_signed_cert" "example" {
  key_algorithm   = "RSA"
  private_key_pem = tls_private_key.example.private_key_pem
  dns_names       = ["foo.%s"]
  subject {
    common_name  = "foo.%s"
    organization = "%s"
  }

  validity_period_hours = 24

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
  ]
}

resource "abrha_spaces_bucket" "space" {
  name   = "%s"
  region = "sfo3"
}

resource "abrha_certificate" "spaces_cert" {
  name             = "%s"
  type             = "custom"
  private_key      = tls_private_key.example.private_key_pem
  leaf_certificate = tls_self_signed_cert.example.cert_pem

  lifecycle {
    create_before_destroy = true
  }
}

resource abrha_domain "domain" {
  name = "%s"
}

resource abrha_record "record" {
  domain = abrha_domain.domain.name
  type   = "CNAME"
  name   = "foo"
  value  = "${abrha_spaces_bucket.space.bucket_domain_name}."
}

resource "abrha_cdn" "space_cdn" {
  depends_on = [
    abrha_spaces_bucket.space,
    abrha_certificate.spaces_cert,
    abrha_record.record
  ]

  origin           = abrha_spaces_bucket.space.bucket_domain_name
  ttl              = 600
  certificate_name = abrha_certificate.spaces_cert.name
  custom_domain    = "foo.%s"
}`, domain, domain, certName, spaceName, certName, domain, domain)
}
