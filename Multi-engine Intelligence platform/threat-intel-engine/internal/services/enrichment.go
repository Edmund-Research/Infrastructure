package services

import (
	"log"
	"strings"
	"time"

	"threat-intel-engine/internal/models"
	"threat-intel-engine/internal/repository"

	"github.com/google/uuid"
)

type EnrichmentService struct {
	repo *repository.FindingRepository
}

func NewEnrichmentService(repo *repository.FindingRepository) *EnrichmentService {
	return &EnrichmentService{repo: repo}
}

func (e *EnrichmentService) ProcessEnrichment(finding models.Finding) {
	go func() {
		// Simulate enrichment processing time
		time.Sleep(2 * time.Second)

		enrichment := models.Enrichment{
			ID:               uuid.New(),
			FindingID:        finding.ID,
			PipelineRunID:    finding.PipelineRunID,
			CVEList:          e.lookupCVEs(finding),
			CWEList:          e.lookupCWEs(finding),
			ExploitScore:     e.calculateExploitScore(finding),
			EnrichedAt:       time.Now(),
			EnrichmentSource: "ThreatIntelDB",
		}

		if err := e.repo.CreateEnrichment(&enrichment); err != nil {
			log.Printf("Error creating enrichment: %v", err)
		}
	}()
}

func (e *EnrichmentService) GetEnrichments(findingID uuid.UUID) ([]models.Enrichment, error) {
	return e.repo.GetEnrichmentsByFindingID(findingID)
}

func (e *EnrichmentService) lookupCVEs(finding models.Finding) []string {
	mockCVEs := map[string][]string{
		"SQL_INJECTION":   {"CVE-2021-44228", "CVE-2022-22965"},
		"XSS":             {"CVE-2021-26855", "CVE-2021-34527"},
		"BUFFER_OVERFLOW": {"CVE-2021-3156", "CVE-2021-40438"},
	}

	for pattern, cves := range mockCVEs {
		if strings.Contains(strings.ToUpper(finding.RuleID), pattern) {
			return cves
		}
	}

	return []string{}
}

func (e *EnrichmentService) lookupCWEs(finding models.Finding) []string {
	mockCWEs := map[string][]string{
		"SQL_INJECTION":   {"CWE-89", "CWE-564"},
		"XSS":             {"CWE-79", "CWE-352"},
		"BUFFER_OVERFLOW": {"CWE-120", "CWE-787"},
	}

	for pattern, cwes := range mockCWEs {
		if strings.Contains(strings.ToUpper(finding.RuleID), pattern) {
			return cwes
		}
	}

	return []string{}
}

func (e *EnrichmentService) calculateExploitScore(finding models.Finding) float64 {
	severityScores := map[string]float64{
		"Critical": 9.5,
		"High":     7.5,
		"Medium":   5.0,
		"Low":      2.5,
	}

	return severityScores[finding.Severity]
}
