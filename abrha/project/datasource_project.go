package project

import (
	"context"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/abrhacom/terraform-provider-abrha/abrha/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceAbrhaProject() *schema.Resource {
	recordSchema := projectSchema()

	for _, f := range recordSchema {
		f.Computed = true
	}

	recordSchema["id"].ConflictsWith = []string{"name"}
	recordSchema["id"].Optional = true
	recordSchema["name"].ConflictsWith = []string{"id"}
	recordSchema["name"].Optional = true

	return &schema.Resource{
		ReadContext: dataSourceAbrhaProjectRead,
		Schema:      recordSchema,
	}
}

func dataSourceAbrhaProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	// Load the specified project, otherwise load the default project.
	var foundProject *goApiAbrha.Project
	if projectId, ok := d.GetOk("id"); ok {
		thisProject, _, err := client.Projects.Get(context.Background(), projectId.(string))
		if err != nil {
			return diag.Errorf("Unable to load project ID %s: %s", projectId, err)
		}
		foundProject = thisProject
	} else if name, ok := d.GetOk("name"); ok {
		projects, err := getAbrhaProjects(meta, nil)
		if err != nil {
			return diag.Errorf("Unable to load projects: %s", err)
		}

		var projectsWithName []goApiAbrha.Project
		for _, p := range projects {
			project := p.(goApiAbrha.Project)
			if project.Name == name.(string) {
				projectsWithName = append(projectsWithName, project)
			}
		}
		if len(projectsWithName) == 0 {
			return diag.Errorf("No projects found with name '%s'", name)
		} else if len(projectsWithName) > 1 {
			return diag.Errorf("Multiple projects found with name '%s'", name)
		}

		// Single result so choose that project.
		foundProject = &projectsWithName[0]
	} else {
		defaultProject, _, err := client.Projects.GetDefault(context.Background())
		if err != nil {
			return diag.Errorf("Unable to load default project: %s", err)
		}
		foundProject = defaultProject
	}

	if foundProject == nil {
		return diag.Errorf("No project found.")
	}

	flattenedProject, err := flattenAbrhaProject(*foundProject, meta, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := util.SetResourceDataFromMap(d, flattenedProject); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(foundProject.ID)
	return nil
}
