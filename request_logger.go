package main

import "time"

type RequestLogEntry struct {
	Worker      int           `json:"worker"`
	StatusCode  int           `json:"status_code"`
	RequestTime time.Duration `json:"request_time"`
	Error       error         `json:"error"`
}

type RequestLogger interface {
	Log(entry RequestLogEntry)
}

type RequestLoggerManager struct {
	loggers []RequestLogger
}

func NewRequestLoggerManager() *RequestLoggerManager {
	return &RequestLoggerManager{}
}

func (l *RequestLoggerManager) AddLogger(logger RequestLogger) {
	l.loggers = append(l.loggers, logger)
}

func (l *RequestLoggerManager) Log(entry RequestLogEntry) {
	for _, logger := range l.loggers {
		go logger.Log(entry)
	}
}
