package monitoring

import (
    "github.com/prometheus/client_golang/prometheus"
)

type PrometheusMetrics struct {
    incidentGauge   prometheus.Gauge
    clusterGauge    prometheus.Gauge
    eventCounter    prometheus.Counter
    scrapeHistogram prometheus.Histogram
}

func NewPrometheusMetrics() *PrometheusMetrics {
    incidentGauge := prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "kubewatch_incident_count",
        Help: "Current number of active incidents detected by KubeWatch AI.",
    })
    clusterGauge := prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "kubewatch_cluster_resources",
        Help: "Resource count for cluster objects monitored by KubeWatch AI.",
    })
    eventCounter := prometheus.NewCounter(prometheus.CounterOpts{
        Name: "kubewatch_incident_events_total",
        Help: "Total number of incident events emitted by KubeWatch AI.",
    })
    scrapeHistogram := prometheus.NewHistogram(prometheus.HistogramOpts{
        Name:    "kubewatch_scrape_duration_seconds",
        Help:    "Duration of Kubernetes cluster scans.",
        Buckets: prometheus.DefBuckets,
    })
    prometheus.MustRegister(incidentGauge, clusterGauge, eventCounter, scrapeHistogram)
    return &PrometheusMetrics{incidentGauge: incidentGauge, clusterGauge: clusterGauge, eventCounter: eventCounter, scrapeHistogram: scrapeHistogram}
}

func (m *PrometheusMetrics) UpdateIncidentCounts(count int) {
    m.incidentGauge.Set(float64(count))
    m.eventCounter.Add(1)
}

func (m *PrometheusMetrics) UpdateOverviewMetrics(nodes, pods, deployments, services int) {
    total := nodes + pods + deployments + services
    m.clusterGauge.Set(float64(total))
}

func (m *PrometheusMetrics) ObserveScrapeDuration(seconds float64) {
    m.scrapeHistogram.Observe(seconds)
}
