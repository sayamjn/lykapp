package handlers

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/sayamjn/lykapp/internal/models"
)

type mockStore struct {
    ads    []models.Ad
    clicks []models.AdClick
}

func (m *mockStore) GetAds() ([]models.Ad, error) {
    return m.ads, nil
}

func (m *mockStore) RecordClick(click models.AdClick) error {
    m.clicks = append(m.clicks, click)
    return nil
}

func (m *mockStore) GetClicks() ([]models.AdClick, error) {
    return m.clicks, nil
}

func (m *mockStore) RefreshAds() error {
    return nil
}

func TestGetAds(t *testing.T) {
    store := &mockStore{
        ads: []models.Ad{
            {ID: "1", ImageURL: "test.jpg"},
        },
    }
    handler := NewHandler(store)

    req := httptest.NewRequest("GET", "/api/ads", nil)
    w := httptest.NewRecorder()

    handler.GetAds(w, req)

    if w.Code != http.StatusOK {
        t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
    }

    var response []models.Ad
    err := json.NewDecoder(w.Body).Decode(&response)
    if err != nil {
        t.Fatalf("Failed to decode response: %v", err)
    }

    if len(response) != 1 {
        t.Errorf("Expected 1 ad, got %d", len(response))
    }
}

func TestRecordClick(t *testing.T) {
    store := &mockStore{}
    handler := NewHandler(store)

    clickData := map[string]interface{}{
        "adId":            "1",
        "videoPlaybackTs": 10.5,
    }
    body, _ := json.Marshal(clickData)

    req := httptest.NewRequest("POST", "/api/ads/click", bytes.NewBuffer(body))
    w := httptest.NewRecorder()

    handler.RecordClick(w, req)

    if w.Code != http.StatusCreated {
        t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
    }

    if len(store.clicks) != 1 {
        t.Errorf("Expected 1 click to be recorded, got %d", len(store.clicks))
    }
}