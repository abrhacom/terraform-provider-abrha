package vpc

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/abrhacom/terraform-provider-abrha/internal/mutexkv"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var mutexKV = mutexkv.NewMutexKV()

func ResourceAbrhaVPC() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbrhaVPCCreate,
		ReadContext:   resourceAbrhaVPCRead,
		UpdateContext: resourceAbrhaVPCUpdate,
		DeleteContext: resourceAbrhaVPCDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the VPC",
				ValidateFunc: validation.NoZeroValues,
			},
			"region": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Abrha region slug for the VPC's location",
				ValidateFunc: validation.NoZeroValues,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "A free-form description for the VPC",
				ValidateFunc: validation.StringLenBetween(0, 255),
			},
			"ip_range": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The range of IP addresses for the VPC in CIDR notation",
				//ValidateFunc: validation.IsCIDR,
			},

			// Computed attributes
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The uniform resource name (URN) for the VPC",
			},
			"default": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether or not the VPC is the default one for the region",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time of when the VPC was created",
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
	}
}

func resourceAbrhaVPCCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	region := d.Get("region").(string)
	vpcRequest := &goApiAbrha.VPCCreateRequest{
		Name:       d.Get("name").(string),
		RegionSlug: region,
	}

	if v, ok := d.GetOk("description"); ok {
		vpcRequest.Description = v.(string)
	}

	if v, ok := d.GetOk("ip_range"); ok {
		vpcRequest.IPRange = v.(string)
	}

	// Prevent parallel creation of VPCs in the same region to protect
	// against race conditions in IP range assignment.
	key := fmt.Sprintf("resource_abrha_vpc/%s", region)
	mutexKV.Lock(key)
	defer mutexKV.Unlock(key)

	log.Printf("[DEBUG] VPC create request: %#v", vpcRequest)
	vpc, _, err := client.VPCs.Create(context.Background(), vpcRequest)
	if err != nil {
		return diag.Errorf("Error creating VPC: %s", err)
	}

	d.SetId(vpc.ID)
	log.Printf("[INFO] VPC created, ID: %s", d.Id())

	return resourceAbrhaVPCRead(ctx, d, meta)
}

func resourceAbrhaVPCRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	vpc, resp, err := client.VPCs.Get(context.Background(), d.Id())

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[DEBUG] VPC  (%s) was not found - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error reading VPC: %s", err)
	}

	d.SetId(vpc.ID)
	d.Set("name", vpc.Name)
	d.Set("region", vpc.RegionSlug)
	d.Set("description", vpc.Description)
	d.Set("ip_range", vpc.IPRange)
	d.Set("urn", vpc.URN)
	d.Set("default", vpc.Default)
	d.Set("created_at", vpc.CreatedAt)

	return nil
}

func resourceAbrhaVPCUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	if d.HasChanges("name", "description") {
		vpcUpdateRequest := &goApiAbrha.VPCUpdateRequest{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			Default:     goApiAbrha.PtrTo(d.Get("default").(bool)),
		}
		_, _, err := client.VPCs.Update(context.Background(), d.Id(), vpcUpdateRequest)

		if err != nil {
			return diag.Errorf("Error updating VPC : %s", err)
		}
		log.Printf("[INFO] Updated VPC")
	}

	return resourceAbrhaVPCRead(ctx, d, meta)
}

func resourceAbrhaVPCDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()
	vpcID := d.Id()

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		resp, err := client.VPCs.Delete(context.Background(), vpcID)
		if err != nil {
			// Retry if VPC still contains member resources to prevent race condition
			if resp.StatusCode == http.StatusForbidden {
				return retry.RetryableError(err)
			} else {
				return retry.NonRetryableError(fmt.Errorf("Error deleting VPC: %s", err))
			}
		}

		d.SetId("")
		log.Printf("[INFO] VPC deleted, ID: %s", vpcID)

		return nil
	})

	if err != nil {
		return diag.FromErr(err)
	} else {
		return nil
	}
}
