package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"threat-intel-engine/internal/models"
)

type FindingRepository struct {
	findingsCollection    *mongo.Collection
	enrichmentCollection  *mongo.Collection
	correlationCollection *mongo.Collection
}

func NewFindingRepository(db *mongo.Database) *FindingRepository {
	return &FindingRepository{
		findingsCollection:    db.Collection("findings"),
		enrichmentCollection:  db.Collection("enrichments"),
		correlationCollection: db.Collection("correlations"),
	}
}

func (r *FindingRepository) Create(finding *models.Finding) error {
	_, err := r.findingsCollection.InsertOne(context.Background(), finding)
	return err
}

func (r *FindingRepository) FindByQuery(params models.QueryParams) ([]models.Finding, error) {
	filter := bson.M{}

	if params.Since != "" {
		sinceTime, err := time.Parse(time.RFC3339, params.Since)
		if err != nil {
			return nil, err
		}
		filter["timestamp"] = bson.M{"$gte": sinceTime}
	}

	if params.Severity != "" {
		filter["severity"] = params.Severity
	}

	if params.Source != "" {
		filter["source"] = params.Source
	}

	if params.Limit <= 0 {
		params.Limit = 50
	}
	if params.Offset < 0 {
		params.Offset = 0
	}

	opts := options.Find().
		SetLimit(int64(params.Limit)).
		SetSkip(int64(params.Offset)).
		SetSort(bson.D{{Key: "timestamp", Value: -1}})

	cursor, err := r.findingsCollection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var findings []models.Finding
	if err = cursor.All(context.Background(), &findings); err != nil {
		return nil, err
	}

	return findings, nil
}

func (r *FindingRepository) CheckDuplicate(req models.FindingRequest) (bool, error) {
	filter := bson.M{
		"pipelineRunId": req.PipelineRunID,
		"source":        req.Source,
		"ruleId":        req.RuleID,
		"filePath":      req.FilePath,
		"endpoint":      req.Endpoint,
	}

	count, err := r.findingsCollection.CountDocuments(context.Background(), filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *FindingRepository) CreateEnrichment(enrichment *models.Enrichment) error {
	_, err := r.enrichmentCollection.InsertOne(context.Background(), enrichment)
	return err
}

func (r *FindingRepository) GetEnrichmentsByFindingID(findingID uuid.UUID) ([]models.Enrichment, error) {
	filter := bson.M{"findingId": findingID}
	cursor, err := r.enrichmentCollection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var enrichments []models.Enrichment
	if err = cursor.All(context.Background(), &enrichments); err != nil {
		return nil, err
	}

	return enrichments, nil
}

func (r *FindingRepository) CreateCorrelation(correlation *models.ThreatCorrelation) error {
	_, err := r.correlationCollection.InsertOne(context.Background(), correlation)
	return err
}
