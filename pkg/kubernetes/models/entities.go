package models

import (
	"fmt"
	"time"
)

// Cluster represents a Kubernetes cluster configuration
type Cluster struct {
	Name       string        `yaml:"name" json:"name"`
	Context    string        `yaml:"context" json:"context"`
	Kubeconfig string        `yaml:"kubeconfig" json:"kubeconfig"`
	Server     string        `json:"server,omitempty"`
	Version    string        `json:"version,omitempty"`
	Status     ClusterStatus `json:"status"`
}

// ClusterStatus represents cluster connection status
type ClusterStatus string

const (
	ClusterStatusConnected    ClusterStatus = "connected"
	ClusterStatusDisconnected ClusterStatus = "disconnected"
	ClusterStatusError        ClusterStatus = "error"
	ClusterStatusUnknown      ClusterStatus = "unknown"
	ClusterStatusReady        ClusterStatus = "ready" // Alias for connected
)

// Node represents a Kubernetes node
type Node struct {
	Name       string          `json:"name"`
	Status     NodeStatus      `json:"status"`
	Roles      []string        `json:"roles"`
	Age        string          `json:"age"`
	Version    string          `json:"version"`
	InternalIP string          `json:"internal_ip"`
	ExternalIP string          `json:"external_ip,omitempty"`
	Conditions []NodeCondition `json:"conditions"`
	Resources  NodeResources   `json:"resources"`
	CreatedAt  time.Time       `json:"created_at"`
}

// NodeStatus represents node operational status
type NodeStatus string

const (
	NodeStatusReady    NodeStatus = "Ready"
	NodeStatusNotReady NodeStatus = "NotReady"
	NodeStatusUnknown  NodeStatus = "Unknown"
)

// NodeCondition represents a node condition
type NodeCondition struct {
	Type               string    `json:"type"`
	Status             string    `json:"status"`
	Reason             string    `json:"reason,omitempty"`
	Message            string    `json:"message,omitempty"`
	LastTransitionTime time.Time `json:"last_transition_time"`
}

// NodeResources represents node resource information
type NodeResources struct {
	Capacity    ResourceList `json:"capacity"`
	Allocatable ResourceList `json:"allocatable"`
	Used        ResourceList `json:"used,omitempty"`
}

// ResourceList represents resource quantities
type ResourceList struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
	Pods   string `json:"pods,omitempty"`
}

// Namespace represents a Kubernetes namespace
type Namespace struct {
	Name      string            `json:"name"`
	Status    NamespaceStatus   `json:"status"`
	Labels    map[string]string `json:"labels,omitempty"`
	Age       string            `json:"age"`
	CreatedAt time.Time         `json:"created_at"`
}

// NamespaceStatus represents namespace status
type NamespaceStatus string

const (
	NamespaceStatusActive      NamespaceStatus = "Active"
	NamespaceStatusTerminating NamespaceStatus = "Terminating"
)

// Deployment represents a Kubernetes deployment
type Deployment struct {
	Name           string                `json:"name"`
	Namespace      string                `json:"namespace"`
	PodInfo        PodInfo               `json:"pod_info"`
	Replicas       DeploymentReplicas    `json:"replicas"`
	Age            string                `json:"age"`
	CreatedAt      time.Time             `json:"created_at"`
	Conditions     []DeploymentCondition `json:"conditions"`
	ServiceContext *ServiceContext       `json:"service_context,omitempty"`
}

// PodInfo represents pod information summary
type PodInfo struct {
	Running int `json:"running"`
	Pending int `json:"pending"`
	Failed  int `json:"failed"`
	Total   int `json:"total"`
}

// DeploymentReplicas represents deployment replica information
type DeploymentReplicas struct {
	Desired   int32 `json:"desired"`
	Current   int32 `json:"current"`
	Ready     int32 `json:"ready"`
	Available int32 `json:"available"`
}

// DeploymentCondition represents a deployment condition
type DeploymentCondition struct {
	Type           string    `json:"type"`
	Status         string    `json:"status"`
	Reason         string    `json:"reason,omitempty"`
	Message        string    `json:"message,omitempty"`
	LastUpdateTime time.Time `json:"last_update_time"`
}

// Pod represents a Kubernetes pod
type Pod struct {
	Name       string         `json:"name"`
	Namespace  string         `json:"namespace"`
	Status     PodStatus      `json:"status"`
	Phase      string         `json:"phase"`
	Node       string         `json:"node"`
	Age        string         `json:"age"`
	Restarts   int32          `json:"restarts"`
	Ready      string         `json:"ready"`
	IP         string         `json:"ip,omitempty"`
	Containers []Container    `json:"containers"`
	Conditions []PodCondition `json:"conditions"`
	CreatedAt  time.Time      `json:"created_at"`
	QoSClass   string         `json:"qos_class,omitempty"`
}

// PodStatus represents pod operational status
type PodStatus string

const (
	PodStatusRunning   PodStatus = "Running"
	PodStatusPending   PodStatus = "Pending"
	PodStatusSucceeded PodStatus = "Succeeded"
	PodStatusFailed    PodStatus = "Failed"
	PodStatusUnknown   PodStatus = "Unknown"
)

