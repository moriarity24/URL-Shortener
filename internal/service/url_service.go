package service

import (
	"context"
	"fmt"
	"time"

	"github.com/teris-io/shortid"

	"github.com/moriarity24/url-shortener/internal/models"
	"github.com/moriarity24/url-shortener/internal/repository"
)

type URLService struct {
	repo    *repository.URLRepository
	baseURL string
}

func NewURLService(repo *repository.URLRepository, baseURL string) *URLService {
	return &URLService{
		repo:    repo,
		baseURL: baseURL,
	}
}

// ShortenURL creates a short URL
func (s *URLService) ShortenURL(ctx context.Context, originalURL string) (*models.CreateURLResponse, error) {
	// Check if URL already shortened
	existing, err := s.repo.FindByOriginalURL(originalURL)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		// Return existing short code
		return &models.CreateURLResponse{
			ShortCode:   existing.ShortCode,
			ShortURL:    fmt.Sprintf("%s/%s", s.baseURL, existing.ShortCode),
			OriginalURL: existing.OriginalURL,
			IsExisting:  true,
		}, nil
	}

	// Generate unique short code
	shortCode, err := shortid.Generate()
	if err != nil {
		return nil, fmt.Errorf("failed to generate short code: %w", err)
	}

	// Create URL record
	url := &models.URL{
		OriginalURL: originalURL,
		ShortCode:   shortCode,
		Clicks:      0,
		ExpiresAt:   time.Now().AddDate(0, 0, 30), // Expires in 30 days
	}

	if err := s.repo.Create(url); err != nil {
		return nil, fmt.Errorf("failed to save URL: %w", err)
	}

	return &models.CreateURLResponse{
		ShortCode:   url.ShortCode,
		ShortURL:    fmt.Sprintf("%s/%s", s.baseURL, url.ShortCode),
		OriginalURL: url.OriginalURL,
	}, nil
}

// GetOriginalURL retrieves original URL and increments clicks
func (s *URLService) GetOriginalURL(ctx context.Context, shortCode string) (string, error) {
	url, err := s.repo.FindByShortCode(ctx, shortCode)
	if err != nil {
		return "", err
	}

	// Increment click counter asynchronously (don't wait)
	go s.repo.IncrementClicks(shortCode)

	return url.OriginalURL, nil
}
