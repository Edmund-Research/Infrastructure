package models

import (
	"time"

	"github.com/google/uuid"
)

type Finding struct {
	ID                 uuid.UUID         `json:"id" bson:"_id"`
	PipelineRunID      uuid.UUID         `json:"pipelineRunId" bson:"pipelineRunId"`
	Source             string            `json:"source" bson:"source"`
	Severity           string            `json:"severity" bson:"severity"`
	Timestamp          time.Time         `json:"timestamp" bson:"timestamp"`
	FilePath           string            `json:"filePath,omitempty" bson:"filePath,omitempty"`
	Endpoint           string            `json:"endpoint,omitempty" bson:"endpoint,omitempty"`
	RuleID             string            `json:"ruleId" bson:"ruleId"`
	Description        string            `json:"description" bson:"description"`
	Evidence           string            `json:"evidence" bson:"evidence"`
	DataClassification string            `json:"dataClassification" bson:"dataClassification"`
	CorrelationKeys    map[string]string `json:"correlationKeys" bson:"correlationKeys"`
}

type Enrichment struct {
	ID               uuid.UUID `json:"id" bson:"_id"`
	FindingID        uuid.UUID `json:"findingId" bson:"findingId"`
	PipelineRunID    uuid.UUID `json:"pipelineRunId" bson:"pipelineRunId"`
	CVEList          []string  `json:"cveList" bson:"cveList"`
	CWEList          []string  `json:"cweList" bson:"cweList"`
	ExploitScore     float64   `json:"exploitabilityScore" bson:"exploitabilityScore"`
	EnrichedAt       time.Time `json:"enrichedAt" bson:"enrichedAt"`
	EnrichmentSource string    `json:"enrichmentSource" bson:"enrichmentSource"`
}

type ThreatCorrelation struct {
	ID                uuid.UUID `json:"id" bson:"_id"`
	FindingID         uuid.UUID `json:"findingId" bson:"findingId"`
	CorrelationType   string    `json:"correlationType" bson:"correlationType"`
	RelatedEntityID   string    `json:"relatedEntityId" bson:"relatedEntityId"`
	RelatedEntityType string    `json:"relatedEntityType" bson:"relatedEntityType"`
	Confidence        float64   `json:"confidence" bson:"confidence"`
}

type FindingRequest struct {
	PipelineRunID      uuid.UUID         `json:"pipelineRunId" binding:"required"`
	Source             string            `json:"source" binding:"required,oneof=SAST DAST Dependency"`
	Severity           string            `json:"severity" binding:"required,oneof=Low Medium High Critical"`
	FilePath           string            `json:"filePath,omitempty"`
	Endpoint           string            `json:"endpoint,omitempty"`
	RuleID             string            `json:"ruleId" binding:"required"`
	Description        string            `json:"description" binding:"required"`
	Evidence           string            `json:"evidence"`
	DataClassification string            `json:"dataClassification"`
	CorrelationKeys    map[string]string `json:"correlationKeys"`
}

type QueryParams struct {
	Since    string `form:"since"`
	Severity string `form:"severity"`
	Source   string `form:"source"`
	Limit    int    `form:"limit"`
	Offset   int    `form:"offset"`
}
