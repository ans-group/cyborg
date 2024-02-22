package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type RequestLoggerStdout struct{}

func (l *RequestLoggerStdout) Log(entry RequestLogEntry) {
	if entry.Error != nil {
		log.Printf("[Worker %d] Request failure: %s", entry.Worker, entry.Error)
	} else {
		log.Printf("[Worker %d] Got result with status code [%s] in [%s]", entry.Worker, formatStatusCodeString(entry.StatusCode), entry.RequestTime)
	}
}

func formatStatusCodeString(code int) string {
	if !viper.GetBool("colour") {
		return fmt.Sprintf("%d", code)
	}

	color := 31
	if code >= 200 && code <= 399 {
		color = 32
	}

	return fmt.Sprintf("\033[1;%dm%d\033[0m", color, code)
}
