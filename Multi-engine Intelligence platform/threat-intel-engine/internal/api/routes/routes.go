package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"threat-intel-engine/internal/api/handlers"
	"threat-intel-engine/internal/config"
	"threat-intel-engine/internal/repository"
	"threat-intel-engine/internal/services"
	"threat-intel-engine/pkg/database"
)

func SetupRoutes(db *database.MongoDB) *gin.Engine {
	cfg, _ := config.Load()

	// Initialize repository
	findingRepo := repository.NewFindingRepository(db.Database)

	// Initialize services
	validationSvc := services.NewValidationService(findingRepo)
	normalizationSvc := services.NewNormalizationService()
	enrichmentSvc := services.NewEnrichmentService(findingRepo)
	correlationSvc := services.NewCorrelationService(findingRepo)
	alertingSvc := services.NewAlertingService(cfg.SOARWebhook)

	// Initialize handlers
	findingsHandler := handlers.NewFindingsHandler(
		validationSvc,
		normalizationSvc,
		enrichmentSvc,
		correlationSvc,
		alertingSvc,
		findingRepo,
	)

	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// API routes
	api := r.Group("/api/v1")
	{
		api.POST("/findings", findingsHandler.CreateFinding)
		api.GET("/findings", findingsHandler.GetFindings)
		api.GET("/findings/:id/enrichments", findingsHandler.GetEnrichments)
	}

	return r
}
