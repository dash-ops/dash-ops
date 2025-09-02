package kubernetes

import "time"

// ClusterResponse represents cluster information response
type ClusterResponse struct {
	Name    string `json:"name"`
	Context string `json:"context"`
	Server  string `json:"server,omitempty"`
	Version string `json:"version,omitempty"`
	Status  string `json:"status"`
}

// ClusterInfoResponse represents comprehensive cluster information response
type ClusterInfoResponse struct {
	Cluster     ClusterResponse        `json:"cluster"`
	Nodes       []NodeResponse         `json:"nodes"`
	Namespaces  []NamespaceResponse    `json:"namespaces"`
	Summary     ClusterSummaryResponse `json:"summary"`
	LastUpdated time.Time              `json:"last_updated"`
}

// ClusterSummaryResponse represents cluster summary response
type ClusterSummaryResponse struct {
	TotalNodes       int `json:"total_nodes"`
	ReadyNodes       int `json:"ready_nodes"`
	TotalNamespaces  int `json:"total_namespaces"`
	TotalDeployments int `json:"total_deployments"`
	TotalPods        int `json:"total_pods"`
	RunningPods      int `json:"running_pods"`
}

// NodeResponse represents node information response
type NodeResponse struct {
	Name       string                  `json:"name"`
	Status     string                  `json:"status"`
	Roles      []string                `json:"roles"`
	Age        string                  `json:"age"`
	Version    string                  `json:"version"`
	InternalIP string                  `json:"internal_ip"`
	ExternalIP string                  `json:"external_ip,omitempty"`
	Conditions []NodeConditionResponse `json:"conditions"`
	Resources  NodeResourcesResponse   `json:"resources"`
	CreatedAt  time.Time               `json:"created_at"`
}

// NodeConditionResponse represents node condition response
type NodeConditionResponse struct {
	Type               string    `json:"type"`
	Status             string    `json:"status"`
	Reason             string    `json:"reason,omitempty"`
	Message            string    `json:"message,omitempty"`
	LastTransitionTime time.Time `json:"last_transition_time"`
}

// NodeResourcesResponse represents node resources response
type NodeResourcesResponse struct {
	Capacity    ResourceListResponse `json:"capacity"`
	Allocatable ResourceListResponse `json:"allocatable"`
	Used        ResourceListResponse `json:"used,omitempty"`
}

// ResourceListResponse represents resource list response
type ResourceListResponse struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
	Pods   string `json:"pods,omitempty"`
}

// NamespaceResponse represents namespace information response
type NamespaceResponse struct {
	Name      string            `json:"name"`
	Status    string            `json:"status"`
	Labels    map[string]string `json:"labels,omitempty"`
	Age       string            `json:"age"`
	CreatedAt time.Time         `json:"created_at"`
}

// DeploymentResponse represents deployment information response
type DeploymentResponse struct {
	Name                string                        `json:"name"`
	Namespace           string                        `json:"namespace"`
	PodInfo             PodInfoResponse               `json:"pod_info"`
	Replicas            DeploymentReplicasResponse    `json:"replicas"`
	Age                 string                        `json:"age"`
	CreatedAt           time.Time                     `json:"created_at"`
	Conditions          []DeploymentConditionResponse `json:"conditions"`
	ServiceContext      *ServiceContextResponse       `json:"service_context,omitempty"`
	HealthStatus        string                        `json:"health_status,omitempty"`
	AvailabilityPercent float64                       `json:"availability_percent"`
}

// PodInfoResponse represents pod information response
type PodInfoResponse struct {
	Running int `json:"running"`
	Pending int `json:"pending"`
	Failed  int `json:"failed"`
	Total   int `json:"total"`
}

// DeploymentReplicasResponse represents deployment replicas response
type DeploymentReplicasResponse struct {
	Desired   int32 `json:"desired"`
	Current   int32 `json:"current"`
	Ready     int32 `json:"ready"`
	Available int32 `json:"available"`
}

