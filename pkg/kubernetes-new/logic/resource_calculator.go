package kubernetes

import (
	"fmt"
	"strconv"
	"strings"

	k8sModels "github.com/dash-ops/dash-ops/pkg/kubernetes-new/models"
)

// ResourceCalculator provides resource calculation and parsing logic
type ResourceCalculator struct{}

// NewResourceCalculator creates a new resource calculator
func NewResourceCalculator() *ResourceCalculator {
	return &ResourceCalculator{}
}

// ParseCPU parses CPU resource string to millicores
func (rc *ResourceCalculator) ParseCPU(cpu string) (int64, error) {
	if cpu == "" {
		return 0, nil
	}

	cpu = strings.TrimSpace(cpu)

	// Handle millicore format (e.g., "100m")
	if strings.HasSuffix(cpu, "m") {
		milliStr := strings.TrimSuffix(cpu, "m")
		milli, err := strconv.ParseInt(milliStr, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid CPU millicore format: %s", cpu)
		}
		return milli, nil
	}

	// Handle core format (e.g., "1", "0.5")
	if strings.Contains(cpu, ".") {
		cores, err := strconv.ParseFloat(cpu, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid CPU core format: %s", cpu)
		}
		return int64(cores * 1000), nil
	}

	// Handle integer cores
	cores, err := strconv.ParseInt(cpu, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid CPU format: %s", cpu)
	}
	return cores * 1000, nil
}

// ParseMemory parses memory resource string to bytes
func (rc *ResourceCalculator) ParseMemory(memory string) (int64, error) {
	if memory == "" {
		return 0, nil
	}

	memory = strings.TrimSpace(memory)

	// Define multipliers
	multipliers := map[string]int64{
		"Ki": 1024,
		"Mi": 1024 * 1024,
		"Gi": 1024 * 1024 * 1024,
		"Ti": 1024 * 1024 * 1024 * 1024,
		"K":  1000,
		"M":  1000 * 1000,
		"G":  1000 * 1000 * 1000,
		"T":  1000 * 1000 * 1000 * 1000,
	}

	// Find suffix
	for suffix, multiplier := range multipliers {
		if strings.HasSuffix(memory, suffix) {
			valueStr := strings.TrimSuffix(memory, suffix)
			value, err := strconv.ParseInt(valueStr, 10, 64)
			if err != nil {
				return 0, fmt.Errorf("invalid memory format: %s", memory)
			}
			return value * multiplier, nil
		}
	}

	// No suffix, assume bytes
	bytes, err := strconv.ParseInt(memory, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid memory format: %s", memory)
	}
	return bytes, nil
}

// FormatCPU formats millicores to human-readable string
func (rc *ResourceCalculator) FormatCPU(millicores int64) string {
	if millicores == 0 {
		return "0"
	}

	if millicores < 1000 {
		return fmt.Sprintf("%dm", millicores)
	}

	cores := float64(millicores) / 1000
	if cores == float64(int64(cores)) {
		return fmt.Sprintf("%.0f", cores)
	}
	return fmt.Sprintf("%.1f", cores)
}

// FormatMemory formats bytes to human-readable string
func (rc *ResourceCalculator) FormatMemory(bytes int64) string {
	if bytes == 0 {
		return "0"
	}

	const (
		Ki = 1024
		Mi = Ki * 1024
		Gi = Mi * 1024
		Ti = Gi * 1024
	)

	switch {
	case bytes >= Ti:
		return fmt.Sprintf("%.1fTi", float64(bytes)/Ti)
	case bytes >= Gi:
		return fmt.Sprintf("%.1fGi", float64(bytes)/Gi)
	case bytes >= Mi:
		return fmt.Sprintf("%.0fMi", float64(bytes)/Mi)
	case bytes >= Ki:
		return fmt.Sprintf("%.0fKi", float64(bytes)/Ki)
	default:
		return fmt.Sprintf("%d", bytes)
	}
}

// CalculateResourceUtilization calculates resource utilization percentage
func (rc *ResourceCalculator) CalculateResourceUtilization(used, total int64) float64 {
	if total == 0 {
		return 0
	}
	return float64(used) / float64(total) * 100
}

// CalculateNodeResourceSummary calculates node resource summary
func (rc *ResourceCalculator) CalculateNodeResourceSummary(node *k8sModels.Node) (*NodeResourceSummary, error) {
	if node == nil {
		return nil, fmt.Errorf("node cannot be nil")
	}

	// Parse CPU resources
	capacityCPU, err := rc.ParseCPU(node.Resources.Capacity.CPU)
	if err != nil {
		return nil, fmt.Errorf("failed to parse capacity CPU: %w", err)
	}

	allocatableCPU, err := rc.ParseCPU(node.Resources.Allocatable.CPU)
	if err != nil {
		return nil, fmt.Errorf("failed to parse allocatable CPU: %w", err)
	}

	usedCPU, _ := rc.ParseCPU(node.Resources.Used.CPU) // Used might not be available

	// Parse Memory resources
	capacityMemory, err := rc.ParseMemory(node.Resources.Capacity.Memory)
	if err != nil {
		return nil, fmt.Errorf("failed to parse capacity memory: %w", err)
	}

	allocatableMemory, err := rc.ParseMemory(node.Resources.Allocatable.Memory)
	if err != nil {
		return nil, fmt.Errorf("failed to parse allocatable memory: %w", err)
	}

	usedMemory, _ := rc.ParseMemory(node.Resources.Used.Memory) // Used might not be available

	return &NodeResourceSummary{
		NodeName: node.Name,
		CPU: ResourceSummary{
			Capacity:           capacityCPU,
			Allocatable:        allocatableCPU,
			Used:               usedCPU,
			UtilizationPercent: rc.CalculateResourceUtilization(usedCPU, allocatableCPU),
		},
		Memory: ResourceSummary{
			Capacity:           capacityMemory,
			Allocatable:        allocatableMemory,
			Used:               usedMemory,
			UtilizationPercent: rc.CalculateResourceUtilization(usedMemory, allocatableMemory),
		},
	}, nil
}

// NodeResourceSummary represents node resource summary
type NodeResourceSummary struct {
	NodeName string          `json:"node_name"`
	CPU      ResourceSummary `json:"cpu"`
	Memory   ResourceSummary `json:"memory"`
}

// ResourceSummary represents resource summary for a specific resource type
type ResourceSummary struct {
	Capacity           int64   `json:"capacity"`
	Allocatable        int64   `json:"allocatable"`
	Used               int64   `json:"used"`
	UtilizationPercent float64 `json:"utilization_percent"`
}
