package analysis

import (
	"fmt"
	"golang-projects/models"
	"strconv"
	"strings"
)

func LogAnalysis(entries []models.LogEntry) models.LogAnalysis {
	analysis := models.LogAnalysis{
		LevelCounts:     make(map[string]int),
		MethodCounts:    make(map[string]int),
		MethodDurations: make(map[string]float64),
	}

	for _, entry := range entries {
		if entry.Level != "" {
			analysis.LevelCounts[entry.Level]++
		}

		msgParts := strings.Split(entry.Message, ", ")
		var method string
		var duration float64

		for _, part := range msgParts {
			if strings.HasPrefix(part, "method: ") {
				method = strings.TrimPrefix(part, "method: ")
				analysis.MethodCounts[method]++
				analysis.TotalRequests++

			} else if strings.HasPrefix(part, "requestDuration: ") {
				if method != "" {
					durationStr := strings.TrimPrefix(part, "requestDuration: ")
					duration, _ = strconv.ParseFloat(durationStr, 64)
					analysis.MethodDurations[method] += duration
				}
			}
		}
	}

	// 평균 duration 계산
	for method, totalDuration := range analysis.MethodDurations {
		count := analysis.MethodCounts[method]
		if count > 0 {
			analysis.MethodDurations[method] = totalDuration / float64(count)
		}
	}

	return analysis
}

func PrintAnalysis(analysis models.LogAnalysis) {
	fmt.Println("Log Level Counts:")
	for level, count := range analysis.LevelCounts {
		fmt.Printf("%s: %d\n", level, count)
	}

	fmt.Println("\nMethod Request Counts:")
	for method, count := range analysis.MethodCounts {
		fmt.Printf("%s: %d\n", method, count)
	}

	fmt.Println("\nAverage Request Duration by Method (in milliseconds):")
	for method, avgDuration := range analysis.MethodDurations {
		fmt.Printf("%s: %.2f\n", method, avgDuration)
	}

	fmt.Printf("\nTotal Requests: %d\n", analysis.TotalRequests)
}
