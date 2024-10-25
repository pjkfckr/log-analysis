package main

import (
	"fmt"
	"golang-projects/analysis"
	"golang-projects/utils"
)

func main() {
	dirPath := "logs"
	entries, err := utils.MultipleReadGzipFile(dirPath)
	if err != nil {
		fmt.Printf("Error processing log files: %v\n", err)
		return
	}

	results := analysis.LogAnalysis(entries)
	analysis.PrintAnalysis(results)

}
