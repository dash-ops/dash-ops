package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	commonsHttp "github.com/dash-ops/dash-ops/pkg/commons-new/adapters/http"
	k8sAdapters "github.com/dash-ops/dash-ops/pkg/kubernetes-new/adapters/http"
	kubernetes "github.com/dash-ops/dash-ops/pkg/kubernetes-new/controllers"
	k8sModels "github.com/dash-ops/dash-ops/pkg/kubernetes-new/models"
	k8sWire "github.com/dash-ops/dash-ops/pkg/kubernetes-new/wire"
)

// HTTPHandler handles HTTP requests for Kubernetes module
type HTTPHandler struct {
	controller      *kubernetes.KubernetesController
	k8sAdapter      *k8sAdapters.KubernetesAdapter
	responseAdapter *commonsHttp.ResponseAdapter
	requestAdapter  *commonsHttp.RequestAdapter
}

// NewHTTPHandler creates a new HTTP handler
func NewHTTPHandler(
	controller *kubernetes.KubernetesController,
	k8sAdapter *k8sAdapters.KubernetesAdapter,
	responseAdapter *commonsHttp.ResponseAdapter,
	requestAdapter *commonsHttp.RequestAdapter,
) *HTTPHandler {
	return &HTTPHandler{
		controller:      controller,
		k8sAdapter:      k8sAdapter,
		responseAdapter: responseAdapter,
		requestAdapter:  requestAdapter,
	}
}

// RegisterRoutes registers all Kubernetes routes
func (h *HTTPHandler) RegisterRoutes(router *mux.Router) {
	// Cluster operations
	router.HandleFunc("/clusters/{context}", h.getClusterInfoHandler).Methods("GET")
	router.HandleFunc("/clusters/{context}/health", h.getClusterHealthHandler).Methods("GET")

	// Node operations
	router.HandleFunc("/clusters/{context}/nodes", h.listNodesHandler).Methods("GET")
	router.HandleFunc("/clusters/{context}/nodes/{name}", h.getNodeHandler).Methods("GET")

	// Namespace operations
	router.HandleFunc("/clusters/{context}/namespaces", h.listNamespacesHandler).Methods("GET")
	router.HandleFunc("/clusters/{context}/namespaces", h.createNamespaceHandler).Methods("POST")
	router.HandleFunc("/clusters/{context}/namespaces/{name}", h.getNamespaceHandler).Methods("GET")
	router.HandleFunc("/clusters/{context}/namespaces/{name}", h.deleteNamespaceHandler).Methods("DELETE")

	// Deployment operations
	router.HandleFunc("/clusters/{context}/deployments", h.listDeploymentsHandler).Methods("GET")
	router.HandleFunc("/clusters/{context}/namespaces/{namespace}/deployments", h.listDeploymentsByNamespaceHandler).Methods("GET")
	router.HandleFunc("/clusters/{context}/namespaces/{namespace}/deployments/{name}", h.getDeploymentHandler).Methods("GET")
	router.HandleFunc("/clusters/{context}/namespaces/{namespace}/deployments/{name}/scale", h.scaleDeploymentHandler).Methods("PUT")
	router.HandleFunc("/clusters/{context}/namespaces/{namespace}/deployments/{name}/restart", h.restartDeploymentHandler).Methods("POST")

	// Pod operations
	router.HandleFunc("/clusters/{context}/pods", h.listPodsHandler).Methods("GET")
	router.HandleFunc("/clusters/{context}/namespaces/{namespace}/pods", h.listPodsByNamespaceHandler).Methods("GET")
	router.HandleFunc("/clusters/{context}/namespaces/{namespace}/pods/{name}", h.getPodHandler).Methods("GET")
	router.HandleFunc("/clusters/{context}/namespaces/{namespace}/pods/{name}", h.deletePodHandler).Methods("DELETE")
	router.HandleFunc("/clusters/{context}/namespaces/{namespace}/pods/{name}/logs", h.getPodLogsHandler).Methods("GET")
}

// getClusterInfoHandler handles GET /clusters/{context}
func (h *HTTPHandler) getClusterInfoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	context := vars["context"]

	if context == "" {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Context is required")
		return
	}

	clusterInfo, err := h.controller.GetClusterInfo(r.Context(), context)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, "Failed to get cluster info: "+err.Error())
		return
	}

	response := h.k8sAdapter.ClusterInfoToResponse(clusterInfo)
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}

