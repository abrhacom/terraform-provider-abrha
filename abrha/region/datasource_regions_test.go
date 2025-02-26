package region_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAbrhaRegions_Basic(t *testing.T) {
	configNoFilter := `
data "abrha_regions" "all" {
}
`
	configAvailableFilter := `
data "abrha_regions" "filtered" {
  filter {
    key    = "available"
    values = ["true"]
  }
  sort {
    key = "slug"
  }
}
`

	configFeaturesFilter := `
data "abrha_regions" "filtered" {
  filter {
    key    = "features"
    values = ["private_networking", "backups"]
  }
  sort {
    key       = "available"
    direction = "desc"
  }
}
`

	configAllFilters := `
data "abrha_regions" "filtered" {
  filter {
    key    = "available"
    values = ["true"]
  }
  filter {
    key    = "features"
    values = ["private_networking", "backups"]
  }
}
`
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: configNoFilter,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.abrha_regions.all", "regions.#"),
					resource.TestCheckResourceAttrSet("data.abrha_regions.all", "regions.#"),
					acceptance.TestResourceInstanceState("data.abrha_regions.all", func(is *terraform.InstanceState) error {
						n, err := strconv.Atoi(is.Attributes["regions.#"])
						if err != nil {
							return err
						}

						for i := 0; i < n; i++ {
							key := fmt.Sprintf("regions.%d.slug", i)
							if _, ok := is.Attributes[key]; !ok {
								return fmt.Errorf("missing key in instance state for %s in %s", key, "data.abrha_regions.all")
							}
						}

						return nil
					}),
				),
			},
			{
				Config: configAvailableFilter,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.abrha_regions.filtered", "regions.#"),
				),
			},
			{
				Config: configFeaturesFilter,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.abrha_regions.filtered", "regions.#"),
				),
			},
			{
				Config: configAllFilters,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.abrha_regions.filtered", "regions.#"),
				),
			},
		},
	})
}
