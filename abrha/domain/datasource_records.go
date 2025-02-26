package domain

import (
	"github.com/abrhacom/terraform-provider-abrha/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceAbrhaRecords() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema:        recordsSchema(),
		ResultAttributeName: "records",
		ExtraQuerySchema: map[string]*schema.Schema{
			"domain": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		FlattenRecord: flattenAbrhaRecord,
		GetRecords:    getAbrhaRecords,
	}

	return datalist.NewResource(dataListConfig)
}
