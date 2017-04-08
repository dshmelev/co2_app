package main

import (
	"bytes"
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	co2Queued = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "co2",
			Subsystem: "metric",
			Name:      "queued",
			Help:      "CO2 Realtime metrics",
		},
		[]string{
			"metric",
			"SSID",
			"MAC",
		},
	)
)

type co2_struct struct {
	Counter, Temp, Humidity, Ppm, FreeRAM, Rssi float64
	Mac, SSID                                   string
}

var co2_data co2_struct

func send(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(req.Body)
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	bodyString := string(bodyBytes)

	log.Println(bodyString)
	err := json.Unmarshal(bodyBytes, &co2_data)
	if err != nil {
		panic(err)
	}
	co2Queued.WithLabelValues("FreeRAM", co2_data.SSID, co2_data.Mac).Set(co2_data.FreeRAM)
	co2Queued.WithLabelValues("Temp", co2_data.SSID, co2_data.Mac).Set(co2_data.Temp)
	co2Queued.WithLabelValues("Humidity", co2_data.SSID, co2_data.Mac).Set(co2_data.Humidity)
	co2Queued.WithLabelValues("PPM", co2_data.SSID, co2_data.Mac).Set(co2_data.Ppm)
	co2Queued.WithLabelValues("RSSI", co2_data.SSID, co2_data.Mac).Set(co2_data.Rssi)
	co2Queued.WithLabelValues("Counter", co2_data.SSID, co2_data.Mac).Add(1)
}

func init() {
	prometheus.MustRegister(co2Queued)
}

func main() {
	log.Println("Running on port :8082")
	http.Handle("/metrics", prometheus.Handler())
	http.HandleFunc("/send", send)
	log.Fatal(http.ListenAndServe(":8082", nil))
}
