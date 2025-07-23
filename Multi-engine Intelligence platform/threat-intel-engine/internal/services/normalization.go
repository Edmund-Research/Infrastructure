package services

import (
	"strings"
	"time"

	"threat-intel-engine/internal/models"

	"github.com/google/uuid"
)

type NormalizationService struct{}

func NewNormalizationService() *NormalizationService {
	return &NormalizationService{}
}

func (n *NormalizationService) NormalizeRequest(req models.FindingRequest) models.Finding {
	finding := models.Finding{
		ID:                 uuid.New(),
		PipelineRunID:      req.PipelineRunID,
		Source:             strings.ToUpper(req.Source),
		Severity:           n.normalizeSeverity(req.Severity),
		Timestamp:          time.Now(),
		FilePath:           req.FilePath,
		Endpoint:           req.Endpoint,
		RuleID:             req.RuleID,
		Description:        req.Description,
		Evidence:           req.Evidence,
		DataClassification: n.normalizeDataClassification(req.DataClassification),
		CorrelationKeys:    n.normalizeCorrelationKeys(req.CorrelationKeys),
	}

	return finding
}

func (n *NormalizationService) normalizeSeverity(severity string) string {
	return strings.Title(strings.ToLower(severity))
}

func (n *NormalizationService) normalizeDataClassification(classification string) string {
	if classification == "" {
		return "Internal"
	}
	return classification
}

func (n *NormalizationService) normalizeCorrelationKeys(keys map[string]string) map[string]string {
	if keys == nil {
		return make(map[string]string)
	}
	return keys
}
