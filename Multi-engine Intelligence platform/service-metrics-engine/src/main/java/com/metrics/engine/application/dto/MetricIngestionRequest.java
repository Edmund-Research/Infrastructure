package com.metrics.engine.application.dto;

import java.util.List;
import java.util.UUID;

public class MetricIngestionRequest {
    private final List<MetricData> metrics;
    private final UUID pipelineRunId;
    private final String testRunId;

    public MetricIngestionRequest(List<MetricData> metrics, UUID pipelineRunId, String testRunId) {
    this.metrics = metrics;
    this.pipelineRunId = pipelineRunId;
    this.testRunId = testRunId;
    }

    public static class MetricData {
    private final String name;
    private final float value;
    private final String environment;
    private final String testType;

    public MetricData(String name, float value, String environment, String testType) {
    this.name = name;
    this.value = value;
    this.environment = environment;
    this.testType = testType;
    }

    // Getters
    public String getName() { return name; }
    public float getValue() { return value; }
    public String getEnvironment() { return environment; }
    public String getTestType() { return testType; }
    }

    // Getters
    public List<MetricData> getMetrics() { return metrics; }
    public UUID getPipelineRunId() { return pipelineRunId; }
    public String getTestRunId() { return testRunId; }
}
