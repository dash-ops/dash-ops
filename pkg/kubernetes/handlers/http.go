package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	commonsHttp "github.com/dash-ops/dash-ops/pkg/commons/adapters/http"
	k8sAdapters "github.com/dash-ops/dash-ops/pkg/kubernetes/adapters/http"
	"github.com/dash-ops/dash-ops/pkg/kubernetes/controllers"
	"github.com/dash-ops/dash-ops/pkg/kubernetes/integrations/external/kubernetes"
	k8sModels "github.com/dash-ops/dash-ops/pkg/kubernetes/models"
	k8sPorts "github.com/dash-ops/dash-ops/pkg/kubernetes/ports"
	"github.com/dash-ops/dash-ops/pkg/kubernetes/repositories"
	k8sWire "github.com/dash-ops/dash-ops/pkg/kubernetes/wire"
)

// HTTPHandler handles HTTP requests for Kubernetes module
type HTTPHandler struct {
	nodesController       *controllers.NodesController
	deploymentsController *controllers.DeploymentsController
	podsController        *controllers.PodsController
	namespacesController  *controllers.NamespacesController
	k8sClient             *kubernetes.KubernetesClient
	responseAdapter       *commonsHttp.ResponseAdapter
	requestAdapter        *commonsHttp.RequestAdapter
}

// NewHTTPHandler creates a new HTTP handler
func NewHTTPHandler(
	k8sClient *kubernetes.KubernetesClient,
	responseAdapter *commonsHttp.ResponseAdapter,
	requestAdapter *commonsHttp.RequestAdapter,
) *HTTPHandler {
	// Initialize repositories with the provided client
	nodesRepo := repositories.NewNodesRepository(k8sClient)
	deploymentsRepo := repositories.NewDeploymentsRepository(k8sClient)
	podsRepo := repositories.NewPodsRepository(k8sClient)
	namespacesRepo := repositories.NewNamespacesRepository(k8sClient)

	// Initialize controllers with repositories
	nodesController := controllers.NewNodesController(nodesRepo)
	deploymentsController := controllers.NewDeploymentsController(deploymentsRepo)
	podsController := controllers.NewPodsController(podsRepo)
	namespacesController := controllers.NewNamespacesController(namespacesRepo)

	return &HTTPHandler{
		nodesController:       nodesController,
		deploymentsController: deploymentsController,
		podsController:        podsController,
		namespacesController:  namespacesController,
		k8sClient:             k8sClient,
		responseAdapter:       responseAdapter,
		requestAdapter:        requestAdapter,
	}
}

// SetServiceContextResolver sets the service context resolver for the deployments controller
func (h *HTTPHandler) SetServiceContextResolver(resolver k8sPorts.ServiceContextResolver) {
	if h.deploymentsController != nil {
		h.deploymentsController.SetServiceContextResolver(resolver)
	}
}

// RegisterRoutes registers all Kubernetes routes
func (h *HTTPHandler) RegisterRoutes(router *mux.Router) {
	// Cluster operations
	router.HandleFunc("/clusters", h.listClustersHandler).Methods("GET")
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

	// TODO: Implement cluster info aggregation from multiple controllers
	// For now, return a basic cluster info structure
	clusterInfo := &k8sModels.ClusterInfo{
		Cluster: k8sModels.Cluster{
			Name:    context,
			Context: context,
			Status:  k8sModels.ClusterStatusConnected,
		},
	}
	// no error path for the stubbed response above

	response := k8sAdapters.ClusterInfoToResponse(clusterInfo)
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

	nodes, err := h.nodesController.ListNodes(r.Context(), context)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, "Failed to list nodes: "+err.Error())
		return
	}

	response := k8sAdapters.NodesToResponse(nodes)
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

	node, err := h.nodesController.GetNode(r.Context(), context, nodeName)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusNotFound, "Node not found: "+err.Error())
		return
	}

	response := k8sAdapters.NodeToResponse(node)
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

	namespaces, err := h.namespacesController.ListNamespaces(r.Context(), context)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, "Failed to list namespaces: "+err.Error())
		return
	}

	response := k8sAdapters.NamespacesToResponse(namespaces)
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

	err := h.deploymentsController.ScaleDeployment(r.Context(), context, namespace, deploymentName, req.Replicas)
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

	err := h.deploymentsController.RestartDeployment(r.Context(), context, namespace, deploymentName)
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

