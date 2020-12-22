package metric

import (
	"github.com/prometheus/client_golang/prometheus"
)

type PrometheusResult struct {
	PromDesc      *prometheus.Desc
	PromValueType prometheus.ValueType
	Value         float64
	LabelValues   []string
}

// Metric struct
type Metric struct {
	// PromDesc       PromDesc   `json:"promDesc"`
	// PromType       string     `json:"promType"`
	// ResultKey      string     `json:"resultKey"`
	// OkValue        string     `json:"okValue"`
	// ResultPath     string     `json:"resultPath"`
	// Page           string     `json:"page"`
	// Service        string     `json:"service"`
	// Action         string     `json:"action"`
	// ActionArgument *ActionArg `json:"actionArgument"`

	// Desc        *prometheus.Desc

	// Value       float64
	// labelValues []string

	HueType      string
	Labels       []string
	MetricResult []map[string]interface{} //filled during collect
	ResultKey    string
	FqName       string

	PromType   prometheus.ValueType
	PromDesc   *prometheus.Desc
	PromResult []*PrometheusResult
}
