package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	n := flag.Int("n", 10, "number of label values (0 to N-1)")
	port := flag.Int("port", 8080, "port to serve metrics on")
	flag.Parse()

	gauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "example_gauge",
			Help: "An example gauge with dynamic labels",
		},
		[]string{"label_id"},
	)

	prometheus.MustRegister(gauge)

	// Update gauge values every second
	go func() {
		for {
			for i := 0; i < *n; i++ {
				gauge.WithLabelValues(strconv.Itoa(i)).Set(rand.Float64() * 100)
			}
			time.Sleep(time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Serving metrics on %s/metrics with %d label values", addr, *n)
	log.Fatal(http.ListenAndServe(addr, nil))
}
