package go_api_abrha

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/abrhacom/go-api-abrha/metrics"
)

const (
	monitoringBasePath          = "api/public/v1/monitoring"
	alertPolicyBasePath         = monitoringBasePath + "/alerts"
	vmMetricsBasePath           = monitoringBasePath + "/metrics/vm"
	loadBalancerMetricsBasePath = monitoringBasePath + "/metrics/load_balancer"

	VmCPUUtilizationPercent        = "v1/insights/vm/cpu"
	VmMemoryUtilizationPercent     = "v1/insights/vm/memory_utilization_percent"
	VmDiskUtilizationPercent       = "v1/insights/vm/disk_utilization_percent"
	VmPublicOutboundBandwidthRate  = "v1/insights/vm/public_outbound_bandwidth"
	VmPublicInboundBandwidthRate   = "v1/insights/vm/public_inbound_bandwidth"
	VmPrivateOutboundBandwidthRate = "v1/insights/vm/private_outbound_bandwidth"
	VmPrivateInboundBandwidthRate  = "v1/insights/vm/private_inbound_bandwidth"
	VmDiskReadRate                 = "v1/insights/vm/disk_read"
	VmDiskWriteRate                = "v1/insights/vm/disk_write"
	VmOneMinuteLoadAverage         = "v1/insights/vm/load_1"
	VmFiveMinuteLoadAverage        = "v1/insights/vm/load_5"
	VmFifteenMinuteLoadAverage     = "v1/insights/vm/load_15"

	LoadBalancerCPUUtilizationPercent                = "v1/insights/lbaas/avg_cpu_utilization_percent"
	LoadBalancerConnectionUtilizationPercent         = "v1/insights/lbaas/connection_utilization_percent"
	LoadBalancerVmHealth                             = "v1/insights/lbaas/vm_health"
	LoadBalancerTLSUtilizationPercent                = "v1/insights/lbaas/tls_connections_per_second_utilization_percent"
	LoadBalancerIncreaseInHTTPErrorRatePercentage5xx = "v1/insights/lbaas/increase_in_http_error_rate_percentage_5xx"
	LoadBalancerIncreaseInHTTPErrorRatePercentage4xx = "v1/insights/lbaas/increase_in_http_error_rate_percentage_4xx"
	LoadBalancerIncreaseInHTTPErrorRateCount5xx      = "v1/insights/lbaas/increase_in_http_error_rate_count_5xx"
	LoadBalancerIncreaseInHTTPErrorRateCount4xx      = "v1/insights/lbaas/increase_in_http_error_rate_count_4xx"
	LoadBalancerHighHttpResponseTime                 = "v1/insights/lbaas/high_http_request_response_time"
	LoadBalancerHighHttpResponseTime50P              = "v1/insights/lbaas/high_http_request_response_time_50p"
	LoadBalancerHighHttpResponseTime95P              = "v1/insights/lbaas/high_http_request_response_time_95p"
	LoadBalancerHighHttpResponseTime99P              = "v1/insights/lbaas/high_http_request_response_time_99p"

	DbaasFifteenMinuteLoadAverage = "v1/dbaas/alerts/load_15_alerts"
	DbaasMemoryUtilizationPercent = "v1/dbaas/alerts/memory_utilization_alerts"
	DbaasDiskUtilizationPercent   = "v1/dbaas/alerts/disk_utilization_alerts"
	DbaasCPUUtilizationPercent    = "v1/dbaas/alerts/cpu_alerts"
)

