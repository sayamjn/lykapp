package handlers

import (
    "encoding/json"
    "net/http"
    "time"

    "github.com/google/uuid"
    "github.com/sayamjn/lykapp/internal/middleware"
    "github.com/sayamjn/lykapp/internal/models"
    "github.com/sayamjn/lykapp/internal/store"
)

type Handler struct {
    store  store.Store
    logger *middleware.Logger
}

func NewHandler(store store.Store, logger *middleware.Logger) *Handler {
    return &Handler{
        store:  store,
        logger: logger,
    }
}

func (h *Handler) GetAds(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        h.logger.LogEvent("error", map[string]string{
            "type":   "method_not_allowed",
            "method": r.Method,
            "path":   r.URL.Path,
        })
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    ads, err := h.store.GetAds()
    if err != nil {
        h.logger.LogEvent("error", map[string]string{
            "type":  "fetch_ads_failed",
            "error": err.Error(),
        })
        http.Error(w, "Failed to fetch ads", http.StatusInternalServerError)
        return
    }

    h.logger.LogEvent("ads_fetched", map[string]interface{}{
        "count": len(ads),
    })

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(ads); err != nil {
        h.logger.LogEvent("error", map[string]string{
            "type":  "response_encode_failed",
            "error": err.Error(),
        })
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
        return
    }
}

func (h *Handler) RecordClick(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        h.logger.LogEvent("error", map[string]string{
            "type":   "method_not_allowed",
            "method": r.Method,
            "path":   r.URL.Path,
        })
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var click struct {
        AdID            string  `json:"adId"`
        VideoPlaybackTs float64 `json:"videoPlaybackTs"`
    }

    if err := json.NewDecoder(r.Body).Decode(&click); err != nil {
        h.logger.LogEvent("error", map[string]string{
            "type":  "invalid_request",
            "error": err.Error(),
        })
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if click.AdID == "" {
        h.logger.LogEvent("error", map[string]string{
            "type": "invalid_request",
            "error": "adId is required",
        })
        http.Error(w, "adId is required", http.StatusBadRequest)
        return
    }

    sessionID := ""
    if cookie, err := r.Cookie("session_id"); err == nil {
        sessionID = cookie.Value
    } else {
        sessionID = uuid.New().String()
        http.SetCookie(w, &http.Cookie{
            Name:     "session_id",
            Value:    sessionID,
            Path:     "/",
            HttpOnly: true,
            Secure:   r.TLS != nil,
            MaxAge:   86400, // 24 hours
            SameSite: http.SameSiteStrictMode,
        })
        h.logger.LogEvent("session_created", map[string]string{
            "session_id": sessionID,
        })
    }

    newClick := models.AdClick{
        ID:              uuid.New().String(),
        AdID:            click.AdID,
        Timestamp:       time.Now(),
        IPAddress:       r.RemoteAddr,
        VideoPlaybackTs: click.VideoPlaybackTs,
    }

    if err := h.store.RecordClick(newClick); err != nil {
        h.logger.LogEvent("error", map[string]string{
            "type":      "record_click_failed",
            "error":     err.Error(),
            "click_id":  newClick.ID,
            "ad_id":     newClick.AdID,
        })
        http.Error(w, "Failed to record click", http.StatusInternalServerError)
        return
    }

    h.logger.LogEvent("click_recorded", map[string]interface{}{
        "click_id":         newClick.ID,
        "ad_id":           newClick.AdID,
        "video_timestamp": newClick.VideoPlaybackTs,
        "session_id":      sessionID,
    })

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{
        "id":      newClick.ID,
        "message": "Click recorded successfully",
    })
}

func (h *Handler) GetClicks(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        h.logger.LogEvent("error", map[string]string{
            "type":   "method_not_allowed",
            "method": r.Method,
            "path":   r.URL.Path,
        })
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    clicks, err := h.store.GetClicks()
    if err != nil {
        h.logger.LogEvent("error", map[string]string{
            "type":  "fetch_clicks_failed",
            "error": err.Error(),
        })
        http.Error(w, "Failed to fetch clicks", http.StatusInternalServerError)
        return
    }

    h.logger.LogEvent("clicks_fetched", map[string]interface{}{
        "count": len(clicks),
    })

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(clicks); err != nil {
        h.logger.LogEvent("error", map[string]string{
            "type":  "response_encode_failed",
            "error": err.Error(),
        })
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
        return
    }
}