package store

import (
    "encoding/json"
    "fmt"
    "net/http"
    "sync"
    "time"

    "github.com/sayamjn/lykapp/internal/models"
)

type MemoryStore struct {
    ads          map[string]models.Ad
    clicks       map[string]models.AdClick
    adsMux       sync.RWMutex
    clickMux     sync.RWMutex
    unsplashKey  string
    client       *http.Client
}

type UnsplashPhoto struct {
    ID          string `json:"id"`
    URLs        struct {
        Regular string `json:"regular"`
    } `json:"urls"`
    Links struct {
        HTML string `json:"html"`
    } `json:"links"`
    User struct {
        Name string `json:"name"`
    } `json:"user"`
}

func NewMemoryStore(unsplashKey string) Store {
    store := &MemoryStore{
        ads:         make(map[string]models.Ad),
        clicks:      make(map[string]models.AdClick),
        unsplashKey: unsplashKey,
        client:      &http.Client{Timeout: 10 * time.Second},
    }
    
    if err := store.RefreshAds(); err != nil {
        fmt.Printf("Failed to fetch initial ads: %v\n", err)
    }

    return store
}

func (s *MemoryStore) fetchUnsplashPhotos() ([]UnsplashPhoto, error) {
    req, err := http.NewRequest("GET", "https://api.unsplash.com/photos/random", nil)
    if err != nil {
        return nil, err
    }

    q := req.URL.Query()
    q.Add("count", "5") 
    q.Add("orientation", "landscape")
    q.Add("query", "product")
    req.URL.RawQuery = q.Encode()

    req.Header.Set("Authorization", fmt.Sprintf("Client-ID %s", s.unsplashKey))

    resp, err := s.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unsplash API returned status code: %d", resp.StatusCode)
    }

    var photos []UnsplashPhoto
    if err := json.NewDecoder(resp.Body).Decode(&photos); err != nil {
        return nil, err
    }

    return photos, nil
}

func (s *MemoryStore) RefreshAds() error {
    photos, err := s.fetchUnsplashPhotos()
    if err != nil {
        return err
    }

    positions := []string{"top-right", "top-left", "bottom-right", "bottom-left"}
    newAds := make(map[string]models.Ad)

    for i, photo := range photos {
        ad := models.Ad{
            ID:        photo.ID,
            ImageURL:  photo.URLs.Regular,
            TargetURL: photo.Links.HTML,
            Position:  positions[i%len(positions)],
            CreatedAt: time.Now(),
            Title:     fmt.Sprintf("Photo by %s", photo.User.Name),
        }
        newAds[ad.ID] = ad
    }

    s.adsMux.Lock()
    s.ads = newAds
    s.adsMux.Unlock()

    return nil
}

func (s *MemoryStore) GetAds() ([]models.Ad, error) {
    s.adsMux.RLock()
    defer s.adsMux.RUnlock()

    ads := make([]models.Ad, 0, len(s.ads))
    for _, ad := range s.ads {
        ads = append(ads, ad)
    }
    return ads, nil
}

func (s *MemoryStore) RecordClick(click models.AdClick) error {
    s.clickMux.Lock()
    defer s.clickMux.Unlock()

    s.clicks[click.ID] = click
    return nil
}

func (s *MemoryStore) GetClicks() ([]models.AdClick, error) {
    s.clickMux.RLock()
    defer s.clickMux.RUnlock()

    clicks := make([]models.AdClick, 0, len(s.clicks))
    for _, click := range s.clicks {
        clicks = append(clicks, click)
    }
    return clicks, nil
}