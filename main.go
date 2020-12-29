package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/namsral/flag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/aexel90/hue_exporter/collector"
	"github.com/aexel90/hue_exporter/hue"
	"github.com/aexel90/hue_exporter/metric"
)

var (
	flagBridgeURL   = flag.String("hue-url", "", "The URL of the bridge")
	flagUsername    = flag.String("username", "", "The username token having bridge access")
	flagAddress     = flag.String("listen-address", "127.0.0.1:9773", "The address to listen on for HTTP requests.")
	flagMetricsFile = flag.String("metrics-file", "metrics.json", "The JSON file with the metric definitions.")

	flagTest        = flag.Bool("test", false, "test configured metrics")
	flagCollect     = flag.Bool("collect", false, "test configured metrics")
	flagCollectFile = flag.String("collect-file", "", "The JSON file where to store collect results")
)

func main() {

	flag.Parse()

	var metricsFile *metric.MetricsFile

	// collect mode
	if *flagCollect {
		hue.CollectAll(*flagBridgeURL, *flagUsername, *flagCollectFile)
		return
	}

	err := readAndParseFile(*flagMetricsFile, &metricsFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	hueCollector, err := collector.NewHueCollector(*flagBridgeURL, *flagUsername, metricsFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	if *flagTest {
		hueCollector.Test()
	} else {
		prometheus.MustRegister(hueCollector)
		http.Handle("/metrics", promhttp.Handler())
		fmt.Printf("metrics available at http://%s/metrics\n", *flagAddress)
		log.Fatal(http.ListenAndServe(*flagAddress, nil))
	}
}

func readAndParseFile(file string, v interface{}) error {
	jsonData, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("error reading metric file: %v", err)
	}

	err = json.Unmarshal(jsonData, v)
	if err != nil {
		return fmt.Errorf("error parsing JSON: %v", err)
	}
	return nil
}
