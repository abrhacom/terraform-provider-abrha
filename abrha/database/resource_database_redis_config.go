package database

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

func ResourceAbrhaDatabaseRedisConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbrhaDatabaseRedisConfigCreate,
		ReadContext:   resourceAbrhaDatabaseRedisConfigRead,
		UpdateContext: resourceAbrhaDatabaseRedisConfigUpdate,
		DeleteContext: resourceAbrhaDatabaseRedisConfigDelete,
		Importer: &schema.ResourceImporter{
			State: resourceAbrhaDatabaseRedisConfigImport,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"maxmemory_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"pubsub_client_output_buffer_limit": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"number_of_databases": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"io_threads": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"lfu_log_factor": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"lfu_decay_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"ssl": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"notify_keyspace_events": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"persistence": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice(
					[]string{
						"off",
						"rdb",
					},
					true,
				),
			},

			"acl_channels_default": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice(
					[]string{
						"allchannels",
						"resetchannels",
					},
					true,
				),
			},
		},
	}
}

func resourceAbrhaDatabaseRedisConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()
	clusterID := d.Get("cluster_id").(string)

	err := updateRedisConfig(ctx, d, client)
	if err != nil {
		return diag.Errorf("Error updating Redis configuration: %s", err)
	}

	d.SetId(makeDatabaseRedisConfigID(clusterID))

	return resourceAbrhaDatabaseRedisConfigRead(ctx, d, meta)
}

func resourceAbrhaDatabaseRedisConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()
	err := updateRedisConfig(ctx, d, client)
	if err != nil {
		return diag.Errorf("Error updating Redis configuration: %s", err)
	}

	return resourceAbrhaDatabaseRedisConfigRead(ctx, d, meta)
}

func updateRedisConfig(ctx context.Context, d *schema.ResourceData, client *goApiAbrha.Client) error {
	clusterID := d.Get("cluster_id").(string)

	opts := &goApiAbrha.RedisConfig{}

	if v, ok := d.GetOk("maxmemory_policy"); ok {
		opts.RedisMaxmemoryPolicy = goApiAbrha.PtrTo(v.(string))
	}

	if v, ok := d.GetOk("pubsub_client_output_buffer_limit"); ok {
		opts.RedisPubsubClientOutputBufferLimit = goApiAbrha.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("number_of_databases"); ok {
		opts.RedisNumberOfDatabases = goApiAbrha.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("io_threads"); ok {
		opts.RedisIOThreads = goApiAbrha.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("lfu_log_factor"); ok {
		opts.RedisLFULogFactor = goApiAbrha.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("lfu_decay_time"); ok {
		opts.RedisLFUDecayTime = goApiAbrha.PtrTo(v.(int))
	}

	if v, ok := d.GetOkExists("ssl"); ok {
		opts.RedisSSL = goApiAbrha.PtrTo(v.(bool))
	}

	if v, ok := d.GetOkExists("timeout"); ok {
		opts.RedisTimeout = goApiAbrha.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("notify_keyspace_events"); ok {
		opts.RedisNotifyKeyspaceEvents = goApiAbrha.PtrTo(v.(string))
	}

	if v, ok := d.GetOk("persistence"); ok {
		opts.RedisPersistence = goApiAbrha.PtrTo(v.(string))
	}

	if v, ok := d.GetOk("acl_channels_default"); ok {
		opts.RedisACLChannelsDefault = goApiAbrha.PtrTo(v.(string))
	}

	log.Printf("[DEBUG] Redis configuration: %s", goApiAbrha.Stringify(opts))
	_, err := client.Databases.UpdateRedisConfig(ctx, clusterID, opts)
	if err != nil {
		return err
	}

	return nil
}

func resourceAbrhaDatabaseRedisConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()
	clusterID := d.Get("cluster_id").(string)

	config, resp, err := client.Databases.GetRedisConfig(ctx, clusterID)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving Redis configuration: %s", err)
	}

	d.Set("maxmemory_policy", config.RedisMaxmemoryPolicy)
	d.Set("pubsub_client_output_buffer_limit", config.RedisPubsubClientOutputBufferLimit)
	d.Set("number_of_databases", config.RedisNumberOfDatabases)
	d.Set("io_threads", config.RedisIOThreads)
	d.Set("lfu_log_factor", config.RedisLFULogFactor)
	d.Set("lfu_decay_time", config.RedisLFUDecayTime)
	d.Set("ssl", config.RedisSSL)
	d.Set("timeout", config.RedisTimeout)
	d.Set("notify_keyspace_events", config.RedisNotifyKeyspaceEvents)
	d.Set("persistence", config.RedisPersistence)
	d.Set("acl_channels_default", config.RedisACLChannelsDefault)

	return nil
}

func resourceAbrhaDatabaseRedisConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	warn := []diag.Diagnostic{
		{
			Severity: diag.Warning,
			Summary:  "abrha_database_redis_config removed from state",
			Detail:   "Database configurations are only removed from state when destroyed. The remote configuration is not unset.",
		},
	}
	return warn
}

func resourceAbrhaDatabaseRedisConfigImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	clusterID := d.Id()
	d.SetId(makeDatabaseRedisConfigID(clusterID))
	d.Set("cluster_id", clusterID)

	return []*schema.ResourceData{d}, nil
}

func makeDatabaseRedisConfigID(clusterID string) string {
	return fmt.Sprintf("%s/redis-config", clusterID)
}
