package kubernetes

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Deployment Struct representing an k8s deployment
type Deployment struct {
	Name       string                `json:"name"`
	Namespace  string                `json:"namespace"`
	PodInfo    PodInfo               `json:"pod_info"`
	Replicas   DeploymentReplicas    `json:"replicas"`
	Age        string                `json:"age"`
	CreatedAt  time.Time             `json:"created_at"`
	Conditions []DeploymentCondition `json:"conditions"`
}

// PodInfo Struct
type PodInfo struct {
	Current int32 `json:"current"`
	Desired int32 `json:"desired"`
}

// DeploymentReplicas represents replica information
type DeploymentReplicas struct {
	Ready     int32 `json:"ready"`
	Available int32 `json:"available"`
	Updated   int32 `json:"updated"`
	Desired   int32 `json:"desired"`
}

// DeploymentCondition represents a deployment condition
type DeploymentCondition struct {
	Type    string `json:"type"`
	Status  string `json:"status"`
	Reason  string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`
}

type deploymentFilter struct {
	Namespace string
}

func (kc client) GetDeployments(filter deploymentFilter) ([]Deployment, error) {
	var deployments []Deployment

	if filter.Namespace == "" {
		filter.Namespace = apiv1.NamespaceAll
	}

	deploys, err := kc.clientSet.
		AppsV1().
		Deployments(filter.Namespace).
		List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get deployments: %s", err)
	}

	for _, deploy := range deploys.Items {
		conditions := getDeploymentConditions(deploy)
		age := calculateDeploymentAge(deploy.CreationTimestamp.Time)
		replicas := getDeploymentReplicas(deploy)

		deployments = append(deployments, Deployment{
			Name:      deploy.GetName(),
			Namespace: deploy.GetNamespace(),
			PodInfo: PodInfo{
				Current: deploy.Status.Replicas,
				Desired: *deploy.Spec.Replicas,
			},
			Replicas:   replicas,
			Age:        age,
			CreatedAt:  deploy.CreationTimestamp.Time,
			Conditions: conditions,
		})
	}

	return deployments, nil
}

func getDeploymentConditions(deploy appsv1.Deployment) []DeploymentCondition {
	var conditions []DeploymentCondition
	for _, condition := range deploy.Status.Conditions {
		conditions = append(conditions, DeploymentCondition{
			Type:    string(condition.Type),
			Status:  string(condition.Status),
			Reason:  condition.Reason,
			Message: condition.Message,
		})
	}
	return conditions
}

func getDeploymentReplicas(deploy appsv1.Deployment) DeploymentReplicas {
	return DeploymentReplicas{
		Ready:     deploy.Status.ReadyReplicas,
		Available: deploy.Status.AvailableReplicas,
		Updated:   deploy.Status.UpdatedReplicas,
		Desired:   *deploy.Spec.Replicas,
	}
}

func calculateDeploymentAge(createdAt time.Time) string {
	now := time.Now()
	duration := now.Sub(createdAt)

	days := int(duration.Hours() / 24)
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60

	if days > 0 {
		if hours > 0 {
			return fmt.Sprintf("%dd %dh", days, hours)
		}
		return fmt.Sprintf("%dd", days)
	} else if hours > 0 {
		if minutes > 0 {
			return fmt.Sprintf("%dh %dm", hours, minutes)
		}
		return fmt.Sprintf("%dh", hours)
	} else {
		return fmt.Sprintf("%dm", minutes)
	}
}

func (kc client) Scale(name string, ns string, replicas int32) error {
	deploy, err := kc.clientSet.AppsV1().Deployments(ns).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get deploy %s on ns %s: %s", name, ns, err)
	}
	deploy.Spec.Replicas = &replicas
	_, err = kc.clientSet.AppsV1().Deployments(deploy.GetNamespace()).Update(context.TODO(), deploy, metav1.UpdateOptions{})
	return err
}

func (kc client) RestartDeployment(name string, ns string) error {
	deploy, err := kc.clientSet.AppsV1().Deployments(ns).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get deploy %s on ns %s: %s", name, ns, err)
	}

	// Restart by updating the template annotation
	if deploy.Spec.Template.Annotations == nil {
		deploy.Spec.Template.Annotations = make(map[string]string)
	}
	deploy.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)

	_, err = kc.clientSet.AppsV1().Deployments(deploy.GetNamespace()).Update(context.TODO(), deploy, metav1.UpdateOptions{})
	return err
}
