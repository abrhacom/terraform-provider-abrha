package abrha

import (
	"context"

	"github.com/abrhacom/terraform-provider-abrha/abrha/account"
	"github.com/abrhacom/terraform-provider-abrha/abrha/app"
	"github.com/abrhacom/terraform-provider-abrha/abrha/cdn"
	"github.com/abrhacom/terraform-provider-abrha/abrha/certificate"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/abrhacom/terraform-provider-abrha/abrha/database"
	"github.com/abrhacom/terraform-provider-abrha/abrha/domain"
	"github.com/abrhacom/terraform-provider-abrha/abrha/firewall"
	"github.com/abrhacom/terraform-provider-abrha/abrha/image"
	"github.com/abrhacom/terraform-provider-abrha/abrha/kubernetes"
	"github.com/abrhacom/terraform-provider-abrha/abrha/loadbalancer"
	"github.com/abrhacom/terraform-provider-abrha/abrha/monitoring"
	"github.com/abrhacom/terraform-provider-abrha/abrha/project"
	"github.com/abrhacom/terraform-provider-abrha/abrha/region"
	"github.com/abrhacom/terraform-provider-abrha/abrha/registry"
	"github.com/abrhacom/terraform-provider-abrha/abrha/reservedip"
	"github.com/abrhacom/terraform-provider-abrha/abrha/reservedipv6"
	"github.com/abrhacom/terraform-provider-abrha/abrha/size"
	"github.com/abrhacom/terraform-provider-abrha/abrha/snapshot"
	"github.com/abrhacom/terraform-provider-abrha/abrha/spaces"
	"github.com/abrhacom/terraform-provider-abrha/abrha/sshkey"
	"github.com/abrhacom/terraform-provider-abrha/abrha/tag"
	"github.com/abrhacom/terraform-provider-abrha/abrha/uptime"
	"github.com/abrhacom/terraform-provider-abrha/abrha/vm"
	"github.com/abrhacom/terraform-provider-abrha/abrha/vmautoscale"
	"github.com/abrhacom/terraform-provider-abrha/abrha/volume"
	"github.com/abrhacom/terraform-provider-abrha/abrha/vpc"
	"github.com/abrhacom/terraform-provider-abrha/abrha/vpcpeering"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns a schema.Provider for Abrha.
func Provider() *schema.Provider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"ABRHA_TOKEN",
					"ABRHA_ACCESS_TOKEN",
				}, nil),
				Description: "The token key for API operations.",
			},
			"api_endpoint": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ABRHA_API_URL", "https://my.abrha.net/cserver/api"),
				Description: "The URL to use for the Abrha API.",
			},
			"spaces_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SPACES_ENDPOINT_URL", "https://{{.Region}}.my.abrha.com"),
				Description: "The URL to use for the Abrha Spaces API.",
			},
			"spaces_access_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SPACES_ACCESS_KEY_ID", nil),
				Description: "The access key ID for Spaces API operations.",
			},
			"spaces_secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SPACES_SECRET_ACCESS_KEY", nil),
				Description: "The secret access key for Spaces API operations.",
			},
			"requests_per_second": {
				Type:        schema.TypeFloat,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ABRHA_REQUESTS_PER_SECOND", 0.0),
				Description: "The rate of requests per second to limit the HTTP client.",
			},
			"http_retry_max": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ABRHA_HTTP_RETRY_MAX", 4),
				Description: "The maximum number of retries on a failed API request.",
			},
			"http_retry_wait_min": {
				Type:        schema.TypeFloat,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ABRHA_HTTP_RETRY_WAIT_MIN", 1.0),
				Description: "The minimum wait time (in seconds) between failed API requests.",
			},
			"http_retry_wait_max": {
				Type:        schema.TypeFloat,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ABRHA_HTTP_RETRY_WAIT_MAX", 30.0),
				Description: "The maximum wait time (in seconds) between failed API requests.",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"abrha_account":                  account.DataSourceAbrhaAccount(),
			"abrha_app":                      app.DataSourceAbrhaApp(),
			"abrha_certificate":              certificate.DataSourceAbrhaCertificate(),
			"abrha_container_registry":       registry.DataSourceAbrhaContainerRegistry(),
			"abrha_database_cluster":         database.DataSourceAbrhaDatabaseCluster(),
			"abrha_database_connection_pool": database.DataSourceAbrhaDatabaseConnectionPool(),
			"abrha_database_ca":              database.DataSourceAbrhaDatabaseCA(),
			"abrha_database_replica":         database.DataSourceAbrhaDatabaseReplica(),
			"abrha_database_user":            database.DataSourceAbrhaDatabaseUser(),
			"abrha_domain":                   domain.DataSourceAbrhaDomain(),
			"abrha_domains":                  domain.DataSourceAbrhaDomains(),
			"abrha_vm":                       vm.DataSourceAbrhaVm(),
			"abrha_vm_autoscale":             vmautoscale.DataSourceAbrhaVmAutoscale(),
			"abrha_vms":                      vm.DataSourceAbrhaVms(),
			"abrha_vm_snapshot":              snapshot.DataSourceAbrhaVmSnapshot(),
			"abrha_firewall":                 firewall.DataSourceAbrhaFirewall(),
			"abrha_floating_ip":              reservedip.DataSourceAbrhaFloatingIP(),
			"abrha_image":                    image.DataSourceAbrhaImage(),
			"abrha_images":                   image.DataSourceAbrhaImages(),
			"abrha_kubernetes_cluster":       kubernetes.DataSourceAbrhaKubernetesCluster(),
			"abrha_kubernetes_versions":      kubernetes.DataSourceAbrhaKubernetesVersions(),
			"abrha_loadbalancer":             loadbalancer.DataSourceAbrhaLoadbalancer(),
			"abrha_project":                  project.DataSourceAbrhaProject(),
			"abrha_projects":                 project.DataSourceAbrhaProjects(),
			"abrha_record":                   domain.DataSourceAbrhaRecord(),
			"abrha_records":                  domain.DataSourceAbrhaRecords(),
			"abrha_region":                   region.DataSourceAbrhaRegion(),
			"abrha_regions":                  region.DataSourceAbrhaRegions(),
			"abrha_reserved_ip":              reservedip.DataSourceAbrhaReservedIP(),
			"abrha_reserved_ipv6":            reservedipv6.DataSourceAbrhaReservedIPV6(),
			"abrha_sizes":                    size.DataSourceAbrhaSizes(),
			"abrha_spaces_bucket":            spaces.DataSourceAbrhaSpacesBucket(),
			"abrha_spaces_buckets":           spaces.DataSourceAbrhaSpacesBuckets(),
			"abrha_spaces_bucket_object":     spaces.DataSourceAbrhaSpacesBucketObject(),
			"abrha_spaces_bucket_objects":    spaces.DataSourceAbrhaSpacesBucketObjects(),
			"abrha_ssh_key":                  sshkey.DataSourceAbrhaSSHKey(),
			"abrha_ssh_keys":                 sshkey.DataSourceAbrhaSSHKeys(),
			"abrha_tag":                      tag.DataSourceAbrhaTag(),
			"abrha_tags":                     tag.DataSourceAbrhaTags(),
			"abrha_volume_snapshot":          snapshot.DataSourceAbrhaVolumeSnapshot(),
			"abrha_volume":                   volume.DataSourceAbrhaVolume(),
			"abrha_vpc":                      vpc.DataSourceAbrhaVPC(),
			"abrha_vpc_peering":              vpcpeering.DataSourceAbrhaVPCPeering(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"abrha_app":                app.ResourceAbrhaApp(),
			"abrha_certificate":        certificate.ResourceAbrhaCertificate(),
			"abrha_container_registry": registry.ResourceAbrhaContainerRegistry(),
			"abrha_container_registry_docker_credentials": registry.ResourceAbrhaContainerRegistryDockerCredentials(),
			"abrha_cdn":                              cdn.ResourceAbrhaCDN(),
			"abrha_database_cluster":                 database.ResourceAbrhaDatabaseCluster(),
			"abrha_database_connection_pool":         database.ResourceAbrhaDatabaseConnectionPool(),
			"abrha_database_db":                      database.ResourceAbrhaDatabaseDB(),
			"abrha_database_firewall":                database.ResourceAbrhaDatabaseFirewall(),
			"abrha_database_replica":                 database.ResourceAbrhaDatabaseReplica(),
			"abrha_database_user":                    database.ResourceAbrhaDatabaseUser(),
			"abrha_database_redis_config":            database.ResourceAbrhaDatabaseRedisConfig(),
			"abrha_database_postgresql_config":       database.ResourceAbrhaDatabasePostgreSQLConfig(),
			"abrha_database_mysql_config":            database.ResourceAbrhaDatabaseMySQLConfig(),
			"abrha_database_mongodb_config":          database.ResourceAbrhaDatabaseMongoDBConfig(),
			"abrha_database_kafka_config":            database.ResourceAbrhaDatabaseKafkaConfig(),
			"abrha_database_opensearch_config":       database.ResourceAbrhaDatabaseOpensearchConfig(),
			"abrha_database_kafka_topic":             database.ResourceAbrhaDatabaseKafkaTopic(),
			"abrha_domain":                           domain.ResourceAbrhaDomain(),
			"abrha_vm":                               vm.ResourceAbrhaVm(),
			"abrha_vm_autoscale":                     vmautoscale.ResourceAbrhaVmAutoscale(),
			"abrha_vm_snapshot":                      snapshot.ResourceAbrhaVmSnapshot(),
			"abrha_firewall":                         firewall.ResourceAbrhaFirewall(),
			"abrha_floating_ip":                      reservedip.ResourceAbrhaFloatingIP(),
			"abrha_floating_ip_assignment":           reservedip.ResourceAbrhaFloatingIPAssignment(),
			"abrha_kubernetes_cluster":               kubernetes.ResourceAbrhaKubernetesCluster(),
			"abrha_kubernetes_node_pool":             kubernetes.ResourceAbrhaKubernetesNodePool(),
			"abrha_loadbalancer":                     loadbalancer.ResourceAbrhaLoadbalancer(),
			"abrha_monitor_alert":                    monitoring.ResourceAbrhaMonitorAlert(),
			"abrha_project":                          project.ResourceAbrhaProject(),
			"abrha_project_resources":                project.ResourceAbrhaProjectResources(),
			"abrha_record":                           domain.ResourceAbrhaRecord(),
			"abrha_reserved_ip":                      reservedip.ResourceAbrhaReservedIP(),
			"abrha_reserved_ip_assignment":           reservedip.ResourceAbrhaReservedIPAssignment(),
			"abrha_reserved_ipv6":                    reservedipv6.ResourceAbrhaReservedIPV6(),
			"abrha_reserved_ipv6_assignment":         reservedipv6.ResourceAbrhaReservedIPV6Assignment(),
			"abrha_spaces_bucket":                    spaces.ResourceAbrhaBucket(),
			"abrha_spaces_bucket_cors_configuration": spaces.ResourceAbrhaBucketCorsConfiguration(),
			"abrha_spaces_bucket_object":             spaces.ResourceAbrhaSpacesBucketObject(),
			"abrha_spaces_bucket_policy":             spaces.ResourceAbrhaSpacesBucketPolicy(),
			"abrha_ssh_key":                          sshkey.ResourceAbrhaSSHKey(),
			"abrha_tag":                              tag.ResourceAbrhaTag(),
			"abrha_uptime_check":                     uptime.ResourceAbrhaUptimeCheck(),
			"abrha_uptime_alert":                     uptime.ResourceAbrhaUptimeAlert(),
			"abrha_volume":                           volume.ResourceAbrhaVolume(),
			"abrha_volume_attachment":                volume.ResourceAbrhaVolumeAttachment(),
			"abrha_volume_snapshot":                  snapshot.ResourceAbrhaVolumeSnapshot(),
			"abrha_vpc":                              vpc.ResourceAbrhaVPC(),
			"abrha_vpc_peering":                      vpcpeering.ResourceAbrhaVPCPeering(),
			"abrha_custom_image":                     image.ResourceAbrhaCustomImage(),
		},
	}

	p.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

		var diags diag.Diagnostics

		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Billing Warning",
			Detail:   "Please note that if the virtual machine is deleted during its initial billing period (the pre-paid duration), no refunds will be issued for the cost already charged.",
		})

		terraformVersion := p.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		client, err := providerConfigure(d, terraformVersion)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		return client, diags
	}

	return p
}

func providerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, error) {
	conf := config.Config{
		Token:             d.Get("token").(string),
		APIEndpoint:       d.Get("api_endpoint").(string),
		AccessID:          d.Get("spaces_access_id").(string),
		SecretKey:         d.Get("spaces_secret_key").(string),
		RequestsPerSecond: d.Get("requests_per_second").(float64),
		HTTPRetryMax:      d.Get("http_retry_max").(int),
		HTTPRetryWaitMin:  d.Get("http_retry_wait_min").(float64),
		HTTPRetryWaitMax:  d.Get("http_retry_wait_max").(float64),
		TerraformVersion:  terraformVersion,
	}

	if endpoint, ok := d.GetOk("spaces_endpoint"); ok {
		conf.SpacesAPIEndpoint = endpoint.(string)
	}

	return conf.Client()
}
