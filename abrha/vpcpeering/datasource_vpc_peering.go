package vpcpeering

import (
	"context"
	"fmt"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceAbrhaVPCPeering() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAbrhaVPCPeeringRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The ID of the VPC Peering",
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The name of the VPC Peering",
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
			},
			"vpc_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "The list of VPCs to be peered",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAbrhaVPCPeeringRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()
	var foundVPCPeering *goApiAbrha.VPCPeering

	if id, ok := d.GetOk("id"); ok {
		vpcPeering, _, err := client.VPCs.GetVPCPeering(context.Background(), id.(string))
		if err != nil {
			return diag.Errorf("Error retrieving VPC Peering: %s", err)
		}

		foundVPCPeering = vpcPeering
	} else if name, ok := d.GetOk("name"); ok {
		vpcPeerings, err := listVPCPeerings(client)
		if err != nil {
			return diag.Errorf("Error retrieving VPC Peering: %s", err)
		}

		vpcPeering, err := findVPCPeeringByName(vpcPeerings, name.(string))
		if err != nil {
			return diag.Errorf("Error retrieving VPC Peering: %s", err)
		}

		foundVPCPeering = vpcPeering
	}

	if foundVPCPeering == nil {
		return diag.Errorf("Bad Request: %s", fmt.Errorf("'name' or 'id' must be provided"))
	}

	d.SetId(foundVPCPeering.ID)
	d.Set("name", foundVPCPeering.Name)
	d.Set("vpc_ids", foundVPCPeering.VPCIDs)
	d.Set("status", foundVPCPeering.Status)
	d.Set("created_at", foundVPCPeering.CreatedAt.UTC().String())

	return nil
}

func listVPCPeerings(client *goApiAbrha.Client) ([]*goApiAbrha.VPCPeering, error) {
	peeringsList := []*goApiAbrha.VPCPeering{}
	opts := &goApiAbrha.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	for {
		peerings, resp, err := client.VPCs.ListVPCPeerings(context.Background(), opts)

		if err != nil {
			return peeringsList, fmt.Errorf("error retrieving VPC Peerings: %s", err)
		}

		peeringsList = append(peeringsList, peerings...)

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return peeringsList, fmt.Errorf("error retrieving VPC Peerings: %s", err)
		}

		opts.Page = page + 1
	}

	return peeringsList, nil
}

func findVPCPeeringByName(vpcPeerings []*goApiAbrha.VPCPeering, name string) (*goApiAbrha.VPCPeering, error) {
	for _, v := range vpcPeerings {
		if v.Name == name {
			return v, nil
		}
	}

	return nil, fmt.Errorf("no VPC Peerings found with name %s", name)
}
