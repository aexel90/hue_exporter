package collector

import (
	"fmt"
	"strings"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/aexel90/hue_exporter/hue"
	"github.com/aexel90/hue_exporter/metric"
)

// Collector instance
type Collector struct {
	exporter *hue.Exporter
	metrics  []*metric.Metric
}

// NewHueCollector initialization
func NewHueCollector(URL string, username string) (*Collector, error) {

	hueExporter := hue.Exporter{
		BaseURL:  URL,
		Username: username,
	}

	return &Collector{&hueExporter, nil}, nil
}

// Describe for prometheus
func (collector *Collector) Describe(ch chan<- *prometheus.Desc) {

	collector.metrics = collector.exporter.InitMetrics()
	collector.initDescAndType()
}

// Collect for prometheus
func (collector *Collector) Collect(ch chan<- prometheus.Metric) {

	err := collector.collect()
	if err != nil {
		fmt.Println("Error: ", err)
	}

	for _, m := range collector.metrics {
		for _, promResult := range m.PromResult {
			ch <- prometheus.MustNewConstMetric(promResult.PromDesc, promResult.PromValueType, promResult.Value, promResult.LabelValues...)
		}
	}
}

//Test collector metrics
func (collector *Collector) Test() {

	collector.metrics = collector.exporter.InitMetrics()
	collector.initDescAndType()

	err := collector.collect()
	if err != nil {
		fmt.Println("Error: ", err)
	}

	err = collector.printResult()
}

func (collector *Collector) printResult() error {

	for _, m := range collector.metrics {
		fmt.Printf("Metric: %v\n", m.FqName)
		fmt.Printf(" - Exporter Result:\n")

		for i, result := range m.MetricResult {
			fmt.Printf("   - Exporter Result %v:\n", i)
			for key, value := range result {
				fmt.Printf("     - %s=\"%v\"\n", key, value)
			}
		}

		for _, promResult := range m.PromResult {

			fmt.Printf("   - prom desc: %v\n", promResult.PromDesc)
			fmt.Printf("     - prom metric type: %v\n", promResult.PromValueType)
			fmt.Printf("     - prom metric value: %v\n", uint64(promResult.Value))
			fmt.Printf("     - prom label values: %v\n", promResult.LabelValues)
		}
	}

	return nil
}

func (collector *Collector) collect() (err error) {

	err = collector.exporter.Collect(collector.metrics)
	if err != nil {
		return err
	}

	err = collector.getResult()
	if err != nil {
		return err
	}
	return nil
}

func (collector *Collector) getResult() (err error) {

	for _, m := range collector.metrics {
		m.PromResult = nil
		for _, metricResult := range m.MetricResult {

			labelValues, err := getLabelValues(m.Labels, metricResult)
			if err != nil {
				return err
			}
			resultValue, err := getResultValue(m.ResultKey, metricResult)
			if err != nil {
				return err
			}

			result := metric.PrometheusResult{PromDesc: m.PromDesc, PromValueType: m.PromType, Value: resultValue, LabelValues: labelValues}
			m.PromResult = append(m.PromResult, &result)
		}
	}
	return nil
}

func (collector *Collector) initDescAndType() {

	for _, metric := range collector.metrics {

		var help string

		switch metric.HueType {
		case hue.TypeLight:
			metric.FqName = "hue_light"
			help = "status of lights registered at hue bridge"
			metric.PromType = prometheus.GaugeValue
		}

		labels := []string{}
		for _, label := range metric.Labels {
			labels = append(labels, strings.ToLower(label))
		}

		metric.PromDesc = prometheus.NewDesc(metric.FqName, help, labels, nil)
	}
}

func getResultValue(resultKey string, result map[string]interface{}) (float64, error) {

	value := result[resultKey]
	var floatValue float64

	switch tval := value.(type) {
	case float64:
		floatValue = tval
	case int:
		floatValue = float64(tval)
	case uint64:
		floatValue = float64(tval)
	case bool:
		if tval {
			floatValue = 1
		} else {
			floatValue = 0
		}
	default:
		return 0, fmt.Errorf("[getResultValue] %v in %v - unknown type: %T", resultKey, result, value)
	}
	return floatValue, nil
}

func getLabelValues(labelNames []string, result map[string]interface{}) ([]string, error) {

	labelValues := []string{}
	for _, labelname := range labelNames {
		labelValue := fmt.Sprintf("%v", result[labelname])
		labelValue = strings.ToLower(labelValue)
		labelValues = append(labelValues, labelValue)
	}
	return labelValues, nil
}
