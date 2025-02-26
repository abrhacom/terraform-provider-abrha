package app

import (
	"context"
	"fmt"
	"log"
	"time"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceAbrhaApp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbrhaAppCreate,
		ReadContext:   resourceAbrhaAppRead,
		UpdateContext: resourceAbrhaAppUpdate,
		DeleteContext: resourceAbrhaAppDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"spec": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "A Abrha App Platform Spec",
				Elem: &schema.Resource{
					Schema: appSpecSchema(true),
				},
			},

			"project_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			// Computed attributes
			"default_ingress": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The default URL to access the App",
			},

			"dedicated_ips": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "The dedicated egress IP addresses associated with the app.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: "The IP address of the dedicated egress IP.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: "The ID of the dedicated egress IP.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: "The status of the dedicated egress IP: 'UNKNOWN', 'ASSIGNING', 'ASSIGNED', or 'REMOVED'",
						},
					},
				},
			},

			"live_url": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"live_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// TODO: The full Deployment should be a data source, not a resource
			// specify the app id for the active deployment, include a deployment
			// id for a specific one
			"active_deployment_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID the App's currently active deployment",
			},

			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The uniform resource identifier for the app",
			},

			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time of when the App was last updated",
			},

			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time of when the App was created",
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},
	}
}

func resourceAbrhaAppCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()
	appCreateRequest := &goApiAbrha.AppCreateRequest{}
	appCreateRequest.Spec = expandAppSpec(d.Get("spec").([]interface{}))

	if v, ok := d.GetOk("project_id"); ok {
		appCreateRequest.ProjectID = v.(string)
	}

	log.Printf("[DEBUG] App create request: %#v", appCreateRequest)
	app, _, err := client.Apps.Create(context.Background(), appCreateRequest)
	if err != nil {
		return diag.Errorf("Error creating App: %s", err)
	}

	d.SetId(app.ID)
	log.Printf("[DEBUG] Waiting for app (%s) deployment to become active", app.ID)
	timeout := d.Timeout(schema.TimeoutCreate)
	err = waitForAppDeployment(client, app.ID, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] App created, ID: %s", d.Id())

	return resourceAbrhaAppRead(ctx, d, meta)
}

func resourceAbrhaAppRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	app, resp, err := client.Apps.Get(context.Background(), d.Id())
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[DEBUG] App (%s) was not found - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error reading App: %s", err)
	}

	d.SetId(app.ID)
	d.Set("default_ingress", app.DefaultIngress)
	d.Set("live_url", app.LiveURL)
	d.Set("live_domain", app.LiveDomain)
	d.Set("updated_at", app.UpdatedAt.UTC().String())
	d.Set("created_at", app.CreatedAt.UTC().String())
	d.Set("urn", app.URN())
	d.Set("project_id", app.ProjectID)

	if app.DedicatedIps != nil {
		d.Set("dedicated_ips", appDedicatedIps(d, app))
	}

	if err := d.Set("spec", flattenAppSpec(d, app.Spec)); err != nil {
		return diag.Errorf("Error setting app spec: %#v", err)
	}

	if app.ActiveDeployment != nil {
		d.Set("active_deployment_id", app.ActiveDeployment.ID)
	} else {
		deploymentWarning := diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  fmt.Sprintf("No active deployment found for app: %s (%s)", app.Spec.Name, app.ID),
		}
		d.Set("active_deployment_id", "")
		return diag.Diagnostics{deploymentWarning}
	}

	return nil
}

func appDedicatedIps(d *schema.ResourceData, app *goApiAbrha.App) []interface{} {
	remote := make([]interface{}, 0, len(app.DedicatedIps))
	for _, change := range app.DedicatedIps {
		rawChange := map[string]interface{}{
			"ip":     change.Ip,
			"id":     change.ID,
			"status": change.Status,
		}
		remote = append(remote, rawChange)
	}
	return remote
}

func resourceAbrhaAppUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	if d.HasChange("spec") {
		appUpdateRequest := &goApiAbrha.AppUpdateRequest{}
		appUpdateRequest.Spec = expandAppSpec(d.Get("spec").([]interface{}))

		app, _, err := client.Apps.Update(context.Background(), d.Id(), appUpdateRequest)
		if err != nil {
			return diag.Errorf("Error updating app (%s): %s", d.Id(), err)
		}

		log.Printf("[DEBUG] Waiting for app (%s) deployment to become active", app.ID)
		timeout := d.Timeout(schema.TimeoutCreate)
		err = waitForAppDeployment(client, app.ID, timeout)
		if err != nil {
			return diag.FromErr(err)
		}

		log.Printf("[INFO] Updated app (%s)", app.ID)
	}

	return resourceAbrhaAppRead(ctx, d, meta)
}

func resourceAbrhaAppDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	log.Printf("[INFO] Deleting App: %s", d.Id())
	_, err := client.Apps.Delete(context.Background(), d.Id())
	if err != nil {
		return diag.Errorf("Error deletingApp: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForAppDeployment(client *goApiAbrha.Client, id string, timeout time.Duration) error {
	tickerInterval := 10 //10s
	timeoutSeconds := int(timeout.Seconds())
	n := 0

	var deploymentID string
	ticker := time.NewTicker(time.Duration(tickerInterval) * time.Second)
	for range ticker.C {
		if n*tickerInterval > timeoutSeconds {
			ticker.Stop()
			break
		}

		if deploymentID == "" {
			// The InProgressDeployment is generally not known and returned as
			// part of the initial response to the request. For config updates
			// (as opposed to updates to the app's source), the "deployment"
			// can complete before the first time we poll the app. We can not
			// know if the InProgressDeployment has not started or if it has
			// already completed. So instead we need to list all of the
			// deployments for the application.
			opts := &goApiAbrha.ListOptions{PerPage: 20}
			deployments, _, err := client.Apps.ListDeployments(context.Background(), id, opts)
			if err != nil {
				return fmt.Errorf("Error trying to read app deployment state: %s", err)
			}

			// We choose the most recent deployment. Note that there is a possibility
			// that the deployment has not been created yet. If that is true,
			// we will do the wrong thing here and test the status of a previously
			// completed deployment and exit. However there is no better way to
			// correlate a deployment with the request that triggered it.
			if len(deployments) > 0 {
				deploymentID = deployments[0].ID
			}
		} else {
			deployment, _, err := client.Apps.GetDeployment(context.Background(), id, deploymentID)
			if err != nil {
				ticker.Stop()
				return fmt.Errorf("Error trying to read app deployment state: %s", err)
			}

			allSuccessful := true
			for _, step := range deployment.Progress.Steps {
				if step.Status != goApiAbrha.DeploymentProgressStepStatus_Success {
					allSuccessful = false
					break
				}
			}

			if allSuccessful {
				ticker.Stop()
				return nil
			}

			if deployment.Progress.ErrorSteps > 0 {
				ticker.Stop()
				return fmt.Errorf("error deploying app (%s) (deployment ID: %s):\n%s", id, deployment.ID, goApiAbrha.Stringify(deployment.Progress))
			}

			log.Printf("[DEBUG] Waiting for app (%s) deployment (%s) to become active. Phase: %s (%d/%d)",
				id, deployment.ID, deployment.Phase, deployment.Progress.SuccessSteps, deployment.Progress.TotalSteps)
		}

		n++
	}

	return fmt.Errorf("timeout waiting for app (%s) deployment", id)
}
