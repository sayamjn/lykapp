package store

import (
    "crypto/rand"
    "database/sql"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "net/http"
    "sync"
    "time"
    _ "github.com/mattn/go-sqlite3"
    "github.com/sayamjn/lykapp/internal/models"
)

type SQLiteStore struct {
    db          *sql.DB
    unsplashKey string
    cache       struct {
        ads    map[string]models.Ad
        adsMux sync.RWMutex
    }
    client *http.Client
}

func NewSQLiteStore(dbPath, unsplashKey string) (Store, error) {
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }

    if err := createTables(db); err != nil {
        return nil, fmt.Errorf("failed to create tables: %w", err)
    }

    store := &SQLiteStore{
        db:          db,
        unsplashKey: unsplashKey,
        client:      &http.Client{Timeout: 10 * time.Second},
    }
    store.cache.ads = make(map[string]models.Ad)

    if err := store.RefreshAds(); err != nil {
        return nil, fmt.Errorf("failed to fetch initial ads: %w", err)
    }

    return store, nil
}

func createTables(db *sql.DB) error {
    _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS sessions (
            id TEXT PRIMARY KEY,
            user_ip TEXT NOT NULL,
            created_at TIMESTAMP NOT NULL,
            last_activity TIMESTAMP NOT NULL
        )
    `)
    if err != nil {
        return err
    }

    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS ad_clicks (
            id TEXT PRIMARY KEY,
            ad_id TEXT NOT NULL,
            session_id TEXT NOT NULL,
            timestamp TIMESTAMP NOT NULL,
            ip_address TEXT NOT NULL,
            video_playback_ts REAL NOT NULL,
            FOREIGN KEY(session_id) REFERENCES sessions(id)
        )
    `)
    if err != nil {
        return err
    }

    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS analytics (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            event_type TEXT NOT NULL,
            event_data TEXT NOT NULL,
            timestamp TIMESTAMP NOT NULL
        )
    `)
    return err
}

func (s *SQLiteStore) RefreshAds() error {
    req, err := http.NewRequest("GET", "https://api.unsplash.com/photos/random", nil)
    if err != nil {
        return err
    }

    q := req.URL.Query()
    q.Add("count", "5")
    q.Add("orientation", "landscape")
    q.Add("query", "product")
    req.URL.RawQuery = q.Encode()

    req.Header.Set("Authorization", fmt.Sprintf("Client-ID %s", s.unsplashKey))

    resp, err := s.client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("unsplash API returned status code: %d", resp.StatusCode)
    }

    var photos []struct {
        ID    string `json:"id"`
        URLs  struct {
            Regular string `json:"regular"`
        } `json:"urls"`
        Links struct {
            HTML string `json:"html"`
        } `json:"links"`
        User struct {
            Name string `json:"name"`
        } `json:"user"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&photos); err != nil {
        return err
    }

    positions := []string{"top-right", "top-left", "bottom-right", "bottom-left"}
    
    s.cache.adsMux.Lock()
    s.cache.ads = make(map[string]models.Ad)
    for i, photo := range photos {
        ad := models.Ad{
            ID:        photo.ID,
            ImageURL:  photo.URLs.Regular,
            TargetURL: photo.Links.HTML,
            Position:  positions[i%len(positions)],
            CreatedAt: time.Now(),
            Title:     fmt.Sprintf("Photo by %s", photo.User.Name),
        }
        s.cache.ads[ad.ID] = ad
    }
    s.cache.adsMux.Unlock()

    return nil
}

func (s *SQLiteStore) GetAds() ([]models.Ad, error) {
    s.logAnalytics("get_ads", "Ads requested")
    
    s.cache.adsMux.RLock()
    defer s.cache.adsMux.RUnlock()

    ads := make([]models.Ad, 0, len(s.cache.ads))
    for _, ad := range s.cache.ads {
        ads = append(ads, ad)
    }
    return ads, nil
}

func (s *SQLiteStore) RecordClick(click models.AdClick) error {
    tx, err := s.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    sessionID, err := s.getOrCreateSession(tx, click.IPAddress)
    if err != nil {
        return err
    }

    _, err = tx.Exec(`
        INSERT INTO ad_clicks (id, ad_id, session_id, timestamp, ip_address, video_playback_ts)
        VALUES (?, ?, ?, ?, ?, ?)
    `, click.ID, click.AdID, sessionID, click.Timestamp, click.IPAddress, click.VideoPlaybackTs)
    if err != nil {
        return err
    }

    s.logAnalytics("ad_click", fmt.Sprintf("Ad %s clicked by session %s", click.AdID, sessionID))

    return tx.Commit()
}

func (s *SQLiteStore) GetClicks() ([]models.AdClick, error) {
    rows, err := s.db.Query(`
        SELECT id, ad_id, timestamp, ip_address, video_playback_ts
        FROM ad_clicks
        ORDER BY timestamp DESC
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var clicks []models.AdClick
    for rows.Next() {
        var click models.AdClick
        err := rows.Scan(&click.ID, &click.AdID, &click.Timestamp, &click.IPAddress, &click.VideoPlaybackTs)
        if err != nil {
            return nil, err
        }
        clicks = append(clicks, click)
    }

    return clicks, nil
}

func generateSessionID() string {
    b := make([]byte, 16)
    rand.Read(b)
    return hex.EncodeToString(b)
}

func (s *SQLiteStore) getOrCreateSession(tx *sql.Tx, userIP string) (string, error) {
    var sessionID string
    err := tx.QueryRow(`
        SELECT id FROM sessions 
        WHERE user_ip = ? AND last_activity > ?
        ORDER BY last_activity DESC LIMIT 1
    `, userIP, time.Now().Add(-24*time.Hour)).Scan(&sessionID)

    if err == nil {
        _, err = tx.Exec("UPDATE sessions SET last_activity = ? WHERE id = ?",
            time.Now(), sessionID)
        return sessionID, err
    }

    sessionID = generateSessionID()
    _, err = tx.Exec(`
        INSERT INTO sessions (id, user_ip, created_at, last_activity)
        VALUES (?, ?, ?, ?)
    `, sessionID, userIP, time.Now(), time.Now())

    return sessionID, err
}

func (s *SQLiteStore) logAnalytics(eventType, eventData string) error {
    _, err := s.db.Exec(`
        INSERT INTO analytics (event_type, event_data, timestamp)
        VALUES (?, ?, ?)
    `, eventType, eventData, time.Now())
    return err
}