// MonitoringService is an interface for interfacing with the
// monitoring endpoints of the Abrha API
// See: https://docs.parspack.com/api/#tag/Monitoring
type MonitoringService interface {
	ListAlertPolicies(context.Context, *ListOptions) ([]AlertPolicy, *Response, error)
	GetAlertPolicy(context.Context, string) (*AlertPolicy, *Response, error)
	CreateAlertPolicy(context.Context, *AlertPolicyCreateRequest) (*AlertPolicy, *Response, error)
	UpdateAlertPolicy(context.Context, string, *AlertPolicyUpdateRequest) (*AlertPolicy, *Response, error)
	DeleteAlertPolicy(context.Context, string) (*Response, error)

	GetVmBandwidth(context.Context, *VmBandwidthMetricsRequest) (*MetricsResponse, *Response, error)
	GetVmAvailableMemory(context.Context, *VmMetricsRequest) (*MetricsResponse, *Response, error)
	GetVmCPU(context.Context, *VmMetricsRequest) (*MetricsResponse, *Response, error)
	GetVmFilesystemFree(context.Context, *VmMetricsRequest) (*MetricsResponse, *Response, error)
	GetVmFilesystemSize(context.Context, *VmMetricsRequest) (*MetricsResponse, *Response, error)
	GetVmLoad1(context.Context, *VmMetricsRequest) (*MetricsResponse, *Response, error)
	GetVmLoad5(context.Context, *VmMetricsRequest) (*MetricsResponse, *Response, error)
	GetVmLoad15(context.Context, *VmMetricsRequest) (*MetricsResponse, *Response, error)
	GetVmCachedMemory(context.Context, *VmMetricsRequest) (*MetricsResponse, *Response, error)
	GetVmFreeMemory(context.Context, *VmMetricsRequest) (*MetricsResponse, *Response, error)
	GetVmTotalMemory(context.Context, *VmMetricsRequest) (*MetricsResponse, *Response, error)

	GetLoadBalancerFrontendHttpRequestsPerSecond(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error)
	GetLoadBalancerFrontendConnectionsCurrent(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error)
	GetLoadBalancerFrontendConnectionsLimit(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error)
	GetLoadBalancerFrontendCpuUtilization(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error)
	GetLoadBalancerFrontendNetworkThroughputHttp(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error)
	GetLoadBalancerFrontendNetworkThroughputUdp(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error)
	GetLoadBalancerFrontendNetworkThroughputTcp(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error)
	GetLoadBalancerFrontendNlbTcpNetworkThroughput(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error)
	GetLoadBalancerFrontendNlbUdpNetworkThroughput(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error)
	GetLoadBalancerFrontendFirewallDroppedBytes(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error)
	GetLoadBalancerFrontendFirewallDroppedPackets(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error)
	GetLoadBalancerFrontendHttpResponses(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error)
	GetLoadBalancerFrontendTlsConnectionsCurrent(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error)
	GetLoadBalancerFrontendTlsConnectionsLimit(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error)
	GetLoadBalancerFrontendTlsConnectionsExceedingRateLimit(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error)
	GetLoadBalancerVmsHttpSessionDurationAvg(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error)
	GetLoadBalancerVmsHttpSessionDuration50P(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error)
	GetLoadBalancerVmsHttpSessionDuration95P(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error)
	GetLoadBalancerVmsHttpResponseTimeAvg(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error)
	GetLoadBalancerVmsHttpResponseTime50P(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error)
	GetLoadBalancerVmsHttpResponseTime95P(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error)
	GetLoadBalancerVmsHttpResponseTime99P(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error)
	GetLoadBalancerVmsQueueSize(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error)
	GetLoadBalancerVmsHttpResponses(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error)
	GetLoadBalancerVmsConnections(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error)
	GetLoadBalancerVmsHealthChecks(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error)
	GetLoadBalancerVmsDowntime(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error)
}

// MonitoringServiceOp handles communication with monitoring related methods of the
// Abrha API.
type MonitoringServiceOp struct {
	client *Client
}

var _ MonitoringService = &MonitoringServiceOp{}

