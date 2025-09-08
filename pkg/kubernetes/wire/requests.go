package kubernetes

// ScaleDeploymentRequest represents deployment scaling request
type ScaleDeploymentRequest struct {
	Replicas int32 `json:"replicas" validate:"min=0,max=100"`
}

// CreateNamespaceRequest represents namespace creation request
type CreateNamespaceRequest struct {
	Name   string            `json:"name" validate:"required,min=1,max=63"`
	Labels map[string]string `json:"labels,omitempty"`
}

// PodLogsRequest represents pod logs request
type PodLogsRequest struct {
	ContainerName string `json:"container_name,omitempty"`
	Follow        bool   `json:"follow,omitempty"`
	TailLines     *int64 `json:"tail_lines,omitempty" validate:"omitempty,min=1,max=10000"`
	SinceSeconds  *int64 `json:"since_seconds,omitempty" validate:"omitempty,min=1"`
	Previous      bool   `json:"previous,omitempty"`
	Timestamps    bool   `json:"timestamps,omitempty"`
}

// DeploymentFilterRequest represents deployment filtering request
type DeploymentFilterRequest struct {
	Namespace     string `json:"namespace,omitempty"`
	LabelSelector string `json:"label_selector,omitempty"`
	ServiceName   string `json:"service_name,omitempty"`
	Status        string `json:"status,omitempty" validate:"omitempty,oneof=healthy degraded unhealthy"`
	Limit         int    `json:"limit,omitempty" validate:"omitempty,min=1,max=1000"`
}

// PodFilterRequest represents pod filtering request
type PodFilterRequest struct {
	Namespace     string `json:"namespace,omitempty"`
	LabelSelector string `json:"label_selector,omitempty"`
	FieldSelector string `json:"field_selector,omitempty"`
	PodName       string `json:"pod_name,omitempty"`
	ContainerName string `json:"container_name,omitempty"`
	Status        string `json:"status,omitempty" validate:"omitempty,oneof=Running Pending Succeeded Failed Unknown"`
	Node          string `json:"node,omitempty"`
	Limit         int    `json:"limit,omitempty" validate:"omitempty,min=1,max=1000"`
}

// NodeFilterRequest represents node filtering request
type NodeFilterRequest struct {
	LabelSelector string `json:"label_selector,omitempty"`
	Role          string `json:"role,omitempty" validate:"omitempty,oneof=master control-plane worker"`
	Status        string `json:"status,omitempty" validate:"omitempty,oneof=Ready NotReady Unknown"`
}

// BatchOperationRequest represents batch operation request
type BatchOperationRequest struct {
	Resources []ResourceTarget `json:"resources" validate:"required,min=1,max=50"`
	Operation string           `json:"operation" validate:"required,oneof=delete restart scale"`
	Replicas  *int32           `json:"replicas,omitempty"` // For scale operations
}

// ResourceTarget represents a target resource for batch operations
type ResourceTarget struct {
	Type      string `json:"type" validate:"required,oneof=deployment pod"`
	Namespace string `json:"namespace" validate:"required"`
	Name      string `json:"name" validate:"required"`
}
