package main

import (
	"log"
	"net/http"
	"os"

	"github.com/sayamjn/lykapp/internal/api/handlers"
	"github.com/sayamjn/lykapp/internal/store"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	store := store.NewMemoryStore()
	handler := handlers.NewHandler(store)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/ads", handler.GetAds)
	mux.HandleFunc("/api/ads/click", handler.RecordClick)

	handlerWithCors := corsMiddleware(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, handlerWithCors); err != nil {
		log.Fatal(err)
	}
}