package registry

import (
	"context"
	"fmt"
	"log"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const RegistryHostname = "registry.abrha.com"

func ResourceAbrhaContainerRegistry() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbrhaContainerRegistryCreate,
		ReadContext:   resourceAbrhaContainerRegistryRead,
		UpdateContext: resourceAbrhaContainerRegistryUpdate,
		DeleteContext: resourceAbrhaContainerRegistryDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"subscription_tier_slug": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"starter",
					"basic",
					"professional",
				}, false),
			},
			"region": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"server_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_usage_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAbrhaContainerRegistryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	// Build up our creation options
	opts := &goApiAbrha.RegistryCreateRequest{
		Name:                 d.Get("name").(string),
		SubscriptionTierSlug: d.Get("subscription_tier_slug").(string),
	}

	if region, ok := d.GetOk("region"); ok {
		opts.Region = region.(string)
	}

	log.Printf("[DEBUG] Container Registry create configuration: %#v", opts)
	reg, _, err := client.Registry.Create(context.Background(), opts)
	if err != nil {
		return diag.Errorf("Error creating container registry: %s", err)
	}

	d.SetId(reg.Name)
	log.Printf("[INFO] Container Registry: %s", reg.Name)

	return resourceAbrhaContainerRegistryRead(ctx, d, meta)
}

func resourceAbrhaContainerRegistryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	reg, resp, err := client.Registry.Get(context.Background())
	if err != nil {
		// If the registry is somehow already destroyed, mark as
		// successfully gone
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving container registry: %s", err)
	}

	d.SetId(reg.Name)
	d.Set("name", reg.Name)
	d.Set("region", reg.Region)
	d.Set("endpoint", fmt.Sprintf("%s/%s", RegistryHostname, reg.Name))
	d.Set("server_url", RegistryHostname)
	d.Set("created_at", reg.CreatedAt.UTC().String())
	d.Set("storage_usage_bytes", reg.StorageUsageBytes)

	sub, _, err := client.Registry.GetSubscription(context.Background())
	if err != nil {
		return diag.Errorf("Error retrieving container registry subscription: %s", err)
	}
	d.Set("subscription_tier_slug", sub.Tier.Slug)

	return nil
}

func resourceAbrhaContainerRegistryUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()
	if d.HasChange("subscription_tier_slug") {
		req := &goApiAbrha.RegistrySubscriptionUpdateRequest{
			TierSlug: d.Get("subscription_tier_slug").(string),
		}

		_, _, err := client.Registry.UpdateSubscription(ctx, req)
		if err != nil {
			return diag.Errorf("Error updating container registry subscription: %s", err)
		}
	}
	return resourceAbrhaContainerRegistryRead(ctx, d, meta)
}

func resourceAbrhaContainerRegistryDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	log.Printf("[INFO] Deleting container registry: %s", d.Id())
	_, err := client.Registry.Delete(context.Background())
	if err != nil {
		return diag.Errorf("Error deleting container registry: %s", err)
	}
	d.SetId("")
	return nil
}
