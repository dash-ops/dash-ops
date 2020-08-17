package commons

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type ResponseError struct {
	Error string `json:"error"`
}

// RespondJSON makes the response with payload as json format
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

// RespondError makes the error response with payload as json format
func RespondError(w http.ResponseWriter, code int, message string) {
	RespondJSON(w, code, ResponseError{Error: message})
}

// HasPermission ...
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

// UnderScoreString ...
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