// getClusterHealthHandler handles GET /clusters/{context}/health
func (h *HTTPHandler) getClusterHealthHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	context := vars["context"]

	if context == "" {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Context is required")
		return
	}

	// TODO: Implement cluster health aggregation from multiple controllers
	// For now, return a basic health structure
	health := &k8sModels.ClusterHealth{
		Context:     context,
		Status:      k8sModels.HealthStatus("healthy"),
		LastUpdated: time.Now(),
	}
	// no error path for the stubbed response above

	response := k8sAdapters.ClusterHealthToResponse(health)
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}

// createNamespaceHandler handles POST /clusters/{context}/namespaces
func (h *HTTPHandler) createNamespaceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	context := vars["context"]

	if context == "" {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Context is required")
		return
	}

	var req k8sWire.CreateNamespaceRequest
	if err := h.requestAdapter.ParseJSON(r, &req); err != nil {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	namespace, err := h.namespacesController.CreateNamespace(r.Context(), context, req.Name)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, "Failed to create namespace: "+err.Error())
		return
	}

	response := k8sAdapters.NamespaceToResponse(namespace)
	h.responseAdapter.WriteCreated(w, "/clusters/"+context+"/namespaces/"+namespace.Name, response)
}

// getNamespaceHandler handles GET /clusters/{context}/namespaces/{name}
func (h *HTTPHandler) getNamespaceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	context := vars["context"]
	name := vars["name"]

	if context == "" || name == "" {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Context and namespace name are required")
		return
	}

	namespace, err := h.namespacesController.GetNamespace(r.Context(), context, name)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusNotFound, "Namespace not found: "+err.Error())
		return
	}

	response := k8sAdapters.NamespaceToResponse(namespace)
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}

// deleteNamespaceHandler handles DELETE /clusters/{context}/namespaces/{name}
func (h *HTTPHandler) deleteNamespaceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	context := vars["context"]
	name := vars["name"]

	if context == "" || name == "" {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Context and namespace name are required")
		return
	}

	err := h.namespacesController.DeleteNamespace(r.Context(), context, name)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, "Failed to delete namespace: "+err.Error())
		return
	}

	response := k8sWire.OperationResponse{
		Success:   true,
		Message:   "Namespace deleted successfully",
		Operation: "delete",
		Resource:  name,
	}
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}

// listDeploymentsHandler handles GET /clusters/{context}/deployments
func (h *HTTPHandler) listDeploymentsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	context := vars["context"]

	if context == "" {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Context is required")
		return
	}

	filter := h.parseDeploymentFilter(r)
	deployments, err := h.deploymentsController.ListDeployments(r.Context(), context, filter)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, "Failed to list deployments: "+err.Error())
		return
	}

	response := k8sAdapters.DeploymentsToResponse(deployments.Deployments)
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}

// listDeploymentsByNamespaceHandler handles GET /clusters/{context}/namespaces/{namespace}/deployments
func (h *HTTPHandler) listDeploymentsByNamespaceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	context := vars["context"]
	namespace := vars["namespace"]

	if context == "" || namespace == "" {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Context and namespace are required")
		return
	}

	filter := h.parseDeploymentFilter(r)
	filter.Namespace = namespace

	deployments, err := h.deploymentsController.ListDeployments(r.Context(), context, filter)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, "Failed to list deployments: "+err.Error())
		return
	}

	response := k8sAdapters.DeploymentsToResponse(deployments.Deployments)
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}

// getDeploymentHandler handles GET /clusters/{context}/namespaces/{namespace}/deployments/{name}
func (h *HTTPHandler) getDeploymentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	context := vars["context"]
	namespace := vars["namespace"]
	name := vars["name"]

	if context == "" || namespace == "" || name == "" {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Context, namespace, and deployment name are required")
		return
	}

	deployment, err := h.deploymentsController.GetDeployment(r.Context(), context, namespace, name)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusNotFound, "Deployment not found: "+err.Error())
		return
	}

	response := k8sAdapters.DeploymentToResponse(deployment)
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}

