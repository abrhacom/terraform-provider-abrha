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

func ResourceAbrhaDatabaseMySQLConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbrhaDatabaseMySQLConfigCreate,
		ReadContext:   resourceAbrhaDatabaseMySQLConfigRead,
		UpdateContext: resourceAbrhaDatabaseMySQLConfigUpdate,
		DeleteContext: resourceAbrhaDatabaseMySQLConfigDelete,
		Importer: &schema.ResourceImporter{
			State: resourceAbrhaDatabaseMySQLConfigImport,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"connect_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"default_time_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"innodb_log_buffer_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"innodb_online_alter_log_max_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"innodb_lock_wait_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"interactive_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"max_allowed_packet": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"net_read_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"sort_buffer_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"sql_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sql_require_primary_key": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"wait_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"net_write_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"group_concat_max_len": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"information_schema_stats_expiry": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"innodb_ft_min_token_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"innodb_ft_server_stopword_table": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"innodb_print_all_deadlocks": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"innodb_rollback_on_timeout": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"internal_tmp_mem_storage_engine": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice(
					[]string{
						"TempTable",
						"MEMORY",
					},
					false,
				),
			},
			"max_heap_table_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"tmp_table_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"slow_query_log": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"long_query_time": {
				Type:     schema.TypeFloat,
				Optional: true,
				Computed: true,
			},
			"backup_hour": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"backup_minute": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"binlog_retention_period": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAbrhaDatabaseMySQLConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()
	clusterID := d.Get("cluster_id").(string)

	if err := updateMySQLConfig(ctx, d, client); err != nil {
		return diag.Errorf("Error updating MySQL configuration: %s", err)
	}

	d.SetId(makeDatabaseMySQLConfigID(clusterID))

	return resourceAbrhaDatabaseMySQLConfigRead(ctx, d, meta)
}

func resourceAbrhaDatabaseMySQLConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	if err := updateMySQLConfig(ctx, d, client); err != nil {
		return diag.Errorf("Error updating MySQL configuration: %s", err)
	}

	return resourceAbrhaDatabaseMySQLConfigRead(ctx, d, meta)
}

