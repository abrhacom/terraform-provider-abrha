package domain

import (
	"github.com/abrhacom/terraform-provider-abrha/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceAbrhaDomains() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema:        domainSchema(),
		ResultAttributeName: "domains",
		GetRecords:          getAbrhaDomains,
		FlattenRecord:       flattenAbrhaDomain,
	}

	return datalist.NewResource(dataListConfig)
}
