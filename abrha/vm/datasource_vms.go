package vm

import (
	"github.com/abrhacom/terraform-provider-abrha/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceAbrhaVms() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema:        vmSchema(),
		ResultAttributeName: "vms",
		GetRecords:          getAbrhaVms,
		FlattenRecord:       flattenAbrhaVm,
		ExtraQuerySchema: map[string]*schema.Schema{
			"gpus": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}

	return datalist.NewResource(dataListConfig)
}
