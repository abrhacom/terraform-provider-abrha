package certificate_test

import (
	"testing"

	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAbrhaCertificate_importBasic(t *testing.T) {
	resourceName := "abrha_certificate.foobar"
	name := acceptance.RandomTestName("certificate")
	privateKeyMaterial, leafCertMaterial, certChainMaterial := acceptance.GenerateTestCertMaterial(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaCertificateConfig_basic(name, privateKeyMaterial, leafCertMaterial, certChainMaterial),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"certificate_chain", "leaf_certificate", "private_key"}, // We ignore these as they are not returned by the API

			},
		},
	})
}
