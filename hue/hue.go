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
	typeLight  = "light"
	typeSensor = "sensor"

	labelName             = "name"
	labelType             = "type"
	labelID               = "id"
	labelModelID          = "model_id"
	labelManufacturerName = "manufacturer_name"
	labelSwVersion        = "sw_version"
	labelSwConfigID       = "sw_config_id"
	labelUniqueID         = "unique_id"

	labelStateOn          = "state_on"
	labelStateAlert       = "state_alert"
	labelStateBri         = "state_bri"
	labelStateColorMode   = "state_color_mode"
	labelStateCT          = "state_ct"
	labelStateReachable   = "state_reachable"
	labelStateSaturation  = "state_saturation"
	labelStateButtonEvent = "state_buttonevent"
	labelStateDaylight    = "state_daylight"
	labelStateLastUpdated = "state_lastupdated"
	labelStateTemperature = "state_temperature"
	labelStateLightLevel  = "state_lightlevel"

	labelConfigBattery   = "config_battery"
	labelConfigOn        = "config_on"
	labelConfigReachable = "config_reachable"

	labelAPIVersion                  = "api_version"
	labelBridgeID                    = "bridge_id"
	labelIPAddress                   = "ip_address"
	labelInternetServiceInternet     = "internetservice_internet"
	labelInternetServiceRemoteAccess = "internetservice_remoteaccess"
	labelInternetServiceSwUpdate     = "internetservice_swupdate"
	labelInternetServiceTime         = "internetservice_time"
	labelLocalTime                   = "local_time"
	labelSwUpdate2LastChange         = "sw_update_last_change"
	labelZigbeeChannel               = "zigbee_channel"
)

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

	bridgeData, err := collectBridgeInfo(bridge)
	if err != nil {
		fmt.Sprintln(err)
		return
	}

	for _, sensor := range sensorData {
		jsonContent = append(jsonContent, collectEntry{Type: "sensor", Result: sensor})
	}
	for _, light := range lightData {
		jsonContent = append(jsonContent, collectEntry{Type: "light", Result: light})
	}
	for _, bridge := range bridgeData {
		jsonContent = append(jsonContent, collectEntry{Type: "bridge", Result: bridge})
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
		case typeLight:
			metric.MetricResult = lightData
		case typeSensor:
			metric.MetricResult = sensorData
		default:
			return fmt.Errorf("Type '%v' currently not supported", metric.HueType)
		}

	}

	return nil
}

func collectSensors(bridge *hueAPI.Bridge) (sensorData []map[string]interface{}, err error) {

	sensors, err := bridge.GetSensors()
	if err != nil {
		return nil, fmt.Errorf("[GetAllSensors()] '%v'", err)
	}
	for _, sensor := range sensors {
		result := make(map[string]interface{})
		result[labelName] = sensor.Name
		result[labelID] = sensor.ID
		result[labelType] = sensor.Type
		result[labelModelID] = sensor.ModelID
		result[labelManufacturerName] = sensor.ManufacturerName
		result[labelSwVersion] = sensor.SwVersion
		result[labelUniqueID] = sensor.UniqueID

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
		return nil, fmt.Errorf("[GetAllLights()] '%v'", err)
	}

	for _, light := range lights {

		result := make(map[string]interface{})
		result[labelName] = light.Name
		result[labelID] = light.ID
		result[labelType] = light.Type
		result[labelModelID] = light.ModelID
		result[labelManufacturerName] = light.ManufacturerName
		result[labelSwVersion] = light.SwVersion
		result[labelSwConfigID] = light.SwConfigID
		result[labelUniqueID] = light.UniqueID

		// State
		result[labelStateOn] = light.State.On
		result[labelStateAlert] = light.State.Alert
		result[labelStateBri] = light.State.Bri
		result[labelStateColorMode] = light.State.ColorMode
		result[labelStateCT] = light.State.Ct
		result[labelStateReachable] = light.State.Reachable
		result[labelStateSaturation] = light.State.Sat

		lightData = append(lightData, result)
	}
	return lightData, nil
}

func collectBridgeInfo(bridge *hueAPI.Bridge) (bridgeData []map[string]interface{}, err error) {

	config, err := bridge.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("[GetConfig] '%v'", err)
	}
	result := make(map[string]interface{})
	result[labelName] = config.Name
	result[labelAPIVersion] = config.APIVersion
	result[labelBridgeID] = config.BridgeID
	result[labelIPAddress] = config.IPAddress
	result[labelInternetServiceInternet] = config.InternetService.Internet
	result[labelInternetServiceRemoteAccess] = config.InternetService.RemoteAccess
	result[labelInternetServiceSwUpdate] = config.InternetService.SwUpdate
	result[labelInternetServiceTime] = config.InternetService.Time
	result[labelLocalTime] = config.LocalTime
	result[labelModelID] = config.ModelID
	result[labelSwVersion] = config.SwVersion
	result[labelSwUpdate2LastChange] = config.SwUpdate2.LastChange
	result[labelZigbeeChannel] = config.ZigbeeChannel

	bridgeData = append(bridgeData, result)
	return bridgeData, nil
}
