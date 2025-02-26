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

func ResourceAbrhaDatabaseOpensearchConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbrhaDatabaseOpensearchConfigCreate,
		ReadContext:   resourceAbrhaDatabaseOpensearchConfigRead,
		UpdateContext: resourceAbrhaDatabaseOpensearchConfigUpdate,
		DeleteContext: resourceAbrhaDatabaseOpensearchConfigDelete,
		Importer: &schema.ResourceImporter{
			State: resourceAbrhaDatabaseOpensearchConfigImport,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"ism_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"ism_history_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"ism_history_max_age_hours": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"ism_history_max_docs": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"ism_history_rollover_check_period_hours": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"ism_history_rollover_retention_period_days": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"http_max_content_length_bytes": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"http_max_header_size_bytes": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"http_max_initial_line_length_bytes": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1024),
			},
			"indices_query_bool_max_clause_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(64),
			},
			"search_max_buckets": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"indices_fielddata_cache_size_percentage": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(3),
			},
			"indices_memory_index_buffer_size_percentage": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(3),
			},
			"indices_memory_min_index_buffer_size_mb": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(3),
			},
			"indices_memory_max_index_buffer_size_mb": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(3),
			},
			"indices_queries_cache_size_percentage": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(3),
			},
			"indices_recovery_max_mb_per_sec": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(40),
			},
			"indices_recovery_max_concurrent_file_chunks": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(2),
			},
			"action_auto_create_index_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"action_destructive_requires_name": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"enable_security_audit": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"thread_pool_search_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"thread_pool_search_throttled_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"thread_pool_search_throttled_queue_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(10),
			},
			"thread_pool_search_queue_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(10),
			},
			"thread_pool_get_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"thread_pool_get_queue_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(10),
			},
			"thread_pool_analyze_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"thread_pool_analyze_queue_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(10),
			},
			"thread_pool_write_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"thread_pool_write_queue_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(10),
			},
			"thread_pool_force_merge_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"override_main_response_version": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"script_max_compilations_rate": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cluster_max_shards_per_node": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(100),
			},
			"cluster_routing_allocation_node_concurrent_recoveries": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(2),
			},
			"plugins_alerting_filter_by_backend_roles_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"reindex_remote_whitelist": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
		},
	}
}

func resourceAbrhaDatabaseOpensearchConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()
	clusterID := d.Get("cluster_id").(string)

	if d.HasChangeExcept("cluster_id") {
		if err := updateOpensearchConfig(ctx, d, client); err != nil {
			return diag.Errorf("Error updating Opensearch configuration: %s", err)
		}
	}

	d.SetId(makeDatabaseOpensearchConfigID(clusterID))

	return resourceAbrhaDatabaseOpensearchConfigRead(ctx, d, meta)
}

func resourceAbrhaDatabaseOpensearchConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	if err := updateOpensearchConfig(ctx, d, client); err != nil {
		return diag.Errorf("Error updating Opensearch configuration: %s", err)
	}

	return resourceAbrhaDatabaseOpensearchConfigRead(ctx, d, meta)
}

