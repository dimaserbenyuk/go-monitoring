package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// MetricData represents parsed metric information
type MetricData struct {
	Endpoint string
	Status   string
	Quantile string
	Value    float64
	Count    int64
	Sum      float64
}

// StatsViewer handles fetching and displaying current statistics
type StatsViewer struct {
	baseURL    string
	metricsURL string
}

func NewStatsViewer(baseURL string) *StatsViewer {
	return &StatsViewer{
		baseURL:    baseURL,
		metricsURL: "http://localhost:8082/metrics",
	}
}

func (sv *StatsViewer) ShowCurrentStats() error {
	fmt.Println("📊 Current Load Testing Statistics")
	fmt.Println("==================================")
	fmt.Printf("Target: %s\n", sv.baseURL)
	fmt.Printf("Timestamp: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	// Check if metrics server is running
	resp, err := http.Get(sv.metricsURL)
	if err != nil {
		fmt.Println("❌ No active load test found. Start with: ./client")
		fmt.Println("💡 Tip: Run load test first, then check stats in another terminal")
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read metrics: %v", err)
	}

	metrics := sv.parseMetrics(string(body))
	sv.displayStats(metrics)

	// Also show server-side app stats if available
	sv.showServerStats()

	return nil
}

func (sv *StatsViewer) parseMetrics(metricsText string) map[string]*MetricData {
	metrics := make(map[string]*MetricData)

	// Parse tester_request_duration_seconds metrics
	lines := strings.Split(metricsText, "\n")
	for _, line := range lines {
		if strings.Contains(line, "tester_request_duration_seconds") &&
			!strings.HasPrefix(line, "#") {

			// Extract endpoint, status, quantile, and value
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				metricName := parts[0]
				value, _ := strconv.ParseFloat(parts[1], 64)

				// Parse labels
				endpoint := sv.extractLabel(metricName, "path")
				status := sv.extractLabel(metricName, "status")
				quantile := sv.extractLabel(metricName, "quantile")

				key := endpoint + "|" + status
				if metrics[key] == nil {
					metrics[key] = &MetricData{
						Endpoint: endpoint,
						Status:   status,
					}
				}

				if strings.Contains(metricName, "_count") {
					metrics[key].Count = int64(value)
				} else if strings.Contains(metricName, "_sum") {
					metrics[key].Sum = value
				} else if quantile == "0.9" {
					metrics[key].Quantile = "0.9"
					metrics[key].Value = value
				}
			}
		}
	}

	return metrics
}

func (sv *StatsViewer) extractLabel(metricName, labelName string) string {
	pattern := labelName + `="([^"]+)"`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(metricName)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func (sv *StatsViewer) displayStats(metrics map[string]*MetricData) {
	if len(metrics) == 0 {
		fmt.Println("❌ No metrics found. Load test may not be running.")
		return
	}

	fmt.Println("🎯 Load Test Results:")
	fmt.Println("---------------------")
	fmt.Printf("%-35s %-8s %-10s %-10s %-10s\n",
		"Endpoint", "Requests", "P90(ms)", "Avg(ms)", "RPS")
	fmt.Println(strings.Repeat("-", 80))

	totalRequests := int64(0)
	for _, metric := range metrics {
		if metric.Count > 0 {
			avgMs := (metric.Sum / float64(metric.Count)) * 1000
			p90Ms := metric.Value * 1000

			// Estimate RPS (very rough calculation)
			rps := float64(metric.Count) / 60.0 // Assume 1 minute window

			// Shorten endpoint for display
			endpoint := metric.Endpoint
			if len(endpoint) > 34 {
				endpoint = "..." + endpoint[len(endpoint)-31:]
			}

			fmt.Printf("%-35s %-8d %-10.1f %-10.1f %-10.1f\n",
				endpoint, metric.Count, p90Ms, avgMs, rps)

			totalRequests += metric.Count
		}
	}

	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("Total Requests: %d\n", totalRequests)
}

func (sv *StatsViewer) showServerStats() {
	fmt.Println("\n🖥️  Server Statistics:")
	fmt.Println("---------------------")

	// Try to get server stats
	resp, err := http.Get(sv.baseURL + "/api/stats")
	if err != nil {
		fmt.Println("❌ Server stats not available")
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("❌ Failed to read server stats")
		return
	}

	// Pretty print JSON
	var jsonData interface{}
	if err := json.Unmarshal(body, &jsonData); err == nil {
		prettyJSON, _ := json.MarshalIndent(jsonData, "", "  ")
		fmt.Println(string(prettyJSON))
	} else {
		fmt.Println(string(body))
	}
}