// listPodsHandler handles GET /clusters/{context}/pods
func (h *HTTPHandler) listPodsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	context := vars["context"]

	if context == "" {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Context is required")
		return
	}

	filter := h.parsePodFilter(r)
	pods, err := h.podsController.ListPods(r.Context(), context, filter)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, "Failed to list pods: "+err.Error())
		return
	}

	response := k8sAdapters.PodsToResponse(pods.Pods)
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}

// listPodsByNamespaceHandler handles GET /clusters/{context}/namespaces/{namespace}/pods
func (h *HTTPHandler) listPodsByNamespaceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	context := vars["context"]
	namespace := vars["namespace"]

	if context == "" || namespace == "" {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Context and namespace are required")
		return
	}

	filter := h.parsePodFilter(r)
	filter.Namespace = namespace

	pods, err := h.podsController.ListPods(r.Context(), context, filter)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, "Failed to list pods: "+err.Error())
		return
	}

	response := k8sAdapters.PodsToResponse(pods.Pods)
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}

// getPodHandler handles GET /clusters/{context}/namespaces/{namespace}/pods/{name}
func (h *HTTPHandler) getPodHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	context := vars["context"]
	namespace := vars["namespace"]
	name := vars["name"]

	if context == "" || namespace == "" || name == "" {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Context, namespace, and pod name are required")
		return
	}

	pod, err := h.podsController.GetPod(r.Context(), context, namespace, name)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusNotFound, "Pod not found: "+err.Error())
		return
	}

	response := k8sAdapters.PodToResponse(pod)
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}

// deletePodHandler handles DELETE /clusters/{context}/namespaces/{namespace}/pods/{name}
func (h *HTTPHandler) deletePodHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	context := vars["context"]
	namespace := vars["namespace"]
	name := vars["name"]

	if context == "" || namespace == "" || name == "" {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Context, namespace, and pod name are required")
		return
	}

	err := h.podsController.DeletePod(r.Context(), context, namespace, name)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, "Failed to delete pod: "+err.Error())
		return
	}

	response := k8sWire.OperationResponse{
		Success:   true,
		Message:   "Pod deleted successfully",
		Operation: "delete",
		Resource:  namespace + "/" + name,
	}
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}

// getPodLogsHandler handles GET /clusters/{context}/namespaces/{namespace}/pods/{name}/logs
func (h *HTTPHandler) getPodLogsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	context := vars["context"]
	namespace := vars["namespace"]
	name := vars["name"]

	if context == "" || namespace == "" || name == "" {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Context, namespace, and pod name are required")
		return
	}

	query := r.URL.Query()
	container := query.Get("container")
	tailLines := query.Get("tailLines")

	lines := int64(100) // default
	if tailLines != "" {
		if l, err := strconv.ParseInt(tailLines, 10, 64); err == nil && l > 0 {
			lines = l
		}
	}

	logFilter := &k8sModels.LogFilter{
		PodName:       name,
		Namespace:     namespace,
		ContainerName: container,
		TailLines:     lines,
	}

	logs, err := h.podsController.GetPodLogs(r.Context(), context, logFilter)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, "Failed to get pod logs: "+err.Error())
		return
	}

	// Convert logs to response format
	logResponses := make([]k8sWire.ContainerLogResponse, len(logs))
	for i, log := range logs {
		logResponses[i] = k8sWire.ContainerLogResponse{
			Timestamp: log.Timestamp,
			Message:   log.Message,
			Level:     log.Level,
		}
	}

	response := k8sWire.PodLogsResponse{
		PodName:       name,
		Namespace:     namespace,
		ContainerName: container,
		Logs:          logResponses,
		TotalLines:    len(logResponses),
	}
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}

// listClustersHandler handles GET /clusters
func (h *HTTPHandler) listClustersHandler(w http.ResponseWriter, r *http.Request) {
	// Build cluster list from injected client (single configured context for now)
	version, err := h.k8sClient.GetServerVersion(r.Context())
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, "Failed to get server version: "+err.Error())
		return
	}

	currentContext := h.k8sClient.GetContext()

	clusters := []k8sModels.Cluster{
		{
			Name:    currentContext,
			Context: currentContext,
			Version: version,
			Status:  k8sModels.ClusterStatusConnected,
		},
	}

	response := k8sAdapters.ClusterListToResponse(clusters)
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}
