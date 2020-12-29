# Philips Hue exporter for prometheus

This exporter exports some variables from Philips Hue Bridge 
(https://www.philips-hue.com)
to prometheus.

## Building

    go get github.com/aexel90/hue_exporter/
    cd $GOPATH/src/github.com/aexel90/hue_exporter
    go install

## Running

How to create a user for your bridge is described here: https://developers.meethue.com/develop/get-started-2/

Usage:

    $GOPATH/bin/hue_exporter -h

    Usage of ./hue_exporter:
        -collect
                Collect all available metrics
        -collect-file string
                The JSON file where to store collect results
        -hue-url string
                The URL of the bridge
        -listen-address string
                The address to listen on for HTTP requests. (default "127.0.0.1:9773")
        -metrics-file string
                The JSON file with the metric definitions. (default "hue_metrics.json")
        -test
                Test configured metrics
        -username string
                The username token having bridge access

## Example execution

### Running within prometheus:

    $GOPATH/bin/hue_exporter -hue_url 192.168.xxx.xxx -username ZlEH24zabK2jTpJ...  -metrics-file hue_metrics.json 

    # HELP hue_light_status status of lights registered at hue bridge
    # TYPE hue_light_status gauge
    hue_light_status{manufacturer_name="...",model_id="...",name="...",state_alert="...",state_bri="...",state_ct="...",state_on="...",state_reachable="...",state_saturation="...",sw_version="...",type="...",unique_id="..."} 1
    hue_light_status{manufacturer_name="...",model_id="...",name="...",state_alert="...",state_bri="...",state_ct="...",state_on="...",state_reachable="...",state_saturation="...",sw_version="...",type="...",unique_id="..."} 0
    ...
    # HELP hue_sensor status of sensors registered at hue bridge
    # TYPE hue_sensor gauge
    hue_sensor{config_battery="...",config_on="...",config_reachable="...",manufacturer_name="...",model_id="...",name="...",state_button_event="...",state_daylight="...",state_last_updated="...",state_last_updated_time="...",sw_version="...",type="...",unique_id="..."} 1
    hue_sensor{config_battery="...",config_on="...",config_reachable="...",manufacturer_name="...",model_id="...",name="...",state_button_event="...",state_daylight="...",state_last_updated="...",state_last_updated_time="...",sw_version="...",type="...",unique_id="..."} 0
    ...

### Test exporter:

    $GOPATH/bin/hue_exporter -hue_url 192.168.xxx.xxx -username ZlEH24zabK2jTpJ... -test -metrics-file hue_metrics.json

### Collect:

    $GOPATH/bin/hue_exporter -hue_url 192.168.xxx.xxx -username ZlEH24zabK2jTpJ... -collect -metrics-file hue_metrics.json -collect-file result.json

## Grafana Dashboard

Dashboard ID is 13645.

![Grafana](https://raw.githubusercontent.com/aexel90/hue_exporter/main/grafana/screenshot.jpg)