package image

import (
	"github.com/abrhacom/terraform-provider-abrha/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceAbrhaImages() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema:        imageSchema(),
		ResultAttributeName: "images",
		FlattenRecord:       flattenAbrhaImage,
		GetRecords:          getAbrhaImages,
	}

	return datalist.NewResource(dataListConfig)
}
