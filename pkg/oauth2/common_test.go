package oauth2

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRespondJSON(t *testing.T) {
	var body interface{}
	res := httptest.NewRecorder()
	respondJSON(res, http.StatusOK, body)

	assert.Equal(t, http.StatusOK, res.Code)
}

func TestRespondError(t *testing.T) {
	res := httptest.NewRecorder()
	respondError(res, http.StatusInternalServerError, "Error")

	assert.Equal(t, http.StatusInternalServerError, res.Code)
}
