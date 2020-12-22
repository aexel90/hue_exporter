package hue

import (
	"fmt"
	"log"

	"github.com/aexel90/hue_exporter/metric"
	hue "github.com/collinux/gohue"
)

// Exporter data
type Exporter struct {
	BaseURL  string
	Username string
}

const (
	TypeLight = "light"
	TypeOther = "???"

	LightLabelName             = "Name"
	LightLabelType             = "Type"
	LightLabelModelID          = "Model_ID"
	LightLabelManufacturerName = "Manufacturer_Name"
	LightLabelSWVersion        = "SW_Version"
	LightLabelUniqueID         = "Unique_ID"
	LightLabelStateOn          = "State_On"
	LightLabelStateAlert       = "State_Alert"
	LightLabelStateBri         = "State_Bri"
	LightLabelStateCT          = "State_CT"
	LightLabelStateReachable   = "State_Reachable"
	LightLabelStateSaturation  = "State_Saturation"
)

// InitMetrics func
func (exporter *Exporter) InitMetrics() (metrics []*metric.Metric) {

	metrics = append(metrics, &metric.Metric{
		HueType:   TypeLight,
		Labels:    []string{LightLabelName, LightLabelType, LightLabelModelID, LightLabelManufacturerName, LightLabelSWVersion, LightLabelUniqueID, LightLabelStateOn, LightLabelStateAlert, LightLabelStateBri, LightLabelStateCT, LightLabelStateReachable, LightLabelStateSaturation},
		ResultKey: LightLabelStateOn})

	return metrics
}

// Collect metrics
func (exporter *Exporter) Collect(metrics []*metric.Metric) (err error) {

	bridge := newBridge(exporter.BaseURL)

	err = bridge.Login(exporter.Username)
	if err != nil {
		return fmt.Errorf("[error login] '%v'", err)
	}

	for _, metric := range metrics {

		var err error

		switch metric.HueType {
		case TypeLight:
			err = collectLights(bridge, metric)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func collectLights(bridge *hue.Bridge, metric *metric.Metric) (err error) {

	metric.MetricResult = nil

	lights, err := bridge.GetAllLights()
	if err != nil {
		return fmt.Errorf("[error GetAllLights()] '%v'", err)
	}

	for _, light := range lights {

		result := make(map[string]interface{})
		for _, label := range metric.Labels {

			switch label {
			case LightLabelName:
				result[LightLabelName] = light.Name
			case LightLabelType:
				result[LightLabelType] = light.Type
			case LightLabelModelID:
				result[LightLabelModelID] = light.ModelID
			case LightLabelManufacturerName:
				result[LightLabelManufacturerName] = light.ManufacturerName
			case LightLabelSWVersion:
				result[LightLabelSWVersion] = light.SWVersion
			case LightLabelUniqueID:
				result[LightLabelUniqueID] = light.UniqueID
			case LightLabelStateOn:
				result[LightLabelStateOn] = light.State.On
			case LightLabelStateAlert:
				result[LightLabelStateAlert] = light.State.Alert
			case LightLabelStateBri:
				result[LightLabelStateBri] = light.State.Bri
			case LightLabelStateCT:
				result[LightLabelStateCT] = light.State.CT
			case LightLabelStateReachable:
				result[LightLabelStateReachable] = light.State.Reachable
			case LightLabelStateSaturation:
				result[LightLabelStateSaturation] = light.State.Saturation
			}
		}

		metric.MetricResult = append(metric.MetricResult, result)
	}

	return nil
}

func newBridge(ipAddr string) *hue.Bridge {
	bridge, err := hue.NewBridge(ipAddr)
	if err != nil {
		log.Fatalf("Error connecting to Hue bridge at '%v': '%v'\n", ipAddr, err)
	}
	return bridge
}
