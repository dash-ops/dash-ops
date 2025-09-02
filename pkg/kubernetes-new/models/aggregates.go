package kubernetes

import "time"

// ClusterInfo represents comprehensive cluster information
type ClusterInfo struct {
	Cluster     Cluster        `json:"cluster"`
	Nodes       []Node         `json:"nodes"`
	Namespaces  []Namespace    `json:"namespaces"`
	Summary     ClusterSummary `json:"summary"`
	LastUpdated time.Time      `json:"last_updated"`
}

// ClusterSummary represents cluster summary statistics
type ClusterSummary struct {
	TotalNodes       int `json:"total_nodes"`
	ReadyNodes       int `json:"ready_nodes"`
	TotalNamespaces  int `json:"total_namespaces"`
	TotalDeployments int `json:"total_deployments"`
	TotalPods        int `json:"total_pods"`
	RunningPods      int `json:"running_pods"`
}

// DeploymentList represents a list of deployments with metadata
type DeploymentList struct {
	Deployments []Deployment      `json:"deployments"`
	Total       int               `json:"total"`
	Namespace   string            `json:"namespace,omitempty"`
	Filter      *DeploymentFilter `json:"filter,omitempty"`
}

// DeploymentFilter represents filtering criteria for deployments
type DeploymentFilter struct {
	Namespace     string `json:"namespace,omitempty"`
	LabelSelector string `json:"label_selector,omitempty"`
	FieldSelector string `json:"field_selector,omitempty"`
	ServiceName   string `json:"service_name,omitempty"`
	Status        string `json:"status,omitempty"`
	Limit         int    `json:"limit,omitempty"`
}

// PodList represents a list of pods with metadata
type PodList struct {
	Pods      []Pod      `json:"pods"`
	Total     int        `json:"total"`
	Namespace string     `json:"namespace,omitempty"`
	Filter    *PodFilter `json:"filter,omitempty"`
}

// PodFilter represents filtering criteria for pods
type PodFilter struct {
	Namespace     string `json:"namespace,omitempty"`
	LabelSelector string `json:"label_selector,omitempty"`
	FieldSelector string `json:"field_selector,omitempty"`
	PodName       string `json:"pod_name,omitempty"`
	ContainerName string `json:"container_name,omitempty"`
	Status        string `json:"status,omitempty"`
	Node          string `json:"node,omitempty"`
	Limit         int    `json:"limit,omitempty"`
}

// LogFilter represents filtering criteria for logs
type LogFilter struct {
	Namespace     string     `json:"namespace"`
	PodName       string     `json:"pod_name"`
	ContainerName string     `json:"container_name,omitempty"`
	Follow        bool       `json:"follow,omitempty"`
	TailLines     int64      `json:"tail_lines,omitempty"`
	SinceTime     *time.Time `json:"since_time,omitempty"`
	Previous      bool       `json:"previous,omitempty"`
}

// Methods for ClusterInfo

// CalculateSummary calculates cluster summary statistics
func (ci *ClusterInfo) CalculateSummary() {
	ci.Summary.TotalNodes = len(ci.Nodes)
	ci.Summary.TotalNamespaces = len(ci.Namespaces)

	for _, node := range ci.Nodes {
		if node.IsReady() {
			ci.Summary.ReadyNodes++
		}
	}
}

// GetMasterNodes returns all master/control-plane nodes
func (ci *ClusterInfo) GetMasterNodes() []Node {
	var masters []Node
	for _, node := range ci.Nodes {
		if node.IsMaster() {
			masters = append(masters, node)
		}
	}
	return masters
}

// GetWorkerNodes returns all worker nodes
func (ci *ClusterInfo) GetWorkerNodes() []Node {
	var workers []Node
	for _, node := range ci.Nodes {
		if !node.IsMaster() {
			workers = append(workers, node)
		}
	}
	return workers
}

// Methods for DeploymentList

// FilterByNamespace filters deployments by namespace
func (dl *DeploymentList) FilterByNamespace(namespace string) *DeploymentList {
	if namespace == "" {
		return dl
	}

	var filtered []Deployment
	for _, deployment := range dl.Deployments {
		if deployment.Namespace == namespace {
			filtered = append(filtered, deployment)
		}
	}

	return &DeploymentList{
		Deployments: filtered,
		Total:       len(filtered),
		Namespace:   namespace,
		Filter:      &DeploymentFilter{Namespace: namespace},
	}
}

// FilterByStatus filters deployments by health status
func (dl *DeploymentList) FilterByStatus(healthy bool) *DeploymentList {
	var filtered []Deployment
	for _, deployment := range dl.Deployments {
		if deployment.IsHealthy() == healthy {
			filtered = append(filtered, deployment)
		}
	}

	status := "unhealthy"
	if healthy {
		status = "healthy"
	}

	return &DeploymentList{
		Deployments: filtered,
		Total:       len(filtered),
		Filter:      &DeploymentFilter{Status: status},
	}
}

// GetHealthyDeployments returns only healthy deployments
func (dl *DeploymentList) GetHealthyDeployments() []Deployment {
	var healthy []Deployment
	for _, deployment := range dl.Deployments {
		if deployment.IsHealthy() {
			healthy = append(healthy, deployment)
		}
	}
	return healthy
}

// GetUnhealthyDeployments returns only unhealthy deployments
func (dl *DeploymentList) GetUnhealthyDeployments() []Deployment {
	var unhealthy []Deployment
	for _, deployment := range dl.Deployments {
		if !deployment.IsHealthy() {
			unhealthy = append(unhealthy, deployment)
		}
	}
	return unhealthy
}

// Methods for PodList

// FilterByStatus filters pods by status
func (pl *PodList) FilterByStatus(status PodStatus) *PodList {
	var filtered []Pod
	for _, pod := range pl.Pods {
		if pod.Status == status {
			filtered = append(filtered, pod)
		}
	}

	return &PodList{
		Pods:   filtered,
		Total:  len(filtered),
		Filter: &PodFilter{Status: string(status)},
	}
}

// GetRunningPods returns only running pods
func (pl *PodList) GetRunningPods() []Pod {
	var running []Pod
	for _, pod := range pl.Pods {
		if pod.IsRunning() {
			running = append(running, pod)
		}
	}
	return running
}

// GetProblematicPods returns pods that are not running or not ready
func (pl *PodList) GetProblematicPods() []Pod {
	var problematic []Pod
	for _, pod := range pl.Pods {
		if !pod.IsRunning() || !pod.IsReady() || pod.GetTotalRestarts() > 5 {
			problematic = append(problematic, pod)
		}
	}
	return problematic
}

// PodMetrics represents pod resource metrics
type PodMetrics struct {
	PodName     string             `json:"pod_name"`
	Namespace   string             `json:"namespace"`
	Containers  []ContainerMetrics `json:"containers"`
	LastUpdated time.Time          `json:"last_updated"`
}

// ContainerMetrics represents container resource metrics
type ContainerMetrics struct {
	Name      string             `json:"name"`
	Resources ContainerResources `json:"resources"`
	Usage     ResourceList       `json:"usage"`
}
