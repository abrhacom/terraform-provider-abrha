package region

import (
	"github.com/abrhacom/terraform-provider-abrha/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceAbrhaRegions() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema: map[string]*schema.Schema{
			"slug": {
				Type: schema.TypeString,
			},
			"name": {
				Type: schema.TypeString,
			},
			"sizes": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"features": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"available": {
				Type: schema.TypeBool,
			},
		},
		ResultAttributeName: "regions",
		FlattenRecord:       flattenRegion,
		GetRecords:          getAbrhaRegions,
	}

	return datalist.NewResource(dataListConfig)
}
