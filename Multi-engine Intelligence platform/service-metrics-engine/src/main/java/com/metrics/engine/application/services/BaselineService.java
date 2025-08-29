package com.metrics.engine.application.services;

import com.metrics.engine.domain.entities.MetricSample;
import com.metrics.engine.domain.entities.PerformanceBaseline;
import com.metrics.engine.domain.services.MetricValidationService;
import com.metrics.engine.domain.services.BaselineEvaluationService;
import com.metrics.engine.domain.repositories.BaselineRepository;
import com.metrics.engine.infrastructure.messaging.AlertManager;
import com.metrics.engine.application.dto.BaselineViolation;
import java.time.Instant;
import java.util.List;

public class BaselineService {
    private final BaselineRepository baselineRepository;
    private final MetricValidationService validationService;
    private final BaselineEvaluationService evaluationService;
    private final AlertManager alertManager;

    public BaselineService(
        BaselineRepository baselineRepository,
        MetricValidationService validationService,
        BaselineEvaluationService evaluationService,
        AlertManager alertManager) {
            this.baselineRepository = baselineRepository;
            this.validationService = validationService;
            this.evaluationService = evaluationService;
            this.alertManager = alertManager;
    }

    public void evaluateMetricAgainstBaselines(MetricSample sample) {
    List<PerformanceBaseline> baselines = baselineRepository
    .findByMetricNameAndEnvironment(sample.getMetricName(), sample.getEnvironment());

    BaselineEvaluationService.EvaluationResult result =
    evaluationService.evaluate(sample, baselines);

    if (!result.isPassed()) {
        handleBaselineViolations(sample, result.getViolatedBaselines());
    }
    }

    private void handleBaselineViolations(MetricSample sample, List<PerformanceBaseline> violations) {
        BaselineViolation alert = new BaselineViolation(
            sample.getMetricName(),
            sample.getValue(),
            violations,
            sample.getEnvironment()
        );

        // Send alerts for each violation
        violations.forEach(baseline -> {
            alertManager.sendBaselineViolationAlert(sample, baseline);
            alertManager.logBaselineViolation(sample.getPipelineRunId(),
            sample.getMetricName(), sample.getValue());
        });
    }

    public void renderDashboard(Instant timestamp) {
        System.out.println("BaselineService: Rendering dashboard at " + timestamp);
    }
}
