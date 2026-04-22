package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/RafayKhattak/aegis-iam-backend/internal/config"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

type healthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
	Version string `json:"version"`
}

func main() {
	appConfig := config.LoadConfig()

	router := chi.NewRouter()
	router.Use(chiMiddleware.Logger)
	router.Use(chiMiddleware.Recoverer)

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := healthResponse{
			Status:  "ok",
			Service: "aegis-iam",
			Version: "1.0.0",
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	})

	address := ":" + appConfig.Port
	log.Printf("Starting Aegis-IAM server on port %s...", appConfig.Port)
	if err := http.ListenAndServe(address, router); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
