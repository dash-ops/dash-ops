package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	awsAdapters "github.com/dash-ops/dash-ops/pkg/aws-new/adapters/http"
	aws "github.com/dash-ops/dash-ops/pkg/aws-new/controllers"
	awsModels "github.com/dash-ops/dash-ops/pkg/aws-new/models"
	awsPorts "github.com/dash-ops/dash-ops/pkg/aws-new/ports"
	awsWire "github.com/dash-ops/dash-ops/pkg/aws-new/wire"
	commonsHttp "github.com/dash-ops/dash-ops/pkg/commons-new/adapters/http"
)

// HTTPHandler handles HTTP requests for AWS module
type HTTPHandler struct {
	controller      *aws.AWSController
	awsAdapter      *awsAdapters.AWSAdapter
	responseAdapter *commonsHttp.ResponseAdapter
	requestAdapter  *commonsHttp.RequestAdapter
}

// NewHTTPHandler creates a new HTTP handler
func NewHTTPHandler(
	controller *aws.AWSController,
	awsAdapter *awsAdapters.AWSAdapter,
	responseAdapter *commonsHttp.ResponseAdapter,
	requestAdapter *commonsHttp.RequestAdapter,
) *HTTPHandler {
	return &HTTPHandler{
		controller:      controller,
		awsAdapter:      awsAdapter,
		responseAdapter: responseAdapter,
		requestAdapter:  requestAdapter,
	}
}

// RegisterRoutes registers all AWS routes
func (h *HTTPHandler) RegisterRoutes(router *mux.Router) {
	// Account operations
	router.HandleFunc("/accounts", h.listAccountsHandler).Methods("GET")
	router.HandleFunc("/accounts/{account}", h.getAccountHandler).Methods("GET")
	router.HandleFunc("/accounts/{account}/summary", h.getAccountSummaryHandler).Methods("GET")

	// Instance operations
	router.HandleFunc("/accounts/{account}/regions/{region}/instances", h.listInstancesHandler).Methods("GET")
	router.HandleFunc("/accounts/{account}/regions/{region}/instances/{id}", h.getInstanceHandler).Methods("GET")
	router.HandleFunc("/accounts/{account}/regions/{region}/instances/{id}/start", h.startInstanceHandler).Methods("POST")
	router.HandleFunc("/accounts/{account}/regions/{region}/instances/{id}/stop", h.stopInstanceHandler).Methods("POST")
	router.HandleFunc("/accounts/{account}/regions/{region}/instances/{id}/metrics", h.getInstanceMetricsHandler).Methods("GET")

	// Batch operations
	router.HandleFunc("/accounts/{account}/regions/{region}/instances/batch", h.batchOperationHandler).Methods("POST")

	// Cost analysis
	router.HandleFunc("/accounts/{account}/regions/{region}/cost/savings", h.getCostSavingsHandler).Methods("GET")
	router.HandleFunc("/accounts/{account}/regions/{region}/instances/{id}/cost", h.getInstanceCostEstimateHandler).Methods("GET")
}

// listAccountsHandler handles GET /accounts
func (h *HTTPHandler) listAccountsHandler(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.controller.ListAccounts(r.Context())
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

	account, err := h.controller.GetAccount(r.Context(), accountKey)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusNotFound, "Account not found: "+err.Error())
		return
	}

	response := h.awsAdapter.AccountToResponse(account)
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}

// listInstancesHandler handles GET /accounts/{account}/regions/{region}/instances
func (h *HTTPHandler) listInstancesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountKey := vars["account"]
	region := vars["region"]

	if accountKey == "" || region == "" {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Account and region are required")
		return
	}

	// Parse query parameters
	filter := h.parseInstanceFilter(r)

	// Get user context
	userContext := h.getUserContext(r)

	// List instances
	instanceList, err := h.controller.ListInstances(r.Context(), accountKey, region, filter, userContext)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, "Failed to list instances: "+err.Error())
		return
	}

	response := h.awsAdapter.InstanceListToResponse(instanceList)
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}

// startInstanceHandler handles POST /accounts/{account}/regions/{region}/instances/{id}/start
func (h *HTTPHandler) startInstanceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountKey := vars["account"]
	region := vars["region"]
	instanceID := vars["id"]

	if accountKey == "" || region == "" || instanceID == "" {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Account, region, and instance ID are required")
		return
	}

	// Get user context
	userContext := h.getUserContext(r)

	// Start instance
	operation, err := h.controller.StartInstance(r.Context(), accountKey, region, instanceID, userContext)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, "Failed to start instance: "+err.Error())
		return
	}

	response := h.awsAdapter.InstanceOperationToResponse(operation)
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}

// stopInstanceHandler handles POST /accounts/{account}/regions/{region}/instances/{id}/stop
func (h *HTTPHandler) stopInstanceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountKey := vars["account"]
	region := vars["region"]
	instanceID := vars["id"]

	if accountKey == "" || region == "" || instanceID == "" {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Account, region, and instance ID are required")
		return
	}

	// Get user context
	userContext := h.getUserContext(r)

	// Stop instance
	operation, err := h.controller.StopInstance(r.Context(), accountKey, region, instanceID, userContext)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, "Failed to stop instance: "+err.Error())
		return
	}

	response := h.awsAdapter.InstanceOperationToResponse(operation)
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

	// Get user context
	userContext := h.getUserContext(r)

	// Execute batch operation
	batchOp, err := h.controller.BatchOperation(r.Context(), accountKey, region, req.Operation, req.InstanceIDs, userContext)
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

	// Get user context
	userContext := h.getUserContext(r)

	// Get cost savings
	savings, err := h.controller.GetCostSavings(r.Context(), accountKey, region, userContext)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, "Failed to get cost savings: "+err.Error())
		return
	}

	response := h.awsAdapter.CostSavingsToResponse(savings)
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
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
