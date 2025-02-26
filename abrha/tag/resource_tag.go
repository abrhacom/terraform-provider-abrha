package tag

import (
	"context"
	"log"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceAbrhaTag() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbrhaTagCreate,
		ReadContext:   resourceAbrhaTagRead,
		DeleteContext: resourceAbrhaTagDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateTag,
			},
			"total_resource_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vms_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"images_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"volumes_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"volume_snapshots_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"databases_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceAbrhaTagCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	// Build up our creation options
	opts := &goApiAbrha.TagCreateRequest{
		Name: d.Get("name").(string),
	}

	log.Printf("[DEBUG] Tag create configuration: %#v", opts)
	tag, _, err := client.Tags.Create(context.Background(), opts)
	if err != nil {
		return diag.Errorf("Error creating tag: %s", err)
	}

	d.SetId(tag.Name)
	log.Printf("[INFO] Tag: %s", tag.Name)

	return resourceAbrhaTagRead(ctx, d, meta)
}

func resourceAbrhaTagRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	tag, resp, err := client.Tags.Get(context.Background(), d.Id())
	if err != nil {
		// If the tag is somehow already destroyed, mark as
		// successfully gone
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving tag: %s", err)
	}

	d.Set("name", tag.Name)
	d.Set("total_resource_count", tag.Resources.Count)
	d.Set("vms_count", tag.Resources.Vms.Count)
	d.Set("images_count", tag.Resources.Images.Count)
	d.Set("volumes_count", tag.Resources.Volumes.Count)
	d.Set("volume_snapshots_count", tag.Resources.VolumeSnapshots.Count)
	d.Set("databases_count", tag.Resources.Databases.Count)

	return nil
}

func resourceAbrhaTagDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	log.Printf("[INFO] Deleting tag: %s", d.Id())
	_, err := client.Tags.Delete(context.Background(), d.Id())
	if err != nil {
		return diag.Errorf("Error deleting tag: %s", err)
	}

	d.SetId("")
	return nil
}
