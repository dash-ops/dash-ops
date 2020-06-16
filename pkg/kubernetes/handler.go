package kubernetes

import (
	"fmt"
	"net/http"

	"github.com/dash-ops/dash-ops/pkg/commons"
	"github.com/gorilla/mux"
)

func k8sNamespacesHandler(k8sClient K8sClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deployment, err := k8sClient.GetNamespaces()
		if err != nil {
			commons.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		commons.RespondJSON(w, http.StatusOK, deployment)
	}
}

func k8sDeploymentsHandler(k8sClient K8sClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		v := r.URL.Query()

		namespace := v.Get("namespace")
		deployment, err := k8sClient.GetDeployments(deploymentFilter{
			Namespace: namespace,
		})
		if err != nil {
			commons.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		commons.RespondJSON(w, http.StatusOK, deployment)
	}
}

func k8sDeploymentUpHandler(k8sClient K8sClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		err := k8sClient.Scale(vars["name"], vars["namespace"], 1)
		if err != nil {
			commons.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		fmt.Println("Up Deployment: ", vars["namespace"], vars["name"])
		commons.RespondJSON(w, http.StatusOK, nil)
	}
}

func k8sDeploymentDownHandler(k8sClient K8sClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		err := k8sClient.Scale(vars["name"], vars["namespace"], 0)
		if err != nil {
			commons.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		fmt.Println("Down Deployment: ", vars["namespace"], vars["name"])
		commons.RespondJSON(w, http.StatusOK, nil)
	}
}

// MakeKubernetesHandlers Add kubernetes module endpoints
func MakeKubernetesHandlers(r *mux.Router, fileConfig []byte) {
	dashConfig := loadConfig(fileConfig)

	k8sClient, err := NewK8sClient(dashConfig.Kubernetes)
	if err != nil {
		fmt.Println(err.Error())
	}

	r.HandleFunc("/k8s/namespaces", k8sNamespacesHandler(k8sClient)).
		Methods("GET", "OPTIONS").
		Name("k8sNamespaces")

	r.HandleFunc("/k8s/deployments", k8sDeploymentsHandler(k8sClient)).
		Methods("GET", "OPTIONS").
		Name("k8sDeployments")

	r.HandleFunc("/k8s/deployment/up/{namespace}/{name}", k8sDeploymentUpHandler(k8sClient)).
		Methods("POST", "OPTIONS").
		Name("k8sDeploymentUp")

	r.HandleFunc("/k8s/deployment/down/{namespace}/{name}", k8sDeploymentDownHandler(k8sClient)).
		Methods("POST", "OPTIONS").
		Name("k8sDeploymentDown")
}