func updateMySQLConfig(ctx context.Context, d *schema.ResourceData, client *goApiAbrha.Client) error {
	clusterID := d.Get("cluster_id").(string)

	opts := &goApiAbrha.MySQLConfig{}

	if v, ok := d.GetOk("connect_timeout"); ok {
		opts.ConnectTimeout = goApiAbrha.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("default_time_zone"); ok {
		opts.DefaultTimeZone = goApiAbrha.PtrTo(v.(string))
	}

	if v, ok := d.GetOk("innodb_log_buffer_size"); ok {
		opts.InnodbLogBufferSize = goApiAbrha.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("innodb_online_alter_log_max_size"); ok {
		opts.InnodbOnlineAlterLogMaxSize = goApiAbrha.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("innodb_lock_wait_timeout"); ok {
		opts.InnodbLockWaitTimeout = goApiAbrha.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("interactive_timeout"); ok {
		opts.InteractiveTimeout = goApiAbrha.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("max_allowed_packet"); ok {
		opts.MaxAllowedPacket = goApiAbrha.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("net_read_timeout"); ok {
		opts.NetReadTimeout = goApiAbrha.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("sort_buffer_size"); ok {
		opts.SortBufferSize = goApiAbrha.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("sql_mode"); ok {
		opts.SQLMode = goApiAbrha.PtrTo(v.(string))
	}

	if v, ok := d.GetOkExists("sql_require_primary_key"); ok {
		opts.SQLRequirePrimaryKey = goApiAbrha.PtrTo(v.(bool))
	}

	if v, ok := d.GetOk("wait_timeout"); ok {
		opts.WaitTimeout = goApiAbrha.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("net_write_timeout"); ok {
		opts.NetWriteTimeout = goApiAbrha.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("group_concat_max_len"); ok {
		opts.GroupConcatMaxLen = goApiAbrha.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("information_schema_stats_expiry"); ok {
		opts.InformationSchemaStatsExpiry = goApiAbrha.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("innodb_ft_min_token_size"); ok {
		opts.InnodbFtMinTokenSize = goApiAbrha.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("innodb_ft_server_stopword_table"); ok {
		opts.InnodbFtServerStopwordTable = goApiAbrha.PtrTo(v.(string))
	}

	if v, ok := d.GetOkExists("innodb_print_all_deadlocks"); ok {
		opts.InnodbPrintAllDeadlocks = goApiAbrha.PtrTo(v.(bool))
	}

	if v, ok := d.GetOkExists("innodb_rollback_on_timeout"); ok {
		opts.InnodbRollbackOnTimeout = goApiAbrha.PtrTo(v.(bool))
	}

	if v, ok := d.GetOk("internal_tmp_mem_storage_engine"); ok {
		opts.InternalTmpMemStorageEngine = goApiAbrha.PtrTo(v.(string))
	}

	if v, ok := d.GetOk("max_heap_table_size"); ok {
		opts.MaxHeapTableSize = goApiAbrha.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("tmp_table_size"); ok {
		opts.TmpTableSize = goApiAbrha.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("slow_query_log"); ok {
		opts.SlowQueryLog = goApiAbrha.PtrTo(v.(bool))
	}

	if v, ok := d.GetOk("long_query_time"); ok {
		opts.LongQueryTime = goApiAbrha.PtrTo(float32(v.(float64)))
	}

	if v, ok := d.GetOk("backup_hour"); ok {
		opts.BackupHour = goApiAbrha.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("backup_minute"); ok {
		opts.BackupMinute = goApiAbrha.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("binlog_retention_period"); ok {
		opts.BinlogRetentionPeriod = goApiAbrha.PtrTo(v.(int))
	}

	log.Printf("[DEBUG] MySQL configuration: %s", goApiAbrha.Stringify(opts))

	if _, err := client.Databases.UpdateMySQLConfig(ctx, clusterID, opts); err != nil {
		return err
	}

	return nil
}

func resourceAbrhaDatabaseMySQLConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()
	clusterID := d.Get("cluster_id").(string)

	config, resp, err := client.Databases.GetMySQLConfig(ctx, clusterID)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving MySQL configuration: %s", err)
	}

	d.Set("connect_timeout", config.ConnectTimeout)
	d.Set("default_time_zone", config.DefaultTimeZone)
	d.Set("innodb_log_buffer_size", config.InnodbLogBufferSize)
	d.Set("innodb_online_alter_log_max_size", config.InnodbOnlineAlterLogMaxSize)
	d.Set("innodb_lock_wait_timeout", config.InnodbLockWaitTimeout)
	d.Set("interactive_timeout", config.InteractiveTimeout)
	d.Set("max_allowed_packet", config.MaxAllowedPacket)
	d.Set("net_read_timeout", config.NetReadTimeout)
	d.Set("sort_buffer_size", config.SortBufferSize)
	d.Set("sql_mode", config.SQLMode)
	d.Set("sql_require_primary_key", config.SQLRequirePrimaryKey)
	d.Set("wait_timeout", config.WaitTimeout)
	d.Set("net_write_timeout", config.NetWriteTimeout)
	d.Set("group_concat_max_len", config.GroupConcatMaxLen)
	d.Set("information_schema_stats_expiry", config.InformationSchemaStatsExpiry)
	d.Set("innodb_ft_min_token_size", config.InnodbFtMinTokenSize)
	d.Set("innodb_ft_server_stopword_table", config.InnodbFtServerStopwordTable)
	d.Set("innodb_print_all_deadlocks", config.InnodbPrintAllDeadlocks)
	d.Set("innodb_rollback_on_timeout", config.InnodbRollbackOnTimeout)
	d.Set("internal_tmp_mem_storage_engine", config.InternalTmpMemStorageEngine)
	d.Set("max_heap_table_size", config.MaxHeapTableSize)
	d.Set("tmp_table_size", config.TmpTableSize)
	d.Set("slow_query_log", config.SlowQueryLog)
	d.Set("long_query_time", config.LongQueryTime)
	d.Set("backup_hour", config.BackupHour)
	d.Set("backup_minute", config.BackupMinute)
	d.Set("binlog_retention_period", config.BinlogRetentionPeriod)

	return nil
}

func resourceAbrhaDatabaseMySQLConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")

	warn := []diag.Diagnostic{
		{
			Severity: diag.Warning,
			Summary:  "abrha_database_mysql_config removed from state",
			Detail:   "Database configurations are only removed from state when destroyed. The remote configuration is not unset.",
		},
	}

	return warn
}

func resourceAbrhaDatabaseMySQLConfigImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	clusterID := d.Id()

	d.SetId(makeDatabaseMySQLConfigID(clusterID))
	d.Set("cluster_id", clusterID)

	return []*schema.ResourceData{d}, nil
}

func makeDatabaseMySQLConfigID(clusterID string) string {
	return fmt.Sprintf("%s/mysql-config", clusterID)
}
