package services

import (
	"log"

	"github.com/google/uuid"

	"threat-intel-engine/internal/models"
	"threat-intel-engine/internal/repository"
)

type CorrelationService struct {
	repo *repository.FindingRepository
}

func NewCorrelationService(repo *repository.FindingRepository) *CorrelationService {
	return &CorrelationService{repo: repo}
}

func (c *CorrelationService) ProcessCorrelations(finding models.Finding) {
	go func() {
		correlations := c.generateCorrelations(finding)
		for _, correlation := range correlations {
			if err := c.repo.CreateCorrelation(&correlation); err != nil {
				log.Printf("Error creating correlation: %v", err)
			}
		}
	}()
}

func (c *CorrelationService) generateCorrelations(finding models.Finding) []models.ThreatCorrelation {
	var correlations []models.ThreatCorrelation

	if finding.FilePath != "" {
		correlation := models.ThreatCorrelation{
			ID:                uuid.New(),
			FindingID:         finding.ID,
			CorrelationType:   "file_based",
			RelatedEntityID:   finding.FilePath,
			RelatedEntityType: "file",
			Confidence:        0.8,
		}
		correlations = append(correlations, correlation)
	}

	if finding.Endpoint != "" {
		correlation := models.ThreatCorrelation{
			ID:                uuid.New(),
			FindingID:         finding.ID,
			CorrelationType:   "endpoint_based",
			RelatedEntityID:   finding.Endpoint,
			RelatedEntityType: "endpoint",
			Confidence:        0.9,
		}
		correlations = append(correlations, correlation)
	}

	return correlations
}
