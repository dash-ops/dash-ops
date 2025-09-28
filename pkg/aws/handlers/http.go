package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	awsAdapters "github.com/dash-ops/dash-ops/pkg/aws/adapters/http"
	aws "github.com/dash-ops/dash-ops/pkg/aws/controllers"
	awsModels "github.com/dash-ops/dash-ops/pkg/aws/models"
	awsPorts "github.com/dash-ops/dash-ops/pkg/aws/ports"
	awsRepositories "github.com/dash-ops/dash-ops/pkg/aws/repositories"
	awsWire "github.com/dash-ops/dash-ops/pkg/aws/wire"
	commonsHttp "github.com/dash-ops/dash-ops/pkg/commons/adapters/http"
)

// HTTPHandler handles HTTP requests for AWS module
type HTTPHandler struct {
	accountsController  *aws.AccountsController
	instancesController *aws.InstancesController
	awsAdapter          *awsAdapters.AWSAdapter
	responseAdapter     *commonsHttp.ResponseAdapter
	requestAdapter      *commonsHttp.RequestAdapter
}

// NewHTTPHandler creates a new HTTP handler with DI
func NewHTTPHandler(
	awsClientService awsPorts.AWSClientService,
	accounts []awsModels.AWSAccount,
	responseAdapter *commonsHttp.ResponseAdapter,
	requestAdapter *commonsHttp.RequestAdapter,
) *HTTPHandler {
	// Create repositories
	instanceRepo := awsRepositories.NewInstanceRepository(awsClientService)

	// Create controllers
	accountsController := aws.NewAccountsController(accounts)
	instancesController := aws.NewInstancesController(instanceRepo, accounts)

	// Create HTTP adapter
	awsAdapter := awsAdapters.NewAWSAdapter()

	return &HTTPHandler{
		accountsController:  accountsController,
		instancesController: instancesController,
		awsAdapter:          awsAdapter,
		responseAdapter:     responseAdapter,
		requestAdapter:      requestAdapter,
	}
}

// RegisterRoutes registers all AWS routes
func (h *HTTPHandler) RegisterRoutes(router *mux.Router) {
	// Account operations
	router.HandleFunc("/aws/accounts", h.listAccountsHandler).Methods("GET")
	router.HandleFunc("/aws/{account}/permissions", h.getPermissionsHandler).Methods("GET")

	// Instance operations - matching frontend expectations
	router.HandleFunc("/aws/{account}/ec2/instances", h.listInstancesHandler).Methods("GET")
	router.HandleFunc("/aws/{account}/ec2/instance/start/{instanceId}", h.startInstanceHandler).Methods("POST")
	router.HandleFunc("/aws/{account}/ec2/instance/stop/{instanceId}", h.stopInstanceHandler).Methods("POST")

	// Additional operations (for future use)
	router.HandleFunc("/aws/{account}/regions/{region}/instances", h.listInstancesByRegionHandler).Methods("GET")
	router.HandleFunc("/aws/{account}/regions/{region}/instances/{id}", h.getInstanceHandler).Methods("GET")
	router.HandleFunc("/aws/{account}/regions/{region}/instances/{id}/metrics", h.getInstanceMetricsHandler).Methods("GET")
	router.HandleFunc("/aws/{account}/regions/{region}/instances/{id}/cost", h.getInstanceCostEstimateHandler).Methods("GET")
	router.HandleFunc("/aws/{account}/regions/{region}/instances/batch", h.batchOperationHandler).Methods("POST")
	router.HandleFunc("/aws/{account}/regions/{region}/cost/savings", h.getCostSavingsHandler).Methods("GET")
	router.HandleFunc("/aws/{account}/summary", h.getAccountSummaryHandler).Methods("GET")
}

// listAccountsHandler handles GET /accounts
func (h *HTTPHandler) listAccountsHandler(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.accountsController.ListAccounts(r.Context())
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, "Failed to list accounts: "+err.Error())
		return
	}

	response := h.awsAdapter.AccountsToResponse(accounts)
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}

// getAccountHandler handles GET /accounts/{account}
func (h *HTTPHandler) getAccountHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountKey := vars["account"]

	if accountKey == "" {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Account key is required")
		return
	}

	account, err := h.accountsController.GetAccount(r.Context(), accountKey)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusNotFound, "Account not found: "+err.Error())
		return
	}

	response := h.awsAdapter.AccountToResponse(account)
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}

