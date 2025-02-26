package vm

import (
	"context"
	"fmt"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/abrhacom/terraform-provider-abrha/abrha/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceAbrhaVm() *schema.Resource {
	recordSchema := vmSchema()

	for _, f := range recordSchema {
		f.Computed = true
	}

	recordSchema["id"].ExactlyOneOf = []string{"id", "tag", "name"}
	recordSchema["id"].Optional = true
	recordSchema["name"].ExactlyOneOf = []string{"id", "tag", "name"}
	recordSchema["name"].Optional = true
	recordSchema["gpu"] = &schema.Schema{
		Type:          schema.TypeBool,
		Optional:      true,
		Default:       false,
		ConflictsWith: []string{"tag"},
	}

	recordSchema["tag"] = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "unique tag of the Vm",
		ValidateFunc: validation.NoZeroValues,
		ExactlyOneOf: []string{"id", "tag", "name"},
	}

	return &schema.Resource{
		ReadContext: dataSourceAbrhaVmRead,
		Schema:      recordSchema,
	}
}

func dataSourceAbrhaVmRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	var foundVm goApiAbrha.Vm

	if id, ok := d.GetOk("id"); ok {
		vm, _, err := client.Vms.Get(context.Background(), id.(string))
		if err != nil {
			return diag.FromErr(err)
		}

		foundVm = *vm
	} else if v, ok := d.GetOk("tag"); ok {
		vmList, err := getAbrhaVms(meta, nil)
		if err != nil {
			return diag.FromErr(err)
		}

		vm, err := findVmByTag(vmList, v.(string))
		if err != nil {
			return diag.FromErr(err)
		}

		foundVm = *vm
	} else if v, ok := d.GetOk("name"); ok {
		gpus := d.Get("gpu").(bool)
		extra := make(map[string]interface{})
		if gpus {
			extra["gpus"] = true
		}

		vmList, err := getAbrhaVms(meta, extra)
		if err != nil {
			return diag.FromErr(err)
		}

		vm, err := findVmByName(vmList, v.(string))

		if err != nil {
			return diag.FromErr(err)
		}

		foundVm = *vm
	} else {
		return diag.Errorf("Error: specify either a name, tag, or id to use to look up the vm")
	}

	flattenedVm, err := flattenAbrhaVm(foundVm, meta, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := util.SetResourceDataFromMap(d, flattenedVm); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(foundVm.ID)
	return nil
}

func findVmByName(vms []interface{}, name string) (*goApiAbrha.Vm, error) {
	results := make([]goApiAbrha.Vm, 0)
	for _, v := range vms {
		vm := v.(goApiAbrha.Vm)
		if vm.Name == name {
			results = append(results, vm)
		}
	}
	if len(results) == 1 {
		return &results[0], nil
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no vm found with name %s", name)
	}
	return nil, fmt.Errorf("too many vms found with name %s (found %d, expected 1)", name, len(results))
}

func findVmByTag(vms []interface{}, tag string) (*goApiAbrha.Vm, error) {
	results := make([]goApiAbrha.Vm, 0)
	for _, d := range vms {
		vm := d.(goApiAbrha.Vm)
		for _, t := range vm.Tags {
			if t == tag {
				results = append(results, vm)
			}
		}
	}
	if len(results) == 1 {
		return &results[0], nil
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no vm found with tag %s", tag)
	}
	return nil, fmt.Errorf("too many vms found with tag %s (found %d, expected 1)", tag, len(results))
}
