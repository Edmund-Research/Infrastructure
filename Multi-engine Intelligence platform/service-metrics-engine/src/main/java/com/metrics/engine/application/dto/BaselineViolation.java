package com.metrics.engine.application.dto;

import com.metrics.engine.domain.entities.PerformanceBaseline;
import java.util.List;

public class BaselineViolation {
    private final String metricName;
    private final float currentValue;
    private final List<PerformanceBaseline> violatedBaselines;
    private final String environment;
    private final String severity;

    public BaselineViolation(String metricName, float currentValue, List<PerformanceBaseline> violatedBaselines,
    String environment) {
        this.metricName = metricName;
        this.currentValue = currentValue;
        this.violatedBaselines = violatedBaselines;
        this.environment = environment;
        this.severity = calculateSeverity(violatedBaselines);
    }

    private String calculateSeverity(List<PerformanceBaseline> violations) {
        if (violations.isEmpty()) return "INFO";
        if (violations.size() == 1) return "WARNING";
        return "CRITICAL";
    }

    // Getters
    public String getMetricName() { return metricName; }
    public float getCurrentValue() { return currentValue; }
    public List<PerformanceBaseline> getViolatedBaselines() { return violatedBaselines; }
    public String getEnvironment() { return environment; }
    public String getSeverity() { return severity; }
}
