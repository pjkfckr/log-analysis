package models

type LogAnalysis struct {
	LevelCounts     map[string]int
	MethodCounts    map[string]int
	MethodDurations map[string]float64
	TotalRequests   int
}
