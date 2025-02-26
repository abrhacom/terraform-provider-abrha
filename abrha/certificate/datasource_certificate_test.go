package certificate_test

import (
	"fmt"
	"testing"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/abrhacom/terraform-provider-abrha/abrha/certificate"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceAbrhaCertificate_Basic(t *testing.T) {
	var certificate goApiAbrha.Certificate
	name := acceptance.RandomTestName("certificate")

	privateKeyMaterial, leafCertMaterial, certChainMaterial := acceptance.GenerateTestCertMaterial(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceAbrhaCertificateConfig_basic(name, privateKeyMaterial, leafCertMaterial, certChainMaterial, false),
			},
			{
				Config: testAccCheckDataSourceAbrhaCertificateConfig_basic(name, privateKeyMaterial, leafCertMaterial, certChainMaterial, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceAbrhaCertificateExists("data.abrha_certificate.foobar", &certificate),
					resource.TestCheckResourceAttr(
						"data.abrha_certificate.foobar", "id", name),
					resource.TestCheckResourceAttr(
						"data.abrha_certificate.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"data.abrha_certificate.foobar", "type", "custom"),
				),
			},
		},
	})
}

func testAccCheckDataSourceAbrhaCertificateExists(n string, cert *goApiAbrha.Certificate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No certificate ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GoApiAbrhaClient()

		foundCertificate, err := certificate.FindCertificateByName(client, rs.Primary.ID)
		if err != nil {
			return err
		}

		*cert = *foundCertificate

		return nil
	}
}

func testAccCheckDataSourceAbrhaCertificateConfig_basic(
	name, privateKeyMaterial, leafCert, certChain string,
	includeDataSource bool,
) string {
	config := fmt.Sprintf(`
resource "abrha_certificate" "foo" {
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
`, name, privateKeyMaterial, leafCert, certChain)

	if includeDataSource {
		config += `
data "abrha_certificate" "foobar" {
  name = abrha_certificate.foo.name
}
`
	}

	return config
}
