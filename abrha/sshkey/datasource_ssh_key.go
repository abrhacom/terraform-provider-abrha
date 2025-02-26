package sshkey

import (
	"context"
	"fmt"
	"strconv"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceAbrhaSSHKey() *schema.Resource {
	recordSchema := sshKeySchema()

	for _, f := range recordSchema {
		f.Computed = true
	}

	recordSchema["name"].Required = true
	recordSchema["name"].Computed = false

	return &schema.Resource{
		ReadContext: dataSourceAbrhaSSHKeyRead,
		Schema:      recordSchema,
	}
}

func dataSourceAbrhaSSHKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keyList, err := getAbrhaSshKeys(meta, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	key, err := findSshKeyByName(keyList, d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	flattenedKey, err := flattenAbrhaSshKey(*key, meta, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := util.SetResourceDataFromMap(d, flattenedKey); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(key.ID))

	return nil
}

func findSshKeyByName(keys []interface{}, name string) (*goApiAbrha.Key, error) {
	results := make([]goApiAbrha.Key, 0)
	for _, v := range keys {
		key := v.(goApiAbrha.Key)
		if key.Name == name {
			results = append(results, key)
		}
	}
	if len(results) == 1 {
		return &results[0], nil
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no ssh key found with name %s", name)
	}
	return nil, fmt.Errorf("too many ssh keys found with name %s (found %d, expected 1)", name, len(results))
}
