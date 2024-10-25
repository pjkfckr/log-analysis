package utils

import (
	"fmt"
	"golang-projects/models"
	"golang-projects/processors"
	"path/filepath"
	"sort"
	"sync"
)

func MultipleReadGzipFile(dirPath string) ([]models.LogEntry, error) {
	files, err := filepath.Glob(filepath.Join(dirPath, "*.gz"))
	if err != nil {
		return nil, fmt.Errorf("파일 목록 가져오기를 실패하였습니다: %v", err)
	}

	sort.Strings(files)

	resultChan := make(chan models.LogEntry, 100)
	var wg sync.WaitGroup

	for _, file := range files {
		fmt.Printf("Read File... %v\n", file)
		wg.Add(1)
		go processors.ProcessGzipLogFile(file, resultChan, &wg)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	var allEntries []models.LogEntry
	for entry := range resultChan {
		allEntries = append(allEntries, entry)
	}

	return allEntries, nil
}
