package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/namsral/flag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/aexel90/hue_exporter/collector"
	"github.com/aexel90/hue_exporter/hue"
)

var (
	flagBridgeURL = flag.String("hue-url", "", "The URL of the bridge")
	flagUsername  = flag.String("username", "", "The username token having bridge access")
	flagAddress   = flag.String("listen-address", "127.0.0.1:9773", "The address to listen on for HTTP requests.")

	flagTest        = flag.Bool("test", false, "test configured metrics")
	flagCollect     = flag.Bool("collect", false, "test configured metrics")
	flagCollectFile = flag.String("collect-file", "", "The JSON file where to store collect results")
)

func main() {

	flag.Parse()

	// collect mode
	if *flagCollect {
		hue.CollectAll(*flagBridgeURL, *flagUsername, *flagCollectFile)
		return
	}

	hueCollector, err := collector.NewHueCollector(*flagBridgeURL, *flagUsername)
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
