package com.metrics.engine.application.services;

import org.springframework.boot.ApplicationRunner;
import org.springframework.stereotype.Component;
import com.metrics.engine.interfaces.collectors.CollectorAgent;
import com.metrics.engine.domain.entities.PipelineContext;
import com.metrics.engine.application.dto.MetricIngestionRequest;
import com.metrics.engine.domain.repositories.*;
import com.metrics.engine.infrastructure.external.AuditLogger;
import com.metrics.engine.domain.entities.PerformanceBaseline;
import com.metrics.engine.domain.entities.MetricCorrelation;
import com.metrics.engine.domain.entities.MetricSample;
import java.util.Arrays;

@Component
public class MetricsDemoRunner implements ApplicationRunner {
    private final CollectorAgent collectorAgent;
    private final AuditLogger auditLogger;
    private final MetricRepository metricRepository;
    private final CorrelationRepository correlationRepository;

    public MetricsDemoRunner(
        CollectorAgent collectorAgent,
        AuditLogger auditLogger,
        MetricRepository metricRepository,
        CorrelationRepository correlationRepository
    )  {
        this.collectorAgent = collectorAgent;
        this.auditLogger = auditLogger;
        this.metricRepository = metricRepository;
        this.correlationRepository = correlationRepository;
    }

    @Override
    public void run(org.springframework.boot.ApplicationArguments args) throws Exception {
        PipelineContext context = new PipelineContext("production");

        MetricIngestionRequest normalRequest = new MetricIngestionRequest(
            Arrays.asList(
                new MetricIngestionRequest.MetricData("response_time", 95.5f, "production", "load_test"),
                new MetricIngestionRequest.MetricData("error_rate", 0.008f, "production", "load_test"),
                new MetricIngestionRequest.MetricData("cpu_usage", 68.2f, "production", "load_test")
            ),
            context.getRunId(),
            "test-run-1"
        );

        collectorAgent.remoteWrite(normalRequest, context);
        Thread.sleep(1000);

        MetricIngestionRequest violatingRequest = new MetricIngestionRequest(
            Arrays.asList(
                new MetricIngestionRequest.MetricData("response_time", 150.0f, "production", "load_test"),
                new MetricIngestionRequest.MetricData("error_rate", 0.02f, "production", "load_test"),
                new MetricIngestionRequest.MetricData("cpu_usage", 95.0f, "production", "load_test")
            ),
            context.getRunId(),
            "test-run-2"
        );

        collectorAgent.remoteWrite(violatingRequest, context);

        System.out.println("\n=== AUDIT LOG ===");
        auditLogger.getAuditLog().forEach(System.out::println);

        System.out.println("\n=== STORED SAMPLES ===");
        metricRepository.findByPipelineRunId(context.getRunId())
            .forEach(sample -> System.out.println(sample.getMetricName() + ": " + sample.getValue()));

        System.out.println("\n=== CORRELATIONS ===");
        correlationRepository.findAll()
            .forEach(correlation -> System.out.println("Correlation: " +
            correlation.getCorrelationType() + " (confidence: " + correlation.getConfidence() + ")"));
    }
}
