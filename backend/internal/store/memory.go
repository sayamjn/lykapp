package store

import (
	"sync"
	"time"

	"github.com/sayamjn/lykapp/internal/models"
)

type MemoryStore struct {
	ads      map[string]models.Ad
	clicks   map[string]models.AdClick
	adsMux   sync.RWMutex
	clickMux sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
	store := &MemoryStore{
		ads:    make(map[string]models.Ad),
		clicks: make(map[string]models.AdClick),
	}
	
	sampleAds := []models.Ad{
		{
			ID:        "ad1",
			ImageURL:  "https://example.com/ad1.jpg",
			TargetURL: "https://example.com/product1",
			Position:  "top-right",
			CreatedAt: time.Now(),
		},
		{
			ID:        "ad2",
			ImageURL:  "https://example.com/ad2.jpg",
			TargetURL: "https://example.com/product2",
			Position:  "bottom-left",
			CreatedAt: time.Now(),
		},
	}

	for _, ad := range sampleAds {
		store.ads[ad.ID] = ad
	}

	return store
}

func (s *MemoryStore) GetAds() []models.Ad {
	s.adsMux.RLock()
	defer s.adsMux.RUnlock()

	ads := make([]models.Ad, 0, len(s.ads))
	for _, ad := range s.ads {
		ads = append(ads, ad)
	}
	return ads
}

func (s *MemoryStore) RecordClick(click models.AdClick) error {
	s.clickMux.Lock()
	defer s.clickMux.Unlock()

	s.clicks[click.ID] = click
	return nil
}

func (s *MemoryStore) GetClicks() []models.AdClick {
	s.clickMux.RLock()
	defer s.clickMux.RUnlock()

	clicks := make([]models.AdClick, 0, len(s.clicks))
	for _, click := range s.clicks {
		clicks = append(clicks, click)
	}
	return clicks
}