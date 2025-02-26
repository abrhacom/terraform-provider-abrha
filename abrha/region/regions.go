package region

import (
	"context"
	"fmt"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func getAbrhaRegions(meta interface{}, extra map[string]interface{}) ([]interface{}, error) {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	allRegions := []interface{}{}

	opts := &goApiAbrha.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	for {
		regions, resp, err := client.Regions.List(context.Background(), opts)

		if err != nil {
			return nil, fmt.Errorf("Error retrieving regions: %s", err)
		}

		for _, region := range regions {
			allRegions = append(allRegions, region)
		}

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, fmt.Errorf("Error retrieving regions: %s", err)
		}

		opts.Page = page + 1
	}

	return allRegions, nil
}

func flattenRegion(rawRegion, meta interface{}, extra map[string]interface{}) (map[string]interface{}, error) {
	region := rawRegion.(goApiAbrha.Region)

	flattenedRegion := map[string]interface{}{}
	flattenedRegion["slug"] = region.Slug
	flattenedRegion["name"] = region.Name
	flattenedRegion["available"] = region.Available

	sizesSet := schema.NewSet(schema.HashString, []interface{}{})
	for _, size := range region.Sizes {
		sizesSet.Add(size)
	}
	flattenedRegion["sizes"] = sizesSet

	featuresSet := schema.NewSet(schema.HashString, []interface{}{})
	for _, feature := range region.Features {
		featuresSet.Add(feature)
	}
	flattenedRegion["features"] = featuresSet

	return flattenedRegion, nil
}
