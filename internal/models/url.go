package models
// this is to check the webhook
import (
	"time"
)
type URL struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	OriginalURL string   `gorm:"type:text;not null" json:"original_url"`
	ShortCode  string   `gorm:"type:varchar(10);uniqueIndex;not null" json:"short_code"`
	Clicks int `gorm:"default:0" json:"clicks"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
}

type CreateURLRequest struct {
	URL string `json:"url" binding:"required,url"`
}

// CreateURLResponse - Response after creating short URL
type CreateURLResponse struct {
	ShortCode   string `json:"short_code"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	IsExisting  bool   `json:"is_existing,omitempty"`
}