// AlertPolicy represents a Abrha alert policy
type AlertPolicy struct {
	UUID        string          `json:"uuid"`
	Type        string          `json:"type"`
	Description string          `json:"description"`
	Compare     AlertPolicyComp `json:"compare"`
	Value       float32         `json:"value"`
	Window      string          `json:"window"`
	Entities    []string        `json:"entities"`
	Tags        []string        `json:"tags"`
	Alerts      Alerts          `json:"alerts"`
	Enabled     bool            `json:"enabled"`
}

// Alerts represents the alerts section of an alert policy
type Alerts struct {
	Slack []SlackDetails `json:"slack"`
	Email []string       `json:"email"`
}

// SlackDetails represents the details required to send a slack alert
type SlackDetails struct {
	URL     string `json:"url"`
	Channel string `json:"channel"`
}

// AlertPolicyComp represents an alert policy comparison operation
type AlertPolicyComp string

const (
	// GreaterThan is the comparison >
	GreaterThan AlertPolicyComp = "GreaterThan"
	// LessThan is the comparison <
	LessThan AlertPolicyComp = "LessThan"
)

// AlertPolicyCreateRequest holds the info for creating a new alert policy
type AlertPolicyCreateRequest struct {
	Type        string          `json:"type"`
	Description string          `json:"description"`
	Compare     AlertPolicyComp `json:"compare"`
	Value       float32         `json:"value"`
	Window      string          `json:"window"`
	Entities    []string        `json:"entities"`
	Tags        []string        `json:"tags"`
	Alerts      Alerts          `json:"alerts"`
	Enabled     *bool           `json:"enabled"`
}

// AlertPolicyUpdateRequest holds the info for updating an existing alert policy
type AlertPolicyUpdateRequest struct {
	Type        string          `json:"type"`
	Description string          `json:"description"`
	Compare     AlertPolicyComp `json:"compare"`
	Value       float32         `json:"value"`
	Window      string          `json:"window"`
	Entities    []string        `json:"entities"`
	Tags        []string        `json:"tags"`
	Alerts      Alerts          `json:"alerts"`
	Enabled     *bool           `json:"enabled"`
}

type alertPoliciesRoot struct {
	AlertPolicies []AlertPolicy `json:"policies"`
	Links         *Links        `json:"links"`
	Meta          *Meta         `json:"meta"`
}

type alertPolicyRoot struct {
	AlertPolicy *AlertPolicy `json:"policy,omitempty"`
}

// VmMetricsRequest holds the information needed to retrieve Vm various metrics.
type VmMetricsRequest struct {
	HostID string
	Start  time.Time
	End    time.Time
}

// VmBandwidthMetricsRequest holds the information needed to retrieve Vm bandwidth metrics.
type VmBandwidthMetricsRequest struct {
	VmMetricsRequest
	Interface string
	Direction string
}

// LoadBalancerMetricsRequest holds the information needed to retrieve Load Balancer various metrics.
type LoadBalancerMetricsRequest struct {
	LoadBalancerID string
	Start          time.Time
	End            time.Time
}

// MetricsResponse holds a Metrics query response.
type MetricsResponse struct {
	Status string      `json:"status"`
	Data   MetricsData `json:"data"`
}

// MetricsData holds the data portion of a Metrics response.
type MetricsData struct {
	ResultType string                 `json:"resultType"`
	Result     []metrics.SampleStream `json:"result"`
}

