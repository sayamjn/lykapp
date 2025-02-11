package main

import (
    "log"
    "net/http"
    "os"
    "path/filepath"
    "time"

    "github.com/sayamjn/lykapp/internal/api/handlers"
    "github.com/sayamjn/lykapp/internal/config"
    "github.com/sayamjn/lykapp/internal/middleware"
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

    logDir := getEnvOrDefault("LOG_DIR", "logs")
    if err := os.MkdirAll(logDir, 0755); err != nil {
        log.Fatalf("Failed to create log directory: %v", err)
    }

    logPath := filepath.Join(logDir, "app.log")
    logger, err := middleware.NewLogger(logPath)
    if err != nil {
        log.Fatalf("Failed to initialize logger: %v", err)
    }
    defer logger.Close()

    adStore := store.NewMemoryStore(cfg.UnsplashKey)
    
    handler := handlers.NewHandler(adStore, logger)

    if cfg.RefreshEnabled {
        go func() {
            ticker := time.NewTicker(cfg.CacheTimeout)
            defer ticker.Stop()

            for range ticker.C {
                if err := adStore.RefreshAds(); err != nil {
                    logger.LogEvent("error", map[string]string{
                        "type":    "ad_refresh_failed",
                        "error":   err.Error(),
                    })
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

    handlerChain := corsMiddleware(logger.LoggingMiddleware(mux))

    port := getEnvOrDefault("PORT", "8080")
    serverAddr := ":" + port

    logger.LogEvent("startup", map[string]string{
        "port": port,
        "environment": getEnvOrDefault("ENV", "development"),
    })

    log.Printf("Server starting on port %s", port)
    if err := http.ListenAndServe(serverAddr, handlerChain); err != nil {
        logger.LogEvent("shutdown", map[string]string{
            "error": err.Error(),
        })
        log.Fatal(err)
    }
}

func getEnvOrDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}