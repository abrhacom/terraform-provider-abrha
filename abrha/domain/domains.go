package domain

import (
	"context"
	"fmt"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func domainSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "name of the domain",
		},
		"urn": {
			Type:        schema.TypeString,
			Description: "the uniform resource name for the domain",
		},
		"ttl": {
			Type:        schema.TypeInt,
			Description: "ttl of the domain",
		},
	}
}

func getAbrhaDomains(meta interface{}, extra map[string]interface{}) ([]interface{}, error) {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	opts := &goApiAbrha.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	var allDomains []interface{}

	for {
		domains, resp, err := client.Domains.List(context.Background(), opts)

		if err != nil {
			return nil, fmt.Errorf("Error retrieving domains: %s", err)
		}

		for _, domain := range domains {
			allDomains = append(allDomains, domain)
		}

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, fmt.Errorf("Error retrieving domains: %s", err)
		}

		opts.Page = page + 1
	}

	return allDomains, nil
}

func flattenAbrhaDomain(rawDomain, meta interface{}, extra map[string]interface{}) (map[string]interface{}, error) {
	domain := rawDomain.(goApiAbrha.Domain)

	flattenedDomain := map[string]interface{}{
		"name": domain.Name,
		"urn":  domain.URN(),
		"ttl":  domain.TTL,
	}

	return flattenedDomain, nil
}
