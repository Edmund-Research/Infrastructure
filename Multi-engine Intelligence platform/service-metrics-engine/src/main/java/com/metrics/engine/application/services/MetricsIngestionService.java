package com.metrics.engine.application.services;

import com.metrics.engine.domain.entities.MetricSample;
import com.metrics.engine.domain.entities.PipelineContext;
import com.metrics.engine.domain.repositories.MetricRepository;
import com.metrics.engine.infrastructure.persistence.TimeSeriesDB;
import com.metrics.engine.application.dto.MetricIngestionRequest;
import java.util.List;
import java.util.stream.Collectors;

public class MetricsIngestionService {
    private final MetricRepository metricRepository;
    private final TimeSeriesDB timeSeriesDB;
    private final NormalizationService normalizationService;

    public MetricsIngestionService(MetricRepository metricRepository, TimeSeriesDB timeSeriesDB, NormalizationService normalizationService) {
        this.metricRepository = metricRepository;
        this.timeSeriesDB = timeSeriesDB;
        this.normalizationService = normalizationService;
    }

    public void ingestMetrics(MetricIngestionRequest request, PipelineContext context) {
        List<MetricSample> samples = request.getMetrics().stream()
            .map(metricData -> new MetricSample(
                request.getPipelineRunId(),
                metricData.getName(),
                metricData.getValue(),
                metricData.getEnvironment(),
                metricData.getTestType()
    ))
    .collect(Collectors.toList());

    // Log ingestion
    timeSeriesDB.logIngestion(request.getPipelineRunId(), "metrics", samples.size());

    // Process each sample
    for (MetricSample sample : samples) {
            context.traces(sample);
            metricRepository.save(sample);
            normalizationService.processMetricSample(sample);
        }
    }
}
