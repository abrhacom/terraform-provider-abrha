package sshkey

import (
	"github.com/abrhacom/terraform-provider-abrha/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceAbrhaSSHKeys() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema:        sshKeySchema(),
		ResultAttributeName: "ssh_keys",
		GetRecords:          getAbrhaSshKeys,
		FlattenRecord:       flattenAbrhaSshKey,
	}

	return datalist.NewResource(dataListConfig)
}
