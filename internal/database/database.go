package database

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/moriarity24/url-shortener/internal/config"
	"github.com/moriarity24/url-shortener/internal/models"
)

type Database struct {
	DB    *gorm.DB
	Redis *redis.Client
}

// NewDatabase creates database connections
func NewDatabase(cfg *config.Config) (*Database, error) {
	// DEBUG: Print the DSN to see what we're trying to connect to
	dsn := cfg.GetDSN()
	log.Println("🔍 Attempting to connect with DSN:", dsn)
	
	// Connect to PostgreSQL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Log SQL queries
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("✅ Connected to PostgreSQL")

	// Auto-migrate database schema
	if err := db.AutoMigrate(&models.URL{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("✅ Database tables migrated")

	// Connect to Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		DB:   0, // Use default DB
	})

	// Test Redis connection
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	log.Println("✅ Connected to Redis")

	return &Database{
		DB:    db,
		Redis: redisClient,
	}, nil
}

// Close closes database connections
func (d *Database) Close() {
	if d.Redis != nil {
		d.Redis.Close()
	}
}
