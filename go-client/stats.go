package main

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// RequestStats holds statistics for each endpoint
type RequestStats struct {
	URL          string
	TotalReqs    int64
	SuccessReqs  int64
	ErrorReqs    int64
	MinTime      time.Duration
	MaxTime      time.Duration
	TotalTime    time.Duration
	LastResponse time.Time
}

// StatsCollector manages statistics collection
type StatsCollector struct {
	mu    sync.RWMutex
	stats map[string]*RequestStats
}

func NewStatsCollector() *StatsCollector {
	return &StatsCollector{
		stats: make(map[string]*RequestStats),
	}
}

func (sc *StatsCollector) AddRequest(url string, duration time.Duration, success bool) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if sc.stats[url] == nil {
		sc.stats[url] = &RequestStats{
			URL:     url,
			MinTime: duration,
			MaxTime: duration,
		}
	}

	s := sc.stats[url]
	s.TotalReqs++
	s.TotalTime += duration
	s.LastResponse = time.Now()

	if success {
		s.SuccessReqs++
	} else {
		s.ErrorReqs++
	}

	if duration < s.MinTime {
		s.MinTime = duration
	}
	if duration > s.MaxTime {
		s.MaxTime = duration
	}
}

func (sc *StatsCollector) GetStats() map[string]*RequestStats {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	result := make(map[string]*RequestStats)
	for k, v := range sc.stats {
		// Copy struct to avoid race conditions
		statsCopy := *v
		result[k] = &statsCopy
	}
	return result
}

func (sc *StatsCollector) PrintTable() {
	stats := sc.GetStats()
	if len(stats) == 0 {
		return
	}

	// Clear screen and print header
	fmt.Print("\033[2J\033[H") // Clear screen
	fmt.Println("🔄 Live Load Testing Results")
	fmt.Println("============================================================")
	fmt.Printf("%-25s %-8s %-8s %-8s %-8s %-8s %-8s\n",
		"Endpoint", "Total", "Success", "Errors", "Min(ms)", "Max(ms)", "Avg(ms)")
	fmt.Println("------------------------------------------------------------")

	// Sort endpoints by name
	var endpoints []string
	for endpoint := range stats {
		endpoints = append(endpoints, endpoint)
	}
	sort.Strings(endpoints)

	totalReqs := int64(0)
	totalSuccess := int64(0)
	totalErrors := int64(0)

	for _, endpoint := range endpoints {
		s := stats[endpoint]
		avgTime := time.Duration(0)
		if s.TotalReqs > 0 {
			avgTime = s.TotalTime / time.Duration(s.TotalReqs)
		}

		// Shorten endpoint name for display
		displayEndpoint := endpoint
		if len(displayEndpoint) > 24 {
			displayEndpoint = "..." + displayEndpoint[len(displayEndpoint)-21:]
		}

		fmt.Printf("%-25s %-8d %-8d %-8d %-8.1f %-8.1f %-8.1f\n",
			displayEndpoint,
			s.TotalReqs,
			s.SuccessReqs,
			s.ErrorReqs,
			float64(s.MinTime.Nanoseconds())/1e6,
			float64(s.MaxTime.Nanoseconds())/1e6,
			float64(avgTime.Nanoseconds())/1e6,
		)

		totalReqs += s.TotalReqs
		totalSuccess += s.SuccessReqs
		totalErrors += s.ErrorReqs
	}

	fmt.Println("------------------------------------------------------------")
	fmt.Printf("%-25s %-8d %-8d %-8d\n", "TOTAL", totalReqs, totalSuccess, totalErrors)
	fmt.Printf("\nLast update: %s | Press Ctrl+C to stop\n",
		time.Now().Format("15:04:05"))
}

func (sc *StatsCollector) PrintSummary() {
	stats := sc.GetStats()
	if len(stats) == 0 {
		fmt.Println("No statistics collected yet.")
		return
	}

	fmt.Println("\n📊 Load Testing Summary")
	fmt.Println("=======================")

	totalReqs := int64(0)
	totalSuccess := int64(0)
	totalErrors := int64(0)

	for endpoint, s := range stats {
		avgTime := time.Duration(0)
		if s.TotalReqs > 0 {
			avgTime = s.TotalTime / time.Duration(s.TotalReqs)
		}

		successRate := float64(0)
		if s.TotalReqs > 0 {
			successRate = float64(s.SuccessReqs) / float64(s.TotalReqs) * 100
		}

		fmt.Printf("\n🔗 %s\n", endpoint)
		fmt.Printf("   Requests: %d (Success: %d, Errors: %d)\n",
			s.TotalReqs, s.SuccessReqs, s.ErrorReqs)
		fmt.Printf("   Success Rate: %.1f%%\n", successRate)
		fmt.Printf("   Response Time: Min=%.1fms, Max=%.1fms, Avg=%.1fms\n",
			float64(s.MinTime.Nanoseconds())/1e6,
			float64(s.MaxTime.Nanoseconds())/1e6,
			float64(avgTime.Nanoseconds())/1e6)

		totalReqs += s.TotalReqs
		totalSuccess += s.SuccessReqs
		totalErrors += s.ErrorReqs
	}

	overallSuccessRate := float64(0)
	if totalReqs > 0 {
		overallSuccessRate = float64(totalSuccess) / float64(totalReqs) * 100
	}

	fmt.Printf("\n🎯 Overall Results:\n")
	fmt.Printf("   Total Requests: %d\n", totalReqs)
	fmt.Printf("   Success Rate: %.1f%% (%d/%d)\n",
		overallSuccessRate, totalSuccess, totalReqs)
	fmt.Printf("   Error Rate: %.1f%% (%d/%d)\n",
		100-overallSuccessRate, totalErrors, totalReqs)
}