// getPermissionsHandler handles GET /aws/{account}/permissions
func (h *HTTPHandler) getPermissionsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountKey := vars["account"]

	if accountKey == "" {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Account key is required")
		return
	}

	account, err := h.accountsController.GetAccount(r.Context(), accountKey)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusNotFound, "Account not found: "+err.Error())
		return
	}

	// Return permissions in format expected by frontend
	response := map[string]interface{}{
		"ec2": map[string]interface{}{
			"start": account.Permissions.EC2.Start,
			"stop":  account.Permissions.EC2.Stop,
			"view":  account.Permissions.EC2.View,
		},
	}

	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}

// listInstancesHandler handles GET /aws/{account}/ec2/instances
func (h *HTTPHandler) listInstancesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountKey := vars["account"]

	if accountKey == "" {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Account is required")
		return
	}

	// Get region from query parameter (default to us-east-1)
	region := r.URL.Query().Get("region")
	if region == "" {
		region = "us-east-1"
	}

	// Parse query parameters
	filter := h.parseInstanceFilter(r)

	// List instances
	instanceList, err := h.instancesController.ListInstances(r.Context(), accountKey, region, filter)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, "Failed to list instances: "+err.Error())
		return
	}

	response := h.awsAdapter.InstanceListToResponse(instanceList)
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}

// listInstancesByRegionHandler handles GET /aws/{account}/regions/{region}/instances
func (h *HTTPHandler) listInstancesByRegionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountKey := vars["account"]
	region := vars["region"]

	if accountKey == "" || region == "" {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Account and region are required")
		return
	}

	// Parse query parameters
	filter := h.parseInstanceFilter(r)

	// List instances
	instanceList, err := h.instancesController.ListInstances(r.Context(), accountKey, region, filter)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, "Failed to list instances: "+err.Error())
		return
	}

	response := h.awsAdapter.InstanceListToResponse(instanceList)
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}

// startInstanceHandler handles POST /aws/{account}/ec2/instance/start/{instanceId}
func (h *HTTPHandler) startInstanceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountKey := vars["account"]
	instanceID := vars["instanceId"]

	if accountKey == "" || instanceID == "" {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Account and instance ID are required")
		return
	}

	// Get region from query parameter (default to us-east-1)
	region := r.URL.Query().Get("region")
	if region == "" {
		region = "us-east-1"
	}

	// Start instance
	operation, err := h.instancesController.StartInstance(r.Context(), accountKey, region, instanceID)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, "Failed to start instance: "+err.Error())
		return
	}

	// Return response in format expected by frontend
	response := map[string]string{
		"current_state": operation.CurrentState.Name,
	}
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}

// stopInstanceHandler handles POST /aws/{account}/ec2/instance/stop/{instanceId}
func (h *HTTPHandler) stopInstanceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountKey := vars["account"]
	instanceID := vars["instanceId"]

	if accountKey == "" || instanceID == "" {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Account and instance ID are required")
		return
	}

	// Get region from query parameter (default to us-east-1)
	region := r.URL.Query().Get("region")
	if region == "" {
		region = "us-east-1"
	}

	// Stop instance
	operation, err := h.instancesController.StopInstance(r.Context(), accountKey, region, instanceID)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, "Failed to stop instance: "+err.Error())
		return
	}

	// Return response in format expected by frontend
	response := map[string]string{
		"current_state": operation.CurrentState.Name,
	}
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}

// batchOperationHandler handles POST /accounts/{account}/regions/{region}/instances/batch
func (h *HTTPHandler) batchOperationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountKey := vars["account"]
	region := vars["region"]

	if accountKey == "" || region == "" {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Account and region are required")
		return
	}

	// Parse request
	var req awsWire.BatchOperationRequest
	if err := h.requestAdapter.ParseJSON(r, &req); err != nil {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	// Execute batch operation
	batchOp, err := h.instancesController.BatchOperation(r.Context(), accountKey, region, req.Operation, req.InstanceIDs)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, "Failed to execute batch operation: "+err.Error())
		return
	}

	response := h.awsAdapter.BatchOperationToResponse(batchOp)
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}

