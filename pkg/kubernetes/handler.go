package kubernetes

import (
	"fmt"
	"net/http"

	"github.com/dash-ops/dash-ops/pkg/commons"
	"github.com/gorilla/mux"
)

// Cluster config
type Cluster struct {
	Name    string `json:"name"`
	Context string `json:"context"`
}

func k8sClustersHandler(dashConfig dashYaml) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var clusters []Cluster

		for _, cluster := range dashConfig.Kubernetes {
			clusters = append(clusters, Cluster{
				Name:    cluster.Name,
				Context: cluster.Context,
			})
		}

		commons.RespondJSON(w, http.StatusOK, clusters)
	}
}

func k8sPermissionsHandler(permission k8sPermission) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		commons.RespondJSON(w, http.StatusOK, permission)
	}
}

func k8sNodesHandler(k8sClient K8sClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nodes, err := k8sClient.GetNodes()
		if err != nil {
			commons.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		commons.RespondJSON(w, http.StatusOK, nodes)
	}
}

func k8sNamespacesHandler(k8sClient K8sClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		namespaces, err := k8sClient.GetNamespaces()
		if err != nil {
			commons.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		commons.RespondJSON(w, http.StatusOK, namespaces)
	}
}

func k8sDeploymentsHandler(k8sClient K8sClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		namespace := query.Get("namespace")
		deployments, err := k8sClient.GetDeployments(deploymentFilter{
			Namespace: namespace,
		})
		if err != nil {
			commons.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		commons.RespondJSON(w, http.StatusOK, deployments)
	}
}

func k8sDeploymentUpHandler(k8sClient K8sClient, permission k8sPermission) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userData := r.Context().Value(commons.UserDataKey).(commons.UserData)
		if isValid := commons.HasPermission(permission.Deployments.Start, userData.Groups); !isValid {
			commons.RespondError(w, http.StatusForbidden, "you do not have permission")
			return
		}

		vars := mux.Vars(r)
		err := k8sClient.Scale(vars["name"], vars["namespace"], 1)
		if err != nil {
			commons.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		commons.RespondJSON(w, http.StatusOK, nil)
	}
}

func k8sDeploymentDownHandler(k8sClient K8sClient, permission k8sPermission) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userData := r.Context().Value(commons.UserDataKey).(commons.UserData)
		if isValid := commons.HasPermission(permission.Deployments.Stop, userData.Groups); !isValid {
			commons.RespondError(w, http.StatusForbidden, "you do not have permission")
			return
		}

		vars := mux.Vars(r)
		err := k8sClient.Scale(vars["name"], vars["namespace"], 0)
		if err != nil {
			commons.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		commons.RespondJSON(w, http.StatusOK, nil)
	}
}

func k8sPodsHandler(k8sClient K8sClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		namespace := query.Get("namespace")
		pods, err := k8sClient.GetPods(podFilter{
			Namespace: namespace,
		})
		if err != nil {
			commons.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		commons.RespondJSON(w, http.StatusOK, pods)
	}
}

func k8sPodLogsHandler(k8sClient K8sClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		query := r.URL.Query()

		logs, err := k8sClient.GetPodLogs(podFilter{
			Name:      vars["name"],
			Namespace: query.Get("namespace"),
		})
		if err != nil {
			commons.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		commons.RespondJSON(w, http.StatusOK, logs)
	}
}

// MakeKubernetesHandlers Add kubernetes module endpoints
func MakeKubernetesHandlers(r *mux.Router, fileConfig []byte) {
	dashConfig := loadConfig(fileConfig)

	r.HandleFunc("/k8s/clusters", k8sClustersHandler(dashConfig)).
		Methods("GET", "OPTIONS").
		Name("k8sClusters")

	for _, cluster := range dashConfig.Kubernetes {
		k8sClient, err := NewK8sClient(cluster)
		if err != nil {
			fmt.Println(err.Error())
		}

		contextRoute := r.PathPrefix("/k8s/" + cluster.Context).Subrouter()

		contextRoute.HandleFunc("/permissions", k8sPermissionsHandler(cluster.Permission)).
			Methods("GET", "OPTIONS").
			Name("k8sPermissions")

		contextRoute.HandleFunc("/nodes", k8sNodesHandler(k8sClient)).
			Methods("GET", "OPTIONS").
			Name("k8sNodes")

		contextRoute.HandleFunc("/namespaces", k8sNamespacesHandler(k8sClient)).
			Methods("GET", "OPTIONS").
			Name("k8sNamespaces")

		contextRoute.HandleFunc("/deployments", k8sDeploymentsHandler(k8sClient)).
			Methods("GET", "OPTIONS").
			Name("k8sDeployments")

		contextRoute.HandleFunc("/deployment/up/{namespace}/{name}", k8sDeploymentUpHandler(k8sClient, cluster.Permission)).
			Methods("POST", "OPTIONS").
			Name("k8sDeploymentUp")

		contextRoute.HandleFunc("/deployment/down/{namespace}/{name}", k8sDeploymentDownHandler(k8sClient, cluster.Permission)).
			Methods("POST", "OPTIONS").
			Name("k8sDeploymentDown")

		contextRoute.HandleFunc("/pods", k8sPodsHandler(k8sClient)).
			Methods("GET", "OPTIONS").
			Name("k8sPods")

		contextRoute.HandleFunc("/pod/{name}/logs", k8sPodLogsHandler(k8sClient)).
			Queries("namespace", "{namespace}").
			Methods("GET", "OPTIONS").
			Name("k8sPodLogs")
	}
}
