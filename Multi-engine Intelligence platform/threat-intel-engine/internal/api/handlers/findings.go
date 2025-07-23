package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"threat-intel-engine/internal/models"
	"threat-intel-engine/internal/services"
)

type FindingsHandler struct {
	validationService    *services.ValidationService
	normalizationService *services.NormalizationService
	enrichmentService    *services.EnrichmentService
	correlationService   *services.CorrelationService
	alertingService      *services.AlertingService
	findingRepo          FindingRepository
}

type FindingRepository interface {
	Create(finding *models.Finding) error
	FindByQuery(params models.QueryParams) ([]models.Finding, error)
}

func NewFindingsHandler(
	validationSvc *services.ValidationService,
	normalizationSvc *services.NormalizationService,
	enrichmentSvc *services.EnrichmentService,
	correlationSvc *services.CorrelationService,
	alertingSvc *services.AlertingService,
	repo FindingRepository,
) *FindingsHandler {
	return &FindingsHandler{
		validationService:    validationSvc,
		normalizationService: normalizationSvc,
		enrichmentService:    enrichmentSvc,
		correlationService:   correlationSvc,
		alertingService:      alertingSvc,
		findingRepo:          repo,
	}
}

func (h *FindingsHandler) CreateFinding(c *gin.Context) {
	var req models.FindingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	if err := h.validationService.ValidateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check for duplicates
	isDuplicate, err := h.validationService.CheckDuplicate(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check duplicates"})
		return
	}

	if isDuplicate {
		c.JSON(http.StatusConflict, gin.H{"error": "Duplicate finding"})
		return
	}

	// Normalize and create finding
	finding := h.normalizationService.NormalizeRequest(req)

	// Store finding
	if err := h.findingRepo.Create(&finding); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store finding"})
		return
	}

	// Trigger async processing
	h.enrichmentService.ProcessEnrichment(finding)
	h.correlationService.ProcessCorrelations(finding)
	h.alertingService.ProcessAlert(finding)

	c.JSON(http.StatusCreated, gin.H{
		"id":      finding.ID,
		"message": "Finding created successfully",
	})
}

func (h *FindingsHandler) GetFindings(c *gin.Context) {
	var params models.QueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	findings, err := h.findingRepo.FindByQuery(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query findings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"findings": findings,
		"total":    len(findings),
		"limit":    params.Limit,
		"offset":   params.Offset,
	})
}

func (h *FindingsHandler) GetEnrichments(c *gin.Context) {
	findingID := c.Param("id")

	findingUUID, err := uuid.Parse(findingID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid finding ID"})
		return
	}

	enrichments, err := h.enrichmentService.GetEnrichments(findingUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get enrichments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"findingId":   findingID,
		"enrichments": enrichments,
	})
}
