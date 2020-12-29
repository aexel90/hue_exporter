package metric

import (
	"github.com/prometheus/client_golang/prometheus"
)

// PrometheusResult struct
type PrometheusResult struct {
	PromDesc      *prometheus.Desc
	PromValueType prometheus.ValueType
	Value         float64
	LabelValues   []string
}

// Metric struct
type Metric struct {
	HueType   string   `json:"type"`
	ResultKey string   `json:"resultKey"`
	FqName    string   `json:"fqName"`
	Help      string   `json:"help"`
	Labels    []string `json:"labels"`

	MetricResult []map[string]interface{}

	PromType   prometheus.ValueType
	PromDesc   *prometheus.Desc
	PromResult []*PrometheusResult
}

// MetricsFile struct
type MetricsFile struct {
	Metrics []*Metric `json:"metrics"`
}
