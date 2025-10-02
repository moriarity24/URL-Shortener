package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/moriarity24/url-shortener/internal/models"
)

type URLRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewURLRepository(db *gorm.DB, redisClient *redis.Client) *URLRepository {
	return &URLRepository{
		db:    db,
		redis: redisClient,
	}
}

// Create saves a new URL to database
func (r *URLRepository) Create(url *models.URL) error {
	return r.db.Create(url).Error
}

// FindByShortCode retrieves URL by short code (checks cache first!)
func (r *URLRepository) FindByShortCode(ctx context.Context, shortCode string) (*models.URL, error) {
	// Try Redis cache first
	cacheKey := fmt.Sprintf("url:%s", shortCode)
	cached, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		// Cache hit! Deserialize and return
		var url models.URL
		if err := json.Unmarshal([]byte(cached), &url); err == nil {
			return &url, nil
		}
	}

	// Cache miss - query database
	var url models.URL
	if err := r.db.Where("short_code = ?", shortCode).First(&url).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("URL not found")
		}
		return nil, err
	}

	// Store in cache for 24 hours
	urlJSON, _ := json.Marshal(url)
	r.redis.Set(ctx, cacheKey, urlJSON, 24*time.Hour)

	return &url, nil
}

// IncrementClicks updates click counter
func (r *URLRepository) IncrementClicks(shortCode string) error {
	return r.db.Model(&models.URL{}).
		Where("short_code = ?", shortCode).
		UpdateColumn("clicks", gorm.Expr("clicks + ?", 1)).
		Error
}

// FindByOriginalURL checks if URL already shortened
func (r *URLRepository) FindByOriginalURL(url string) (*models.URL, error) {
	var existingURL models.URL
	err := r.db.Where("original_url = ?", url).First(&existingURL).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil // Not found is OK
	}
	return &existingURL, err
}