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

func clustersHandler(dashConfig dashYaml) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var clusters []Cluster

		for _, cluster := range dashConfig.Kubernetes {
			c := Cluster{
				Name:    cluster.Name,
				Context: cluster.Context,
			}

			if cluster.Context == "" {
				c.Context = "default"
			}

			clusters = append(clusters, c)
		}

		commons.RespondJSON(w, http.StatusOK, clusters)
	}
}

func permissionsHandler(permission permission) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		commons.RespondJSON(w, http.StatusOK, permission)
	}
}

func nodesHandler(client Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nodes, err := client.GetNodes()
		if err != nil {
			commons.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		commons.RespondJSON(w, http.StatusOK, nodes)
	}
}

func namespacesHandler(client Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		namespaces, err := client.GetNamespaces()
		if err != nil {
			commons.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		commons.RespondJSON(w, http.StatusOK, namespaces)
	}
}

func deploymentsHandler(client Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		namespace := query.Get("namespace")
		deployments, err := client.GetDeployments(deploymentFilter{
			Namespace: namespace,
		})
		if err != nil {
			commons.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		commons.RespondJSON(w, http.StatusOK, deployments)
	}
}

func deploymentUpHandler(client Client, permission permission) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userData := r.Context().Value(commons.UserDataKey).(commons.UserData)
		if isValid := commons.HasPermission(permission.Deployments.Start, userData.Groups); !isValid {
			commons.RespondError(w, http.StatusForbidden, "you do not have permission")
			return
		}

		vars := mux.Vars(r)
		if isValid := hasPermissionNamespace(permission.Deployments.Namespaces, vars["namespace"]); !isValid {
			commons.RespondError(w, http.StatusForbidden, "you do not have permission")
			return
		}

		err := client.Scale(vars["name"], vars["namespace"], 1)
		if err != nil {
			commons.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		commons.RespondJSON(w, http.StatusOK, nil)
	}
}

func deploymentDownHandler(client Client, permission permission) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userData := r.Context().Value(commons.UserDataKey).(commons.UserData)
		if isValid := commons.HasPermission(permission.Deployments.Stop, userData.Groups); !isValid {
			commons.RespondError(w, http.StatusForbidden, "you do not have permission")
			return
		}

		vars := mux.Vars(r)
		if isValid := hasPermissionNamespace(permission.Deployments.Namespaces, vars["namespace"]); !isValid {
			commons.RespondError(w, http.StatusForbidden, "you do not have permission")
			return
		}

		err := client.Scale(vars["name"], vars["namespace"], 0)
		if err != nil {
			commons.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		commons.RespondJSON(w, http.StatusOK, nil)
	}
}

func podsHandler(client Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		namespace := query.Get("namespace")
		pods, err := client.GetPods(podFilter{
			Namespace: namespace,
		})
		if err != nil {
			commons.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		commons.RespondJSON(w, http.StatusOK, pods)
	}
}

func podLogsHandler(client Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		query := r.URL.Query()

		logs, err := client.GetPodLogs(podFilter{
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

func hasPermissionNamespace(namespaces []string, namespace string) bool {
	isValid := false

	for _, n := range namespaces {
		if n == namespace {
			isValid = true
		}
	}

	return isValid
}

// MakeKubernetesHandlers Add kubernetes module endpoints
func MakeKubernetesHandlers(r *mux.Router, fileConfig []byte) {
	dashConfig, err := loadConfig(fileConfig)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	r.HandleFunc("/k8s/clusters", clustersHandler(dashConfig)).
		Methods("GET", "OPTIONS").
		Name("k8sClusters")

	for _, cluster := range dashConfig.Kubernetes {
		client, err := NewClient(cluster)
		if err != nil {
			fmt.Println(err.Error())
		}

		contextRoute := r.PathPrefix("/k8s/default").Subrouter()
		if cluster.Context != "" {
			contextRoute = r.PathPrefix("/k8s/" + cluster.Context).Subrouter()
		}

		contextRoute.HandleFunc("/permissions", permissionsHandler(cluster.Permission)).
			Methods("GET", "OPTIONS").
			Name("k8sPermissions")

		contextRoute.HandleFunc("/nodes", nodesHandler(client)).
			Methods("GET", "OPTIONS").
			Name("k8sNodes")

		contextRoute.HandleFunc("/namespaces", namespacesHandler(client)).
			Methods("GET", "OPTIONS").
			Name("k8sNamespaces")

		contextRoute.HandleFunc("/deployments", deploymentsHandler(client)).
			Methods("GET", "OPTIONS").
			Name("k8sDeployments")

		contextRoute.HandleFunc("/deployment/up/{namespace}/{name}", deploymentUpHandler(client, cluster.Permission)).
			Methods("POST", "OPTIONS").
			Name("k8sDeploymentUp")

		contextRoute.HandleFunc("/deployment/down/{namespace}/{name}", deploymentDownHandler(client, cluster.Permission)).
			Methods("POST", "OPTIONS").
			Name("k8sDeploymentDown")

		contextRoute.HandleFunc("/pods", podsHandler(client)).
			Methods("GET", "OPTIONS").
			Name("k8sPods")

		contextRoute.HandleFunc("/pod/{name}/logs", podLogsHandler(client)).
			Queries("namespace", "{namespace}").
			Methods("GET", "OPTIONS").
			Name("k8sPodLogs")
	}
}
