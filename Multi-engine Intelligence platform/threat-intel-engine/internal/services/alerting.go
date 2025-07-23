package services

import (
	"encoding/json"
	"log"

	"threat-intel-engine/internal/models"
)

type AlertingService struct {
	soarWebhookURL string
}

func NewAlertingService(webhookURL string) *AlertingService {
	return &AlertingService{
		soarWebhookURL: webhookURL,
	}
}

func (a *AlertingService) ProcessAlert(finding models.Finding) {
	if finding.Severity == "High" || finding.Severity == "Critical" {
		go func() {
			alert := map[string]interface{}{
				"id":          finding.ID,
				"severity":    finding.Severity,
				"source":      finding.Source,
				"description": finding.Description,
				"timestamp":   finding.Timestamp,
			}

			alertJSON, _ := json.Marshal(alert)
			log.Printf("Sending SOAR alert: %s", string(alertJSON))

			// Implement HTTP Call to SOAR
		}()
	}
}