// getCostSavingsHandler handles GET /accounts/{account}/regions/{region}/cost/savings
func (h *HTTPHandler) getCostSavingsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountKey := vars["account"]
	region := vars["region"]

	if accountKey == "" || region == "" {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Account and region are required")
		return
	}

	// TODO: Implement cost savings functionality
	h.responseAdapter.WriteError(w, http.StatusNotImplemented, "Cost savings functionality not implemented yet")
}

// parseInstanceFilter parses query parameters into InstanceFilter
func (h *HTTPHandler) parseInstanceFilter(r *http.Request) *awsModels.InstanceFilter {
	query := r.URL.Query()

	filter := &awsModels.InstanceFilter{
		State:        query.Get("state"),
		InstanceType: query.Get("instance_type"),
		Platform:     query.Get("platform"),
		Search:       query.Get("search"),
	}

	// Parse limit and offset
	if limit := query.Get("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 {
			filter.Limit = l
		}
	}

	if offset := query.Get("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil && o >= 0 {
			filter.Offset = o
		}
	}

	// Parse tags (format: tag.key=value)
	for key, values := range query {
		if strings.HasPrefix(key, "tag.") && len(values) > 0 {
			tagKey := strings.TrimPrefix(key, "tag.")
			filter.Tags = append(filter.Tags, awsModels.TagFilter{
				Key:   tagKey,
				Value: values[0],
			})
		}
	}

	return filter
}

// getInstanceHandler handles GET /accounts/{account}/regions/{region}/instances/{id}
func (h *HTTPHandler) getInstanceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountKey := vars["account"]
	region := vars["region"]
	instanceID := vars["id"]

	if accountKey == "" || region == "" || instanceID == "" {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Account, region, and instance ID are required")
		return
	}

	// Get instance
	instance, err := h.instancesController.GetInstance(r.Context(), accountKey, region, instanceID)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusNotFound, "Instance not found: "+err.Error())
		return
	}

	response := h.awsAdapter.InstanceToResponse(instance)
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}

// getAccountSummaryHandler handles GET /accounts/{account}/summary
func (h *HTTPHandler) getAccountSummaryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountKey := vars["account"]

	if accountKey == "" {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Account key is required")
		return
	}

	// Get region from query parameter (default to us-east-1)
	region := r.URL.Query().Get("region")
	if region == "" {
		region = "us-east-1"
	}

	// TODO: Implement account summary functionality
	h.responseAdapter.WriteError(w, http.StatusNotImplemented, "Account summary functionality not implemented yet")
}

// getInstanceMetricsHandler handles GET /accounts/{account}/regions/{region}/instances/{id}/metrics
func (h *HTTPHandler) getInstanceMetricsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountKey := vars["account"]
	region := vars["region"]
	instanceID := vars["id"]

	if accountKey == "" || region == "" || instanceID == "" {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Account, region, and instance ID are required")
		return
	}

	// Get period from query parameter (default to 1h)
	period := r.URL.Query().Get("period")
	if period == "" {
		period = "1h"
	}

	// TODO: Implement instance metrics functionality
	h.responseAdapter.WriteError(w, http.StatusNotImplemented, "Instance metrics functionality not implemented yet")
}

// getInstanceCostEstimateHandler handles GET /accounts/{account}/regions/{region}/instances/{id}/cost
func (h *HTTPHandler) getInstanceCostEstimateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountKey := vars["account"]
	region := vars["region"]
	instanceID := vars["id"]

	if accountKey == "" || region == "" || instanceID == "" {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Account, region, and instance ID are required")
		return
	}

	// Get operation from query parameter (default to current state)
	operation := r.URL.Query().Get("operation")
	if operation == "" {
		operation = "current"
	}

	// TODO: Implement cost estimate functionality
	h.responseAdapter.WriteError(w, http.StatusNotImplemented, "Cost estimate functionality not implemented yet")
}

// getUserContext extracts user context from request
func (h *HTTPHandler) getUserContext(r *http.Request) *awsPorts.UserContext {
	// TODO: Implement OAuth2 integration to get user context
	// For now, return a default user context for testing
	return &awsPorts.UserContext{
		Username: "test-user",
		Name:     "Test User",
		Email:    "test@example.com",
		Groups:   []string{"admin"}, // Default to admin for testing
	}
}
