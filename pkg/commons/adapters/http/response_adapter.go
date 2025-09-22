package http

import (
	"encoding/json"
	"log"
	"net/http"

	commonsWire "github.com/dash-ops/dash-ops/pkg/commons/wire"
)

// ResponseAdapter handles HTTP response formatting and writing
type ResponseAdapter struct{}

// NewResponseAdapter creates a new response adapter
func NewResponseAdapter() *ResponseAdapter {
	return &ResponseAdapter{}
}

// WriteJSON writes a JSON response with the given status code and payload
func (r *ResponseAdapter) WriteJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		r.WriteError(w, http.StatusInternalServerError, "Failed to marshal response")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if _, err := w.Write(response); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

// WriteError writes an error response with the given status code and message
func (r *ResponseAdapter) WriteError(w http.ResponseWriter, code int, message string) {
	errorResponse := commonsWire.ErrorResponse{
		Error: message,
	}
	r.WriteJSON(w, code, errorResponse)
}

// WriteSuccess writes a success response with optional data
func (r *ResponseAdapter) WriteSuccess(w http.ResponseWriter, message string, data interface{}) {
	successResponse := commonsWire.SuccessResponse{
		Message: message,
		Data:    data,
	}
	r.WriteJSON(w, http.StatusOK, successResponse)
}

// WriteNoContent writes a 204 No Content response
func (r *ResponseAdapter) WriteNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// WriteCreated writes a 201 Created response with optional location header
func (r *ResponseAdapter) WriteCreated(w http.ResponseWriter, location string, payload interface{}) {
	if location != "" {
		w.Header().Set("Location", location)
	}
	r.WriteJSON(w, http.StatusCreated, payload)
}
