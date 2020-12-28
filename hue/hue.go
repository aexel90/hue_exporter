package hue

import (
	"fmt"
	"log"

	"github.com/aexel90/hue_exporter/metric"
	hue "github.com/shamx9ir/gohue" // github.com/collinux/gohue
)

// Exporter data
type Exporter struct {
	BaseURL  string
	Username string
}

const (
	TypeLight  = "light"
	TypeSensor = "sensor"

	LabelName                 = "Name"
	LabelType                 = "Type"
	LabelIndex                = "Index"
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
	LabelStateTemperature     = "State_Temperature"
	LabelConfigBattery        = "Config_Battery"
	LabelConfigOn             = "Config_On"
	LabelConfigReachable      = "Config_Reachable"
)

// InitMetrics func
func (exporter *Exporter) InitMetrics() (metrics []*metric.Metric) {

	metrics = append(metrics, &metric.Metric{
		HueType: TypeLight,
		FqName:  "hue_light_info",
		Help:    "Non-numeric data, value is always 1",
		Labels:  []string{LabelName, LabelIndex, LabelType, LabelModelID, LabelManufacturerName, LabelSWVersion, LabelUniqueID, LabelStateOn, LabelStateAlert, LabelStateBri, LabelStateCT, LabelStateReachable, LabelStateSaturation},
	})

	metrics = append(metrics, &metric.Metric{
		HueType: TypeSensor,
		FqName:  "hue_sensor_info",
		Help:    "Non-numeric data, value is always 1",
		Labels:  []string{LabelName, LabelIndex, LabelType, LabelModelID, LabelManufacturerName, LabelSWVersion, LabelUniqueID, LabelStateButtonEvent, LabelStateDaylight, LabelStateLastUpdated, LabelStateLastUpdatedTime, LabelStateTemperature, LabelConfigBattery, LabelConfigOn, LabelConfigReachable},
	})

	metrics = append(metrics, &metric.Metric{
		HueType:   TypeSensor,
		FqName:    "hue_sensor_battery",
		Help:      "battery level percentage",
		Labels:    []string{LabelName, LabelIndex, LabelType, LabelModelID, LabelManufacturerName, LabelSWVersion, LabelUniqueID},
		ResultKey: LabelConfigBattery,
	})

	metrics = append(metrics, &metric.Metric{
		HueType:   TypeSensor,
		FqName:    "hue_sensor_temperature",
		Help:      "temperature level celsius degree",
		Labels:    []string{LabelName, LabelIndex, LabelType, LabelModelID, LabelManufacturerName, LabelSWVersion, LabelUniqueID},
		ResultKey: LabelStateTemperature,
	})

	return metrics
}

// Collect metrics
func (exporter *Exporter) Collect(metrics []*metric.Metric) (err error) {

	bridge := newBridge(exporter.BaseURL)

	err = bridge.Login(exporter.Username)
	if err != nil {
		return fmt.Errorf("[error Login()] '%v'", err)
	}

	sensorData, err := collectSensors(bridge)
	if err != nil {
		return err
	}

	lightData, err := collectLights(bridge)
	if err != nil {
		return err
	}

	for _, metric := range metrics {
		switch metric.HueType {
		case TypeLight:
			metric.MetricResult = lightData
		case TypeSensor:
			metric.MetricResult = sensorData
		}
	}

	return nil
}

func collectSensors(bridge *hue.Bridge) (sensorData []map[string]interface{}, err error) {

	sensors, err := bridge.GetAllSensors()
	if err != nil {
		return nil, fmt.Errorf("[error GetAllSensors()] '%v'", err)
	}
	for _, sensor := range sensors {
		result := make(map[string]interface{})
		result[LabelName] = sensor.Name
		result[LabelIndex] = sensor.Index
		result[LabelType] = sensor.Type
		result[LabelModelID] = sensor.ModelID
		result[LabelManufacturerName] = sensor.ManufacturerName
		result[LabelSWVersion] = sensor.SWVersion
		result[LabelUniqueID] = sensor.UniqueID
		result[LabelStateButtonEvent] = float64(sensor.State.ButtonEvent)
		result[LabelStateDaylight] = sensor.State.Daylight
		result[LabelStateLastUpdated] = sensor.State.LastUpdated
		result[LabelStateLastUpdatedTime] = sensor.State.LastUpdated.Time
		result[LabelConfigBattery] = sensor.Config.Battery
		result[LabelConfigOn] = sensor.Config.On
		result[LabelConfigReachable] = sensor.Config.Reachable

		if sensor.Type == "ZLLTemperature" {
			result[LabelStateTemperature] = float64(sensor.State.Temperature)
		}

		sensorData = append(sensorData, result)
	}
	return sensorData, nil
}

func collectLights(bridge *hue.Bridge) (lightData []map[string]interface{}, err error) {

	lights, err := bridge.GetAllLights()
	if err != nil {
		return nil, fmt.Errorf("[error GetAllLights()] '%v'", err)
	}

	for _, light := range lights {

		result := make(map[string]interface{})
		result[LabelName] = light.Name
		result[LabelIndex] = light.Index
		result[LabelType] = light.Type
		result[LabelModelID] = light.ModelID
		result[LabelManufacturerName] = light.ManufacturerName
		result[LabelSWVersion] = light.SWVersion
		result[LabelUniqueID] = light.UniqueID
		result[LabelStateOn] = light.State.On
		result[LabelStateAlert] = light.State.Alert
		result[LabelStateBri] = light.State.Bri
		result[LabelStateCT] = light.State.CT
		result[LabelStateReachable] = light.State.Reachable
		result[LabelStateSaturation] = light.State.Saturation

		lightData = append(lightData, result)
	}
	return lightData, nil
}

func newBridge(ipAddr string) *hue.Bridge {
	bridge, err := hue.NewBridge(ipAddr)
	if err != nil {
		log.Fatalf("Error connecting to Hue bridge with '%v': '%v'\n", ipAddr, err)
	}
	return bridge
}