// Container represents a container within a pod
type Container struct {
	Name         string             `json:"name"`
	Image        string             `json:"image"`
	Ready        bool               `json:"ready"`
	RestartCount int32              `json:"restart_count"`
	State        ContainerState     `json:"state"`
	Resources    ContainerResources `json:"resources,omitempty"`
}

// ContainerState represents container state
type ContainerState struct {
	Running    *ContainerStateRunning    `json:"running,omitempty"`
	Waiting    *ContainerStateWaiting    `json:"waiting,omitempty"`
	Terminated *ContainerStateTerminated `json:"terminated,omitempty"`
}

// ContainerStateRunning represents running container state
type ContainerStateRunning struct {
	StartedAt time.Time `json:"started_at"`
}

// ContainerStateWaiting represents waiting container state
type ContainerStateWaiting struct {
	Reason  string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`
}

// ContainerStateTerminated represents terminated container state
type ContainerStateTerminated struct {
	ExitCode   int32     `json:"exit_code"`
	Reason     string    `json:"reason,omitempty"`
	Message    string    `json:"message,omitempty"`
	StartedAt  time.Time `json:"started_at,omitempty"`
	FinishedAt time.Time `json:"finished_at,omitempty"`
}

// ContainerResources represents container resource information
type ContainerResources struct {
	Requests ResourceList `json:"requests,omitempty"`
	Limits   ResourceList `json:"limits,omitempty"`
}

// PodCondition represents a pod condition
type PodCondition struct {
	Type               string    `json:"type"`
	Status             string    `json:"status"`
	Reason             string    `json:"reason,omitempty"`
	Message            string    `json:"message,omitempty"`
	LastTransitionTime time.Time `json:"last_transition_time"`
}

// ServiceContext represents service information linked to a Kubernetes resource
type ServiceContext struct {
	ServiceName string `json:"service_name,omitempty"`
	ServiceTier string `json:"service_tier,omitempty"`
	Environment string `json:"environment,omitempty"`
	Context     string `json:"context,omitempty"`
	Team        string `json:"team,omitempty"`
	Description string `json:"description,omitempty"`
	Found       bool   `json:"found"`
}

// ContainerLog represents container log information
type ContainerLog struct {
	ContainerName string    `json:"container_name"`
	PodName       string    `json:"pod_name"`
	Namespace     string    `json:"namespace"`
	Timestamp     time.Time `json:"timestamp"`
	Message       string    `json:"message"`
	Level         string    `json:"level,omitempty"`
}

// Domain methods for Cluster

// IsConnected checks if cluster is connected
func (c *Cluster) IsConnected() bool {
	return c.Status == ClusterStatusConnected
}

// Validate validates cluster configuration
func (c *Cluster) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("cluster name is required")
	}
	if c.Context == "" {
		return fmt.Errorf("cluster context is required")
	}
	return nil
}

// Domain methods for Node

// IsReady checks if node is ready
func (n *Node) IsReady() bool {
	return n.Status == NodeStatusReady
}

// HasRole checks if node has a specific role
func (n *Node) HasRole(role string) bool {
	for _, r := range n.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// IsMaster checks if node is a master/control-plane node
func (n *Node) IsMaster() bool {
	return n.HasRole("master") || n.HasRole("control-plane")
}

// Domain methods for Deployment

// IsHealthy checks if deployment is healthy
func (d *Deployment) IsHealthy() bool {
	return d.Replicas.Ready == d.Replicas.Desired && d.Replicas.Desired > 0
}

// GetAvailabilityPercentage returns deployment availability percentage
func (d *Deployment) GetAvailabilityPercentage() float64 {
	if d.Replicas.Desired == 0 {
		return 0
	}
	return float64(d.Replicas.Ready) / float64(d.Replicas.Desired) * 100
}

// HasServiceContext checks if deployment has service context
func (d *Deployment) HasServiceContext() bool {
	return d.ServiceContext != nil && d.ServiceContext.Found
}

// Domain methods for Pod

// IsRunning checks if pod is running
func (p *Pod) IsRunning() bool {
	return p.Status == PodStatusRunning
}

// IsReady checks if all containers in pod are ready
func (p *Pod) IsReady() bool {
	for _, container := range p.Containers {
		if !container.Ready {
			return false
		}
	}
	return len(p.Containers) > 0
}

// GetTotalRestarts returns total restart count for all containers
func (p *Pod) GetTotalRestarts() int32 {
	var total int32
	for _, container := range p.Containers {
		total += container.RestartCount
	}
	return total
}

// Domain methods for Container

// IsReady checks if container is ready
func (c *Container) IsReady() bool {
	return c.Ready
}

// IsRunning checks if container is running
func (c *Container) IsRunning() bool {
	return c.State.Running != nil
}

// IsWaiting checks if container is waiting
func (c *Container) IsWaiting() bool {
	return c.State.Waiting != nil
}

// IsTerminated checks if container is terminated
func (c *Container) IsTerminated() bool {
	return c.State.Terminated != nil
}
