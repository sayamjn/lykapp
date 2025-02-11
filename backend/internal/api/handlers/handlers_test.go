package handlers

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "os"
    "path/filepath"
    "testing"
    "time"

    "github.com/sayamjn/lykapp/internal/middleware"
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

func setupTest(t *testing.T) (*Handler, *mockStore, *middleware.Logger, func()) {
    tmpDir, err := os.MkdirTemp("", "test-logs-*")
    if err != nil {
        t.Fatalf("Failed to create temp directory: %v", err)
    }

    logger, err := middleware.NewLogger(filepath.Join(tmpDir, "test.log"))
    if err != nil {
        os.RemoveAll(tmpDir)
        t.Fatalf("Failed to create logger: %v", err)
    }

    store := &mockStore{
        ads: []models.Ad{
            {
                ID:        "test-ad-1",
                ImageURL:  "test.jpg",
                Position: "top-right",
                Title:    "Test Ad",
                CreatedAt: time.Now(),
            },
        },
        clicks: make([]models.AdClick, 0),
    }

    handler := NewHandler(store, logger)

    cleanup := func() {
        logger.Close()
        os.RemoveAll(tmpDir)
    }

    return handler, store, logger, cleanup
}

func TestGetAds(t *testing.T) {
    handler, _, _, cleanup := setupTest(t)
    defer cleanup()

    tests := []struct {
        name         string
        method      string
        expectedCode int
    }{
        {
            name:         "Success - GET method",
            method:      http.MethodGet,
            expectedCode: http.StatusOK,
        },
        {
            name:         "Failed - Wrong HTTP method",
            method:      http.MethodPost,
            expectedCode: http.StatusMethodNotAllowed,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req := httptest.NewRequest(tt.method, "/api/ads", nil)
            w := httptest.NewRecorder()

            handler.GetAds(w, req)

            if w.Code != tt.expectedCode {
                t.Errorf("Expected status code %d, got %d", tt.expectedCode, w.Code)
            }

            if tt.expectedCode == http.StatusOK {
                var response []models.Ad
                err := json.NewDecoder(w.Body).Decode(&response)
                if err != nil {
                    t.Fatalf("Failed to decode response: %v", err)
                }

                if len(response) != 1 {
                    t.Errorf("Expected 1 ad, got %d", len(response))
                }
            }
        })
    }
}

func TestRecordClick(t *testing.T) {
    handler, store, _, cleanup := setupTest(t)
    defer cleanup()

    tests := []struct {
        name         string
        requestBody interface{}
        expectedCode int
    }{
        {
            name: "Success - Valid click",
            requestBody: map[string]interface{}{
                "adId":            "test-ad-1",
                "videoPlaybackTs": 10.5,
            },
            expectedCode: http.StatusCreated,
        },
        {
            name:         "Failed - Invalid request body",
            requestBody:  "invalid json",
            expectedCode: http.StatusBadRequest,
        },
        {
            name: "Failed - Missing required fields",
            requestBody: map[string]interface{}{
                "videoPlaybackTs": 10.5,
            },
            expectedCode: http.StatusBadRequest,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            var body bytes.Buffer
            if str, ok := tt.requestBody.(string); ok {
                body.WriteString(str)
            } else {
                json.NewEncoder(&body).Encode(tt.requestBody)
            }

            req := httptest.NewRequest(http.MethodPost, "/api/ads/click", &body)
            w := httptest.NewRecorder()

            handler.RecordClick(w, req)

            if w.Code != tt.expectedCode {
                t.Errorf("Expected status code %d, got %d", tt.expectedCode, w.Code)
            }

            if tt.expectedCode == http.StatusCreated {
                if len(store.clicks) != 1 {
                    t.Errorf("Expected 1 click to be recorded, got %d", len(store.clicks))
                }

                var response map[string]string
                err := json.NewDecoder(w.Body).Decode(&response)
                if err != nil {
                    t.Fatalf("Failed to decode response: %v", err)
                }

                if response["message"] != "Click recorded successfully" {
                    t.Errorf("Expected success message, got %s", response["message"])
                }

                cookies := w.Result().Cookies()
                var sessionCookie *http.Cookie
                for _, cookie := range cookies {
                    if cookie.Name == "session_id" {
                        sessionCookie = cookie
                        break
                    }
                }

                if sessionCookie == nil {
                    t.Error("Expected session cookie to be set")
                } else {
                    if !sessionCookie.HttpOnly {
                        t.Error("Expected HttpOnly cookie")
                    }
                    if sessionCookie.MaxAge != 86400 {
                        t.Errorf("Expected MaxAge 86400, got %d", sessionCookie.MaxAge)
                    }
                }
            }
        })
    }
}

func TestGetClicks(t *testing.T) {
    handler, store, _, cleanup := setupTest(t)
    defer cleanup()

    store.clicks = append(store.clicks, models.AdClick{
        ID:              "test-click-1",
        AdID:            "test-ad-1",
        Timestamp:       time.Now(),
        IPAddress:       "127.0.0.1",
        VideoPlaybackTs: 10.5,
    })

    tests := []struct {
        name         string
        method      string
        expectedCode int
    }{
        {
            name:         "Success - GET method",
            method:      http.MethodGet,
            expectedCode: http.StatusOK,
        },
        {
            name:         "Failed - Wrong HTTP method",
            method:      http.MethodPost,
            expectedCode: http.StatusMethodNotAllowed,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req := httptest.NewRequest(tt.method, "/api/ads/clicks", nil)
            w := httptest.NewRecorder()

            handler.GetClicks(w, req)

            if w.Code != tt.expectedCode {
                t.Errorf("Expected status code %d, got %d", tt.expectedCode, w.Code)
            }

            if tt.expectedCode == http.StatusOK {
                var response []models.AdClick
                err := json.NewDecoder(w.Body).Decode(&response)
                if err != nil {
                    t.Fatalf("Failed to decode response: %v", err)
                }

                if len(response) != 1 {
                    t.Errorf("Expected 1 click, got %d", len(response))
                }
            }
        })
    }
}