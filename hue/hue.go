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
	TypeLight  = "light"
	TypeSesnor = "sensor"

	LabelName                 = "Name"
	LabelType                 = "Type"
	LabelModelID              = "Model_ID"
	LabelManufacturerName     = "Manufacturer_Name"
	LabelSWVersion            = "SW_Version"
	LabelUniqueID             = "Unique_ID"
	LabelStateOn              = "State_On"
	LabelStateAlert           = "State_Alert"
	LabelStateBri             = "State_Bri"
	LabelStateCT              = "State_CT"
	LabelStateReachable       = "State_Reachable"
	LabelStateSaturation      = "State_Saturation"
	LabelStateButtonEvent     = "State_Button_Event"
	LabelStateDaylight        = "State_Daylight"
	LabelStateLastUpdated     = "State_Last_Updated"
	LabelStateLastUpdatedTime = "State_Last_Updated_Time"
	LabelConfigBatery         = "Config_Battery"
	LabelConfigOn             = "Config_On"
	LabelConfigReachable      = "Config_Reachable"
)

// InitMetrics func
func (exporter *Exporter) InitMetrics() (metrics []*metric.Metric) {

	metrics = append(metrics, &metric.Metric{
		HueType:   TypeLight,
		Labels:    []string{LabelName, LabelType, LabelModelID, LabelManufacturerName, LabelSWVersion, LabelUniqueID, LabelStateOn, LabelStateAlert, LabelStateBri, LabelStateCT, LabelStateReachable, LabelStateSaturation},
		ResultKey: LabelStateOn})

	metrics = append(metrics, &metric.Metric{
		HueType:   TypeSesnor,
		Labels:    []string{LabelName, LabelType, LabelModelID, LabelManufacturerName, LabelSWVersion, LabelUniqueID, LabelStateButtonEvent, LabelStateDaylight, LabelStateLastUpdated, LabelStateLastUpdatedTime, LabelConfigBatery, LabelConfigOn, LabelConfigReachable},
		ResultKey: LabelConfigOn})

	return metrics
}

// Collect metrics
func (exporter *Exporter) Collect(metrics []*metric.Metric) (err error) {

	bridge := newBridge(exporter.BaseURL)

	err = bridge.Login(exporter.Username)
	if err != nil {
		return fmt.Errorf("[error Login()] '%v'", err)
	}

	for _, metric := range metrics {

		var err error

		switch metric.HueType {
		case TypeLight:
			err = collectLights(bridge, metric)
		case TypeSesnor:
			err = collectSensors(bridge, metric)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func collectSensors(bridge *hue.Bridge, metric *metric.Metric) (err error) {

	metric.MetricResult = nil

	sensors, err := bridge.GetAllSensors()
	if err != nil {
		return fmt.Errorf("[error GetAllSensors()] '%v'", err)
	}

	for _, sensor := range sensors {

		result := make(map[string]interface{})
		for _, label := range metric.Labels {

			switch label {
			case LabelName:
				result[LabelName] = sensor.Name
			case LabelType:
				result[LabelType] = sensor.Type
			case LabelModelID:
				result[LabelModelID] = sensor.ModelID
			case LabelManufacturerName:
				result[LabelManufacturerName] = sensor.ManufacturerName
			case LabelSWVersion:
				result[LabelSWVersion] = sensor.SWVersion
			case LabelUniqueID:
				result[LabelUniqueID] = sensor.UniqueID
			case LabelStateButtonEvent:
				result[LabelStateButtonEvent] = sensor.State.ButtonEvent
			case LabelStateDaylight:
				result[LabelStateDaylight] = sensor.State.Daylight
			case LabelStateLastUpdated:
				result[LabelStateLastUpdated] = sensor.State.LastUpdated
			case LabelStateLastUpdatedTime:
				result[LabelStateLastUpdatedTime] = sensor.State.LastUpdated.Time
			case LabelConfigBatery:
				result[LabelConfigBatery] = sensor.Config.Battery
			case LabelConfigOn:
				result[LabelConfigOn] = sensor.Config.On
			case LabelConfigReachable:
				result[LabelConfigReachable] = sensor.Config.Reachable
			}
		}

		metric.MetricResult = append(metric.MetricResult, result)
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
			case LabelName:
				result[LabelName] = light.Name
			case LabelType:
				result[LabelType] = light.Type
			case LabelModelID:
				result[LabelModelID] = light.ModelID
			case LabelManufacturerName:
				result[LabelManufacturerName] = light.ManufacturerName
			case LabelSWVersion:
				result[LabelSWVersion] = light.SWVersion
			case LabelUniqueID:
				result[LabelUniqueID] = light.UniqueID
			case LabelStateOn:
				result[LabelStateOn] = light.State.On
			case LabelStateAlert:
				result[LabelStateAlert] = light.State.Alert
			case LabelStateBri:
				result[LabelStateBri] = light.State.Bri
			case LabelStateCT:
				result[LabelStateCT] = light.State.CT
			case LabelStateReachable:
				result[LabelStateReachable] = light.State.Reachable
			case LabelStateSaturation:
				result[LabelStateSaturation] = light.State.Saturation
			}
		}

		metric.MetricResult = append(metric.MetricResult, result)
	}

	return nil
}

func newBridge(ipAddr string) *hue.Bridge {
	bridge, err := hue.NewBridge(ipAddr)
	if err != nil {
		log.Fatalf("Error connecting to Hue bridge with '%v': '%v'\n", ipAddr, err)
	}
	return bridge
}
