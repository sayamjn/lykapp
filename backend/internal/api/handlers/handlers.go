package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sayamjn/lykapp/internal/models"
	"github.com/sayamjn/lykapp/internal/store"
)

type Handler struct {
	store *store.MemoryStore
}

func NewHandler(store *store.MemoryStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) GetAds(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ads := h.store.GetAds()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ads)
}

func (h *Handler) RecordClick(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var click struct {
		AdID            string  `json:"adId"`
		VideoPlaybackTs float64 `json:"videoPlaybackTs"`
	}

	if err := json.NewDecoder(r.Body).Decode(&click); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newClick := models.AdClick{
		ID:              uuid.New().String(),
		AdID:            click.AdID,
		Timestamp:       time.Now(),
		IPAddress:       r.RemoteAddr,
		VideoPlaybackTs: click.VideoPlaybackTs,
	}

	if err := h.store.RecordClick(newClick); err != nil {
		http.Error(w, "Failed to record click", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}