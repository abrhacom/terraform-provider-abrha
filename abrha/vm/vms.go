package vm

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/abrhacom/terraform-provider-abrha/abrha/tag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func vmSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Description: "id of the Vm",
		},
		"name": {
			Type:        schema.TypeString,
			Description: "name of the Vm",
		},
		"created_at": {
			Type:        schema.TypeString,
			Description: "the creation date for the Vm",
		},
		"urn": {
			Type:        schema.TypeString,
			Description: "the uniform resource name for the Vm",
		},
		"region": {
			Type:        schema.TypeString,
			Description: "the region that the Vm instance is deployed in",
		},
		"image": {
			Type:        schema.TypeString,
			Description: "the image id or slug of the Vm",
		},
		"size": {
			Type:        schema.TypeString,
			Description: "the current size of the Vm",
		},
		"disk": {
			Type:        schema.TypeInt,
			Description: "the size of the Vms disk in gigabytes",
		},
		"vcpus": {
			Type:        schema.TypeInt,
			Description: "the number of virtual cpus",
		},
		"memory": {
			Type:        schema.TypeInt,
			Description: "memory of the Vm in megabytes",
		},
		"price_hourly": {
			Type:        schema.TypeFloat,
			Description: "the Vms hourly price",
		},
		"price_monthly": {
			Type:        schema.TypeFloat,
			Description: "the Vms monthly price",
		},
		"status": {
			Type:        schema.TypeString,
			Description: "state of the Vm instance",
		},
		"locked": {
			Type:        schema.TypeBool,
			Description: "whether the Vm has been locked",
		},
		"ipv4_address": {
			Type:        schema.TypeString,
			Description: "the Vms public ipv4 address",
		},
		"ipv4_address_private": {
			Type:        schema.TypeString,
			Description: "the Vms private ipv4 address",
		},
		"ipv6_address": {
			Type:        schema.TypeString,
			Description: "the Vms public ipv6 address",
		},
		"ipv6_address_private": {
			Type:        schema.TypeString,
			Description: "the Vms private ipv4 address",
		},
		"backups": {
			Type:        schema.TypeBool,
			Description: "whether the Vm has backups enabled",
		},
		"ipv6": {
			Type:        schema.TypeBool,
			Description: "whether the Vm has ipv6 enabled",
		},
		"private_networking": {
			Type:        schema.TypeBool,
			Description: "whether the Vm has private networking enabled",
		},
		"monitoring": {
			Type:        schema.TypeBool,
			Description: "whether the Vm has monitoring enabled",
		},
		"volume_ids": {
			Type:        schema.TypeSet,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "list of volumes attached to the Vm",
		},
		"tags": tag.TagsDataSourceSchema(),
		"vpc_uuid": {
			Type:        schema.TypeString,
			Description: "UUID of the VPC in which the Vm is located",
		},
	}
}

func getAbrhaVms(meta interface{}, extra map[string]interface{}) ([]interface{}, error) {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	gpus, _ := extra["gpus"].(bool)

	opts := &goApiAbrha.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	var vmList []interface{}

	for {
		var (
			vms  []goApiAbrha.Vm
			resp *goApiAbrha.Response
			err  error
		)
		if gpus {
			vms, resp, err = client.Vms.ListWithGPUs(context.Background(), opts)
		} else {
			vms, resp, err = client.Vms.List(context.Background(), opts)
		}

		if err != nil {
			return nil, fmt.Errorf("Error retrieving vms: %s", err)
		}

		for _, vm := range vms {
			vmList = append(vmList, vm)
		}

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, fmt.Errorf("Error retrieving vms: %s", err)
		}

		opts.Page = page + 1
	}

	return vmList, nil
}

func flattenAbrhaVm(rawVm, meta interface{}, extra map[string]interface{}) (map[string]interface{}, error) {
	vm := rawVm.(goApiAbrha.Vm)

	flattenedVm := map[string]interface{}{
		"id":            vm.ID,
		"name":          vm.Name,
		"urn":           vm.URN(),
		"region":        vm.Region.Slug,
		"size":          vm.Size.Slug,
		"price_hourly":  vm.Size.PriceHourly,
		"price_monthly": vm.Size.PriceMonthly,
		"disk":          vm.Disk,
		"vcpus":         vm.Vcpus,
		"memory":        vm.Memory,
		"status":        vm.Status,
		"locked":        vm.Locked,
		"created_at":    vm.Created,
		"vpc_uuid":      vm.VPCUUID,
	}

	if vm.Image.Slug == "" {
		flattenedVm["image"] = strconv.Itoa(vm.Image.ID)
	} else {
		flattenedVm["image"] = vm.Image.Slug
	}

	if publicIPv4 := FindIPv4AddrByType(&vm, "public"); publicIPv4 != "" {
		flattenedVm["ipv4_address"] = publicIPv4
	}

	if privateIPv4 := FindIPv4AddrByType(&vm, "private"); privateIPv4 != "" {
		flattenedVm["ipv4_address_private"] = privateIPv4
	}

	if publicIPv6 := FindIPv6AddrByType(&vm, "public"); publicIPv6 != "" {
		flattenedVm["ipv6_address"] = strings.ToLower(publicIPv6)
	}

	if privateIPv6 := FindIPv6AddrByType(&vm, "private"); privateIPv6 != "" {
		flattenedVm["ipv6_address_private"] = strings.ToLower(privateIPv6)
	}

	if features := vm.Features; features != nil {
		flattenedVm["backups"] = containsAbrhaVmFeature(features, "backups")
		flattenedVm["ipv6"] = containsAbrhaVmFeature(features, "ipv6")
		flattenedVm["private_networking"] = containsAbrhaVmFeature(features, "private_networking")
		flattenedVm["monitoring"] = containsAbrhaVmFeature(features, "monitoring")
	}

	flattenedVm["volume_ids"] = flattenAbrhaVmVolumeIds(vm.VolumeIDs)

	flattenedVm["tags"] = tag.FlattenTags(vm.Tags)

	return flattenedVm, nil
}
