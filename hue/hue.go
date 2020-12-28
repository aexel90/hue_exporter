package hue

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/aexel90/hue_exporter/metric"
	hueAPI "github.com/amimof/huego"
)

// Exporter data
type Exporter struct {
	BaseURL  string
	Username string
}

type collectEntry struct {
	Type   string                 `json:"type"`
	Result map[string]interface{} `json:"result"`
}

const (
	TypeLight  = "light"
	TypeSensor = "sensor"

	LabelName             = "name"
	LabelType             = "type"
	LabelID               = "id"
	LabelModelID          = "model_id"
	LabelManufacturerName = "manufacturer_name"
	LabelSwVersion        = "sw_version"
	LabelSwConfigID       = "sw_config_id"
	LabelUniqueID         = "unique_id"

	LabelStateOn          = "state_on"
	LabelStateAlert       = "state_alert"
	LabelStateBri         = "state_bri"
	LabelStateColorMode   = "state_color_mode"
	LabelStateCT          = "state_ct"
	LabelStateReachable   = "state_reachable"
	LabelStateSaturation  = "state_saturation"
	LabelStateButtonEvent = "state_buttonevent"
	LabelStateDaylight    = "state_daylight"
	LabelStateLastUpdated = "state_lastupdated"
	LabelStateTemperature = "state_temperature"
	LabelStateLightLevel  = "state_lightlevel"

	LabelConfigBattery   = "config_battery"
	LabelConfigOn        = "config_on"
	LabelConfigReachable = "config_reachable"
)

// InitMetrics func
func (exporter *Exporter) InitMetrics() (metrics []*metric.Metric) {

	metrics = append(metrics, &metric.Metric{
		HueType: TypeLight,
		FqName:  "hue_light_info",
		Help:    "Non-numeric data, value is always 1",
		Labels: []string{
			LabelName,
			LabelID,
			LabelType,
			LabelModelID,
			LabelManufacturerName,
			LabelSwVersion,
			LabelSwConfigID,
			LabelUniqueID,
			LabelStateOn,
			LabelStateAlert,
			LabelStateBri,
			LabelStateCT,
			LabelStateReachable,
			LabelStateSaturation,
		},
	})

	metrics = append(metrics, &metric.Metric{
		HueType: TypeLight,
		FqName:  "hue_light_state",
		Help:    "light status (1=ON, 0=OFF)",
		Labels: []string{
			LabelName,
		},
		ResultKey: LabelStateOn,
	})

	metrics = append(metrics, &metric.Metric{
		HueType: TypeSensor,
		FqName:  "hue_sensor_info",
		Help:    "Non-numeric data, value is always 1",
		Labels: []string{
			LabelName,
			LabelID,
			LabelType,
			LabelModelID,
			LabelManufacturerName,
			LabelSwVersion,
			LabelUniqueID,
			LabelStateButtonEvent,
			LabelStateDaylight,
			LabelStateLastUpdated,
			LabelStateTemperature,
			LabelConfigBattery,
			LabelConfigOn,
			LabelConfigReachable,
		},
	})

	metrics = append(metrics, &metric.Metric{
		HueType: TypeSensor,
		FqName:  "hue_sensor_temperature",
		Help:    "temperature level celsius degree",
		Labels: []string{
			LabelName,
		},
		ResultKey: LabelStateTemperature,
	})

	metrics = append(metrics, &metric.Metric{
		HueType: TypeSensor,
		FqName:  "hue_sensor_lightlevel",
		Help:    "light level",
		Labels: []string{
			LabelName,
		},
		ResultKey: LabelStateLightLevel,
	})

	return metrics
}

// CollectAll available hue metrics
func CollectAll(url string, username string, fileName string) {

	bridge := hueAPI.New(url, username)

	jsonContent := []collectEntry{}

	sensorData, err := collectSensors(bridge)
	if err != nil {
		fmt.Sprintln(err)
		return
	}

	lightData, err := collectLights(bridge)
	if err != nil {
		fmt.Sprintln(err)
		return
	}

	for _, sensor := range sensorData {
		jsonContent = append(jsonContent, collectEntry{
			Type:   "sensor",
			Result: sensor,
		})
	}
	for _, light := range lightData {
		jsonContent = append(jsonContent, collectEntry{
			Type:   "light",
			Result: light,
		})
	}

	jsonString, err := json.MarshalIndent(jsonContent, "", "\t")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(jsonString))
	if fileName != "" {
		err = ioutil.WriteFile(fileName, jsonString, 0644)
		if err != nil {
			fmt.Printf("Failed writing JSON file '%s': %s\n", fileName, err.Error())
		}
	}
}

// Collect metrics
func (exporter *Exporter) Collect(metrics []*metric.Metric) (err error) {

	bridge := hueAPI.New(exporter.BaseURL, exporter.Username)

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

func collectSensors(bridge *hueAPI.Bridge) (sensorData []map[string]interface{}, err error) {

	sensors, err := bridge.GetSensors()
	if err != nil {
		return nil, fmt.Errorf("[error GetAllSensors()] '%v'", err)
	}
	for _, sensor := range sensors {
		result := make(map[string]interface{})
		result[LabelName] = sensor.Name
		result[LabelID] = sensor.ID
		result[LabelType] = sensor.Type
		result[LabelModelID] = sensor.ModelID
		result[LabelManufacturerName] = sensor.ManufacturerName
		result[LabelSwVersion] = sensor.SwVersion
		result[LabelUniqueID] = sensor.UniqueID

		//State
		for stateKey, stateValue := range sensor.State {
			result["state_"+stateKey] = stateValue
		}

		//Config
		for stateKey, stateValue := range sensor.Config {
			result["config_"+stateKey] = stateValue
		}

		sensorData = append(sensorData, result)
	}
	return sensorData, nil
}

func collectLights(bridge *hueAPI.Bridge) (lightData []map[string]interface{}, err error) {

	lights, err := bridge.GetLights()
	if err != nil {
		return nil, fmt.Errorf("[error GetAllLights()] '%v'", err)
	}

	for _, light := range lights {

		result := make(map[string]interface{})
		result[LabelName] = light.Name
		result[LabelID] = light.ID
		result[LabelType] = light.Type
		result[LabelModelID] = light.ModelID
		result[LabelManufacturerName] = light.ManufacturerName
		result[LabelSwVersion] = light.SwVersion
		result[LabelSwConfigID] = light.SwConfigID
		result[LabelUniqueID] = light.UniqueID

		// State
		result[LabelStateOn] = light.State.On
		result[LabelStateAlert] = light.State.Alert
		result[LabelStateBri] = light.State.Bri
		result[LabelStateColorMode] = light.State.ColorMode
		result[LabelStateCT] = light.State.Ct
		result[LabelStateReachable] = light.State.Reachable
		result[LabelStateSaturation] = light.State.Sat

		lightData = append(lightData, result)
	}
	return lightData, nil
}
