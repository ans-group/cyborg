package main

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
