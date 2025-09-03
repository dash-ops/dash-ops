package commons

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRespondJSON(t *testing.T) {
	var body interface{}
	res := httptest.NewRecorder()
	RespondJSON(res, http.StatusOK, body)

	assert.Equal(t, http.StatusOK, res.Code)
}

func TestRespondError(t *testing.T) {
	res := httptest.NewRecorder()
	RespondError(res, http.StatusInternalServerError, "Error")

	assert.Equal(t, http.StatusInternalServerError, res.Code)
}

func TestHasPermission(t *testing.T) {
	isOK := HasPermission([]string{"OrgGH:Team1", "OrgGH:Team2"}, []string{"OrgGH:Team2"})
	isNotPermission := HasPermission([]string{"OrgGH:Team1", "OrgGH:Team2"}, []string{"OrgGH:Team3"})

	assert.True(t, isOK)
	assert.False(t, isNotPermission)
}

func TestUnderScoreString(t *testing.T) {
	myString := "Xpto in My Test"

	result := UnderScoreString(myString)

	assert.Equal(t, "xpto_in_my_test", result)
}