// ListAlertPolicies all alert policies
func (s *MonitoringServiceOp) ListAlertPolicies(ctx context.Context, opt *ListOptions) ([]AlertPolicy, *Response, error) {
	path := alertPolicyBasePath
	path, err := addOptions(path, opt)

	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(alertPoliciesRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if l := root.Links; l != nil {
		resp.Links = l
	}
	if m := root.Meta; m != nil {
		resp.Meta = m
	}
	return root.AlertPolicies, resp, err
}

// GetAlertPolicy gets a single alert policy
func (s *MonitoringServiceOp) GetAlertPolicy(ctx context.Context, uuid string) (*AlertPolicy, *Response, error) {
	path := fmt.Sprintf("%s/%s", alertPolicyBasePath, uuid)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(alertPolicyRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.AlertPolicy, resp, err
}

// CreateAlertPolicy creates a new alert policy
func (s *MonitoringServiceOp) CreateAlertPolicy(ctx context.Context, createRequest *AlertPolicyCreateRequest) (*AlertPolicy, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, alertPolicyBasePath, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(alertPolicyRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.AlertPolicy, resp, err
}

// UpdateAlertPolicy updates an existing alert policy
func (s *MonitoringServiceOp) UpdateAlertPolicy(ctx context.Context, uuid string, updateRequest *AlertPolicyUpdateRequest) (*AlertPolicy, *Response, error) {
	if uuid == "" {
		return nil, nil, NewArgError("uuid", "cannot be empty")
	}
	if updateRequest == nil {
		return nil, nil, NewArgError("updateRequest", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%s", alertPolicyBasePath, uuid)
	req, err := s.client.NewRequest(ctx, http.MethodPut, path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(alertPolicyRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.AlertPolicy, resp, err
}

// DeleteAlertPolicy deletes an existing alert policy
func (s *MonitoringServiceOp) DeleteAlertPolicy(ctx context.Context, uuid string) (*Response, error) {
	if uuid == "" {
		return nil, NewArgError("uuid", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s", alertPolicyBasePath, uuid)
	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)

	return resp, err
}

// GetVmBandwidth retrieves Vm bandwidth metrics.
func (s *MonitoringServiceOp) GetVmBandwidth(ctx context.Context, args *VmBandwidthMetricsRequest) (*MetricsResponse, *Response, error) {
	path := vmMetricsBasePath + "/bandwidth"
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	q := req.URL.Query()
	q.Add("host_id", args.HostID)
	q.Add("interface", args.Interface)
	q.Add("direction", args.Direction)
	q.Add("start", fmt.Sprintf("%d", args.Start.Unix()))
	q.Add("end", fmt.Sprintf("%d", args.End.Unix()))
	req.URL.RawQuery = q.Encode()

	root := new(MetricsResponse)
	resp, err := s.client.Do(ctx, req, root)

	return root, resp, err
}

// GetVmCPU retrieves Vm CPU metrics.
func (s *MonitoringServiceOp) GetVmCPU(ctx context.Context, args *VmMetricsRequest) (*MetricsResponse, *Response, error) {
	return s.getVmMetrics(ctx, "/cpu", args)
}

// GetVmFilesystemFree retrieves Vm filesystem free metrics.
func (s *MonitoringServiceOp) GetVmFilesystemFree(ctx context.Context, args *VmMetricsRequest) (*MetricsResponse, *Response, error) {
	return s.getVmMetrics(ctx, "/filesystem_free", args)
}

// GetVmFilesystemSize retrieves Vm filesystem size metrics.
func (s *MonitoringServiceOp) GetVmFilesystemSize(ctx context.Context, args *VmMetricsRequest) (*MetricsResponse, *Response, error) {
	return s.getVmMetrics(ctx, "/filesystem_size", args)
}

// GetVmLoad1 retrieves Vm load 1 metrics.
func (s *MonitoringServiceOp) GetVmLoad1(ctx context.Context, args *VmMetricsRequest) (*MetricsResponse, *Response, error) {
	return s.getVmMetrics(ctx, "/load_1", args)
}

// GetVmLoad5 retrieves Vm load 5 metrics.
func (s *MonitoringServiceOp) GetVmLoad5(ctx context.Context, args *VmMetricsRequest) (*MetricsResponse, *Response, error) {
	return s.getVmMetrics(ctx, "/load_5", args)
}

// GetVmLoad15 retrieves Vm load 15 metrics.
func (s *MonitoringServiceOp) GetVmLoad15(ctx context.Context, args *VmMetricsRequest) (*MetricsResponse, *Response, error) {
	return s.getVmMetrics(ctx, "/load_15", args)
}

// GetVmCachedMemory retrieves Vm cached memory metrics.
func (s *MonitoringServiceOp) GetVmCachedMemory(ctx context.Context, args *VmMetricsRequest) (*MetricsResponse, *Response, error) {
	return s.getVmMetrics(ctx, "/memory_cached", args)
}

// GetVmFreeMemory retrieves Vm free memory metrics.
func (s *MonitoringServiceOp) GetVmFreeMemory(ctx context.Context, args *VmMetricsRequest) (*MetricsResponse, *Response, error) {
	return s.getVmMetrics(ctx, "/memory_free", args)
}

// GetVmTotalMemory retrieves Vm total memory metrics.
func (s *MonitoringServiceOp) GetVmTotalMemory(ctx context.Context, args *VmMetricsRequest) (*MetricsResponse, *Response, error) {
	return s.getVmMetrics(ctx, "/memory_total", args)
}

// GetVmAvailableMemory retrieves Vm available memory metrics.
func (s *MonitoringServiceOp) GetVmAvailableMemory(ctx context.Context, args *VmMetricsRequest) (*MetricsResponse, *Response, error) {
	return s.getVmMetrics(ctx, "/memory_available", args)
}

func (s *MonitoringServiceOp) getVmMetrics(ctx context.Context, path string, args *VmMetricsRequest) (*MetricsResponse, *Response, error) {
	fullPath := vmMetricsBasePath + path
	req, err := s.client.NewRequest(ctx, http.MethodGet, fullPath, nil)
	if err != nil {
		return nil, nil, err
	}

	q := req.URL.Query()
	q.Add("host_id", args.HostID)
	q.Add("start", fmt.Sprintf("%d", args.Start.Unix()))
	q.Add("end", fmt.Sprintf("%d", args.End.Unix()))
	req.URL.RawQuery = q.Encode()

	root := new(MetricsResponse)
	resp, err := s.client.Do(ctx, req, root)

	return root, resp, err
}

// GetLoadBalancerFrontendHttpRequestsPerSecond retrieves frontend HTTP requests per second for a given load balancer.
func (s *MonitoringServiceOp) GetLoadBalancerFrontendHttpRequestsPerSecond(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error) {
	return s.getLoadBalancerMetrics(ctx, "/frontend_http_requests_per_second", args)
}

// GetLoadBalancerFrontendConnectionsCurrent retrieves frontend total current active connections for a given load balancer.
func (s *MonitoringServiceOp) GetLoadBalancerFrontendConnectionsCurrent(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error) {
	return s.getLoadBalancerMetrics(ctx, "/frontend_connections_current", args)
}

// GetLoadBalancerFrontendConnectionsLimit retrieves frontend max connections limit for a given load balancer.
func (s *MonitoringServiceOp) GetLoadBalancerFrontendConnectionsLimit(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error) {
	return s.getLoadBalancerMetrics(ctx, "/frontend_connections_limit", args)
}

// GetLoadBalancerFrontendCpuUtilization retrieves frontend average percentage cpu utilization for a given load balancer.
func (s *MonitoringServiceOp) GetLoadBalancerFrontendCpuUtilization(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error) {
	return s.getLoadBalancerMetrics(ctx, "/frontend_cpu_utilization", args)
}

// GetLoadBalancerFrontendNetworkThroughputHttp retrieves frontend HTTP throughput for a given load balancer.
func (s *MonitoringServiceOp) GetLoadBalancerFrontendNetworkThroughputHttp(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error) {
	return s.getLoadBalancerMetrics(ctx, "/frontend_network_throughput_http", args)
}

// GetLoadBalancerFrontendNetworkThroughputUdp retrieves frontend UDP throughput for a given load balancer.
func (s *MonitoringServiceOp) GetLoadBalancerFrontendNetworkThroughputUdp(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error) {
	return s.getLoadBalancerMetrics(ctx, "/frontend_network_throughput_udp", args)
}

// GetLoadBalancerFrontendNetworkThroughputTcp retrieves frontend TCP throughput for a given load balancer.
func (s *MonitoringServiceOp) GetLoadBalancerFrontendNetworkThroughputTcp(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error) {
	return s.getLoadBalancerMetrics(ctx, "/frontend_network_throughput_tcp", args)
}

// GetLoadBalancerFrontendNlbTcpNetworkThroughput retrieves frontend TCP throughput for a given network load balancer.
func (s *MonitoringServiceOp) GetLoadBalancerFrontendNlbTcpNetworkThroughput(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error) {
	return s.getLoadBalancerMetrics(ctx, "/frontend_nlb_tcp_network_throughput", args)
}

// GetLoadBalancerFrontendNlbUdpNetworkThroughput retrieves frontend UDP throughput for a given network load balancer.
func (s *MonitoringServiceOp) GetLoadBalancerFrontendNlbUdpNetworkThroughput(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error) {
	return s.getLoadBalancerMetrics(ctx, "/frontend_nlb_udp_network_throughput", args)
}

// GetLoadBalancerFrontendFirewallDroppedBytes retrieves firewall dropped bytes for a given load balancer.
func (s *MonitoringServiceOp) GetLoadBalancerFrontendFirewallDroppedBytes(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error) {
	return s.getLoadBalancerMetrics(ctx, "/frontend_firewall_dropped_bytes", args)
}

// GetLoadBalancerFrontendFirewallDroppedPackets retrieves firewall dropped packets for a given load balancer.
func (s *MonitoringServiceOp) GetLoadBalancerFrontendFirewallDroppedPackets(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error) {
	return s.getLoadBalancerMetrics(ctx, "/frontend_firewall_dropped_packets", args)
}

// GetLoadBalancerFrontendHttpResponses retrieves frontend HTTP rate of response code for a given load balancer.
func (s *MonitoringServiceOp) GetLoadBalancerFrontendHttpResponses(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error) {
	return s.getLoadBalancerMetrics(ctx, "/frontend_http_responses", args)
}

// GetLoadBalancerFrontendTlsConnectionsCurrent retrieves frontend current TLS connections rate for a given load balancer.
func (s *MonitoringServiceOp) GetLoadBalancerFrontendTlsConnectionsCurrent(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error) {
	return s.getLoadBalancerMetrics(ctx, "/frontend_tls_connections_current", args)
}

// GetLoadBalancerFrontendTlsConnectionsLimit retrieves frontend max TLS connections limit for a given load balancer.
func (s *MonitoringServiceOp) GetLoadBalancerFrontendTlsConnectionsLimit(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error) {
	return s.getLoadBalancerMetrics(ctx, "/frontend_tls_connections_limit", args)
}

// GetLoadBalancerFrontendTlsConnectionsExceedingRateLimit retrieves frontend closed TLS connections for exceeded rate limit for a given load balancer.
func (s *MonitoringServiceOp) GetLoadBalancerFrontendTlsConnectionsExceedingRateLimit(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error) {
	return s.getLoadBalancerMetrics(ctx, "/frontend_tls_connections_exceeding_rate_limit", args)
}

// GetLoadBalancerVmsHttpSessionDurationAvg retrieves vm average HTTP session duration for a given load balancer.
func (s *MonitoringServiceOp) GetLoadBalancerVmsHttpSessionDurationAvg(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error) {
	return s.getLoadBalancerMetrics(ctx, "/vms_http_session_duration_avg", args)
}

// GetLoadBalancerVmsHttpSessionDuration50P retrieves vm 50th percentile HTTP session duration for a given load balancer.
func (s *MonitoringServiceOp) GetLoadBalancerVmsHttpSessionDuration50P(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error) {
	return s.getLoadBalancerMetrics(ctx, "/vms_http_session_duration_50p", args)
}

// GetLoadBalancerVmsHttpSessionDuration95P retrieves vm 95th percentile HTTP session duration for a given load balancer.
func (s *MonitoringServiceOp) GetLoadBalancerVmsHttpSessionDuration95P(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error) {
	return s.getLoadBalancerMetrics(ctx, "/vms_http_session_duration_95p", args)
}

// GetLoadBalancerVmsHttpResponseTimeAvg retrieves vm average HTTP response time for a given load balancer.
func (s *MonitoringServiceOp) GetLoadBalancerVmsHttpResponseTimeAvg(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error) {
	return s.getLoadBalancerMetrics(ctx, "/vms_http_response_time_avg", args)
}

// GetLoadBalancerVmsHttpResponseTime50P retrieves vm 50th percentile HTTP response time for a given load balancer.
func (s *MonitoringServiceOp) GetLoadBalancerVmsHttpResponseTime50P(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error) {
	return s.getLoadBalancerMetrics(ctx, "/vms_http_response_time_50p", args)
}

// GetLoadBalancerVmsHttpResponseTime95P retrieves vm 95th percentile HTTP response time for a given load balancer.
func (s *MonitoringServiceOp) GetLoadBalancerVmsHttpResponseTime95P(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error) {
	return s.getLoadBalancerMetrics(ctx, "/vms_http_response_time_95p", args)
}

// GetLoadBalancerVmsHttpResponseTime99P retrieves vm 99th percentile HTTP response time for a given load balancer.
func (s *MonitoringServiceOp) GetLoadBalancerVmsHttpResponseTime99P(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error) {
	return s.getLoadBalancerMetrics(ctx, "/vms_http_response_time_99p", args)
}

// GetLoadBalancerVmsQueueSize retrieves vm queue size for a given load balancer.
func (s *MonitoringServiceOp) GetLoadBalancerVmsQueueSize(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error) {
	return s.getLoadBalancerMetrics(ctx, "/vms_queue_size", args)
}

// GetLoadBalancerVmsHttpResponses retrieves vm HTTP rate of response code for a given load balancer.
func (s *MonitoringServiceOp) GetLoadBalancerVmsHttpResponses(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error) {
	return s.getLoadBalancerMetrics(ctx, "/vms_http_responses", args)
}

// GetLoadBalancerVmsConnections retrieves vm active connections for a given load balancer.
func (s *MonitoringServiceOp) GetLoadBalancerVmsConnections(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error) {
	return s.getLoadBalancerMetrics(ctx, "/vms_connections", args)
}

// GetLoadBalancerVmsHealthChecks retrieves vm health check status for a given load balancer.
func (s *MonitoringServiceOp) GetLoadBalancerVmsHealthChecks(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error) {
	return s.getLoadBalancerMetrics(ctx, "/vms_health_checks", args)
}

// GetLoadBalancerVmsDowntime retrieves vm downtime status for a given load balancer.
func (s *MonitoringServiceOp) GetLoadBalancerVmsDowntime(ctx context.Context, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error) {
	return s.getLoadBalancerMetrics(ctx, "/vms_downtime", args)
}

func (s *MonitoringServiceOp) getLoadBalancerMetrics(ctx context.Context, path string, args *LoadBalancerMetricsRequest) (*MetricsResponse, *Response, error) {
	fullPath := loadBalancerMetricsBasePath + path
	req, err := s.client.NewRequest(ctx, http.MethodGet, fullPath, nil)
	if err != nil {
		return nil, nil, err
	}

	q := req.URL.Query()
	q.Add("lb_id", args.LoadBalancerID)
	q.Add("start", fmt.Sprintf("%d", args.Start.Unix()))
	q.Add("end", fmt.Sprintf("%d", args.End.Unix()))
	req.URL.RawQuery = q.Encode()

	root := new(MetricsResponse)
	resp, err := s.client.Do(ctx, req, root)

	return root, resp, err
}
