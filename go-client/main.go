package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	maxClients    = flag.Int("maxClients", 10, "Maximum number of virtual clients")
	scaleInterval = flag.Int("scaleInterval", 500, "Scale interval in milliseconds")
	randomSleep   = flag.Int("randomSleep", 1000, "Random sleep from 0 to target microseconds")
	baseURL       = flag.String("baseURL", "http://localhost:8000", "Base URL for the target server")
)

func main() {
	// Sleep for 5 seconds before running test
	time.Sleep(5 * time.Second)

	// Parse the command line into the defined flags
	flag.Parse()

	// Create Prometheus registry
	reg := prometheus.NewRegistry()
	m := NewMetrics(reg)

	// Create Prometheus HTTP server to expose metrics
	pMux := http.NewServeMux()
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	pMux.Handle("/metrics", promHandler)

	go func() {
		log.Printf("Starting client metrics server on port 8082")
		log.Fatal(http.ListenAndServe(":8082", pMux))
	}()

	// Create transport and client to reuse connection pool
	tr := &http.Transport{}
	client := &http.Client{Transport: tr}

	// Create job queue
	var ch = make(chan string, *maxClients*2)
	var wg sync.WaitGroup

	// Slowly increase the number of virtual clients
	for clients := 0; clients <= *maxClients; clients++ {
		wg.Add(1)

		for i := 0; i < clients; i++ {
			go func() {
				for {
					url, ok := <-ch
					if !ok {
						wg.Done()
						return
					}
					sendReq(m, client, url)
				}
			}()
		}

		doWork(ch, clients)

		// Sleep for one second and increase number of clients
		time.Sleep(time.Duration(*scaleInterval) * time.Millisecond)
	}
}

func doWork(ch chan string, clients int) {
	// Define different endpoints to test
	endpoints := []string{
		*baseURL + "/health",
		*baseURL + "/api/devices",
		*baseURL + "/api/images",
	}

	if clients == *maxClients {
		for {
			// Randomly select an endpoint
			endpoint := endpoints[rand.Intn(len(endpoints))]
			ch <- endpoint
			sleep(*randomSleep)
		}
	}

	for i := 0; i < clients; i++ {
		// Randomly select an endpoint
		endpoint := endpoints[rand.Intn(len(endpoints))]
		ch <- endpoint
		sleep(*randomSleep)
	}
}

func sleep(us int) {
	r := rand.Intn(us)
	time.Sleep(time.Duration(r) * time.Microsecond)
}
