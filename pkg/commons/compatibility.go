package commons

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"

	commonsModels "github.com/dash-ops/dash-ops/pkg/commons/models"
)

// Legacy compatibility functions for existing code

// ResponseError legacy struct for compatibility
type ResponseError struct {
	Error string `json:"error"`
}

// RespondJSON makes the response with payload as json format - legacy compatibility
func RespondJSON(w http.ResponseWriter, code int, payload interface{}) {
	r, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			log.Fatalln("write failed", err)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write([]byte(r))
	if err != nil {
		log.Fatalln("write failed", err)
	}
}

// RespondError makes the error response with payload as json format - legacy compatibility
func RespondError(w http.ResponseWriter, code int, message string) {
	RespondJSON(w, code, ResponseError{Error: message})
}

// HasPermission checks permissions - legacy compatibility
func HasPermission(featurePermissions []string, groupsPermission []string) bool {
	isValid := false

	for i := 0; i < len(featurePermissions); i++ {
		for _, gP := range groupsPermission {
			if strings.ToLower(featurePermissions[i]) == strings.ToLower(gP) {
				isValid = true
			}
		}
	}

	return isValid
}

// UnderScoreString converts string to underscore format - legacy compatibility
func UnderScoreString(str string) string {
	// convert every letter to lower case
	newStr := strings.ToLower(str)

	// convert all spaces/tab to underscore
	regExp := regexp.MustCompile("[[:space:][:blank:]]")
	newStrByte := regExp.ReplaceAll([]byte(newStr), []byte("_"))

	regExp = regexp.MustCompile("`[^a-z0-9]`i")
	newStrByte = regExp.ReplaceAll(newStrByte, []byte("_"))

	regExp = regexp.MustCompile("[!/']")
	newStrByte = regExp.ReplaceAll(newStrByte, []byte("_"))

	// and remove underscore from beginning and ending
	newStr = strings.TrimPrefix(string(newStrByte), "_")
	newStr = strings.TrimSuffix(newStr, "_")

	return newStr
}

// Legacy type aliases for compatibility
type key = commonsModels.ContextKey

// Legacy constants for compatibility
const (
	TokenKey    = commonsModels.TokenKey
	UserDataKey = commonsModels.UserDataKey
)

// UserData legacy type alias for compatibility
type UserData = commonsModels.UserData
