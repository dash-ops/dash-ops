package commons

import (
	"encoding/json"
	"net/http"

	"github.com/apex/log"
)

// RespondJSON makes the response with payload as json format
func RespondJSON(w http.ResponseWriter, code int, payload interface{}) {
	r, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, errw := w.Write([]byte(err.Error()))
		if errw != nil {
			log.WithError(err).Fatal("write failed")
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, errw := w.Write([]byte(r))
	if errw != nil {
		log.WithError(err).Fatal("write failed")
	}
}

// RespondError makes the error response with payload as json format
func RespondError(w http.ResponseWriter, code int, message string) {
	RespondJSON(w, code, map[string]string{"error": message})
}
