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
	flagMetricsFile = flag.String("metrics-file", "hue_metrics.json", "The JSON file with the metric definitions.")
	authUser        = flag.String("auth.user", "", "Username for basic auth.")
	authPass        = flag.String("auth.pass", "", "Password for basic auth. Enables basic auth if set.")

	flagTest        = flag.Bool("test", false, "Test configured metrics")
	flagCollect     = flag.Bool("collect", false, "Collect all available metrics")
	flagCollectFile = flag.String("collect-file", "", "The JSON file where to store collect results")
)

type basicAuthHandler struct {
	handler  http.HandlerFunc
	user     string
	password string
}

func (h *basicAuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, password, ok := r.BasicAuth()
	if !ok || password != h.password || user != h.user {
		w.Header().Set("WWW-Authenticate", "Basic realm=\"metrics\"")
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}
	h.handler(w, r)
	return
}

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
		handler := promhttp.Handler()
		if *authUser != "" || *authPass != "" {
			if *authUser == "" || *authPass == "" {
				log.Fatal("You need to specify -auth.user and -auth.pass to enable basic auth")
			}
			handler = &basicAuthHandler{
			        handler:  promhttp.Handler().ServeHTTP,
			        user:     *authUser,
			        password: *authPass,
		        }
	        }
		http.Handle("/metrics", handler)
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
