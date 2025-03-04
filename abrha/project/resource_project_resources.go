package project

import (
	"context"

	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/abrhacom/terraform-provider-abrha/abrha/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceAbrhaProjectResources() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbrhaProjectResourcesUpdate,
		UpdateContext: resourceAbrhaProjectResourcesUpdate,
		ReadContext:   resourceAbrhaProjectResourcesRead,
		DeleteContext: resourceAbrhaProjectResourcesDelete,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "project ID",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"resources": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "the resources associated with the project",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAbrhaProjectResourcesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	projectId := d.Get("project").(string)

	_, resp, err := client.Projects.Get(context.Background(), projectId)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			// Project does not exist. Mark this resource as not existing.
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error while retrieving project %s: %v", projectId, err)
	}

	if d.HasChange("resources") {
		oldURNs, newURNs := d.GetChange("resources")
		remove, add := util.GetSetChanges(oldURNs.(*schema.Set), newURNs.(*schema.Set))

		if remove.Len() > 0 {
			_, err = assignResourcesToDefaultProject(client, remove)
			if err != nil {
				return diag.Errorf("Error assigning resources to default project: %s", err)
			}
		}

		if add.Len() > 0 {
			_, err = assignResourcesToProject(client, projectId, add)
			if err != nil {
				return diag.Errorf("Error assigning resources to project %s: %s", projectId, err)
			}
		}

		if err = d.Set("resources", newURNs); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(projectId)

	return resourceAbrhaProjectResourcesRead(ctx, d, meta)
}

func resourceAbrhaProjectResourcesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	projectId := d.Id()

	_, resp, err := client.Projects.Get(context.Background(), projectId)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			// Project does not exist. Mark this resource as not existing.
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error while retrieving project: %v", err)
	}

	if err = d.Set("project", projectId); err != nil {
		return diag.FromErr(err)
	}

	apiURNs, err := LoadResourceURNs(client, projectId)
	if err != nil {
		return diag.Errorf("Error while retrieving project resources: %s", err)
	}

	var newURNs []string

	configuredURNs := d.Get("resources").(*schema.Set).List()
	for _, rawConfiguredURN := range configuredURNs {
		configuredURN := rawConfiguredURN.(string)

		for _, apiURN := range *apiURNs {
			if configuredURN == apiURN {
				newURNs = append(newURNs, configuredURN)
			}
		}
	}

	if err = d.Set("resources", newURNs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceAbrhaProjectResourcesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	projectId := d.Get("project").(string)
	urns := d.Get("resources").(*schema.Set)

	_, resp, err := client.Projects.Get(context.Background(), projectId)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			// Project does not exist. Mark this resource as not existing.
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error while retrieving project: %s", err)
	}

	if urns.Len() > 0 {
		if _, err = assignResourcesToDefaultProject(client, urns); err != nil {
			return diag.Errorf("Error assigning resources to default project: %s", err)
		}
	}

	d.SetId("")
	return nil
}
