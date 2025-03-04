package snapshot

import (
	"context"
	"log"
	"strings"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceAbrhaVmSnapshot() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAbrhaVmSnapshotRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"name_regex": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validation.StringIsValidRegExp,
				ConflictsWith: []string{"name"},
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				StateFunc: func(val interface{}) string {
					// DO API V2 region slug is always lowercase
					return strings.ToLower(val.(string))
				},
				ValidateFunc: validation.NoZeroValues,
			},
			"most_recent": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			// Computed values.
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"min_disk_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"regions": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vm_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
		},
	}
}

// dataSourceDoSnapshotRead performs the Snapshot lookup.
func dataSourceAbrhaVmSnapshotRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	name, hasName := d.GetOk("name")
	nameRegex, hasNameRegex := d.GetOk("name_regex")
	region, hasRegion := d.GetOk("region")

	if !hasName && !hasNameRegex {
		return diag.Errorf("One of `name` or `name_regex` must be assigned")
	}

	opts := &goApiAbrha.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	var snapshotList []goApiAbrha.Snapshot

	for {
		snapshots, resp, err := client.Snapshots.ListVm(context.Background(), opts)

		if err != nil {
			return diag.Errorf("Error retrieving Vm snapshots: %s", err)
		}

		snapshotList = append(snapshotList, snapshots...)

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return diag.Errorf("Error retrieving Vm snapshots: %s", err)
		}

		opts.Page = page + 1
	}

	// Go through all the possible filters
	if hasName {
		snapshotList = filterSnapshotsByName(snapshotList, name.(string))
	} else {
		snapshotList = filterSnapshotsByNameRegex(snapshotList, nameRegex.(string))
	}
	if hasRegion {
		snapshotList = filterSnapshotsByRegion(snapshotList, region.(string))
	}

	// Get the queried snapshot or fail if it can't be determined
	var snapshot *goApiAbrha.Snapshot
	if len(snapshotList) == 0 {
		return diag.Errorf("No VM snapshot found with name %s", name)
	}
	if len(snapshotList) > 1 {
		recent := d.Get("most_recent").(bool)
		if recent {
			snapshot = findMostRecentSnapshot(snapshotList)
		} else {
			return diag.Errorf("too many Vm snapshots found with name %s (found %d, expected 1)", name, len(snapshotList))
		}
	} else {
		snapshot = &snapshotList[0]
	}

	log.Printf("[DEBUG] do_snapshot - Single Vm Snapshot found: %s", snapshot.ID)

	d.SetId(snapshot.ID)
	d.Set("name", snapshot.Name)
	d.Set("created_at", snapshot.Created)
	d.Set("min_disk_size", snapshot.MinDiskSize)
	d.Set("regions", snapshot.Regions)
	d.Set("vm_id", snapshot.ResourceID)
	d.Set("size", snapshot.SizeGigaBytes)

	return nil
}
