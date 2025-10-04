package main

import (
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// sendReq sends request to the server and collects statistics
func sendReq(m *metrics, statsCollector *StatsCollector, client *http.Client, url string) {
	// Sleep to avoid sending requests at the same time.
	rn := rand.Intn(*scaleInterval)
	time.Sleep(time.Duration(rn) * time.Millisecond)

	// Get timestamp for histogram
	now := time.Now()

	// Send a request to the server
	res, err := client.Get(url)
	duration := time.Since(now)

	if err != nil {
		// Record error in both metrics and statistics
		m.duration.With(prometheus.Labels{"path": url, "status": "500"}).Observe(duration.Seconds())
		statsCollector.AddRequest(url, duration, false)
		log.Printf("client.Get failed: %v", err)
		return
	}

	// Read until the response is complete to reuse connection
	io.ReadAll(res.Body)

	// Close the body to reuse connection
	res.Body.Close()

	// Determine if request was successful
	success := res.StatusCode >= 200 && res.StatusCode < 300

	// Record request duration in Prometheus metrics
	m.duration.With(prometheus.Labels{"path": url, "status": strconv.Itoa(res.StatusCode)}).Observe(duration.Seconds())

	// Record request in statistics collector
	statsCollector.AddRequest(url, duration, success)
}
