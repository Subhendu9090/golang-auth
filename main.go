package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/subhendu/go-auth/routes"
)

func main() {
	r := mux.NewRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	routes.AuthRoutes(r)
	routes.UserRoutes(r)

	r.HandleFunc("/api-1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Access get from api-1")
	})

	r.HandleFunc("/api-2", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Access get from api-2")
	})

	http.ListenAndServe(":"+port, r)
}
