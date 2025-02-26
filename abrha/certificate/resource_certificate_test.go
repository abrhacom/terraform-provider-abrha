package certificate_test

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/abrhacom/terraform-provider-abrha/abrha/certificate"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/abrhacom/terraform-provider-abrha/abrha/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func testCertificateStateDataV0() map[string]interface{} {
	return map[string]interface{}{
		"name": "test",
		"id":   "aaa-bbb-123-ccc",
	}
}

func testCertificateStateDataV1() map[string]interface{} {
	v0 := testCertificateStateDataV0()
	return map[string]interface{}{
		"name": v0["name"],
		"uuid": v0["id"],
		"id":   v0["name"],
	}
}

func TestResourceExampleInstanceStateUpgradeV0(t *testing.T) {
	expected := testCertificateStateDataV1()
	actual, err := certificate.MigrateCertificateStateV0toV1(context.Background(), testCertificateStateDataV0(), nil)
	if err != nil {
		t.Fatalf("error migrating state: %s", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("\n\nexpected:\n\n%#v\n\ngot:\n\n%#v\n\n", actual, expected)
	}
}

func TestAccAbrhaCertificate_Basic(t *testing.T) {
	var cert goApiAbrha.Certificate
	name := acceptance.RandomTestName("certificate")
	privateKeyMaterial, leafCertMaterial, certChainMaterial := acceptance.GenerateTestCertMaterial(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaCertificateConfig_basic(name, privateKeyMaterial, leafCertMaterial, certChainMaterial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaCertificateExists("abrha_certificate.foobar", &cert),
					resource.TestCheckResourceAttr(
						"abrha_certificate.foobar", "id", name),
					resource.TestCheckResourceAttr(
						"abrha_certificate.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"abrha_certificate.foobar", "private_key", util.HashString(fmt.Sprintf("%s\n", privateKeyMaterial))),
					resource.TestCheckResourceAttr(
						"abrha_certificate.foobar", "leaf_certificate", util.HashString(fmt.Sprintf("%s\n", leafCertMaterial))),
					resource.TestCheckResourceAttr(
						"abrha_certificate.foobar", "certificate_chain", util.HashString(fmt.Sprintf("%s\n", certChainMaterial))),
				),
			},
		},
	})
}

func TestAccAbrhaCertificate_ExpectedErrors(t *testing.T) {
	name := acceptance.RandomTestName("certificate")
	privateKeyMaterial, leafCertMaterial, certChainMaterial := acceptance.GenerateTestCertMaterial(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckAbrhaCertificateConfig_customNoLeaf(name, privateKeyMaterial, certChainMaterial),
				ExpectError: regexp.MustCompile("`leaf_certificate` is required for when type is `custom` or empty"),
			},
			{
				Config:      testAccCheckAbrhaCertificateConfig_customNoKey(name, leafCertMaterial, certChainMaterial),
				ExpectError: regexp.MustCompile("`private_key` is required for when type is `custom` or empty"),
			},
			{
				Config:      testAccCheckAbrhaCertificateConfig_noDomains(name),
				ExpectError: regexp.MustCompile("`domains` is required for when type is `lets_encrypt`"),
			},
		},
	})
}

func testAccCheckAbrhaCertificateDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "abrha_certificate" {
			continue
		}

		_, err := certificate.FindCertificateByName(client, rs.Primary.ID)

		if err != nil && !strings.Contains(err.Error(), "not found") {
			return fmt.Errorf(
				"Error waiting for certificate (%s) to be destroyed: %s",
				rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckAbrhaCertificateExists(n string, cert *goApiAbrha.Certificate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Certificate ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

		c, err := certificate.FindCertificateByName(client, rs.Primary.ID)
		if err != nil {
			return err
		}

		*cert = *c

		return nil
	}
}

func testAccCheckAbrhaCertificateConfig_basic(name, privateKeyMaterial, leafCert, certChain string) string {
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
}`, name, privateKeyMaterial, leafCert, certChain)
}

func testAccCheckAbrhaCertificateConfig_customNoLeaf(name, privateKeyMaterial, certChain string) string {
	return fmt.Sprintf(`
resource "abrha_certificate" "foobar" {
  name              = "%s"
  private_key       = <<EOF
%s
EOF
  certificate_chain = <<EOF
%s
EOF
}`, name, privateKeyMaterial, certChain)
}

func testAccCheckAbrhaCertificateConfig_customNoKey(name, leafCert, certChain string) string {
	return fmt.Sprintf(`
resource "abrha_certificate" "foobar" {
  name              = "%s"
  leaf_certificate  = <<EOF
%s
EOF
  certificate_chain = <<EOF
%s
EOF
}`, name, leafCert, certChain)
}

func testAccCheckAbrhaCertificateConfig_noDomains(name string) string {
	return fmt.Sprintf(`
resource "abrha_certificate" "foobar" {
  name = "%s"
  type = "lets_encrypt"
}`, name)
}
