package com.metrics.engine.application.services;

import com.metrics.engine.domain.entities.MetricSample;
import java.util.Map;
import java.util.HashMap;

public class NormalizationService {
    private final BaselineService baselineService;
    private final CorrelationService correlationService;

    public NormalizationService(
        BaselineService baselineService,
        CorrelationService correlationService) {
            this.baselineService = baselineService;
            this.correlationService = correlationService;
    }

    public void processMetricSample(MetricSample sample) {
    // Enrich with additional metadata
    enrichSample(sample);

    // Register for correlation analysis
    correlationService.registerMetricForCorrelation(sample.getPipelineRunId(), sample);

    // Analyze correlations
    correlationService.analyzeCorrelations(sample.getPipelineRunId());

    // Evaluate against baselines
    baselineService.evaluateMetricAgainstBaselines(sample);

    System.out.println("NormalizationService: Processed sample " + sample.getMetricName());
    }

    private void enrichSample(MetricSample sample) {
        Map<String, String> enrichedLabels = new HashMap<>(sample.getLabels());
        enrichedLabels.put("processing_timestamp", java.time.Instant.now().toString());
        enrichedLabels.put("processor", "NormalizationService");
        sample.setLabels(enrichedLabels);
    }

}
