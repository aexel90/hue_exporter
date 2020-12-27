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
	HueType      string
	Labels       []string
	MetricResult []map[string]interface{}
	ResultKey    string
	FqName       string

	PromType   prometheus.ValueType
	PromDesc   *prometheus.Desc
	PromResult []*PrometheusResult
}
