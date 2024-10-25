package processors

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"golang-projects/models"
	"golang-projects/parser"
	"os"
	"sync"
)

func customSplit(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		return i + 1, data[0:i], nil
	}
	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}

func ProcessGzipLogFile(filename string, resultChan chan<- models.LogEntry, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", filename, err)
		return
	}
	defer func(file *os.File) {
		closeErr := file.Close()
		if closeErr != nil {
			fmt.Printf("Error closing file %s: %v\n", filename, closeErr)
			panic(closeErr)
		}
	}(file)

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		fmt.Printf("Error creating gzip reader for %s: %v\n", filename, err)
		return
	}
	defer func(gzipReader *gzip.Reader) {
		readerErr := gzipReader.Close()
		if readerErr != nil {
			fmt.Printf("Error closing gzip reader for %s: %v\n", filename, readerErr)
			panic(readerErr)
		}
	}(gzipReader)

	scanner := bufio.NewScanner(gzipReader)
	scanner.Split(customSplit)
	bufferSize := 10 * 1024 * 1024 // 10MB 제한
	scanner.Buffer(make([]byte, bufferSize), bufferSize)

	for scanner.Scan() {
		entry, isFound := parser.ParseLogLine(scanner.Text())
		if isFound {
			resultChan <- entry
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file %s: %v\n", filename, err)
	}
}
