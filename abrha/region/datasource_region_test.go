package region_test

import (
	"regexp"
	"testing"

	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAbrhaRegion_Basic(t *testing.T) {
	config := `
data "abrha_region" "lon1" {
  slug = "lon1"
}
`
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.abrha_region.lon1", "slug", "lon1"),
					resource.TestCheckResourceAttrSet("data.abrha_region.lon1", "name"),
					resource.TestCheckResourceAttrSet("data.abrha_region.lon1", "available"),
					resource.TestCheckResourceAttrSet("data.abrha_region.lon1", "sizes.#"),
					resource.TestCheckResourceAttrSet("data.abrha_region.lon1", "features.#"),
				),
			},
		},
	})
}

func TestAccAbrhaRegion_MissingSlug(t *testing.T) {
	config := `
data "abrha_region" "xyz5" {
  slug = "xyz5"
}
`
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("Region does not exist: xyz5"),
			},
		},
	})
}