// listNodesHandler handles GET /clusters/{context}/nodes
func (h *HTTPHandler) listNodesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	context := vars["context"]

	if context == "" {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Context is required")
		return
	}

	nodes, err := h.controller.ListNodes(r.Context(), context)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, "Failed to list nodes: "+err.Error())
		return
	}

	response := h.k8sAdapter.NodesToResponse(nodes)
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}

// getNodeHandler handles GET /clusters/{context}/nodes/{name}
func (h *HTTPHandler) getNodeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	context := vars["context"]
	nodeName := vars["name"]

	if context == "" {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Context is required")
		return
	}

	if nodeName == "" {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Node name is required")
		return
	}

	node, err := h.controller.GetNode(r.Context(), context, nodeName)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusNotFound, "Node not found: "+err.Error())
		return
	}

	response := h.k8sAdapter.NodeToResponse(node)
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}

// listNamespacesHandler handles GET /clusters/{context}/namespaces
func (h *HTTPHandler) listNamespacesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	context := vars["context"]

	if context == "" {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Context is required")
		return
	}

	namespaces, err := h.controller.ListNamespaces(r.Context(), context)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, "Failed to list namespaces: "+err.Error())
		return
	}

	response := h.k8sAdapter.NamespacesToResponse(namespaces)
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}

// scaleDeploymentHandler handles PUT /clusters/{context}/namespaces/{namespace}/deployments/{name}/scale
func (h *HTTPHandler) scaleDeploymentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	context := vars["context"]
	namespace := vars["namespace"]
	deploymentName := vars["name"]

	if context == "" || namespace == "" || deploymentName == "" {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Context, namespace, and deployment name are required")
		return
	}

	var req k8sWire.ScaleDeploymentRequest
	if err := h.requestAdapter.ParseJSON(r, &req); err != nil {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	err := h.controller.ScaleDeployment(r.Context(), context, namespace, deploymentName, req.Replicas)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, "Failed to scale deployment: "+err.Error())
		return
	}

	response := k8sWire.OperationResponse{
		Success:   true,
		Message:   "Deployment scaled successfully",
		Operation: "scale",
		Resource:  namespace + "/" + deploymentName,
	}
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}

// restartDeploymentHandler handles POST /clusters/{context}/namespaces/{namespace}/deployments/{name}/restart
func (h *HTTPHandler) restartDeploymentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	context := vars["context"]
	namespace := vars["namespace"]
	deploymentName := vars["name"]

	if context == "" || namespace == "" || deploymentName == "" {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Context, namespace, and deployment name are required")
		return
	}

	err := h.controller.RestartDeployment(r.Context(), context, namespace, deploymentName)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, "Failed to restart deployment: "+err.Error())
		return
	}

	response := k8sWire.OperationResponse{
		Success:   true,
		Message:   "Deployment restart initiated",
		Operation: "restart",
		Resource:  namespace + "/" + deploymentName,
	}
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}

// parseDeploymentFilter parses query parameters into DeploymentFilter
func (h *HTTPHandler) parseDeploymentFilter(r *http.Request) *k8sModels.DeploymentFilter {
	query := r.URL.Query()

	filter := &k8sModels.DeploymentFilter{
		Namespace:     query.Get("namespace"),
		LabelSelector: query.Get("label_selector"),
		ServiceName:   query.Get("service_name"),
		Status:        query.Get("status"),
	}

	if limit := query.Get("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 {
			filter.Limit = l
		}
	}

	return filter
}

// parsePodFilter parses query parameters into PodFilter
func (h *HTTPHandler) parsePodFilter(r *http.Request) *k8sModels.PodFilter {
	query := r.URL.Query()

	filter := &k8sModels.PodFilter{
		Namespace:     query.Get("namespace"),
		LabelSelector: query.Get("label_selector"),
		FieldSelector: query.Get("field_selector"),
		PodName:       query.Get("pod_name"),
		ContainerName: query.Get("container_name"),
		Status:        query.Get("status"),
		Node:          query.Get("node"),
	}

	if limit := query.Get("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 {
			filter.Limit = l
		}
	}

	return filter
}