func updateOpensearchConfig(ctx context.Context, d *schema.ResourceData, client *goApiAbrha.Client) error {
	clusterID := d.Get("cluster_id").(string)

	opts := &goApiAbrha.OpensearchConfig{}

	if d.HasChanges("ism_enabled", "ism_history_enabled", "ism_history_max_age_hours",
		"ism_history_max_docs", "ism_history_rollover_check_period_hours",
		"ism_history_rollover_retention_period_days") {

		if v, ok := d.GetOkExists("ism_enabled"); ok {
			opts.IsmEnabled = goApiAbrha.PtrTo(v.(bool))
		}

		if v, ok := d.GetOkExists("ism_history_enabled"); ok {
			opts.IsmHistoryEnabled = goApiAbrha.PtrTo(v.(bool))
		}

		if v, ok := d.GetOk("ism_history_max_age_hours"); ok {
			opts.IsmHistoryMaxAgeHours = goApiAbrha.PtrTo(v.(int))
		}

		if v, ok := d.GetOk("ism_history_max_docs"); ok {
			opts.IsmHistoryMaxDocs = goApiAbrha.PtrTo(int64(v.(int)))
		}

		if v, ok := d.GetOk("ism_history_rollover_check_period_hours"); ok {
			opts.IsmHistoryRolloverCheckPeriodHours = goApiAbrha.PtrTo(v.(int))
		}

		if v, ok := d.GetOk("ism_history_rollover_retention_period_days"); ok {
			opts.IsmHistoryRolloverRetentionPeriodDays = goApiAbrha.PtrTo(v.(int))
		}
	}

	if d.HasChange("http_max_content_length_bytes") {
		if v, ok := d.GetOk("http_max_content_length_bytes"); ok {
			opts.HttpMaxContentLengthBytes = goApiAbrha.PtrTo(v.(int))
		}
	}

	if d.HasChange("http_max_header_size_bytes") {
		if v, ok := d.GetOk("http_max_header_size_bytes"); ok {
			opts.HttpMaxHeaderSizeBytes = goApiAbrha.PtrTo(v.(int))
		}
	}

	if d.HasChange("http_max_initial_line_length_bytes") {
		if v, ok := d.GetOk("http_max_initial_line_length_bytes"); ok {
			opts.HttpMaxInitialLineLengthBytes = goApiAbrha.PtrTo(v.(int))
		}
	}

	if d.HasChange("indices_query_bool_max_clause_count") {
		if v, ok := d.GetOk("indices_query_bool_max_clause_count"); ok {
			opts.IndicesQueryBoolMaxClauseCount = goApiAbrha.PtrTo(v.(int))
		}
	}

	if d.HasChange("search_max_buckets") {
		if v, ok := d.GetOk("search_max_buckets"); ok {
			opts.SearchMaxBuckets = goApiAbrha.PtrTo(v.(int))
		}
	}

	if d.HasChange("indices_fielddata_cache_size_percentage") {
		if v, ok := d.GetOk("indices_fielddata_cache_size_percentage"); ok {
			opts.IndicesFielddataCacheSizePercentage = goApiAbrha.PtrTo(v.(int))
		}
	}

	if d.HasChange("indices_memory_index_buffer_size_percentage") {
		if v, ok := d.GetOk("indices_memory_index_buffer_size_percentage"); ok {
			opts.IndicesMemoryIndexBufferSizePercentage = goApiAbrha.PtrTo(v.(int))
		}
	}

	if d.HasChange("indices_memory_min_index_buffer_size_mb") {
		if v, ok := d.GetOk("indices_memory_min_index_buffer_size_mb"); ok {
			opts.IndicesMemoryMinIndexBufferSizeMb = goApiAbrha.PtrTo(v.(int))
		}
	}

	if d.HasChange("indices_memory_max_index_buffer_size_mb") {
		if v, ok := d.GetOk("indices_memory_max_index_buffer_size_mb"); ok {
			opts.IndicesMemoryMaxIndexBufferSizeMb = goApiAbrha.PtrTo(v.(int))
		}
	}

	if d.HasChange("indices_queries_cache_size_percentage") {
		if v, ok := d.GetOk("indices_queries_cache_size_percentage"); ok {
			opts.IndicesQueriesCacheSizePercentage = goApiAbrha.PtrTo(v.(int))
		}
	}

	if d.HasChange("indices_recovery_max_mb_per_sec") {
		if v, ok := d.GetOk("indices_recovery_max_mb_per_sec"); ok {
			opts.IndicesRecoveryMaxMbPerSec = goApiAbrha.PtrTo(v.(int))
		}
	}

	if d.HasChange("indices_recovery_max_concurrent_file_chunks") {
		if v, ok := d.GetOk("indices_recovery_max_concurrent_file_chunks"); ok {
			opts.IndicesRecoveryMaxConcurrentFileChunks = goApiAbrha.PtrTo(v.(int))
		}
	}

	if d.HasChange("action_auto_create_index_enabled") {
		if v, ok := d.GetOkExists("action_auto_create_index_enabled"); ok {
			opts.ActionAutoCreateIndexEnabled = goApiAbrha.PtrTo(v.(bool))
		}
	}

	if d.HasChange("action_destructive_requires_name") {
		if v, ok := d.GetOkExists("action_destructive_requires_name"); ok {
			opts.ActionDestructiveRequiresName = goApiAbrha.PtrTo(v.(bool))
		}
	}

	if d.HasChange("enable_security_audit") {
		if v, ok := d.GetOkExists("enable_security_audit"); ok {
			opts.EnableSecurityAudit = goApiAbrha.PtrTo(v.(bool))
		}
	}

	if d.HasChange("thread_pool_search_size") {
		if v, ok := d.GetOk("thread_pool_search_size"); ok {
			opts.ThreadPoolSearchSize = goApiAbrha.PtrTo(v.(int))
		}
	}

	if d.HasChange("thread_pool_search_throttled_size") {
		if v, ok := d.GetOk("thread_pool_search_throttled_size"); ok {
			opts.ThreadPoolSearchThrottledSize = goApiAbrha.PtrTo(v.(int))
		}
	}

	if d.HasChange("thread_pool_search_throttled_queue_size") {
		if v, ok := d.GetOk("thread_pool_search_throttled_queue_size"); ok {
			opts.ThreadPoolSearchThrottledQueueSize = goApiAbrha.PtrTo(v.(int))
		}
	}

	if d.HasChange("thread_pool_search_queue_size") {
		if v, ok := d.GetOk("thread_pool_search_queue_size"); ok {
			opts.ThreadPoolSearchQueueSize = goApiAbrha.PtrTo(v.(int))
		}
	}

	if d.HasChange("thread_pool_get_size") {
		if v, ok := d.GetOk("thread_pool_get_size"); ok {
			opts.ThreadPoolGetSize = goApiAbrha.PtrTo(v.(int))
		}
	}

	if d.HasChange("thread_pool_get_queue_size") {
		if v, ok := d.GetOk("thread_pool_get_queue_size"); ok {
			opts.ThreadPoolGetQueueSize = goApiAbrha.PtrTo(v.(int))
		}
	}

	if d.HasChange("thread_pool_analyze_size") {
		if v, ok := d.GetOk("thread_pool_analyze_size"); ok {
			opts.ThreadPoolAnalyzeSize = goApiAbrha.PtrTo(v.(int))
		}
	}

	if d.HasChange("thread_pool_analyze_queue_size") {
		if v, ok := d.GetOk("thread_pool_analyze_queue_size"); ok {
			opts.ThreadPoolAnalyzeQueueSize = goApiAbrha.PtrTo(v.(int))
		}
	}

	if d.HasChange("thread_pool_write_size") {
		if v, ok := d.GetOk("thread_pool_write_size"); ok {
			opts.ThreadPoolWriteSize = goApiAbrha.PtrTo(v.(int))
		}
	}

	if d.HasChange("thread_pool_write_queue_size") {
		if v, ok := d.GetOk("thread_pool_write_queue_size"); ok {
			opts.ThreadPoolWriteQueueSize = goApiAbrha.PtrTo(v.(int))
		}
	}

	if d.HasChange("thread_pool_force_merge_size") {
		if v, ok := d.GetOk("thread_pool_force_merge_size"); ok {
			opts.ThreadPoolForceMergeSize = goApiAbrha.PtrTo(v.(int))
		}
	}

	if d.HasChange("override_main_response_version") {
		if v, ok := d.GetOkExists("override_main_response_version"); ok {
			opts.OverrideMainResponseVersion = goApiAbrha.PtrTo(v.(bool))
		}
	}

	if d.HasChange("script_max_compilations_rate") {
		if v, ok := d.GetOk("script_max_compilations_rate"); ok {
			opts.ScriptMaxCompilationsRate = goApiAbrha.PtrTo(v.(string))
		}
	}

	if d.HasChange("cluster_max_shards_per_node") {
		if v, ok := d.GetOk("cluster_max_shards_per_node"); ok {
			opts.ClusterMaxShardsPerNode = goApiAbrha.PtrTo(v.(int))
		}
	}

	if d.HasChange("cluster_routing_allocation_node_concurrent_recoveries") {
		if v, ok := d.GetOk("cluster_routing_allocation_node_concurrent_recoveries"); ok {
			opts.ClusterRoutingAllocationNodeConcurrentRecoveries = goApiAbrha.PtrTo(v.(int))
		}
	}

	if d.HasChange("plugins_alerting_filter_by_backend_roles_enabled") {
		if v, ok := d.GetOkExists("plugins_alerting_filter_by_backend_roles_enabled"); ok {
			opts.PluginsAlertingFilterByBackendRolesEnabled = goApiAbrha.PtrTo(v.(bool))
		}
	}

	if d.HasChange("reindex_remote_whitelist") {
		if v, ok := d.GetOk("reindex_remote_whitelist"); ok {
			if exampleSet, ok := v.(*schema.Set); ok {
				var items []string
				for _, item := range exampleSet.List() {
					if str, ok := item.(string); ok {
						items = append(items, str)
					} else {
						return fmt.Errorf("non-string item found in set")
					}
				}
				opts.ReindexRemoteWhitelist = items
			}
		}
	}

	log.Printf("[DEBUG] Opensearch configuration: %s", goApiAbrha.Stringify(opts))

	if _, err := client.Databases.UpdateOpensearchConfig(ctx, clusterID, opts); err != nil {
		return err
	}

	return nil
}

func resourceAbrhaDatabaseOpensearchConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()
	clusterID := d.Get("cluster_id").(string)

	config, resp, err := client.Databases.GetOpensearchConfig(ctx, clusterID)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving Opensearch configuration: %s", err)
	}

	d.Set("ism_enabled", config.IsmEnabled)
	d.Set("ism_history_enabled", config.IsmHistoryEnabled)
	d.Set("ism_history_max_age_hours", config.IsmHistoryMaxAgeHours)
	d.Set("ism_history_max_docs", config.IsmHistoryMaxDocs)
	d.Set("ism_history_rollover_check_period_hours", config.IsmHistoryRolloverCheckPeriodHours)
	d.Set("ism_history_rollover_retention_period_days", config.IsmHistoryRolloverRetentionPeriodDays)
	d.Set("http_max_content_length_bytes", config.HttpMaxContentLengthBytes)
	d.Set("http_max_header_size_bytes", config.HttpMaxHeaderSizeBytes)
	d.Set("http_max_initial_line_length_bytes", config.HttpMaxInitialLineLengthBytes)
	d.Set("indices_query_bool_max_clause_count", config.IndicesQueryBoolMaxClauseCount)
	d.Set("search_max_buckets", config.SearchMaxBuckets)
	d.Set("indices_fielddata_cache_size_percentage", config.IndicesFielddataCacheSizePercentage)
	d.Set("indices_memory_index_buffer_size_percentage", config.IndicesMemoryIndexBufferSizePercentage)
	d.Set("indices_memory_min_index_buffer_size_mb", config.IndicesMemoryMinIndexBufferSizeMb)
	d.Set("indices_memory_max_index_buffer_size_mb", config.IndicesMemoryMaxIndexBufferSizeMb)
	d.Set("indices_queries_cache_size_percentage", config.IndicesQueriesCacheSizePercentage)
	d.Set("indices_recovery_max_mb_per_sec", config.IndicesRecoveryMaxMbPerSec)
	d.Set("indices_recovery_max_concurrent_file_chunks", config.IndicesRecoveryMaxConcurrentFileChunks)
	d.Set("action_auto_create_index_enabled", config.ActionAutoCreateIndexEnabled)
	d.Set("action_destructive_requires_name", config.ActionDestructiveRequiresName)
	d.Set("enable_security_audit", config.EnableSecurityAudit)
	d.Set("thread_pool_search_size", config.ThreadPoolSearchSize)
	d.Set("thread_pool_search_throttled_size", config.ThreadPoolSearchThrottledSize)
	d.Set("thread_pool_search_throttled_queue_size", config.ThreadPoolSearchThrottledQueueSize)
	d.Set("thread_pool_search_queue_size", config.ThreadPoolSearchQueueSize)
	d.Set("thread_pool_get_size", config.ThreadPoolGetSize)
	d.Set("thread_pool_get_queue_size", config.ThreadPoolGetQueueSize)
	d.Set("thread_pool_analyze_size", config.ThreadPoolAnalyzeSize)
	d.Set("thread_pool_analyze_queue_size", config.ThreadPoolAnalyzeQueueSize)
	d.Set("thread_pool_write_size", config.ThreadPoolWriteSize)
	d.Set("thread_pool_write_queue_size", config.ThreadPoolWriteQueueSize)
	d.Set("thread_pool_force_merge_size", config.ThreadPoolForceMergeSize)
	d.Set("override_main_response_version", config.OverrideMainResponseVersion)
	d.Set("script_max_compilations_rate", config.ScriptMaxCompilationsRate)
	d.Set("cluster_max_shards_per_node", config.ClusterMaxShardsPerNode)
	d.Set("cluster_routing_allocation_node_concurrent_recoveries", config.ClusterRoutingAllocationNodeConcurrentRecoveries)
	d.Set("plugins_alerting_filter_by_backend_roles_enabled", config.PluginsAlertingFilterByBackendRolesEnabled)
	d.Set("reindex_remote_whitelist", config.ReindexRemoteWhitelist)

	return nil
}

func resourceAbrhaDatabaseOpensearchConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")

	warn := []diag.Diagnostic{
		{
			Severity: diag.Warning,
			Summary:  "abrha_database_opensearch_config removed from state",
			Detail:   "Database configurations are only removed from state when destroyed. The remote configuration is not unset.",
		},
	}

	return warn
}

func resourceAbrhaDatabaseOpensearchConfigImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	clusterID := d.Id()

	d.SetId(makeDatabaseOpensearchConfigID(clusterID))
	d.Set("cluster_id", clusterID)

	return []*schema.ResourceData{d}, nil
}

func makeDatabaseOpensearchConfigID(clusterID string) string {
	return fmt.Sprintf("%s/opensearch-config", clusterID)
}
