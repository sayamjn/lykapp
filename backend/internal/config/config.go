package config

import (
    "os"
    "time"
)

type Config struct {
    UseMemoryStore  bool
    UnsplashKey     string
    CacheTimeout    time.Duration
    RefreshEnabled  bool
    AllowedOrigins  string
}

func LoadConfig() *Config {
    return &Config{
        UnsplashKey:     getEnvOrDefault("UNSPLASH_ACCESS_KEY", "HbVp0qrtx6GAVd8Hat1gclyRLSNIBh9-Kzk0Rq_JlVs"),
        CacheTimeout:    getEnvDurationOrDefault("AD_CACHE_TIMEOUT", 5*time.Minute),
        RefreshEnabled:  getEnvBoolOrDefault("AD_REFRESH_ENABLED", true),
        AllowedOrigins:  getEnvOrDefault("CORS_ALLOWED_ORIGINS", "http://localhost:5173"),
    }
}

func getEnvOrDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvDurationOrDefault(key string, defaultValue time.Duration) time.Duration {
    if value := os.Getenv(key); value != "" {
        if duration, err := time.ParseDuration(value); err == nil {
            return duration
        }
    }
    return defaultValue
}

func getEnvBoolOrDefault(key string, defaultValue bool) bool {
    if value := os.Getenv(key); value != "" {
        return value == "true" || value == "1"
    }
    return defaultValue
}