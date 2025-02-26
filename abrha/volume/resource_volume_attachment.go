package volume

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/abrhacom/terraform-provider-abrha/abrha/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceAbrhaVolumeAttachment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbrhaVolumeAttachmentCreate,
		ReadContext:   resourceAbrhaVolumeAttachmentRead,
		DeleteContext: resourceAbrhaVolumeAttachmentDelete,

		Schema: map[string]*schema.Schema{
			"vm_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"volume_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
		},
	}
}

func resourceAbrhaVolumeAttachmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	vmId := d.Get("vm_id").(string)
	volumeId := d.Get("volume_id").(string)

	volume, _, err := client.Storage.GetVolume(context.Background(), volumeId)
	if err != nil {
		return diag.Errorf("Error retrieving volume: %s", err)
	}

	if len(volume.VmIDs) == 0 || volume.VmIDs[0] != vmId {

		// Only one volume can be attached at one time to a single vm.
		err := retry.RetryContext(ctx, 5*time.Minute, func() *retry.RetryError {

			log.Printf("[DEBUG] Attaching Volume (%s) to Vm (%s)", volumeId, vmId)
			action, _, err := client.StorageActions.Attach(context.Background(), volumeId, vmId)
			if err != nil {
				if util.IsAbrhaError(err, 422, "Vm already has a pending event.") {
					log.Printf("[DEBUG] Received %s, retrying attaching volume to vm", err)
					return retry.RetryableError(err)
				}

				return retry.NonRetryableError(
					fmt.Errorf("[WARN] Error attaching volume (%s) to Vm (%s): %s", volumeId, vmId, err))
			}

			log.Printf("[DEBUG] Volume attach action id: %d", action.ID)
			if err = util.WaitForAction(client, action); err != nil {
				return retry.NonRetryableError(
					fmt.Errorf("[DEBUG] Error waiting for attach volume (%s) to Vm (%s) to finish: %s", volumeId, vmId, err))
			}

			return nil
		})

		if err != nil {
			return diag.Errorf("Error attaching volume to vm after retry timeout: %s", err)
		}
	}

	d.SetId(id.PrefixedUniqueId(fmt.Sprintf("%s-%s-", vmId, volumeId)))

	return nil
}

func resourceAbrhaVolumeAttachmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	vmId := d.Get("vm_id")
	volumeId := d.Get("volume_id").(string)

	volume, resp, err := client.Storage.GetVolume(context.Background(), volumeId)
	if err != nil {
		// If the volume is already destroyed, mark as
		// successfully removed
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving volume: %s", err)
	}

	if len(volume.VmIDs) == 0 || volume.VmIDs[0] != vmId {
		log.Printf("[DEBUG] Volume Attachment (%s) not found, removing from state", d.Id())
		d.SetId("")
	}

	return nil
}

func resourceAbrhaVolumeAttachmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	vmId := d.Get("vm_id").(string)
	volumeId := d.Get("volume_id").(string)

	// Only one volume can be detached at one time to a single vm.
	err := retry.RetryContext(ctx, 5*time.Minute, func() *retry.RetryError {

		log.Printf("[DEBUG] Detaching Volume (%s) from Vm (%s)", volumeId, vmId)
		action, _, err := client.StorageActions.DetachByVmID(context.Background(), volumeId, vmId)
		if err != nil {
			if util.IsAbrhaError(err, 422, "Vm already has a pending event.") {
				log.Printf("[DEBUG] Received %s, retrying detaching volume from vm", err)
				return retry.RetryableError(err)
			}

			return retry.NonRetryableError(
				fmt.Errorf("[WARN] Error detaching volume (%s) from Vm (%s): %s", volumeId, vmId, err))
		}

		log.Printf("[DEBUG] Volume detach action id: %d", action.ID)
		if err = util.WaitForAction(client, action); err != nil {
			return retry.NonRetryableError(
				fmt.Errorf("error waiting for detach volume (%s) from Vm (%s) to finish: %s", volumeId, vmId, err))
		}

		return nil
	})

	if err != nil {
		return diag.Errorf("Error detaching volume from vm after retry timeout: %s", err)
	}

	return nil
}
