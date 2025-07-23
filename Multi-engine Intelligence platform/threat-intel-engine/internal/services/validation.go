package services

import (
	"fmt"

	"github.com/google/uuid"

	"threat-intel-engine/internal/models"
	"threat-intel-engine/internal/repository"
)

type ValidationService struct {
	repo *repository.FindingRepository
}

func NewValidationService(repo *repository.FindingRepository) *ValidationService {
	return &ValidationService{repo: repo}
}

func (v *ValidationService) ValidateRequest(finding models.FindingRequest) error {
	if finding.PipelineRunID == uuid.Nil {
		return fmt.Errorf("pipelineRunId is required")
	}
	if finding.Source == "" {
		return fmt.Errorf("source is required")
	}
	if !v.isValidSeverity(finding.Severity) {
		return fmt.Errorf("invalid severity level")
	}
	if finding.RuleID == "" {
		return fmt.Errorf("ruleId is required")
	}
	if finding.Description == "" {
		return fmt.Errorf("description is required")
	}
	return nil
}

func (v *ValidationService) CheckDuplicate(finding models.FindingRequest) (bool, error) {
	return v.repo.CheckDuplicate(finding)
}

func (v *ValidationService) isValidSeverity(severity string) bool {
	validSeverities := map[string]bool{
		"Low":      true,
		"Medium":   true,
		"High":     true,
		"Critical": true,
	}
	return validSeverities[severity]
}
