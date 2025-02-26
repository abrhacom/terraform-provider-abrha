package spaces

import (
	"github.com/abrhacom/terraform-provider-abrha/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceAbrhaSpacesBuckets() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema:        spacesBucketSchema(),
		ResultAttributeName: "buckets",
		FlattenRecord:       flattenSpacesBucket,
		GetRecords:          getAbrhaBuckets,
	}

	return datalist.NewResource(dataListConfig)
}
