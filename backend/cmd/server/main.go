package main

import (
    "log"
    "net/http"
    "os"
    "time"

    "github.com/sayamjn/lykapp/internal/api/handlers"
    "github.com/sayamjn/lykapp/internal/config"
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
    cfg := config.LoadConfig()

    if cfg.UnsplashKey == "" {
        log.Fatal("UNSPLASH_ACCESS_KEY is required")
    }

    adStore := store.NewMemoryStore(cfg.UnsplashKey)
    handler := handlers.NewHandler(adStore)

    if cfg.RefreshEnabled {
        go func() {
            ticker := time.NewTicker(cfg.CacheTimeout)
            defer ticker.Stop()

            for range ticker.C {
                if err := adStore.RefreshAds(); err != nil {
                    log.Printf("Error refreshing ads: %v", err)
                }
            }
        }()
    }

    mux := http.NewServeMux()

    mux.HandleFunc("/api/ads", handler.GetAds)
    mux.HandleFunc("/api/ads/click", handler.RecordClick)
    mux.HandleFunc("/api/ads/clicks", handler.GetClicks)

    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("healthy"))
    })

    handlerWithCors := corsMiddleware(mux)

    port := getEnvOrDefault("PORT", "8080")
    log.Printf("Server starting on port %s", port)
    if err := http.ListenAndServe(":"+port, handlerWithCors); err != nil {
        log.Fatal(err)
    }
}

func getEnvOrDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}