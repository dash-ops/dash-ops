package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/dash-ops/dash-ops/pkg/spa"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	spaHandler := spa.SpaHandler{StaticPath: "front/build", IndexPath: "index.html"}
	router.PathPrefix("/").Handler(spaHandler)

	srv := &http.Server{
		Handler: router,
		Addr:    "localhost:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
