package store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/sayamjn/lykapp/internal/models"
)

type APIStore struct {
	apiBaseURL string
	apiKey     string
	client     *http.Client
	cache      struct {
		ads      map[string]models.Ad
		clicks   map[string]models.AdClick
		adsMux   sync.RWMutex
		clickMux sync.RWMutex
	}
}

func NewAPIStore(apiBaseURL, apiKey string) *APIStore {
	store := &APIStore{
		apiBaseURL: apiBaseURL,
		apiKey:     apiKey,
		client:     &http.Client{Timeout: 10 * time.Second},
	}
	store.cache.ads = make(map[string]models.Ad)
	store.cache.clicks = make(map[string]models.AdClick)
	return store
}

func (s *APIStore) FetchAdsFromAPI() error {
	url := fmt.Sprintf("%s/ads?key=%s", s.apiBaseURL, s.apiKey)
	
	resp, err := s.client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch ads: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	var ads []models.Ad
	if err := json.NewDecoder(resp.Body).Decode(&ads); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	s.cache.adsMux.Lock()
	s.cache.ads = make(map[string]models.Ad)
	for _, ad := range ads {
		s.cache.ads[ad.ID] = ad
	}
	s.cache.adsMux.Unlock()

	return nil
}

func (s *APIStore) GetAds() ([]models.Ad, error) {
	s.cache.adsMux.RLock()
	cacheEmpty := len(s.cache.ads) == 0
	s.cache.adsMux.RUnlock()

	if cacheEmpty {
		if err := s.FetchAdsFromAPI(); err != nil {
			return nil, err
		}
	}

	s.cache.adsMux.RLock()
	defer s.cache.adsMux.RUnlock()

	ads := make([]models.Ad, 0, len(s.cache.ads))
	for _, ad := range s.cache.ads {
		ads = append(ads, ad)
	}
	return ads, nil
}

func (s *APIStore) RecordClick(click models.AdClick) error {
	url := fmt.Sprintf("%s/clicks?key=%s", s.apiBaseURL, s.apiKey)
	
	clickData, err := json.Marshal(click)
	if err != nil {
		return fmt.Errorf("failed to marshal click data: %w", err)
	}

	resp, err := s.client.Post(url, "application/json", bytes.NewBuffer(clickData))
	if err != nil {
		return fmt.Errorf("failed to send click data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	s.cache.clickMux.Lock()
	defer s.cache.clickMux.Unlock()
	s.cache.clicks[click.ID] = click

	return nil
}

func (s *APIStore) GetClicks() ([]models.AdClick, error) {
	url := fmt.Sprintf("%s/clicks?key=%s", s.apiBaseURL, s.apiKey)
	
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch clicks: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	var clicks []models.AdClick
	if err := json.NewDecoder(resp.Body).Decode(&clicks); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return clicks, nil
}

func (s *APIStore) RefreshAds() error {
	return s.FetchAdsFromAPI()
}