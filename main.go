package main

import (
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
	"strconv"
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
			"Id",
			"SSID",
			"MAC",
		},
	)
)

type co2_struct struct {
	Id, Temp, Humidity, Ppm, FreeRAM int
	Mac, SSID                        string
}

var co2_data co2_struct

func send(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&co2_data)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()
	co2Queued.WithLabelValues("FreeRAM", strconv.Itoa(co2_data.Id), co2_data.SSID, co2_data.Mac).Add(float64(co2_data.FreeRAM))
	co2Queued.WithLabelValues("Temp", strconv.Itoa(co2_data.Id), co2_data.SSID, co2_data.Mac).Add(float64(co2_data.Temp))
	co2Queued.WithLabelValues("Humidity", strconv.Itoa(co2_data.Id), co2_data.SSID, co2_data.Mac).Add(float64(co2_data.Humidity))
	co2Queued.WithLabelValues("PPM", strconv.Itoa(co2_data.Id), co2_data.SSID, co2_data.Mac).Add(float64(co2_data.Ppm))
	log.Printf("%+v\n", co2_data)
}

func init() {
	prometheus.MustRegister(co2Queued)
}

func main() {
	http.Handle("/metrics", prometheus.Handler())
	http.HandleFunc("/send", send)
	log.Fatal(http.ListenAndServe(":8082", nil))
}
