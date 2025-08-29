package com.metrics.engine.application.services;

import com.metrics.engine.domain.entities.MetricSample;
import com.metrics.engine.domain.entities.MetricCorrelation;
import com.metrics.engine.domain.services.CorrelationAnalysisService;
import com.metrics.engine.domain.repositories.CorrelationRepository;
import com.metrics.engine.domain.repositories.MetricRepository;
import java.util.List;
import java.util.UUID;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

public class CorrelationService {
    private final CorrelationRepository correlationRepository;
    private final MetricRepository metricRepository;
    private final CorrelationAnalysisService analysisService;
    private final Map<UUID, List<MetricSample>> pipelineMetrics;

    public CorrelationService(
        CorrelationRepository correlationRepository,
        MetricRepository metricRepository,
        CorrelationAnalysisService analysisService) {
            this.correlationRepository = correlationRepository;
            this.metricRepository = metricRepository;
            this.analysisService = analysisService;
            this.pipelineMetrics = new ConcurrentHashMap<>();
    }

    public void registerMetricForCorrelation(UUID pipelineRunId, MetricSample sample) {
        pipelineMetrics.computeIfAbsent(pipelineRunId, k ->
        metricRepository.findByPipelineRunId(pipelineRunId));

        System.out.println("CorrelationService: Registered metric " + sample.getMetricName() +
        " for pipeline " + pipelineRunId);
    }

    public List<MetricCorrelation> analyzeCorrelations(UUID pipelineRunId) {
        List<MetricSample> samples = metricRepository.findByPipelineRunId(pipelineRunId);

        List<MetricCorrelation> temporalCorrelations =
        analysisService.analyzeTemporalCorrelations(samples);

        List<MetricCorrelation> valueCorrelations =
        analysisService.analyzeValueCorrelations(samples);

        // Save correlations
        temporalCorrelations.forEach(correlationRepository::save);
        valueCorrelations.forEach(correlationRepository::save);

        System.out.println("CorrelationService: Found " + temporalCorrelations.size() +
        " temporal correlations for pipeline " + pipelineRunId);

        return temporalCorrelations;
    }

    public List<MetricCorrelation> getAllCorrelations() {
        return correlationRepository.findAll();
    }
}
