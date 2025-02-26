package project

import (
	"github.com/abrhacom/terraform-provider-abrha/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceAbrhaProjects() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema:        projectSchema(),
		ResultAttributeName: "projects",
		FlattenRecord:       flattenAbrhaProject,
		GetRecords:          getAbrhaProjects,
	}

	return datalist.NewResource(dataListConfig)
}
