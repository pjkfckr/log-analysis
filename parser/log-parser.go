package parser

import (
	"golang-projects/models"
	"strings"
)

func ParseLogLine(line string) (models.LogEntry, bool) {
	entry := models.LogEntry{}

	// 타임스탬프 추출
	timestampEnd := strings.Index(line, "]")
	if timestampEnd != -1 {
		entry.Timestamp = strings.TrimSpace(line[:timestampEnd])
	}

	// 나머지 부분 처리
	if timestampEnd != -1 && len(line) > timestampEnd+1 {
		remainder := line[timestampEnd+1:]

		// 로그 레벨 추출
		logLevels := []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL", "TRACE"}
		for _, level := range logLevels {
			if idx := strings.Index(remainder, level); idx != -1 {
				entry.Level = level
				entry.Message = strings.TrimSpace(remainder[idx+len(level):])
				break
			}
		}

		if entry.Level == "" {
			entry.Message = strings.TrimSpace(remainder)
		}
	}

	if entry.Message == "" {
		return entry, false
	}

	return entry, true
}
