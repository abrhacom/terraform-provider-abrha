package database

import (
	"context"
	"fmt"
	"log"
	"strings"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceAbrhaDatabaseMongoDBConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbrhaDatabaseMongoDBConfigCreate,
		ReadContext:   resourceAbrhaDatabaseMongoDBConfigRead,
		UpdateContext: resourceAbrhaDatabaseMongoDBConfigUpdate,
		DeleteContext: resourceAbrhaDatabaseMongoDBConfigDelete,
		Importer: &schema.ResourceImporter{
			State: resourceAbrhaDatabaseMongoDBConfigImport,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"default_read_concern": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice(
					[]string{
						"local",
						"available",
						"majority",
					},
					true,
				),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new)
				},
			},
			"default_write_concern": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"transaction_lifetime_limit_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"slow_op_threshold_ms": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"verbosity": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAbrhaDatabaseMongoDBConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()
	clusterID := d.Get("cluster_id").(string)

	if err := updateMongoDBConfig(ctx, d, client); err != nil {
		return diag.Errorf("Error updating MongoDB configuration: %s", err)
	}

	d.SetId(makeDatabaseMongoDBConfigID(clusterID))

	return resourceAbrhaDatabaseMongoDBConfigRead(ctx, d, meta)
}

func resourceAbrhaDatabaseMongoDBConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	if err := updateMongoDBConfig(ctx, d, client); err != nil {
		return diag.Errorf("Error updating MongoDB configuration: %s", err)
	}

	return resourceAbrhaDatabaseMongoDBConfigRead(ctx, d, meta)
}

func updateMongoDBConfig(ctx context.Context, d *schema.ResourceData, client *goApiAbrha.Client) error {
	clusterID := d.Get("cluster_id").(string)

	opts := &goApiAbrha.MongoDBConfig{}

	if v, ok := d.GetOk("default_read_concern"); ok {
		opts.DefaultReadConcern = goApiAbrha.PtrTo(v.(string))
	}

	if v, ok := d.GetOk("default_write_concern"); ok {
		opts.DefaultWriteConcern = goApiAbrha.PtrTo(v.(string))
	}

	if v, ok := d.GetOk("transaction_lifetime_limit_seconds"); ok {
		opts.TransactionLifetimeLimitSeconds = goApiAbrha.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("slow_op_threshold_ms"); ok {
		opts.SlowOpThresholdMs = goApiAbrha.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("verbosity"); ok {
		opts.Verbosity = goApiAbrha.PtrTo(v.(int))
	}

	log.Printf("[DEBUG] MongoDB configuration: %s", goApiAbrha.Stringify(opts))

	if _, err := client.Databases.UpdateMongoDBConfig(ctx, clusterID, opts); err != nil {
		return err
	}

	return nil
}

func resourceAbrhaDatabaseMongoDBConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()
	clusterID := d.Get("cluster_id").(string)

	config, resp, err := client.Databases.GetMongoDBConfig(ctx, clusterID)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving MongoDB configuration: %s", err)
	}

	d.Set("default_read_concern", config.DefaultReadConcern)
	d.Set("default_write_concern", config.DefaultWriteConcern)
	d.Set("transaction_lifetime_limit_seconds", config.TransactionLifetimeLimitSeconds)
	d.Set("slow_op_threshold_ms", config.SlowOpThresholdMs)
	d.Set("verbosity", config.Verbosity)

	return nil
}

func resourceAbrhaDatabaseMongoDBConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")

	warn := []diag.Diagnostic{
		{
			Severity: diag.Warning,
			Summary:  "abrha_database_mongodb_config removed from state",
			Detail:   "Database configurations are only removed from state when destroyed. The remote configuration is not unset.",
		},
	}

	return warn
}

func resourceAbrhaDatabaseMongoDBConfigImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	clusterID := d.Id()

	d.SetId(makeDatabaseMongoDBConfigID(clusterID))
	d.Set("cluster_id", clusterID)

	return []*schema.ResourceData{d}, nil
}

func makeDatabaseMongoDBConfigID(clusterID string) string {
	return fmt.Sprintf("%s/mongodb-config", clusterID)
}
