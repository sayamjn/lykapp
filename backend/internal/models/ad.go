package models

import "time"

type Ad struct {
	ID        string    `json:"id"`
	ImageURL  string    `json:"imageUrl"`
	TargetURL string    `json:"targetUrl"`
	Position  string    `json:"position"`
	CreatedAt time.Time `json:"createdAt"`
}

type AdClick struct {
	ID              string    `json:"id"`
	AdID            string    `json:"adId"`
	Timestamp       time.Time `json:"timestamp"`
	IPAddress       string    `json:"ipAddress"`
	VideoPlaybackTs float64   `json:"videoPlaybackTs"`
}
