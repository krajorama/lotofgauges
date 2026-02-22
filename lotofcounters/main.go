package main

import (
	"flag"
	"fmt"
	"log"
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

	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "example_counter",
			Help: "An example counter with dynamic labels",
		},
		[]string{"label_id"},
	)

	prometheus.MustRegister(counter)

	// Update counter values every second
	go func() {
		for {
			for i := 0; i < *n; i++ {
				counter.WithLabelValues(strconv.Itoa(i)).Inc()
			}
			time.Sleep(time.Second)
		}
	}()

	handler := promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{
		EnableOpenMetrics:                   true,
		EnableOpenMetricsTextCreatedSamples: true,
	})

	http.Handle("/metrics", handler)
	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Serving metrics on %s/metrics with %d label values", addr, *n)
	log.Fatal(http.ListenAndServe(addr, nil))
}
