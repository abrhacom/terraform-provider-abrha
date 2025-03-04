package image

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/abrhacom/terraform-provider-abrha/abrha/tag"
	"github.com/abrhacom/terraform-provider-abrha/abrha/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Ref: https://developers.parspack.com/documentation/v2/#retrieve-an-existing-image-by-id
const (
	ImageAvailableStatus = "available"
	ImageDeletedStatus   = "deleted"
)

func ResourceAbrhaCustomImage() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceAbrhaCustomImageRead,
		CreateContext: resourceAbrhaCustomImageCreate,
		UpdateContext: resourceAbrhaCustomImageUpdate,
		DeleteContext: resourceAbrhaCustomImageDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"url": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"regions": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"distribution": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Unknown",
				ValidateFunc: validation.StringInSlice(validImageDistributions(), false),
			},
			"tags": tag.TagsSchema(),
			"image_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"slug": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"min_disk_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"size_gigabytes": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
		},

		// Images can not currently be removed from a region.
		CustomizeDiff: customdiff.ForceNewIfChange("regions", func(ctx context.Context, old, new, meta interface{}) bool {
			remove, _ := util.GetSetChanges(old.(*schema.Set), new.(*schema.Set))
			return len(remove.List()) > 0
		}),
	}
}

func resourceAbrhaCustomImageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	// We import the image to the first region. We can distribute it to others once it is available.
	regions := d.Get("regions").(*schema.Set).List()
	region := regions[0].(string)

	imageCreateRequest := goApiAbrha.CustomImageCreateRequest{
		Name:   d.Get("name").(string),
		Url:    d.Get("url").(string),
		Region: region,
	}

	if desc, ok := d.GetOk("description"); ok {
		imageCreateRequest.Description = desc.(string)
	}

	if dist, ok := d.GetOk("distribution"); ok {
		imageCreateRequest.Distribution = dist.(string)
	}

	if tags, ok := d.GetOk("tags"); ok {
		imageCreateRequest.Tags = tag.ExpandTags(tags.(*schema.Set).List())
	}

	imageResponse, _, err := client.Images.Create(ctx, &imageCreateRequest)
	if err != nil {
		return diag.Errorf("Error creating custom image: %s", err)
	}

	id := strconv.Itoa(imageResponse.ID)
	d.SetId(id)

	_, err = waitForImage(ctx, d, ImageAvailableStatus, imagePendingStatuses(), "status", meta)
	if err != nil {
		return diag.Errorf("Error waiting for image (%s) to become ready: %s", d.Id(), err)
	}

	if len(regions) > 1 {
		// Remove the first region from the slice as the image is already there.
		regions[0] = regions[len(regions)-1]
		regions[len(regions)-1] = ""
		regions = regions[:len(regions)-1]
		log.Printf("[INFO] Image available in: %s Distributing to: %v", region, regions)
		err = distributeImageToRegions(client, imageResponse.ID, regions)
		if err != nil {
			return diag.Errorf("Error distributing image (%s) to additional regions: %s", d.Id(), err)
		}
	}

	return resourceAbrhaCustomImageRead(ctx, d, meta)
}

func resourceAbrhaCustomImageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	imageID := d.Id()

	id, err := strconv.Atoi(imageID)
	if err != nil {
		return diag.Errorf("Error converting id %s to string: %s", imageID, err)
	}

	imageResponse, _, err := client.Images.GetByID(ctx, id)
	if err != nil {
		return diag.Errorf("Error retrieving image with id %s: %s", imageID, err)
	}
	// Set status as deleted if image is deleted
	if imageResponse.Status == ImageDeletedStatus {
		d.SetId("")
		return nil
	}
	d.Set("image_id", imageResponse.ID)
	d.Set("name", imageResponse.Name)
	d.Set("type", imageResponse.Type)
	d.Set("distribution", imageResponse.Distribution)
	d.Set("slug", imageResponse.Slug)
	d.Set("public", imageResponse.Public)
	d.Set("regions", imageResponse.Regions)
	d.Set("min_disk_size", imageResponse.MinDiskSize)
	d.Set("size_gigabytes", imageResponse.SizeGigaBytes)
	d.Set("created_at", imageResponse.Created)
	d.Set("description", imageResponse.Description)
	if err := d.Set("tags", tag.FlattenTags(imageResponse.Tags)); err != nil {
		return diag.Errorf("Error setting `tags`: %+v", err)
	}
	d.Set("status", imageResponse.Status)
	return nil
}

func resourceAbrhaCustomImageUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	imageID := d.Id()

	id, err := strconv.Atoi(imageID)
	if err != nil {
		return diag.Errorf("Error converting id %s to string: %s", imageID, err)
	}

	if d.HasChanges("name", "description", "distribution") {
		imageName := d.Get("name").(string)
		imageUpdateRequest := &goApiAbrha.ImageUpdateRequest{
			Name:         imageName,
			Distribution: d.Get("distribution").(string),
			Description:  d.Get("description").(string),
		}

		_, _, err := client.Images.Update(ctx, id, imageUpdateRequest)
		if err != nil {
			return diag.Errorf("Error updating image %s, name %s: %s", imageID, imageName, err)
		}
	}

	if d.HasChange("regions") {
		old, new := d.GetChange("regions")
		_, add := util.GetSetChanges(old.(*schema.Set), new.(*schema.Set))
		err = distributeImageToRegions(client, id, add.List())
		if err != nil {
			return diag.Errorf("Error distributing image (%s) to additional regions: %s", d.Id(), err)
		}
	}

	return resourceAbrhaCustomImageRead(ctx, d, meta)
}

func resourceAbrhaCustomImageDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	imageID := d.Id()

	id, err := strconv.Atoi(imageID)
	if err != nil {
		return diag.Errorf("Error converting id %s to string: %s", imageID, err)
	}
	_, err = client.Images.Delete(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "Image Can not delete an already deleted image.") {
			log.Printf("[INFO] Image %s no longer exists", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting image id %s: %s", imageID, err)
	}
	return nil
}

func waitForImage(ctx context.Context, d *schema.ResourceData, target string, pending []string, attribute string, meta interface{}) (interface{}, error) {
	log.Printf("[INFO] Waiting for image (%s) to have %s of %s", d.Id(), attribute, target)
	stateConf := &retry.StateChangeConf{
		Pending:    pending,
		Target:     []string{target},
		Refresh:    imageStateRefreshFunc(ctx, d, attribute, meta),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 60 * time.Second,
	}

	return stateConf.WaitForStateContext(ctx)
}

func imageStateRefreshFunc(ctx context.Context, d *schema.ResourceData, state string, meta interface{}) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

		imageID := d.Id()

		id, err := strconv.Atoi(imageID)
		if err != nil {
			return nil, "", err
		}

		imageResponse, _, err := client.Images.GetByID(ctx, id)
		if err != nil {
			return nil, "", err
		}

		if imageResponse.Status == ImageDeletedStatus {
			return nil, "", fmt.Errorf(imageResponse.ErrorMessage)
		}

		return imageResponse, imageResponse.Status, nil
	}
}

func distributeImageToRegions(client *goApiAbrha.Client, imageId int, regions []interface{}) (err error) {
	for _, region := range regions {
		transferRequest := &goApiAbrha.ActionRequest{
			"type":   "transfer",
			"region": region.(string),
		}

		log.Printf("[INFO] Transferring image (%d) to: %s", imageId, region)
		action, _, err := client.ImageActions.Transfer(context.TODO(), imageId, transferRequest)
		if err != nil {
			return err
		}

		err = util.WaitForAction(client, action)
		if err != nil {
			return err
		}
	}

	return nil
}

// Ref: https://developers.parspack.com/documentation/v2/#retrieve-an-existing-image-by-id
func imagePendingStatuses() []string {
	return []string{"NEW", "pending"}
}

// Ref:https://developers.parspack.com/documentation/v2/#create-a-custom-image
func validImageDistributions() []string {
	return []string{
		"Arch Linux",
		"CentOS",
		"CoreOS",
		"Debian",
		"Fedora",
		"Fedora Atomic",
		"FreeBSD",
		"Gentoo",
		"openSUSE",
		"RancherOS",
		"Ubuntu",
		"Unknown",
		"Unknown OS",
	}
}
