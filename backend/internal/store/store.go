package store

import "github.com/sayamjn/lykapp/internal/models"

type Store interface {
    GetAds() ([]models.Ad, error)
    RecordClick(click models.AdClick) error
    GetClicks() ([]models.AdClick, error)
    RefreshAds() error
}