package main
// testing webhooks
import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/moriarity24/url-shortener/internal/config"
	"github.com/moriarity24/url-shortener/internal/database"
	"github.com/moriarity24/url-shortener/internal/handlers"
	"github.com/moriarity24/url-shortener/internal/repository"
	"github.com/moriarity24/url-shortener/internal/service"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize database connections
	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Initialize layers
	urlRepo := repository.NewURLRepository(db.DB, db.Redis)
	urlService := service.NewURLService(urlRepo, cfg.BaseURL)
	urlHandler := handlers.NewURLHandler(urlService)

	// Setup Gin router
	router := gin.Default()

	// CORS middleware (allow frontend to call API)
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Routes
	router.GET("/health", urlHandler.HealthCheck)
	router.POST("/api/shorten", urlHandler.CreateShortURL)
	router.GET("/:shortCode", urlHandler.RedirectToOriginal)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("🚀 Server starting on http://localhost%s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}