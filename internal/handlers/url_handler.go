package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/moriarity24/url-shortener/internal/models"
	"github.com/moriarity24/url-shortener/internal/service"
)

type URLHandler struct {
	service *service.URLService
}

func NewURLHandler(service *service.URLService) *URLHandler {
	return &URLHandler{
		service: service,
	}
}

// CreateShortURL handles POST /api/shorten
func (h *URLHandler) CreateShortURL(c *gin.Context) {
	var req models.CreateURLRequest

	// Validate request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: " + err.Error(),
		})
		return
	}

	// Call service to shorten URL
	response, err := h.service.ShortenURL(c.Request.Context(), req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to shorten URL: " + err.Error(),
		})
		return
	}

	// Return 201 for new URLs, 200 for existing ones
	status := http.StatusCreated
	if response.IsExisting {
		status = http.StatusOK
	}
	c.JSON(status, response)
}

// RedirectToOriginal handles GET /:shortCode
func (h *URLHandler) RedirectToOriginal(c *gin.Context) {
	shortCode := c.Param("shortCode")

	// Get original URL
	originalURL, err := h.service.GetOriginalURL(c.Request.Context(), shortCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Short URL not found",
		})
		return
	}

	// Redirect with 301 (Moved Permanently)
	c.Redirect(http.StatusMovedPermanently, originalURL)
}

// HealthCheck handles GET /health
func (h *URLHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
	})
}