// DeploymentConditionResponse represents deployment condition response
type DeploymentConditionResponse struct {
	Type           string    `json:"type"`
	Status         string    `json:"status"`
	Reason         string    `json:"reason,omitempty"`
	Message        string    `json:"message,omitempty"`
	LastUpdateTime time.Time `json:"last_update_time"`
}

// ServiceContextResponse represents service context response
type ServiceContextResponse struct {
	ServiceName string `json:"service_name,omitempty"`
	ServiceTier string `json:"service_tier,omitempty"`
	Environment string `json:"environment,omitempty"`
	Context     string `json:"context,omitempty"`
	Team        string `json:"team,omitempty"`
	Description string `json:"description,omitempty"`
	Found       bool   `json:"found"`
}

// PodResponse represents pod information response
type PodResponse struct {
	Name       string                 `json:"name"`
	Namespace  string                 `json:"namespace"`
	Status     string                 `json:"status"`
	Phase      string                 `json:"phase"`
	Node       string                 `json:"node"`
	Age        string                 `json:"age"`
	Restarts   int32                  `json:"restarts"`
	Ready      string                 `json:"ready"`
	IP         string                 `json:"ip,omitempty"`
	Containers []ContainerResponse    `json:"containers"`
	Conditions []PodConditionResponse `json:"conditions"`
	CreatedAt  time.Time              `json:"created_at"`
}

// ContainerResponse represents container information response
type ContainerResponse struct {
	Name         string                     `json:"name"`
	Image        string                     `json:"image"`
	Ready        bool                       `json:"ready"`
	RestartCount int32                      `json:"restart_count"`
	State        ContainerStateResponse     `json:"state"`
	Resources    ContainerResourcesResponse `json:"resources,omitempty"`
}

// ContainerStateResponse represents container state response
type ContainerStateResponse struct {
	Running    *ContainerStateRunningResponse    `json:"running,omitempty"`
	Waiting    *ContainerStateWaitingResponse    `json:"waiting,omitempty"`
	Terminated *ContainerStateTerminatedResponse `json:"terminated,omitempty"`
}

// ContainerStateRunningResponse represents running container state response
type ContainerStateRunningResponse struct {
	StartedAt time.Time `json:"started_at"`
}

