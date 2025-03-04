package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceAbrhaDatabaseDB() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbrhaDatabaseDBCreate,
		ReadContext:   resourceAbrhaDatabaseDBRead,
		DeleteContext: resourceAbrhaDatabaseDBDelete,
		Importer: &schema.ResourceImporter{
			State: resourceAbrhaDatabaseDBImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
		},
	}
}

func resourceAbrhaDatabaseDBCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()
	clusterID := d.Get("cluster_id").(string)

	opts := &goApiAbrha.DatabaseCreateDBRequest{
		Name: d.Get("name").(string),
	}

	log.Printf("[DEBUG] Database DB create configuration: %#v", opts)
	db, _, err := client.Databases.CreateDB(context.Background(), clusterID, opts)
	if err != nil {
		return diag.Errorf("Error creating Database DB: %s", err)
	}

	d.SetId(makeDatabaseDBID(clusterID, db.Name))
	log.Printf("[INFO] Database DB Name: %s", db.Name)

	return resourceAbrhaDatabaseDBRead(ctx, d, meta)
}

func resourceAbrhaDatabaseDBRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()
	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)

	// Check if the database DB still exists
	_, resp, err := client.Databases.GetDB(context.Background(), clusterID, name)
	if err != nil {
		// If the database DB is somehow already destroyed, mark as
		// successfully gone
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving Database DB: %s", err)
	}

	return nil
}

func resourceAbrhaDatabaseDBDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()
	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)

	log.Printf("[INFO] Deleting Database DB: %s", d.Id())
	_, err := client.Databases.DeleteDB(context.Background(), clusterID, name)
	if err != nil {
		return diag.Errorf("Error deleting Database DB: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceAbrhaDatabaseDBImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if strings.Contains(d.Id(), ",") {
		s := strings.Split(d.Id(), ",")
		d.SetId(makeDatabaseDBID(s[0], s[1]))
		d.Set("cluster_id", s[0])
		d.Set("name", s[1])
	} else {
		return nil, errors.New("must use the ID of the source database cluster and the name of the database joined with a comma (e.g. `id,name`)")
	}

	return []*schema.ResourceData{d}, nil
}

func makeDatabaseDBID(clusterID string, name string) string {
	return fmt.Sprintf("%s/database/%s", clusterID, name)
}