// ContainerStateWaitingResponse represents waiting container state response
type ContainerStateWaitingResponse struct {
	Reason  string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`
}

// ContainerStateTerminatedResponse represents terminated container state response
type ContainerStateTerminatedResponse struct {
	ExitCode   int32     `json:"exit_code"`
	Reason     string    `json:"reason,omitempty"`
	Message    string    `json:"message,omitempty"`
	StartedAt  time.Time `json:"started_at,omitempty"`
	FinishedAt time.Time `json:"finished_at,omitempty"`
}

// ContainerResourcesResponse represents container resources response
type ContainerResourcesResponse struct {
	Requests ResourceListResponse `json:"requests,omitempty"`
	Limits   ResourceListResponse `json:"limits,omitempty"`
}

// PodConditionResponse represents pod condition response
type PodConditionResponse struct {
	Type               string    `json:"type"`
	Status             string    `json:"status"`
	Reason             string    `json:"reason,omitempty"`
	Message            string    `json:"message,omitempty"`
	LastTransitionTime time.Time `json:"last_transition_time"`
}

// DeploymentListResponse represents deployment list response
type DeploymentListResponse struct {
	Deployments []DeploymentResponse `json:"deployments"`
	Total       int                  `json:"total"`
	Namespace   string               `json:"namespace,omitempty"`
	Filter      interface{}          `json:"filter,omitempty"`
}

// PodListResponse represents pod list response
type PodListResponse struct {
	Pods      []PodResponse `json:"pods"`
	Total     int           `json:"total"`
	Namespace string        `json:"namespace,omitempty"`
	Filter    interface{}   `json:"filter,omitempty"`
}

// PodLogsResponse represents pod logs response
type PodLogsResponse struct {
	PodName       string                 `json:"pod_name"`
	Namespace     string                 `json:"namespace"`
	ContainerName string                 `json:"container_name,omitempty"`
	Logs          []ContainerLogResponse `json:"logs"`
	TotalLines    int                    `json:"total_lines"`
	Truncated     bool                   `json:"truncated,omitempty"`
}

// ContainerLogResponse represents container log response
type ContainerLogResponse struct {
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
	Level     string    `json:"level,omitempty"`
}

// OperationResponse represents operation result response
type OperationResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	Operation string `json:"operation"`
	Resource  string `json:"resource"`
}

// BatchOperationResponse represents batch operation result response
type BatchOperationResponse struct {
	TotalRequested int                 `json:"total_requested"`
	Successful     int                 `json:"successful"`
	Failed         int                 `json:"failed"`
	Results        []OperationResponse `json:"results"`
}

// ClusterHealthResponse represents cluster health response
type ClusterHealthResponse struct {
	Context     string                    `json:"context"`
	Status      string                    `json:"status"`
	Nodes       []NodeHealthResponse      `json:"nodes"`
	Namespaces  []NamespaceHealthResponse `json:"namespaces,omitempty"`
	Summary     ClusterSummaryResponse    `json:"summary"`
	LastUpdated time.Time                 `json:"last_updated"`
}

// NodeHealthResponse represents node health response
type NodeHealthResponse struct {
	Name        string                  `json:"name"`
	Status      string                  `json:"status"`
	Conditions  []NodeConditionResponse `json:"conditions"`
	Resources   ResourceHealthResponse  `json:"resources"`
	LastUpdated time.Time               `json:"last_updated"`
}

// NamespaceHealthResponse represents namespace health response
type NamespaceHealthResponse struct {
	Name        string                     `json:"name"`
	Status      string                     `json:"status"`
	Deployments []DeploymentHealthResponse `json:"deployments"`
	Pods        []PodHealthResponse        `json:"pods"`
	LastUpdated time.Time                  `json:"last_updated"`
}

// DeploymentHealthResponse represents deployment health response
type DeploymentHealthResponse struct {
	Name                string                        `json:"name"`
	Namespace           string                        `json:"namespace"`
	Status              string                        `json:"status"`
	Replicas            DeploymentReplicasResponse    `json:"replicas"`
	Conditions          []DeploymentConditionResponse `json:"conditions"`
	AvailabilityPercent float64                       `json:"availability_percent"`
	LastUpdated         time.Time                     `json:"last_updated"`
}

// PodHealthResponse represents pod health response
type PodHealthResponse struct {
	Name        string                    `json:"name"`
	Namespace   string                    `json:"namespace"`
	Status      string                    `json:"status"`
	Phase       string                    `json:"phase"`
	Ready       bool                      `json:"ready"`
	Restarts    int32                     `json:"restarts"`
	Containers  []ContainerHealthResponse `json:"containers"`
	LastUpdated time.Time                 `json:"last_updated"`
}

// ContainerHealthResponse represents container health response
type ContainerHealthResponse struct {
	Name         string    `json:"name"`
	Ready        bool      `json:"ready"`
	RestartCount int32     `json:"restart_count"`
	State        string    `json:"state"`
	LastUpdated  time.Time `json:"last_updated"`
}

// ResourceHealthResponse represents resource health response
type ResourceHealthResponse struct {
	CPU    ResourceHealthDetailResponse `json:"cpu"`
	Memory ResourceHealthDetailResponse `json:"memory"`
	Pods   ResourceHealthDetailResponse `json:"pods,omitempty"`
}

// ResourceHealthDetailResponse represents detailed resource health response
type ResourceHealthDetailResponse struct {
	Used               int64   `json:"used"`
	Available          int64   `json:"available"`
	Total              int64   `json:"total"`
	UtilizationPercent float64 `json:"utilization_percent"`
	Status             string  `json:"status"`
}
